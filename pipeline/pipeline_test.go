package pipeline

import (
	"fmt"
	"testing"
)

type testStruct1 struct {
	T string
}

func (t *testStruct1) Call() string { return t.T }

func (t *testStruct1) Call2() string { return t.T }

type testStruct2 struct {
	T string
}

func (t *testStruct2) Call() string { return t.T }

type testStruct3 struct {
	T string
}

type testInterface interface {
	Call() string
	Call2() string
}

type testConfusedInterface interface {
	Call() string
}

func invoke(a *testStruct3) {
	fmt.Println("invoke", a.T)
}

func testFunc1() (b *testStruct1) {
	fmt.Println("testFunc1")
	return &testStruct1{
		T: "111",
	}
}

func testFunc2(a testStruct1) (b testStruct2) {
	fmt.Println("testFunc2", a.T)
	return testStruct2{
		T: a.T + "222",
	}
}

func testFunc3(a *testStruct1, b *testStruct2) (c *testStruct3) {
	fmt.Println("testFunc3", a.T, b.T)
	return &testStruct3{
		T: a.T + b.T + "333",
	}
}

func TestPipeline(t *testing.T) {
	NewPipeline().
		Invoke(invoke).
		Provide(testFunc3).
		Provide(testFunc1).
		Provide(testFunc2).
		Prepare().
		Run()
}

func testFuncInterface1() (b *testStruct1) {
	fmt.Println("testFunc1")
	return &testStruct1{
		T: "111",
	}
}

func testFuncInterface2(a testInterface) (b testStruct2) {
	fmt.Println("testFunc2", a.Call())
	return testStruct2{
		T: a.Call() + "222",
	}
}

func testFuncConfusedInterface2(a testConfusedInterface) (b testStruct2) {
	fmt.Println("testFunc2", a.Call())
	return testStruct2{
		T: a.Call() + "222",
	}
}

func testFuncInterface3(a testInterface, b *testStruct2) (c *testStruct3) {
	fmt.Println("testFunc3", a.Call(), b.T)
	return &testStruct3{
		T: a.Call() + b.T + "333",
	}
}

func TestPipelineInterface(t *testing.T) {
	NewPipeline().
		Invoke(invoke).
		Provide(testFuncInterface1).
		Provide(testFuncInterface2).
		Provide(testFuncInterface3).
		Prepare().
		Run()
}

func TestPipelineConfusedInterface(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
		} else {
			t.Failed()
		}
	}()
	NewPipeline().
		Invoke(invoke).
		Provide(testFuncInterface1).
		Provide(testFuncConfusedInterface2).
		Provide(testFuncInterface3).
		Prepare().
		Run()
}
