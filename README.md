# stack-info-errors

## example
main.go 代码：
```
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
	err := syserrors.New(errorInfo)
	return errors.New(err)
}
```

输出：

```
err:here is an error.
stack info: 
1.file:/Users/zengnifeng/go/src/stack-info-errors/main.go:33, func:main.f4
2.file:/Users/zengnifeng/go/src/stack-info-errors/main.go:27, func:main.f3
3.file:/Users/zengnifeng/go/src/stack-info-errors/main.go:22, func:main.f2
  params:
        ids:[1 2 3]
        name:another name
4.file:/Users/zengnifeng/go/src/stack-info-errors/main.go:17, func:main.f1
  params:
        ids:[4 5 6]
        name:a name
5.file:/Users/zengnifeng/go/src/stack-info-errors/main.go:12, func:main.main

```

## 名词约定

* go errors包下的errorString类型，简称为`errorString类型`。

* 带参数的栈信息错误日志输出，简称为`栈信息错误日志`；error接口的Error()方法返回的错误信息，简称为`原始错误日志`。

* `stackError类型`为本项目中实际error类型，对外不可见。

## 接口列表


#### 1.func (e *stackError) Error() string

与errorString类型的实现一致。打印原始error string。


#### 2.func Error(err error) string 

* 若err为stackError类型，则返回栈信息错误日志。

* 否则，则返回原始错误日志。

#### 3.func With(err error, k string, v interface{}) *stackError

* 若err不是stackError类型，则会先生成stackError并带上调用栈信息。

* 给err带上参数信息，并对应到当前栈。

#### 4.func (e *stackError)With(k string, v interface{}) *stackError

与接口3类似，用于*stackError内部类型的链式调用。

#### 5.New(err error) *stackError 

将error转化为stackError类型时使用。若需同时带上参数信息，请直接调用接口3，无需先New。

## 特性

### 生成完整的带参数栈信息

日志格式：
```
err:here is an error.
stack info: 
(调用栈层次).file:(文件路径), func:(函数名), line:(行数)
  params:
        (参数名1):(参数值1)
        (参数名2):(参数值2)
        ...

```

此处函数名准确地说应该是selector，带有包名前缀，是一个可唯一确定的函数符号。


### 栈信息以最深层次的生成为准

以下两个方法在无栈信息时会以当前栈调用开始生成调用栈信息：

* New

* With（接口3）

若在调用栈的深层和浅层均有生成栈信息的方法调用，将会以深层为准只生成一次调用栈信息。

> 第三方框架`github.com/pkg/errors`会在每次withStack时重复生成，导致调用栈信息大量冗余。

### 无侵入式

与errorString保持兼容，不影响原有error日志，可渐进式修改现有error逻辑。

对于errorString类型的err：

```
	err := syserrors.New(errorInfo) //syserrors是go的errors包
```

1.参考example可见，要获取`栈信息错误日志`需显式调用`errors.Error(err)`。直接打印err将返回`原始错误日志`。

2.显式调用`errors.Error(err)`时，若err并非stackError类型，则返回`原始错误日志`。

总之，stackError类型对外不可见，对外仍可视作普通的errorString类型，不影响原有的任何输出和处理逻辑。


