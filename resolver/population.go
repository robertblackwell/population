package resolver

import (
	"errors"
	"fmt"
	"forecast_model/mockdb"
	"math"
	"slices"
)

// copied from graphql project
type Population struct {
	AgeRange        string `db:"age_range"`
	Sex             string `db:"sex"`
	TotalPopulation int    `db:"total_population"`
	Year            int    `db:"year"`
}

// copied from graphql project
type LadPopulationProjection struct {
	Code            string `db:"code"`
	Type            string `db:"type"`
	AgeRange        string `db:"age_range"`
	TotalPopulation int    `db:"total_population"`
	Year            int    `db:"year"`
}

type GrowthRate = struct {
	Year             int
	RateFromBaseYear float64
}

type EstimatedPopulation = struct {
	Year       int
	Population int
}

// Get the base year projected population for a list of codes
func GetProjectedPopulationsByCodeForBaseYear(ctx mockdb.Context, codes []string, baseYear int, minAge int, maxAge int) (map[string]LadPopulationProjection, error) {
	result := make(map[string]LadPopulationProjection, 0)
	asSlice, e := GetProjectedPopulationsByCodeForYears(ctx, codes, []int{baseYear}, minAge, maxAge)
	if e != nil {
		return result, e
	}
	for k, _ := range asSlice {
		result[k] = asSlice[k][0]
	}
	return result, nil
}

// Get projected population data for a set of codes for a set of years
func GetProjectedPopulationsByCodeForYears(ctx mockdb.Context, codes []string, years []int, minAge int, maxAge int) (map[string][]LadPopulationProjection, error) {
	// get projected pop data for all years and then filter out the ones we want
	all, err := GetAllProjectedPopulationsByCodes(ctx, codes, minAge, maxAge)
	result := make(map[string][]LadPopulationProjection, 0)
	if err != nil {
		return result, err
	}
	for k, v := range all {
		for _, v := range v {
			if slices.Contains(years, v.Year) {
				_, e := result[k]
				if !e {
					result[k] = make([]LadPopulationProjection, 0)
				}
				result[k] = append(result[k], v)
			}
		}
	}
	return result, nil
}

