package mockdb

import (
	"encoding/json"
	"fmt"
	"os"
)

type JsonRecord struct {
	Value  int    `json:"value"`
	Code   string `json:"code"`
	Type   string `json:"type"`
	Date   string `json:"date"`
	Age    int    `json:"age"`
	Gender string `json:"gender"`
}

type JsonList struct {
	Records []JsonRecord `json:"records"`
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
func LoadMockDb() map[string]map[int]map[string]JsonRecord {
	home_path := os.Getenv("HOME")
	p := fmt.Sprintf("%s/%s/%s/%s", home_path, "Projects/popmodel", "mockdb", "prop-db.json")
	fileName := p
	b, err := os.ReadFile(fileName)
	if err != nil {
		panic("failed to open file")
	}
	var target []JsonRecord
	json.Unmarshal([]byte(b), &target)
	fmt.Println("Json marshall complete", len(target))
	//            code       age        date
	m := make(map[string]map[int]map[string]JsonRecord, 0)
	for _, jr := range target {
		code := jr.Code
		_, ok := m[code]
		if !ok {
			m[code] = make(map[int]map[string]JsonRecord, 0)
		}
		_, ok = m[code][jr.Age]
		if !ok {
			m[code][jr.Age] = make(map[string]JsonRecord, 0)
		}
		_, ok = m[code][jr.Age][jr.Date]
		if ok {
			panic("something went wrong")
			// m[code][jr.Age][jr.Date] = make([]JsonRecord, 0)
		}
		m[code][jr.Age][jr.Date] = jr
	}
	return m
}
