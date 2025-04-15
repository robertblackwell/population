package mockdb

import (
	"encoding/json"
	"fmt"
	"os"
)

type JsonProjectedPopulationRecord struct {
	Value  int    `json:"value"`
	Code   string `json:"code"`
	Type   string `json:"type"`
	Date   string `json:"date"`
	Age    int    `json:"age"`
	Gender string `json:"gender"`
}
type JsonCurrentPopulationRecord struct {
	Value  int    `json:"value"`
	Code   string `json:"code"`
	Type   string `json:"type"`
	Year   int    `json:"year"`
	Age    int    `json:"age"`
	Gender string `json:"gender"`
}

type JsonProjectedList struct {
	Records []JsonProjectedPopulationRecord `json:"records"`
}
type JsonCurrentList struct {
	Records []JsonProjectedPopulationRecord `json:"records"`
}

type Context struct {
	Dummy int
}

// Loads data from a json file (prop-db.json) in this directory and transforms the data into a form like
//
//	map[string]map[string][]JsonRecord
//
// where the string index of the 1st map is 'Code' values and second index id date
//
//	code      age      date
func LoadProjectedMockDb(fName string) map[string]map[int]map[string]JsonProjectedPopulationRecord {
	home_path := os.Getenv("HOME")
	p := fmt.Sprintf("%s/%s/%s/%s", home_path, "Projects/popmodel", "mockdb", fName)
	fileName := p
	b, err := os.ReadFile(fileName)
	if err != nil {
		panic(fmt.Errorf("failed to open file: %s", p))
	}
	var target []JsonProjectedPopulationRecord
	json.Unmarshal([]byte(b), &target)
	fmt.Println("Json marshall complete", len(target))
	//            code       age        date
	m := make(map[string]map[int]map[string]JsonProjectedPopulationRecord, 0)
	for _, jr := range target {
		code := jr.Code
		_, ok := m[code]
		if !ok {
			m[code] = make(map[int]map[string]JsonProjectedPopulationRecord, 0)
		}
		_, ok = m[code][jr.Age]
		if !ok {
			m[code][jr.Age] = make(map[string]JsonProjectedPopulationRecord, 0)
		}
		_, ok = m[code][jr.Age][jr.Date]
		if ok {
			panic(fmt.Errorf("something went wrong code:%s age:%d date:%s allready exists", jr.Code, jr.Age, jr.Date))
			// m[code][jr.Age][jr.Date] = make([]JsonRecord, 0)
		}
		m[code][jr.Age][jr.Date] = jr
	}
	return m
}
func LoadCurrentMockDb(fName string) map[string]map[int]map[int]JsonCurrentPopulationRecord {
	home_path := os.Getenv("HOME")
	p := fmt.Sprintf("%s/%s/%s/%s", home_path, "Projects/popmodel", "mockdb", fName)
	fileName := p
	b, err := os.ReadFile(fileName)
	if err != nil {
		panic(fmt.Errorf("failed to open file: %s", p))
	}
	var target []JsonCurrentPopulationRecord
	json.Unmarshal([]byte(b), &target)
	fmt.Println("Json marshall complete", len(target))
	//            code       age        date
	m := make(map[string]map[int]map[int]JsonCurrentPopulationRecord, 0)
	for _, jr := range target {
		code := jr.Code
		_, ok := m[code]
		if !ok {
			m[code] = make(map[int]map[int]JsonCurrentPopulationRecord, 0)
		}
		_, ok = m[code][jr.Age]
		if !ok {
			m[code][jr.Age] = make(map[int]JsonCurrentPopulationRecord, 0)
		}
		_, ok = m[code][jr.Age][jr.Year]
		if ok {
			panic(fmt.Errorf("something went wrong code:%s age:%d year:%d allready exists", jr.Code, jr.Age, jr.Year))
			// m[code][jr.Age][jr.Date] = make([]JsonRecord, 0)
		}
		m[code][jr.Age][jr.Year] = jr
	}
	return m
}
