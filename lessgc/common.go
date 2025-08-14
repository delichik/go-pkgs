package lessgc

const segmentCapacity = 32

type segment[T any] []T

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}
