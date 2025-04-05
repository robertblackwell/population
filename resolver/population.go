package resolver

import (
	"fmt"
	"math"
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

type GrowthRate struct {
	Year             int
	RateFromBaseYear float64
}

type EstimatedPopulation struct {
	Year       int
	Population int
}

type PopVec = []LadPopulationProjection

// first key is Code 2nd key is AgeRange as string
type MapsOfPopVecs = map[string]map[string][]LadPopulationProjection

func MapOfPopVecs_at(m MapsOfPopVecs, code string, ageRange string) ([]LadPopulationProjection, bool) {
	a, ok := m[code]
	if !ok {
		return []LadPopulationProjection{}, false
	}
	b, ok := a[ageRange]
	if !ok {
		return []LadPopulationProjection{}, false
	}
	return b, true
}

// first key is Code 2nd key is AgeRange as string
type MapsOfGrowthRates = map[string]map[string][]GrowthRate

// first key is Code 2nd key is AgeRange as string
type MapsOfEstimatedPopulations = map[string]map[string][]EstimatedPopulation

// validates that a population slice is internally consistent
func IsValidPopVec(p PopVec) bool {
	code := p[0].Code
	typ := p[0].Type
	ageRange := p[0].AgeRange
	for _, v := range p {
		if v.Code != code || v.Type != typ || v.AgeRange != ageRange {
			return false
		}
	}
	return true
}

// validates that a map of population slices is internally consistent
func IsValidMapOfPopVec(mp MapsOfPopVecs) bool {
	for k1, v1 := range mp {
		for k2, v2 := range v1 {
			if !(IsValidPopVec(v2)) {
				return false
			}
			for _, v3 := range v2 {
				if v3.Code != k1 || v3.AgeRange != k2 {
					return false
				}
			}
		}
	}
	return true
}

func Index2LevelMap[T any](m map[string]map[string][]T, k1 string, k2 string) ([]T, bool) {
	a, ok := m[k1]
	if !ok {
		return []T{}, false
	}
	b, ok := a[k2]
	if !ok {
		return []T{}, false
	}
	return b, true
}

// Apply growth rates to a population for the base year. Do this for a number of codes.
// Filter the output so that the returned value only has forrecast populations for the specified list of years (between FirstYear and LastYear)
func CalculateEstimatedPopulationsForSomeYearsMultipleCodes(growthRates MapsOfGrowthRates, baseYearPopulation map[string]map[string]int) (MapsOfEstimatedPopulations, error) {
	result := make(map[string]map[string][]EstimatedPopulation, 0)
	for code, v1 := range growthRates {
		for ageRange, v2 := range v1 {
			basePop, ok := baseYearPopulation[code][ageRange]
			if !ok {
				return result, fmt.Errorf("a base year population was not provided for code: %s", code)
			}
			ep := CalculateEstimatedPopulationsForSomeYears(v2, basePop)
			result[code][ageRange] = ep
		}
	}
	return result, nil
}

// For a single code
// Apply growth rates to the baseYear population to get population forecasts.
// Then filter the result to get population forecasts for the years of interest/requiredYears
func CalculateEstimatedPopulationsForSomeYears(growthRates []GrowthRate, baseYearPopulation int) []EstimatedPopulation {

	result := make([]EstimatedPopulation, 0)
	for _, gr := range growthRates {
		p := int(math.Round(float64(baseYearPopulation) * gr.RateFromBaseYear))
		result = append(result, EstimatedPopulation{Year: gr.Year, Population: p})
	}
	return result
}
