package plugin

type Handler func(name string, data []byte) ([]byte, error)

type Parasite interface {
	Init() error
	UnInit()
	Handle(data []byte) ([]byte, error)
}

type Executor interface {
	OnCall(call string, data []byte) ([]byte, error)
}
