package errors

import (
	"errors"
	"fmt"
	"testing"
)

const errorInfo = "here is an error."

// 此处不是真正的单元测试，仅用以查看输出
func TestErrors(t *testing.T) {
	err := f1()
	fmt.Println(Error(err))
}

func f1() error {
	err := f2()
	return New(err).With("ids", []string{"4", "5", "6"}).With("name", "a name")
}

func f2() error {
	err := f3()
	return With(err, "ids", []string{"1","2","3"}).With("name", "another name")
}

func f3() error {
	err := f4()
	return err
}

func f4() error {
	err := errors.New(errorInfo)
	return New(err)
}