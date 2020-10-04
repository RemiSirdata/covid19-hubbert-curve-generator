package pkg

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
	"time"
)

func TestGetFranceRawDataCsvUrl(t *testing.T) {
	r := rawFranceData{}
	u, err := r.getRawCsvUrl()
	assert.Nil(t, err, "get raw csv url shouldn't return err")
	assert.NotNil(t, u, "return csv url shouldn't be nil")
	assert.NotEmpty(t, u, "return csv url shouldn't be empty")
	_, err = url.Parse(*u)
	assert.Nil(t, err, "parse raw csv url shouldn't return err")
}

func TestGetFranceRawData(t *testing.T) {
	r, err := newRawFranceData()
	assert.Nil(t, err, "init shouldn't return error")
	if err != nil {
		return
	}
	newCases, err := r.GetDailyNewHospitalisations(time.Now().AddDate(0, 0, -21), time.Now().AddDate(0, 0, -7))
	assert.Nil(t, err, "GetDailyNewHospitalisations shouldn't return error")
	assert.NotNil(t, newCases, "GetDailyNewHospitalisations newCases shouldn't be nil")
	assert.NotEmpty(t, newCases, "GetDailyNewHospitalisations new cases shouldn't be empty")

	newCases, err = r.GetDailyNewHospitalisationsForRegion("75", time.Now().AddDate(0, 0, -21), time.Now().AddDate(0, 0, -7))
	assert.Nil(t, err, "GetDailyNewHospitalisationsForRegion shouldn't return error")
	assert.NotNil(t, newCases, "GetDailyNewHospitalisationsForRegion newCases shouldn't be nil")
	assert.NotEmpty(t, newCases, "GetDailyNewHospitalisationsForRegion new cases shouldn't be empty")
}
