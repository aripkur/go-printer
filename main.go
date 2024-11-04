package main

import (
	"os"
	"printer/app"
)

func main() {
	s := &app.Service{
		Name: "aaaPDFprint",
		Port: "8888",
	}

	if len(os.Args) > 1 && os.Args[1] == "d" {
		s.RunDev()
	} else {
		s.RunProd()
	}
}
