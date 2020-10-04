package main

import (
	"flag"
	"github.com/RemiSirdata/covid19-hubbert-curve-generator/cmd"
	"log"
	"os"
)

func main() {
	flag.Parse()
	if len(os.Args) < 2 {
		log.Fatalf("args is required")
	}
	switch os.Args[1] {
	case "webserver":
		cmd.Webserver()
		break
	default:
		log.Fatalf("command not found")
	}
}
