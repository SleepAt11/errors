package errors

import (
	"fmt"
	"runtime"
)

type stackError struct {
	s      string
	stacks []*stack
	params map[stackID][]*stackErrorParam
}

type stackID string

type stackErrorParam struct {
	k string
	v interface{}
}

type stack struct {
	file     string
	funcName string
	line     int
}

func (e *stackError) Error() string {
	return e.s
}

func Error(err error) string {
	e, ok := err.(*stackError)
	if !ok {
		return err.Error()
	}
	return fmt.Sprintf("err:%s\n", e.s) + e.stackInfo()
}

func With(err error, k string, v interface{}) *stackError {
	e, ok := err.(*stackError)
	if !ok {
		e = New(err)
	}
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return e
	}
	f := runtime.FuncForPC(pc)
	s := stack{
		file:     file,
		funcName: f.Name(),
		line:     line,
	}
	stackId := e.getStackID(&s)
	param := stackErrorParam{
		k: k,
		v: v,
	}
	currentParams := e.params[stackId]
	currentParams = append(currentParams, &param)
	e.params[stackId] = currentParams
	return e
}

func (e *stackError)With(k string, v interface{}) *stackError {
	// skip 1 to get With's caller
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return e
	}
	f := runtime.FuncForPC(pc)
	s := stack{
		file:     file,
		funcName: f.Name(),
		line:     line,
	}
	stackId := e.getStackID(&s)
	param := stackErrorParam{
		k: k,
		v: v,
	}
	currentParams := e.params[stackId]
	currentParams = append(currentParams, &param)
	e.params[stackId] = currentParams
	return e
}

func New(err error) *stackError {
	currentErr, ok := err.(*stackError)
	if ok {
		return currentErr
	}
	e := stackError{s: err.Error()}
	e.params = map[stackID][]*stackErrorParam{}
	var stacks []*stack
	// i range from 1，omit function stack-info-errors.New itself.
	for i := 1; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		f := runtime.FuncForPC(pc)
		s := stack{}
		notSysCall := f.Name() != "runtime.main" && f.Name() != "runtime.goexit"
		// New function may be called by With function, when err param's type is not a *stackError.
		notWithCall := f.Name() != "stack-info-errors/errors.With"
		if notSysCall && notWithCall {
			s.file = file
			s.line = line
			s.funcName = f.Name()
		}
		if s.funcName != "" {
			stacks = append(stacks, &s)
		}
	}
	e.stacks = stacks
	return &e
}

func (e *stackError)stackInfo() string {
	info := "stack info: \n"
	for idx, stack := range e.stacks {
		info += fmt.Sprintf("%d.file:%s:%d, func:%s\n", idx + 1, stack.file, stack.line, stack.funcName)
		stackId := e.getStackID(stack)
		params, ok := e.params[stackId]
		if !ok {
			continue
		}
		info += "  params:\n"
		for idx, param := range params {
			info += fmt.Sprintf("\t%s:%+v", param.k, param.v)
			if idx != len(params) - 1 {
				info += "\n"
			}
		}
		info += "\n"
	}
	return info
}

func (e *stackError)getStackID(s *stack) stackID {
	if s == nil {
		return ""
	}
	id := fmt.Sprintf("%s-%s", s.file, s.funcName)
	return stackID(id)
}