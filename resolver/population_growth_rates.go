package resolver

import (
	"fmt"
	"math"
)

// projected_pops - the key is a Code and the values must be valid PopVecs
// baseYear       - the key is a Code (same as projected_pop) and the value is a LapopulationProjection - all values must have the same year - 2023
func CalculateGrowthRates(projected_pops MapsOfPopVecs, baseYear map[string]map[string]LadPopulationProjection) MapsOfGrowthRates {
	result := make(map[string]map[string][]GrowthRate, 0)
	for code, v1 := range projected_pops {
		for ageRange, v2 := range v1 {
			gr := CalculateGrowthRatesBaseCase(v2, baseYear[code][ageRange])
			result[code] = map[string][]GrowthRate{ageRange: gr}
		}
	}
	return result
}
func CalculateGrowthRatesAsFloats(projected_pops MapsOfPopVecs, baseYear map[string]map[string]LadPopulationProjection) map[string]map[string][]float64 {
	result := make(map[string]map[string][]float64, 0)
	for code, v1 := range projected_pops {
		for ageRange, v2 := range v1 {
			gr := CalculateGrowthRatesAsFloatsBaseCase(v2, baseYear[code][ageRange])
			result[code] = map[string][]float64{ageRange: gr}
		}
	}
	return result
}

// This is a better interface as the need for a baseYear and a population for the baseYear is explicit
func CalculateGrowthRatesBaseCase(projected_pops []LadPopulationProjection, baseYear LadPopulationProjection) []GrowthRate {
	growthRates := make([]GrowthRate, 0)
	for _, pp := range projected_pops {
		var gr float64
		if baseYear.Year == pp.Year {
			gr = 1.0
		} else {
			gr = float64(pp.TotalPopulation) / float64(baseYear.TotalPopulation)
		}
		growthRates = append(growthRates, GrowthRate{Year: pp.Year, RateFromBaseYear: gr})
	}
	return growthRates

}

// This is a better interface as the need for a baseYear and a population for the baseYear is explicit
func CalculateGrowthRatesAsFloatsBaseCase(projected_pops []LadPopulationProjection, baseYear LadPopulationProjection) []float64 {
	growthRates := make([]float64, 0)
	for _, pp := range projected_pops {
		var gr float64
		if baseYear.Year == pp.Year {
			gr = 1.0
		} else {
			gr = float64(pp.TotalPopulation) / float64(baseYear.TotalPopulation)
		}
		growthRates = append(growthRates, gr)
	}
	return growthRates

}

// Apply growth rates to a population for the base year. Do this for a number of codes.
// Filter the output so that the returned value only has forrecast populations for the specified list of years (between FirstYear and LastYear)
func CalculateEstimatedPopulationsFromInts(growthRates MapsOfGrowthRates, baseYearPopulation map[string]map[string]int) (MapsOfEstimatedPopulations, error) {
	result := make(map[string]map[string][]EstimatedPopulation, 0)
	for code, v1 := range growthRates {
		for ageRange, v2 := range v1 {
			basePop, ok := baseYearPopulation[code][ageRange]
			if !ok {
				return result, fmt.Errorf("a base year population was not provided for code: %s", code)
			}
			ep := CalculateEstimatedPopulationsBaseCase(v2, basePop)
			result[code][ageRange] = ep
		}
	}
	return result, nil
}
func CalculateEstimatedPopulations(growthRates MapsOfGrowthRates, baseYearPopulation map[string]map[string]LadPopulationProjection) (MapsOfEstimatedPopulations, error) {
	result := make(map[string]map[string][]EstimatedPopulation, 0)
	for code, v1 := range growthRates {
		for ageRange, v2 := range v1 {
			basePop, ok := baseYearPopulation[code][ageRange]
			if !ok {
				return result, fmt.Errorf("a base year population was not provided for code: %s", code)
			}
			ep := CalculateEstimatedPopulationsBaseCase(v2, basePop.TotalPopulation)
			result[code][ageRange] = ep
		}
	}
	return result, nil
}

// For a single code
// Apply growth rates to the baseYear population to get population forecasts.
// Then filter the result to get population forecasts for the years of interest/requiredYears
func CalculateEstimatedPopulationsBaseCase(growthRates []GrowthRate, baseYearPopulation int) []EstimatedPopulation {

	result := make([]EstimatedPopulation, 0)
	for _, gr := range growthRates {
		p := int(math.Round(float64(baseYearPopulation) * gr.RateFromBaseYear))
		result = append(result, EstimatedPopulation{Year: gr.Year, Population: p})
	}
	return result
}

// // Calculates the growth rates from the data in the projected_populations_v2 table
// // For each code a separate calculation is performed.
// // The returned map is indexed by code and the value for a code is a slice/array of LaPopulationProjection object
// // In these arrays/slices there is a growth rate for every year between FirstYear(2018) and LastYear(2035)
// // Growth rate values are numbers like 1.02 for a 2% positive growth and 0.98 for a 2% reduction.
// // For each year the growth rate is the projected population for that year divided by the projected population for the baseYear (usually 2022)
// // The growth rate for baseYear is 1.0
// func CalculateGrowthRatesAllYearsMultipleCodes(projected_pops map[string][]LadPopulationProjection, baseYear int) map[string][]GrowthRate {
//     result := make(map[string][]GrowthRate, 0)
//     for k, v := range projected_pops {
//         gr := CalculateGrowthRatesAllYears(v, baseYear)
//         result[k] = gr
//     }
//     return result
// }

// // Calculate the growth rates of each years projected population relative to the projected population of the baseYear.
// //
// // NOTE: the projected_pops array MUST contain an entry for the baseYear
// func CalculateGrowthRatesAllYears(projected_pops []LadPopulationProjection, baseYear int) []GrowthRate {
//     base_year_index, err := findBaseYearIndexInPopulationProjectsions(projected_pops, 2023)
//     if err != nil {
//         panic("CalculateGrowthRatesAllYears: could not find base year in projected_pops")
//     }
//     return CalculateGrowthRatesRelativeToBaseYear(projected_pops, projected_pops[base_year_index])
//     // growthRates := make([]GrowthRate, 0)
//     // for ix, pp := range projected_pops {
//     //     var p float64
//     //     if ix == base_year_index {
//     //         p = 1.0
//     //     } else if ix < base_year_index {
//     //         p = float64(pp.TotalPopulation) / float64(projected_pops[base_year_index].TotalPopulation)
//     //     } else {
//     //         p = float64(pp.TotalPopulation) / float64(projected_pops[base_year_index].TotalPopulation)
//     //     }
//     //     growthRates = append(growthRates, GrowthRate{Year: pp.Year, RateFromBaseYear: p})
//     // }
//     // return growthRates
// }

// func findBaseYearIndexInPopulationProjectsions(pop_projections []LadPopulationProjection, base_year int) (int, error) {
//     for i, el := range pop_projections {
//         if el.Year == base_year {
//             return i, nil
//         }
//     }
//     return 0, errors.New("could not find 2023")
// }

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
