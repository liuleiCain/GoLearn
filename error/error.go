package learn_error

import (
	"fmt"
	"reflect"
)

type MyError struct {
	error
}

func returnsError() error {
	var e2 *MyError = nil
	return e2
}

func MyErrorEqualsNil() {
	// 内置error的nil
	var e error = nil
	fmt.Printf("isEqual: %v type: %s error: %v \n", e == nil, reflect.TypeOf(e), e)
	// 自定义error的nil
	var e1 *MyError = nil
	fmt.Printf("isEqual: %v type: %s error: %v \n", e1 == nil, reflect.TypeOf(e1), e1)
	// 自定义error的类型存在，值为nil
	e2 := returnsError()
	fmt.Printf("isEqual: %v type: %s error2: %v \n", e2 == nil, reflect.TypeOf(e2), e2)
	// 自定义error的类型存在，值也存在，内部变量为nil
	e3 := &MyError{error: nil}
	fmt.Printf("isEqual: %v type: %s error: %v \n", e3 == nil, reflect.TypeOf(e3), e3)
}
