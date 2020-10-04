package web

import (
	"github.com/RemiSirdata/covid19-hubbert-curve-generator/pkg"
	"github.com/gin-gonic/gin"
	"log"
)

var (
	rawData       = &pkg.RawData{}
	franceRawData pkg.RawCountryData
)

func NewHttpServer(port string) {
	var err error
	franceRawData, err = rawData.GetCountry("france")
	if err != nil {
		log.Fatalf("fail to retreive data: %s", err.Error())
	}

	r := gin.Default()
	r.GET("/country/:country", handlerGetNewCases)
	r.GET("/country/:country/region/:region", handlerGetNewCases)
	if err := r.Run(port); err != nil {
		log.Fatalf("fail to start web server on %s: %s", port, err.Error())
	}
}
