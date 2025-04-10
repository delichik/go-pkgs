package pipeline

import (
	"fmt"
	"testing"
)

type testAutoFillStruct struct {
	s1 testStruct1
	s2 *testStruct2
	s3 testStruct3
}

func invokeAutoFill(a *testAutoFillStruct) {
	fmt.Println("invoke", a.s3.T)
}

func TestWrappers(t *testing.T) {
	v := &testStruct1{
		T: "111",
	}
	NewPipeline().
		Invoke(invokeAutoFill).
		ProvideValue(v).
		ProvideStruct(testAutoFillStruct{}).
		Provide(testFunc1).
		Provide(testFunc2).
		Provide(testFunc3).
		Prepare().
		Run()
}
