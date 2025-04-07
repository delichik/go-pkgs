package pipeline

import "testing"

func TestPureParam(t *testing.T) {
	v := &testStruct1{
		T: "111",
	}
	NewPipeline().
		Invoke(invoke).
		Provide(PureParam(v)).
		Provide(testFunc2).
		Provide(testFunc3).
		Prepare().
		Run()
}
