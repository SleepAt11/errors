package main

import (
	syserrors "errors"
	"fmt"
	"stack-info-errors/errors"
)

const errorInfo = "here is an error."

func main()  {
	err := f1()
	fmt.Println(errors.Error(err))
}

func f1() error {
	err := f2()
	return errors.New(err).With("ids", []string{"4", "5", "6"}).With("name", "a name")
}

func f2() error {
	err := f3()
	return errors.With(err, "ids", []string{"1","2","3"}).With("name", "another name")
}

func f3() error {
	err := f4()
	return err
}

func f4() error {
	err := syserrors.New(errorInfo)
	return errors.New(err)
}