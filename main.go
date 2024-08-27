package main

import (
	"os"
	"printer/app"
)

func main() {
	s := &app.Service{
		Name: "aaa_service_printV2",
		Port: "8888",
	}

	if len(os.Args) > 1 && os.Args[1] == "d" {
		s.RunDev()
	} else {
		s.RunProd()
	}
}
