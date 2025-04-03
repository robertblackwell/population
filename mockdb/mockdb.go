package mockdb

import (
	"encoding/json"
	"fmt"
	"os"
)

// // Years are the years for which population forcasts have been made and from which we will build
// // a population forecasting model. Currently they are 2018 .. 2035
// type Year int

// const FirstPastYear = 2018
// const LastPastYear = 2035

// func IsValidYear(y int) bool {
// 	return y >= 2018 && y <= 2035
// }
// func YearFromInt(y int) Year {
// 	if IsValidYear(y) {
// 		return Year(y)
// 	}
// 	panic(fmt.Sprintf("YearFromInt y: %d", y))
// }

// type Age int

// func IsValidAge(a int) bool {
// 	return a >= 0 && a <= 100
// }

// func int2Age(a int) (Age, error) {
// 	if IsValidAge(a) {
// 		return Age(a), nil
// 	}
// 	return Age(0), errors.New("invalid age")
// }

// type AgeRange = struct {
// 	Start Age
// 	End   Age
// }

// func createAgeRange(start int, end int) (AgeRange, error) {
// 	if start <= end {
// 		s, err1 := int2Age(start)
// 		e, err2 := int2Age(end)
// 		if err1 == nil && err2 == nil {
// 			return AgeRange{Start: s, End: e}, nil
// 		}
// 	}
// 	return AgeRange{}, errors.New("invalid ages or range")
// }
// func AgeRangeToString(ar AgeRange) string {
// 	return fmt.Sprintf("%d-%d", int(ar.Start), int(ar.End))
// }

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

// type Population struct {
// 	AgeRange        string `db:"age_range"`
// 	Sex             string `db:"sex"`
// 	TotalPopulation int    `db:"total_population"`
// 	Year            int    `db:"year"`
// }

// type LadPopulationProjection struct {
// 	Code            string `db:"code"`
// 	Type            string `db:"type"`
// 	AgeRange        string `db:"age_range"`
// 	TotalPopulation int    `db:"total_population"`
// 	Year            int    `db:"year"`
// }

type Context struct {
	Dummy int
}

// Loads data from a json file (prop-db.json) in this directory and transforms the data into a form like
//
//	map[string][]JsonRecord
//
// where the string index of the map are 'Code' values
func LoadMockDb() map[string][]JsonRecord {
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
	m := make(map[string][]JsonRecord, 0)
	for _, jr := range target {
		code := jr.Code
		_, ok := m[code]
		if !ok {
			m[code] = make([]JsonRecord, 0)
		}
		m[code] = append(m[code], jr)
	}
	return m
}

// Gets the projected population total for the specified age range, from the projected_populations_v2 table for all years between 2018 and 2035.
// amalgamate (ie sum over ages within the range)
// amalgamate over codes - this is a mock function only one code is allowed and is ignored
// startYear, rangesize, futureOffset, includeIntermediaries are ignored.
// func GetProjectedPopulationByCodes(ctx Context, codes []string, startYear, rangeSize, minAge, maxAge, futureOffset int, includeIntermediates bool) ([]LadPopulationProjection, error) {
// 	target := LoadMockDb()
// 	valid_age_range := IsValidAge(minAge) && IsValidAge(maxAge) && minAge < maxAge
// 	result := make([]LadPopulationProjection, 0)
// 	if !valid_age_range {
// 		return result, errors.New("invalid age range")
// 	}
// 	age_range_string := fmt.Sprintf("%d-%d", minAge, maxAge)
// 	var byYear map[string][]JsonRecord = make(map[string][]JsonRecord, 1)
// 	fmt.Println("Json marshall complete", len(target))
// 	for _, v := range target {
// 		// fmt.Printf("k: %d v: %v\n", k, v)
// 		_, ok := byYear[v.Date]
// 		if !ok {
// 			byYear[v.Date] = make([]JsonRecord, 0)
// 			byYear[v.Date] = append(byYear[v.Date], v)
// 		} else {
// 			byYear[v.Date] = append(byYear[v.Date], v)
// 		}
// 	}
// 	sortedKeys := SortedMapKeys(byYear)

// 	for _, k := range sortedKeys {
// 		v, ok := byYear[k]
// 		if ok {
// 			ytmp := YearFromDate(k)
// 			fmt.Printf("k: %s v: \n", k)
// 			pop := 0
// 			for _, r := range v {
// 				fmt.Printf("\t %v \n", r)
// 				if r.Age >= minAge && r.Age <= maxAge {
// 					pop = pop + r.Value
// 				}
// 				// tmp := YearAgePopulationProjection{
// 				// 	year:            int(ytmp),
// 				// 	age:             r.Age,
// 				// 	totalPopulation: r.Value,
// 				// }
// 				// if r.Age > yage.maxAge {
// 				// 	yage.maxAge = r.Age
// 				// }
// 				// yage.popByAge = append(yage.popByAge, tmp)
// 			}
// 			pp := LadPopulationProjection{Code: codes[0], Type: v[0].Type, TotalPopulation: pop, Year: int(ytmp), AgeRange: age_range_string}
// 			result = append(result, pp)
// 		} else {
// 			panic("something went wrong")
// 		}
// 	}
// 	return result, nil
// }

// // Takes a string of the form 2024-01-01 and extracts the year value as an int
// // Checks that the year number is in the range 2018 .. 2023
// //
// // panics on error
// func YearFromDate(dateStr string) Year {
// 	bits := strings.Split(dateStr, "-")
// 	if len(bits) != 3 {
// 		fmt.Printf("YearFromDate %s\n", dateStr)
// 		panic("YearFromDate failed")
// 	}
// 	y := bits[0]
// 	if ynum, ok := strconv.Atoi(y); ok == nil {
// 		if ynum >= 2018 && ynum <= 2035 {
// 			return Year(ynum)
// 		} else {
// 			fmt.Printf("YearFromDate %d\n", ynum)
// 			panic("YearFromDate out of range failed")
// 		}
// 	} else {
// 		panic("YearFromDate Atoi failed")
// 	}
// }

// // Generic - Extracts the string keys from a value of type map[string]T
// // and sorts the keys
// func SortedMapKeys[T any](m map[string]T) []string {
// 	keys := slices.Collect(maps.Keys(m))
// 	sort.Strings(keys)
// 	return keys
// }
