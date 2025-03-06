package main

import (
	"fmt"
	"os"

	"http_tester/http_server/controller"
	"http_tester/log"
)

func main() {
	defer func() {
		if log.BaseLogger != nil {
			log.BaseLogger.Sync()
		}
	}()
	if err := controller.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
