package csvdata

import (
	"encoding/csv"
	"io"
	"os"
	"strings"
)

type Covid struct {
	Positive      string `json:"Positive"`
	Tests  	      string `json:"Tests"`
	Date          string `json:"Date"`
	Discharged    string `json:"Discharged"`
	Expired       string `json:"Expired"`
	Admitted      string `json:"Admitted"`
	Region        string `json:"Region"`
}

type CovidDataRequest struct {
	Get string `json:"get"`
}

type CovidDataError struct {
	Error string `json:"covid_error"`
}


func Load(path string) []Covid {
	table := make([]Covid, 0)
	csvFile, err := os.Open(path)
	if err != nil {
		panic(err.Error())
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err.Error())
		}
		c := Covid{
			Positive:      row[0],
			Tests: 	       row[1],
			Date:          row[2],
			Discharged:    row[3],
			Expired:       row[4],
			Admitted:      row[5],
			Region:        row[6],
		}
		table = append(table, c)
	}
	return table
}

func Find(table []Covid, filter string) []Covid {
	if filter == "" || filter == "*" {
		return table
	}
	result := make([]Covid, 0)
	filter = strings.ToUpper(filter)
	for _, cov := range table {
		if cov.Region == filter ||
			cov.Date == filter ||
			strings.Contains(strings.ToUpper(cov.Region), filter) ||
			strings.Contains(strings.ToUpper(cov.Date), filter) {
			result = append(result, cov)
		}
	}
	return result
}


