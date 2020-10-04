package web

import (
	"github.com/RemiSirdata/covid19-hubbert-curve-generator/pkg"
	"github.com/gin-gonic/gin"
	"log"
)

var (
	franceRawData = pkg.RawData.GetCountry("france")
)

func NewHttpServer(port string) {
	r := gin.Default()
	r.GET("/country/:country", handlerGetNewCases)
	r.GET("/country/:country/region/:region", handlerGetNewCases)
	if err := r.Run(port); err != nil {
		log.Fatalf("fail to start web server on %s: %s", port, err.Error())
	}
}
