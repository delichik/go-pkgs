package pipeline

import (
	"fmt"
	"testing"
)

type testStruct1 struct {
	T string
}

type testStruct2 struct {
	T string
}

type testStruct3 struct {
	T string
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
