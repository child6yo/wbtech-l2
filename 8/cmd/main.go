package main

import (
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

func getNetworkTime(ntpAddr string) (time.Time, error) {
	response, err := ntp.Query(ntpAddr)
	if err != nil {
		return time.Time{}, err
	}

	time := response.Time

	return time, nil
}

func main() {
	time, err := getNetworkTime("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(time.Local())
}
