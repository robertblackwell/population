package resolver

import (
	"fmt"
	"forecast_model/mockdb"
	"forecast_model/models"
	"math"
	"sort"
	"strconv"
	"strings"
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

func /*(repo *ResolverRepository)*/ GetPopulationByLadCodeFromDb(ctx mockdb.Context, ladCode string, startYear, rangeSize, minAge, maxAge, futureOffset int, includeIntermediates bool) (*models.PopulationProjections, error) {

	requestedOffset := futureOffset

	// IDEA: Why default this to 10?
	if futureOffset <= 0 {
		futureOffset = 10
	}

	// The most current year we have population data for is 2023
	currentPopulationYear := 2023

	// gets the lad, region and country codes
	parents, err := /*repo.*/ GetParentCodes(ctx, ladCode)
	if err != nil {
		return nil, fmt.Errorf("error getting parent codes: %w", err)
	}

	// get the current (2023) static populations by LAD
	allPopulations, err := /*repo.*/ GetCurrentPopulationsByCodes(ctx, parents, rangeSize, minAge, maxAge, currentPopulationYear)
	if err != nil {
		return nil, fmt.Errorf("error getting current population by lad code: %w", err)
	}
	currentBaseYearPopulation, err := CollateBaseYearProjectedPopulationsByCode(allPopulations)
	if err != nil {
		return nil, err
	}
	// Get the projections for each LAD (2018-2035)
	// IDEA: get all the projections for all years
	// IDEA: always get the intermediates for the growth rates
	projections, err := /*repo.*/ GetPopulationByCodes(ctx, parents, startYear, rangeSize, minAge, maxAge, futureOffset, includeIntermediates)
	if err != nil {
		return nil, fmt.Errorf("error getting population by codes: %w", err)
	}

	// Sort them in to a map of geo code to age range to population
	// geoPopulations := make(map[string]map[string][]int)
	// for _, pop := range projections {
	// 	if geoPopulations[pop.Code] == nil {
	// 		geoPopulations[pop.Code] = make(map[string][]int)
	// 	}
	// 	if geoPopulations[pop.Code][pop.AgeRange] == nil {
	// 		geoPopulations[pop.Code][pop.AgeRange] = []int{}
	// 	}
	// 	geoPopulations[pop.Code][pop.AgeRange] = append(geoPopulations[pop.Code][pop.AgeRange], pop.TotalPopulation)
	// }
	geoPopulations := CollateProjectedPopulationsByCode(projections)
	growthRates := CalculateGrowthRatesAsFloats(geoPopulations, currentBaseYearPopulation)
	// Calculate the growth rate for each LAD and age range
	// growthRates := make(map[string]map[string][]float64)
	// for code, pop := range geoPopulations {
	// 	for ageRange, values := range pop {
	// 		// Growth rate is calculated from the year to the year ahead (not previous)
	// 		growthRateArr := CalculateGrowthRates(values)
	// 		for _, v := range values {
	// 			values = append(values, v)
	// 		}
	// 		if growthRates[code] == nil {
	// 			growthRates[code] = make(map[string][]float64)
	// 		}
	// 		growthRates[code][ageRange] = growthRateArr
	// 	}
	// }

	return processPopulationResponse(allPopulations, growthRates, startYear, requestedOffset, includeIntermediates)
}

type Result struct {
	Code        string
	AgeRange    string
	Current     int
	Projected   []int
	GrowthRates []float64
}

// allPopulations currently (Not the projections),
// growthRates (percentage), (we get one even if the offset is 0)
// startingYear e.g. 2020, (can be in past or future)
// requestedOffset e.g. 10 (2030),
// includeIntermediates - add all years in between or not (boolean). If there is a requestedOffset

// Best if returns a growth rate for each year

// Growth rate e.g. 2020-21,2021-22,2022-23,2023-24,2024-25,2025-26,2026-27,2027-28,2028-29,2029-30
// Current year e.g. 2023

func processPopulationResponse(allCurrentPopulations []LadPopulationProjection, growthRates map[string]map[string][]float64, startYear, requestedOffset int, includeIntermediates bool) (*models.PopulationProjections, error) {

	// Calculate the projected population for each LAD and age range
	projectedPopulation := []Result{}
	for _, pop := range allCurrentPopulations {
		code := pop.Code
		ageRange := pop.AgeRange
		currentPopulation := pop.TotalPopulation // for 2023

		if growthRateArr, ok := growthRates[code][ageRange]; ok {
			// collect the projected population for each year
			var projected = make([]int, len(growthRateArr))
			// Calculate the projected population for each year
			for i, growthRate := range growthRateArr {
				futurePopulation := int(math.Round(float64(currentPopulation) * (growthRate)))
				projected[i] = futurePopulation
			}

			projectedPopulation = append(projectedPopulation, Result{
				Code:        code,
				AgeRange:    ageRange,
				Current:     currentPopulation,
				Projected:   projected,
				GrowthRates: growthRateArr,
			})
		}
	}

	ageRangeMap := make(map[string][]models.PopulationAgeRange)
	for _, pop := range projectedPopulation {
		values := []int{pop.Current}
		if requestedOffset > 0 {
			values = append(values, pop.Projected...)
		}

		// if includeIntermediates is false, filter out the values that are not in the requested range
		ageRangeMap[pop.Code] = append(ageRangeMap[pop.Code], models.PopulationAgeRange{
			AgeRange:    pop.AgeRange,
			Values:      values,
			GrowthRates: pop.GrowthRates,
		})
	}

	geographies := []models.PopulationGeography{}

	for code, ageRanges := range ageRangeMap {
		geographies = append(geographies, models.PopulationGeography{
			Code:      code,
			AgeRanges: ageRanges,
		})
	}

	for _, geography := range geographies {
		sort.Slice(geography.AgeRanges, func(i, j int) bool {
			iAge, _ := strconv.Atoi(strings.Split(geography.AgeRanges[i].AgeRange, "-")[0])
			jAge, _ := strconv.Atoi(strings.Split(geography.AgeRanges[j].AgeRange, "-")[0])
			return iAge < jAge
		})
	}
	sort.Slice(geographies, func(i, j int) bool {
		return geographies[i].Code < geographies[j].Code
	})

	years := []int{startYear}

	if requestedOffset > 0 {
		if includeIntermediates {
			for i := 1; i <= requestedOffset; i++ {
				years = append(years, startYear+i)
			}
		} else {
			years = append(years, startYear+requestedOffset)
		}
	}

	// filter years if include intermediates is false <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
	// filter age ranges if include intermediates is false <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

	return &models.PopulationProjections{
		Geographies: geographies,
		Years:       years,
	}, nil
}
