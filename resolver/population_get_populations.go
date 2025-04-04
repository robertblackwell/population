package resolver

import (
	"errors"
	"fmt"
	"forecast_model/mockdb"
	"slices"
)

// Get the base year projected population for a list of codes
func GetProjectedPopulationsByCodeForBaseYear(ctx mockdb.Context, codes []string, baseYear int, minAge int, maxAge int) (map[string]LadPopulationProjection, error) {
	result := make(map[string]LadPopulationProjection, 0)
	asSlice, e := GetProjectedPopulationsByCodeForYears(ctx, codes, []int{baseYear}, minAge, maxAge)
	if e != nil {
		return result, e
	}
	for k, v := range asSlice {
		result[k] = v[0]
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
