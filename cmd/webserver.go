package cmd

import (
	"flag"
	"github.com/RemiSirdata/covid19-hubbert-curve-generator/web"
)

func Webserver() {
	port := flag.String("port", ":8080", "Webserver port")
	flag.Parse()
	web.NewHttpServer(*port)
}
