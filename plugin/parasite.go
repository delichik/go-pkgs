package plugin

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/delichik/daf/logger"
	"go.uber.org/zap"
)

type Options struct {
	Name               string
	Version            string
	HostName           string
	HostMinimalVersion string
}

var registeredParasites = make(map[string]Parasite)

func RegisterHandler(name string, parasite Parasite) {
	registeredParasites[name] = parasite
}

func RunParasite(options *Options) {
	if options.Name == "" {
		panic("parasite name is required")
	}

	if options.Version == "" {
		panic("parasite version is required")
	}

	if options.HostName == "" {
		panic("parasite host name is required")
	}

	if options.HostMinimalVersion == "" {
		panic("parasite host minimal version is required")
	}

	handshake := ""
	flag.StringVar(&handshake, "h", "", "")
	flag.Parse()
	if !checkHandshake(handshake, options) {
		fmt.Printf("This executable binary is a parasite for %s %s+, Do not run it alone\n",
			options.HostName, options.HostMinimalVersion)
		os.Exit(1)
	}

	for name, parasite := range registeredParasites {
		fmt.Printf("Parasite %s is starting...\n", name)
		parasite.Init()
		fmt.Printf("Parasite %s is started\n", name)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		buf := bufio.NewReader(os.Stdin)
		for {
			req, err := read(buf)
			if err != nil {
				return
			}
			if ctx.Err() != nil {
				return
			}
			parasite, ok := registeredParasites[req.call]
			if !ok {
				send(os.Stdout, &sendObject{
					id:      req.id,
					call:    req.call + callReply,
					content: []byte(""),
				})
				continue
			}
			rsp, err := parasite.Handle(req.content)
			if err != nil {
				logger.Error("parasite handle failed", zap.String("call", req.call), zap.Error(err))
				send(os.Stdout, &sendObject{
					id:      req.id,
					call:    req.call + callReply,
					content: []byte(err.Error()),
				})
				continue
			}
			send(os.Stdout, &sendObject{
				id:      req.id,
				call:    req.call + callReply,
				content: rsp,
			})
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGQUIT)
	<-signalChan
	cancel()
	signal.Stop(signalChan)

	for _, parasite := range registeredParasites {
		parasite.UnInit()
	}
}

type logWriter struct {
	writer io.Writer
}

func (w *logWriter) Write(data []byte) (n int, err error) {
	req := &sendObject{
		id:      0,
		call:    callLogger,
		content: data,
	}
	err = send(w.writer, req)
	return len(data), err
}

func (w *logWriter) Sync() error {
	return nil
}

func InitParasiteLogger() {
	logger.InitDefaultManual(&logger.Config{
		Level:     "debug",
		Format:    "json",
		LogDriver: &logWriter{writer: os.Stdout},
	})
}
