package resolver

import (
	"errors"
	"fmt"
	"forecast_model/mockdb"
	"maps"
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
		return map_result, errors.New("invalid age range")
	}
	ageRangeString := fmt.Sprintf("%d-%d", minAge, maxAge)
	fmt.Printf("Iam here")
	for k1, v1 := range target {
		if slices.Contains(codes, k1) {
			keys := slices.Sorted(maps.Keys(v1))
			for _, k2 := range keys {
				v2 := v1[k2]
				var p = LadPopulationProjection{
					Code: k1, Type: v2[0].Type, AgeRange: ageRangeString, Year: int(YearFromDate(k2)), TotalPopulation: 0,
				}
				for _, j := range v2 {
					if j.Age >= minAge && j.Age <= maxAge {
						p.TotalPopulation = p.TotalPopulation + j.Value
					}
				}
				_, ok := map_result[k1]
				if !ok {
					map_result[k1] = make([]LadPopulationProjection, 0)
				}
				map_result[k1] = append(map_result[k1], p)
			}
		}
	}
	fmt.Printf("I am here again")
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
	ageRangeString := fmt.Sprintf("%d-%d", minAge, maxAge)
	fmt.Printf("Iam here")
	for k1, v1 := range target {
		if slices.Contains(codes, k1) {
			keys := slices.Sorted(maps.Keys(v1))
			for _, k2 := range keys {
				v2 := v1[k2]
				var p = LadPopulationProjection{
					Code: k1, Type: v2[0].Type, AgeRange: ageRangeString, Year: int(YearFromDate(k2)), TotalPopulation: 0,
				}
				for _, j := range v2 {
					if j.Age >= minAge && j.Age <= maxAge {
						p.TotalPopulation = p.TotalPopulation + j.Value
					}
				}
				_, ok := map_result[k1]
				if !ok {
					map_result[k1] = make([]LadPopulationProjection, 0)
				}
				map_result[k1] = append(map_result[k1], p)
			}
		}
	}
	fmt.Printf("I am here again")
	return map_result, nil
}
