package model

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
)

// Takes a string of the form 2024-01-01 and extracts the year value as an int
// Checks that the year number is in the range 2018 .. 2023
//
// panics on error
func PastYearFromDate(dateStr string) PastYear {
	bits := strings.Split(dateStr, "-")
	if len(bits) != 3 {
		fmt.Printf("PastYearFromDate %s\n", dateStr)
		panic("PastYearFromDate failed")
	}
	y := bits[0]
	if ynum, ok := strconv.Atoi(y); ok == nil {
		if ynum >= 2018 && ynum <= 2023 {
			return PastYear(ynum)
		} else {
			fmt.Printf("PastYearFromDate %d\n", ynum)
			panic("PastYearFromDate out of range failed")
		}
	} else {
		panic("PastYearFromDate Atoi failed")
	}
}

// Generic - Extracts the string keys from a value of type map[string]T
// and sorts the keys
func SortedMapKeys[T any](m map[string]T) []string {
	keys := slices.Collect(maps.Keys(m))
	sort.Strings(keys)
	return keys
}

// Load Population projection data into an instance of type PopulationProjectionInputData
func PopulationForecastInputData_LoadFromFile(fileName string) PopulationForecastInputData {
	b, err := os.ReadFile(fileName)
	if err != nil {
		panic("failed to open file")
	}
	// var jsonList JsonList
	var target []JsonRecord
	json.Unmarshal([]byte(b), &target)
	var byYear map[string][]JsonRecord = make(map[string][]JsonRecord, 1)
	fmt.Println("Json marshall complete", len(target))
	for _, v := range target {
		// fmt.Printf("k: %d v: %v\n", k, v)
		_, ok := byYear[v.Date]
		if !ok {
			byYear[v.Date] = make([]JsonRecord, 0)
			byYear[v.Date] = append(byYear[v.Date], v)
		} else {
			byYear[v.Date] = append(byYear[v.Date], v)
		}
	}
	sortedKeys := SortedMapKeys(byYear)
	forecastData := PopulationForecastInputData_CreateEmpty()

	for _, k := range sortedKeys {
		v, ok := byYear[k]
		if ok {
			ytmp := PastYearFromDate(k)
			yearData := YearPopulationProjection_CreateEmpty(PastYear(ytmp))
			fmt.Printf("k: %s v: \n", k)
			for _, r := range v {
				fmt.Printf("\t %v \n", r)
				yearData.AddYearAge(ytmp, r.Age, r.Value)
				// tmp := YearAgePopulationProjection{
				// 	year:            int(ytmp),
				// 	age:             r.Age,
				// 	totalPopulation: r.Value,
				// }
				// if r.Age > yage.maxAge {
				// 	yage.maxAge = r.Age
				// }
				// yage.popByAge = append(yage.popByAge, tmp)
			}
			forecastData.AddYearData(yearData)
			// forecastData.projectionByYear = append(forecastData.projectionByYear, yearData)
		} else {
			panic("something went wrong")
		}
	}
	return forecastData
}
