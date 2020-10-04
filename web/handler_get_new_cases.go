package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func handlerGetNewCases(c *gin.Context) {
	country := c.Param("country")
	region := c.Param("region")
	start := time.Now().AddDate(0, -2, 0)
	end := time.Now().AddDate(0, 0, -9)
	newCases, err := franceRawData.GetDailyNewHospitalisations(start, end)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(newCases)
}
