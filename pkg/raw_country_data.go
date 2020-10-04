package pkg

import "time"

type RawCountryData interface {
	GetDailyNewHospitalisations(start time.Time, end time.Time) (*DailyNewCases, error)
	GetDailyNewHospitalisationsForRegion(region string, start time.Time, end time.Time) (*DailyNewCases, error)
}

type Gender string

var (
	genderMale   Gender = "MALE"
	genderFemale Gender = "FEMALE"
)

type DailyNewCases struct {
	Start time.Time
	End   time.Time
	Data  []DailyNewCase
}

type DailyNewCase struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type DailyNewCaseByRegion struct {
	Date   string `json:"date"`
	Region string `json:"region"`
	Count  int    `json:"count"`
}
