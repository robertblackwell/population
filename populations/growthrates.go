package populations

import (
	"fmt"
	"math"
	"popmodel/cayvalues"
)

// projected_pops - the key is a Code and the values must be valid PopVecs
// baseYear       - the key is a Code (same as projected_pop) and the value is a LapopulationProjection - all values must have the same year - 2023
func CalculateGrowthRates(projected_pops cayvalues.CayValues[int], baseYear cayvalues.CayValues[int]) (cayvalues.CayValues[float64], error) {
	f := func(code string, ageRange string, year int, population int) (float64, error) {
		baseYearPop, ok := baseYear.At(code, ageRange, 2023)
		if !ok {
			return 0.0, fmt.Errorf("CalcuateGrowthRates did not find a base year population for %s %s %d", code, ageRange, 2023)
		}
		gr := float64(population) / float64(baseYearPop)
		return gr, nil
	}
	r, err := cayvalues.Map(projected_pops, f)
	if err != nil {
		return r, err
	}
	return r, nil
}

// Apply growth rates to a population for the base year. Do this for a number of codes.
// Filter the output so that the returned value only has forrecast populations for the specified list of years (between FirstYear and LastYear)
func CalculateEstimatedPopulations(growthRates cayvalues.CayValues[float64], baseYearPopulation cayvalues.CayValues[int]) (cayvalues.CayValues[int], error) {
	f := func(code string, ageRange string, year int, growthRate float64) (int, error) {
		baseYearPop, ok := baseYearPopulation.At(code, ageRange, 2023)
		if !ok {
			return 0.0, fmt.Errorf("CalcuateGrowthRates did not find a base year population for %s %s %d", code, ageRange, 2023)
		}
		estimatedPop := int(math.Round(growthRate * float64(baseYearPop)))
		return estimatedPop, nil
	}
	return cayvalues.Map(growthRates, f)
}

// // find the projected population for the base year
// func FindBaseYearProjectedPopulation(projected_pops []LadPopulationProjection, baseYear int) (int, error) {
// 	for _, v := range projected_pops {
// 		if v.Year == baseYear {
// 			return v.TotalPopulation, nil
// 		}
// 	}
// 	return 0, fmt.Errorf("failed to find base year %d", baseYear)
// }

// // find the projected population of the base year for a number of codes
// func FindBaseYearProjectedPopulationMultiCodes(projected_pops map[string][]LadPopulationProjection, baseYear int) (map[string]int, error) {
// 	m := make(map[string]int)
// 	for k, pops := range projected_pops {
// 		p, ok := FindBaseYearProjectedPopulation(pops, baseYear)
// 		if ok != nil {
// 			return m, ok
// 		}
// 		m[k] = p
// 	}
// 	return m, nil
// }
