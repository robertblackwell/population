package model

import (
	"fmt"
	"math"
	"slices"
	"sort"
)

type PastYear int

const FirstPastYear = 2018
const LastPastYear = 2023

func IsValidPastYear(y int) bool {
	return y >= 2018 && y <= 2023
}
func PastYearFromInt(y int) PastYear {
	if IsValidPastYear(y) {
		return PastYear(y)
	}
	panic(fmt.Sprintf("PastYearFromInt y: %d", y))
}

type FutureYear int

func IsValidFutureYear(y FutureYear) bool {
	return y >= 2024
}

func FirstFutureYear() FutureYear { return 2024 }

type PopulationType string

type Gender string

const (
	Female Gender = "female"
	Male   Gender = "male"
)

type YearAgePopulationProjection struct {
	year            int // 2018 .. 2023
	age             int // >= 0
	totalPopulation int // > 0
	// male_population   int // > 0 .. unused at this time not in population_projections_v2
	// female_population int // > 0 .. unused at this time not in population_projections_v2
}

type YearPopulationProjection struct {
	year     int // 2018 .. 2023
	maxAge   int
	popByAge []YearAgePopulationProjection
	/*
	* Valid only if both pop_by_age slice has max_age entries
	* all numbers between 0 and max_age inclusive must an age value for both female and male slices
	* the year field in every pop_by_age slice entry must be the same as this structs year value
	 */
}

func YearPopulationProjection_CreateEmpty(pastYear PastYear) YearPopulationProjection {
	return YearPopulationProjection{
		year:     int(pastYear),
		maxAge:   0,
		popByAge: make([]YearAgePopulationProjection, 0),
	}
}
func (y *YearPopulationProjection) AddYearAge(pastYear PastYear, age int, popValue int) {
	tmp := YearAgePopulationProjection{
		year:            int(pastYear),
		age:             age,
		totalPopulation: popValue,
	}
	if age > y.maxAge {
		y.maxAge = age
	}
	y.popByAge = append(y.popByAge, tmp)
}

// *******************************************************************************************************
// The type PopulationForecastInputData and its related methods implement a model for forecasting
// future populations
// *******************************************************************************************************
type PopulationForecastInputData struct {
	projectionByYear []YearPopulationProjection
}

// Create an empty model
func PopulationForecastInputData_CreateEmpty() PopulationForecastInputData {
	return PopulationForecastInputData{projectionByYear: make([]YearPopulationProjection, 0)}
}

// Add population by age data for a given year to the model
// The forecast of future populations is based on the forecast pop for the
// period 2018-2023
func (pfid *PopulationForecastInputData) AddYearData(yearData YearPopulationProjection) {
	pfid.projectionByYear = append(pfid.projectionByYear, yearData)
}

// Get the population forecast for one of the years between 2018 and 2023 inclusive
func (pfid *PopulationForecastInputData) GetPopulationByYearAndAge(year PastYear, age int) int {
	year_index := int(year) - int(FirstPastYear)
	pop := pfid.projectionByYear[year_index].popByAge[age].totalPopulation
	return pop
}

// Calculate the annual population growth rates from the input population projections for
// a specific age.
// Returns a vector of growthrates such as 1.02 or 0.98
func (pfid *PopulationForecastInputData) AnnualizedGrowthRatesByAge(age int) []float64 {
	result := make([]float64, 0)
	for y := 1; y <= LastPastYear-FirstPastYear; y++ {
		pop2 := pfid.GetPopulationByYearAndAge(PastYear(FirstPastYear+y), age)
		pop1 := pfid.GetPopulationByYearAndAge(PastYear(FirstPastYear+y-1), age)
		gr := float64(pop2) / float64(pop1)
		fmt.Printf("pop2: %d pop1: %d  gr: %f\n", pop2, pop1, gr)
		result = append(result, gr)
	}
	return result
}

type BasicPopulationForecast struct {
	year       int
	population int
}

// Perform a population forcast using the models data.
// The forecast is for a particular age cohort -- age
// The forecast should cover `nbr_years`
// Should start with the assumption that the population at the start of the forecast is `pop_in_2023`
// Returns a slice where eachelement is a year and its corresponding population
func (pfid *PopulationForecastInputData) ForecastPopulationByAgeAndNumberofYears(age int, nbr_years int, pop_in_2023 int) []BasicPopulationForecast {
	gr := pfid.AnnualizedGrowthRatesByAge(age)
	result := make([]BasicPopulationForecast, 0)
	counter := 0
	latest_pop := pop_in_2023
	// apply the growth rates successively until we have used them all
	for _, g := range gr {
		latest_pop = int(math.Round(float64(latest_pop) * g))
		counter++
		result = append(result, BasicPopulationForecast{year: int(LastPastYear) + counter, population: latest_pop})
	}
	// once all the growth rates are used keep applying the last growth rate until we have forecast enough
	// future years
	g := gr[len(gr)-1]
	for counter < nbr_years {
		latest_pop = int(math.Round(float64(latest_pop) * g))
		counter++
		result = append(result, BasicPopulationForecast{year: int(LastPastYear) + counter, population: latest_pop})
	}
	return result
}

// Perform a population forcast using the models data.
// The forecast is for a particular age cohort -- age
// The forecast is provided for a list of horizons or years in the future. In ascending order
// Should start with the assumption that the population at the start of the forecast is `pop_in_2023`
// Returns a slice where eachelement is a year and its corresponding population
func (pfid *PopulationForecastInputData) ForecastPopulationByAgeToHorizon(age int, pop_in_2023 int, horizons []int) []BasicPopulationForecast {
	sort.Ints(horizons)
	nbrOfYears := horizons[len(horizons)-1] - LastPastYear
	r := pfid.ForecastPopulationByAgeAndNumberofYears(age, nbrOfYears, pop_in_2023)
	result := make([]BasicPopulationForecast, 0)
	for _, f := range r {
		y := f.year
		if slices.Contains(horizons, y) {
			result = append(result, f)
		}
	}
	return result
}

// func PopulatioProjectionByAgeAndGenders(pid PopulationForecastInputData, age int, gender Gender) {
// 	result := make([]int, 1)

// 	for i, p := range pid.projectionByYear {
// 		result = append(result, 0)
// 		result[i] += p.popByAge[age].totalPopulation
// 		// if gender == Female {
// 		// 	result[i] += p.pop_by_age[age].female_population
// 		// } else if gender == Male {
// 		// 	result[i] += p.pop_by_age[age].female_population
// 		// }
// 	}
// }

// func (pid PopulationForecastInputData) IsValid() bool {
// 	// 1. should have a YearPopulationProjectionRecord for each year between FirstPastYear and LastPstYear inclusive
// 	// 2. each YearPopulationProjectionRecord should be valid
// 	// 3. all YearPopulationProjectionRecords should have the same max_age
// 	return true
// }
