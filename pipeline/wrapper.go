package pipeline

func PureParam[T any](t T) func() T {
	return func() T {
		return t
	}
}
