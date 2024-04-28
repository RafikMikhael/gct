package main

import "flag"

func main() {
	var portNum = flag.Int("p", 8080, "port number to use for api requests")
	flag.Parse()
	api := &App{}
	api.Initialize(portNum)
	api.Run()
}
