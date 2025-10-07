package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/pflag"
	"l2.17/pkg/client"
)

var (
	// flag -h/--host - хост TCP сервера.
	host string

	// flag -p/--port - порт TCP сервера
	port string

	// flag -t/--timeout - таймаут на подключение в секундах.
	timeout int
)

func init() {
	pflag.StringVarP(&host, "host", "h", "", "TCP server host")
	pflag.StringVarP(&port, "port", "p", "", "TCP server port")
	pflag.IntVarP(&timeout, "timeout", "t", 10, "connection timeout in seconds")
}

func main() {
	pflag.Parse()
	if host == "" {
		fmt.Println("missing value: host")
		os.Exit(1)
	}
	if port == "" {
		fmt.Println("missing value: port")
		os.Exit(1)
	}

	tcpConnection, err := client.ConnectTCP(host, port, time.Duration(timeout)*time.Second)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = tcpConnection.Relay(context.Background(), os.Stdin, os.Stdout)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
