package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/iamtraining/forum/apiserver"
	"github.com/iamtraining/forum/store"
	"github.com/iamtraining/forum/web"
	"github.com/pingcap/log"
)

var (
	address = flag.String(
		"address", ":3000",
		"address the server should bind to",
	)
	configFile = flag.String(
		"config", "configs/apiserver.toml",
		"configuration file for the application",
	)
	debug = flag.Bool(
		"debug", false,
		"enable debug logging",
	)
)

type Config struct {
	Address string
	App     apiserver.Config
}

var (
	config Config
)

func init() {
	flag.Parse()
	_, err := toml.DecodeFile(*configFile, &config)
	if err != nil {
		os.Stderr.WriteString("error: couldnt load configuration: " + err.Error())
		os.Exit(1)
	}
}

func main() {
	store, err := store.NewStore("postgres://postgres:1111@localhost/iamtraining?sslmode=disable")
	if err != nil {
		panic(err)
	}

	h := web.NewHandler(store)

	go func() {
		if err := h.App.Listen(*address); err != nil {
			fmt.Println("server error " + err.Error())
		}
	}()

	time.Sleep(time.Millisecond * 100)
	log.Info("server started")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	fmt.Println("server stopping", "why", <-sig)

	fmt.Println("goodbye")
}
