package main

import (
	"github.com/delichik/daf/logger"

	"github.com/delichik/go-pkgs/plugin"
)

func main() {
	plugin.InitParasiteLogger()
	plugin.RegisterHandler("hello", &ExamplePlugin{})
	plugin.RunParasite(&plugin.Options{
		Name:               "example-plugin-parasite",
		Version:            "0.0.1",
		HostName:           "example-plugin-host",
		HostMinimalVersion: "0.0.1",
	})
}

type ExamplePlugin struct {
}

func (p *ExamplePlugin) Init() error {
	logger.Info("init parasite")
	return nil
}

func (p *ExamplePlugin) UnInit() {
	logger.Info("uninit parasite")
}

func (p *ExamplePlugin) Handle(data []byte) ([]byte, error) {
	return []byte("hello example-plugin-host"), nil
}
