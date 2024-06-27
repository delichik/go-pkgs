package main

import (
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/delichik/daf/logger"

	"github.com/delichik/go-pkgs/plugin"
)

func main() {
	logger.InitDefaultManual(&logger.Config{
		Level:   "debug",
		Format:  "json",
		LogPath: "stdout",
	})
	h := plugin.NewHost("example-plugin-host", "0.0.1", &executor{})
	err := h.Load("./example/parasite")
	if err != nil {
		panic(err)
	}

	go func() {
		time.Sleep(500 * time.Millisecond)
		logger.Info("Sending hello")
		rsp, err := h.Call("hello", []byte("hello example-plugin-parasite"))
		if err != nil {
			logger.Warn("err", zap.Error(err))
		} else {
			logger.Info("info", zap.ByteString("rsp", rsp))
		}
	}()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGQUIT)
	<-signalChan
	signal.Stop(signalChan)
}

type executor struct{}

func (e *executor) OnCall(call string, data []byte) ([]byte, error) {
	switch call {
	case "hello":
		logger.Info("hello from parasite", zap.ByteString("content", data))
		return []byte("hello from example-plugin-host"), nil
	}
	return nil, errors.New("unknown call")
}
