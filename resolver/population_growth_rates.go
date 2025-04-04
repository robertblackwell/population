package resolver

import (
	"errors"
	"fmt"
)

// Calculates the growth rates from the data in the projected_populations_v2 table
// For each code a separate calculation is performed.
// The returned map is indexed by code and the value for a code is a slice/array of LaPopulationProjection object
// In these arrays/slices there is a growth rate for every year between FirstYear(2018) and LastYear(2035)
// Growth rate values are numbers like 1.02 for a 2% positive growth and 0.98 for a 2% reduction.
// For each year the growth rate is the projected population for that year divided by the projected population for the baseYear (usually 2022)
// The growth rate for baseYear is 1.0
func CalculateGrowthRatesAllYearsMultipleCodes(projected_pops map[string][]LadPopulationProjection, baseYear int) map[string][]GrowthRate {
	result := make(map[string][]GrowthRate, 0)
	for k, v := range projected_pops {
		gr := CalculateGrowthRatesAllYears(v, baseYear)
		result[k] = gr
	}
	return result
}

// Calculate the growth rates of each years projected population relative to the projected population of the baseYear.
//
// NOTE: the projected_pops array MUST contain an entry for the baseYear
func CalculateGrowthRatesAllYears(projected_pops []LadPopulationProjection, baseYear int) []GrowthRate {
	base_year_index, err := findBaseYearIndexInPopulationProjectsions(projected_pops, 2023)
	if err != nil {
		panic("CalculateGrowthRatesAllYears: could not find base year in projected_pops")
	}
	return CalculateGrowthRatesRelativeToBaseYear(projected_pops, projected_pops[base_year_index])
	// growthRates := make([]GrowthRate, 0)
	// for ix, pp := range projected_pops {
	// 	var p float64
	// 	if ix == base_year_index {
	// 		p = 1.0
	// 	} else if ix < base_year_index {
	// 		p = float64(pp.TotalPopulation) / float64(projected_pops[base_year_index].TotalPopulation)
	// 	} else {
	// 		p = float64(pp.TotalPopulation) / float64(projected_pops[base_year_index].TotalPopulation)
	// 	}
	// 	growthRates = append(growthRates, GrowthRate{Year: pp.Year, RateFromBaseYear: p})
	// }
	// return growthRates
}

func CalculateGrowthRatesRelativeToBaseYear(projected_pops []LadPopulationProjection, baseYear LadPopulationProjection) []GrowthRate {
	growthRates := make([]GrowthRate, 0)
	for _, pp := range projected_pops {
		var p float64
		if baseYear.Year == pp.Year {
			p = 1.0
		} else {
			p = float64(pp.TotalPopulation) / float64(baseYear.TotalPopulation)
		}
		growthRates = append(growthRates, GrowthRate{Year: pp.Year, RateFromBaseYear: p})
	}
	return growthRates

}

func findBaseYearIndexInPopulationProjectsions(pop_projections []LadPopulationProjection, base_year int) (int, error) {
	for i, el := range pop_projections {
		if el.Year == base_year {
			return i, nil
		}
	}
	return 0, errors.New("could not find 2023")
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
