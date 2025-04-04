package resolver

import (
	"fmt"
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

type PopVec []LadPopulationProjection

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

type MapOfPopVec map[string][]LadPopulationProjection

// validates that a map of population slices is internally consistent
func IsValidMapOfPopVec(mp MapOfPopVec) bool {
	for k, v := range mp {
		if !(IsValidPopVec(v)) || v[0].Code != k {
			return false
		}
	}
	return true
}

type GrowthRate = struct {
	Year             int
	RateFromBaseYear float64
}

type EstimatedPopulation = struct {
	Year       int
	Population int
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

// For a single code
// Apply growth rates to the baseYear population to get population forecasts.
// Then filter the result to get population forecasts for the years of interest/requiredYears
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
