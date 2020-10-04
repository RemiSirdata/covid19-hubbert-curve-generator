package pkg

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-gota/gota/dataframe"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	franceRawDataSource = "https://www.data.gouv.fr/fr/datasets/donnees-hospitalieres-relatives-a-lepidemie-de-covid-19/"
)

func newRawFranceData() (*rawFranceData, error) {
	r := rawFranceData{}
	if err := r.init(); err != nil {
		return nil, err
	}
	return &r, nil
}

type rawFranceData struct {
	dataframe dataframe.DataFrame
	DB        *sql.DB
}

func (r *rawFranceData) init() error {
	var err error
	r.DB, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Printf("fail to initialize sqlite db: %s", err.Error())
		return err
	}
	if err := r.loadNewHospitalisationData(); err != nil {
		return err
	}
	return nil
}

func (r *rawFranceData) GetDailyNewHospitalisations(start time.Time, end time.Time) (*DailyNewCases, error) {
	rows, err := r.DB.Query("SELECT date(jour) day, SUM(incid_hosp) sum FROM new_hospitalisation WHERE jour >=? AND jour <= ? GROUP BY day ORDER BY day", start, end)
	if err != nil {
		log.Printf("fail to exec GetDailyNewHospitalisations query: %s", err.Error())
		return nil, err
	}
	newCase := []DailyNewCase{}
	for rows.Next() {
		var dnc DailyNewCase
		if err := rows.Scan(&dnc.Date, &dnc.Count); err != nil {
			log.Printf("fail to scan row: %s", err.Error())
		}
		newCase = append(newCase, dnc)
	}
	return &DailyNewCases{Start: start, End: end, Data: newCase}, nil
}

// "departement" is considered as a region for France
func (r *rawFranceData) GetDailyNewHospitalisationsForRegion(region string, start time.Time, end time.Time) (*DailyNewCases, error) {
	rows, err := r.DB.Query("SELECT date(jour) day, SUM(incid_hosp) sum FROM new_hospitalisation WHERE dep=? AND jour >=? AND jour <= ? GROUP BY day ORDER BY day", region, start, end)
	if err != nil {
		log.Printf("fail to exec GetDailyNewHospitalisations query: %s", err.Error())
		return nil, err
	}
	newCase := []DailyNewCase{}
	for rows.Next() {
		var dnc DailyNewCase
		if err := rows.Scan(&dnc.Date, &dnc.Count); err != nil {
			log.Printf("fail to scan row: %s", err.Error())
		}
		newCase = append(newCase, dnc)
	}
	return &DailyNewCases{Start: start, End: end, Data: newCase}, nil
}

func (r *rawFranceData) loadNewHospitalisationData() error {
	if _, err := r.DB.Exec("CREATE TABLE new_hospitalisation(dep INT, jour TEXT, incid_hosp INT, incid_rea INT, incid_dc INT, incid_rad INT);"); err != nil {
		return fmt.Errorf("fail to create table new_hospitalisation: %s", err.Error())
	}
	csvUrl, err := r.getRawCsvUrl()
	httpClient := http.Client{
		Timeout: time.Minute,
	}
	res, err := httpClient.Get(*csvUrl)
	if err != nil {
		return fmt.Errorf("fail to get csv: %s", err.Error())
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("fail to get csv data %s, status code error: %d %s", *csvUrl, res.StatusCode, res.Status)
	}
	csvReader := csv.NewReader(res.Body)
	csvReader.Comma = ';'

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(record) != 6 || record[0] == "dep" {
			// log here?
			continue
		}
		t, err := time.Parse(dateFormat, record[1])
		if err != nil {
			log.Printf("invalid time format %s: %s", record[1], err.Error())
			continue
		}
		record[1] = t.Format(time.RFC3339)
		if _, err := r.DB.Exec("INSERT INTO new_hospitalisation(dep, jour, incid_hosp, incid_rea, incid_dc, incid_rad) VALUES(?, ?, ?, ?, ?, ?)", record[0], record[1], record[2], record[3], record[4], record[5]); err != nil {
			return err
		}
	}
	return nil
}

func (r *rawFranceData) getRawCsvUrl() (*string, error) {
	httpClient := http.Client{
		Timeout: time.Minute,
	}
	res, err := httpClient.Get(franceRawDataSource)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fail to get page %s, status code error: %d %s", franceRawDataSource, res.StatusCode, res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	// Find the review items
	rawJson := doc.Find("#json_ld").First().Text()
	if rawJson == "" {
		return nil, fmt.Errorf("fail to get json data from page %s", franceRawDataSource)
	}
	var rawJsonStruct struct {
		Distribution []struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"distribution"`
	}
	if err := json.Unmarshal([]byte(rawJson), &rawJsonStruct); err != nil {
		return nil, fmt.Errorf("fail to decode raw data json: %s", err.Error())
	}
	for _, file := range rawJsonStruct.Distribution {
		if strings.Contains(file.Name, "donnees-hospitalieres-nouveaux-covid19") {
			return &file.Url, nil
		}
	}
	return nil, fmt.Errorf("raw file not found")
}
