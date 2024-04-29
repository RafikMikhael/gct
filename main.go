package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var portNum = flag.Int("p", 8080, "port number to use for api requests")
	flag.Parse()
	api := &App{}
	errInit := api.Initialize(portNum)
	if errInit != nil {
		fmt.Printf("initialization failed:%s\n", errInit.Error())
		os.Exit(1)
	}
	api.Run()
}
