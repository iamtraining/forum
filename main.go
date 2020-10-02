package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
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
)

type Config struct {
	Address string
	App     apiserver.Config
}

var (
	quit   = make(chan os.Signal, 1)
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

	h := web.NewHandler(store, config.App)

	go func() {
		if err := h.App.Listen(*address); err != nil && err != http.ErrServerClosed {
			fmt.Println("server error " + err.Error())
		}
	}()

	time.Sleep(time.Millisecond * 100)
	log.Info("server started")

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	fmt.Println(" server stopping.", "why?", <-quit)

	if err := h.App.Shutdown(); err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("goodbye")
}