// This is a mock version of a function that gets population project data.
//
// // It has the correct interface
//
// Will need to be re-implemented to use the sql database
func GetAllProjectedPopulationsByCodes(ctx mockdb.Context, codes []string, minAge int, maxAge int) (map[string][]LadPopulationProjection, error) {
	target := mockdb.LoadMockDb()
	valid_age_range := IsValidAge(minAge) && IsValidAge(maxAge) && minAge < maxAge
	map_result := make(map[string][]LadPopulationProjection, 0)
	if !valid_age_range {
		return map_result, fmt.Errorf("invalid age range %d %d", minAge, maxAge)
	}
	age_range_string := fmt.Sprintf("%d-%d", minAge, maxAge)
	var byYear map[string][]mockdb.JsonRecord = make(map[string][]mockdb.JsonRecord, 1)
	fmt.Println("Json marshall complete", len(target))

	for _, kv := range codes {

		code_target, ok := target[kv]
		if !ok {
			return map_result, fmt.Errorf("could not find codes %s in json database", kv)
		}

		_, b := map_result[kv]
		if !b {
			map_result[kv] = make([]LadPopulationProjection, 0)
		}

		for _, v := range code_target {
			// fmt.Printf("k: %d v: %v\n", k, v)
			_, ok := byYear[v.Date]
			if !ok {
				byYear[v.Date] = make([]mockdb.JsonRecord, 0)
				byYear[v.Date] = append(byYear[v.Date], v)
			} else {
				byYear[v.Date] = append(byYear[v.Date], v)
			}
		}
		sortedKeys := SortedMapKeys(byYear)

		for _, k := range sortedKeys {
			v, ok := byYear[k]
			if ok {
				ytmp := YearFromDate(k)
				fmt.Printf("k: %s v: \n", k)
				pop := 0
				for _, r := range v {
					fmt.Printf("\t %v \n", r)
					if r.Age >= minAge && r.Age <= maxAge {
						pop = pop + r.Value
					}
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
				pp := LadPopulationProjection{Code: codes[0], Type: v[0].Type, TotalPopulation: pop, Year: int(ytmp), AgeRange: age_range_string}
				map_result[kv] = append(map_result[kv], pp)
			} else {
				panic("something went wrong")
			}
		}
	}
	return map_result, nil
}

// Calculates the growth rates from the data in the projected_populations_v2 table
// For each code a separate calculation is performed.
// The returned map is indexed by code and the value for a code is a slice/array of LaPopulationProjection object
// In these arrays/slices there is a growth rate for every year between FirstYear(2018) and LastYear(2035)
// Growth rate values are numbers like 1.02 for a 2% positive growth and 0.98 for a 2% reduction.
// For each year the growth rate is the projected population for that year divided by the projected population
// for the baseYear (usually 2022)
// The growth rate for baseYear is 1.0
func GrowthRatesAllYearsMultipleCodes(projected_pops map[string][]LadPopulationProjection, baseYear int) map[string][]GrowthRate {
	result := make(map[string][]GrowthRate, 0)
	for k, v := range projected_pops {
		gr := GrowthRatesAllYears(v, baseYear)
		result[k] = gr
	}
	return result
}

// Apply growth rates to a population for the base year. Do this for a number of codes.
// Filter the output so that the returned value only has forrecast populations for the specified list of years (between FirstYear and LastYear)
func CalculateEstimatedPopulationsForSomeYearsMultipleCodes(growthRates map[string][]GrowthRate, baseYearPopulation map[string]int, requiredYears []int) (map[string][]EstimatedPopulation, error) {
	result := make(map[string][]EstimatedPopulation, 0)
	for k, v := range growthRates {
		basePop, ok := baseYearPopulation[k]
		if !ok {
			return result, fmt.Errorf("a base year population was not provided for code: %s", k)
		}
		ep := CalculateEstimatedPopulationsForSomeYears(v, basePop, requiredYears)
		result[k] = ep
	}
	return result, nil
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// From this point down are function private to this file
// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func GrowthRatesAllYears(projected_pops []LadPopulationProjection, baseYear int) []GrowthRate {
	base_year_index, err := findBaseYearIndexInPopulationProjectsions(projected_pops, 2023)
	if err != nil {
		panic("CalculateGrowthRatesAllYears: could not find base year in projected_pops")
	}
	growthRates := make([]GrowthRate, 0)
	for ix, pp := range projected_pops {
		var p float64
		if ix == base_year_index {
			p = 1.0
		} else if ix < base_year_index {
			p = float64(pp.TotalPopulation) / float64(projected_pops[base_year_index].TotalPopulation)
		} else {
			p = float64(pp.TotalPopulation) / float64(projected_pops[base_year_index].TotalPopulation)
		}
		growthRates = append(growthRates, GrowthRate{Year: pp.Year, RateFromBaseYear: p})
	}
	return growthRates
}

func CalculateEstimatedPopulationsForSomeYears(growthRates []GrowthRate, baseYearPopulation int, requiredYears []int) []EstimatedPopulation {

	result := make([]EstimatedPopulation, 0)
	for _, gr := range growthRates {
		if slices.Contains(requiredYears, gr.Year) {
			p := int(math.Round(float64(baseYearPopulation) * gr.RateFromBaseYear))
			result = append(result, EstimatedPopulation{Year: gr.Year, Population: p})
		}
	}
	return result
}

func findBaseYearIndexInPopulationProjectsions(pop_projections []LadPopulationProjection, base_year int) (int, error) {
	for i, el := range pop_projections {
		if el.Year == base_year {
			return i, nil
		}
	}
	return 0, errors.New("could not find 2023")
}

// This is a mock version of a function written for test purposes.
//
// In the live system the analogous functions quesries the populations_projection_v2 table to get:Gets the projected population total for the specified age range, from the projected_populations_v2 table for all years between 2018 and 2035.
// amalgamate (ie sum over ages within the range)
// amalgamate over codes - this is a mock function only one code is allowed and is ignored
// startYear, rangesize, futureOffset, includeIntermediaries are ignored.
func GetProjectedPopulationByCodes(ctx mockdb.Context, codes []string, startYear, rangeSize, minAge, maxAge, futureOffset int, includeIntermediates bool) (map[string][]LadPopulationProjection, error) {
	target := mockdb.LoadMockDb()
	valid_age_range := IsValidAge(minAge) && IsValidAge(maxAge) && minAge < maxAge
	map_result := make(map[string][]LadPopulationProjection, 0)
	if !valid_age_range {
		return map_result, errors.New("invalid age range")
	}
	age_range_string := fmt.Sprintf("%d-%d", minAge, maxAge)
	var byYear map[string][]mockdb.JsonRecord = make(map[string][]mockdb.JsonRecord, 1)
	fmt.Println("Json marshall complete", len(target))

	for _, kv := range codes {

		code_target, ok := target[kv]
		if !ok {
			return map_result, fmt.Errorf("could not find codes %s in json database", kv)
		}

		_, b := map_result[kv]
		if !b {
			map_result[kv] = make([]LadPopulationProjection, 0)
		}

		for _, v := range code_target {
			// fmt.Printf("k: %d v: %v\n", k, v)
			_, ok := byYear[v.Date]
			if !ok {
				byYear[v.Date] = make([]mockdb.JsonRecord, 0)
				byYear[v.Date] = append(byYear[v.Date], v)
			} else {
				byYear[v.Date] = append(byYear[v.Date], v)
			}
		}
		sortedKeys := SortedMapKeys(byYear)

		for _, k := range sortedKeys {
			v, ok := byYear[k]
			if ok {
				ytmp := YearFromDate(k)
				fmt.Printf("k: %s v: \n", k)
				pop := 0
				for _, r := range v {
					fmt.Printf("\t %v \n", r)
					if r.Age >= minAge && r.Age <= maxAge {
						pop = pop + r.Value
					}
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
				pp := LadPopulationProjection{Code: codes[0], Type: v[0].Type, TotalPopulation: pop, Year: int(ytmp), AgeRange: age_range_string}
				map_result[kv] = append(map_result[kv], pp)
			} else {
				panic("something went wrong")
			}
		}
	}
	return map_result, nil
}

// find the projected population for the base year
func FindBaseYearProjectedPopulation(projected_pops []LadPopulationProjection, baseYear int) (int, error) {
	for _, v := range projected_pops {
		if v.Year == baseYear {
			return v.TotalPopulation, nil
		}
	}
	return 0, fmt.Errorf("failed to find base year %d", baseYear)
}

// find the projected population of the base year for a number of codes
func FindBaseYearProjectedPopulationMultiCodes(projected_pops map[string][]LadPopulationProjection, baseYear int) (map[string]int, error) {
	m := make(map[string]int)
	for k, pops := range projected_pops {
		p, ok := FindBaseYearProjectedPopulation(pops, baseYear)
		if ok != nil {
			return m, ok
		}
		m[k] = p
	}
	return m, nil
}
