package populations

import (
	"fmt"
	"math"
	"popmodel/cayvalues"
)

// projected_pops - A CayValue[int] where the values are projected populations as integers
// baseYear       - A CayValue[int] with the same 1st and 2nd levels keys as projected_pop
//
//	and only one 3rd level key which is 2023 or whatever the baseyear should be
//
// Returns a CayValue[float64] structure where the values are growth rates.
// Returns an error if the 1st and 2nd level keys of baseYear do not match the correspondng keys of projected_pop
// Returns an error if the 3rd level keys of baseYear does not include 2023 (or the baseYear)
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
// Filter the output so that the returned value only has forecast populations for the specified list of years (between FirstYear and LastYear)
// Returns an error if the 1st and 2nd level jeys of baseYearPopulation do not match the 1st and 2nd level keys of growthRate
// Return an error if the 3rd level keys of baseYearPopulation doe n ot include 2023 or whatever the baseYear is
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

type Result struct {
	Code        string
	AgeRange    string
	Current     int
	Projected   []int
	GrowthRates []float64
}
type FlatResult struct {
	Code       string
	AgeRange   string
	Current    int // all these should be the same
	Projected  int
	GrowthRate float64
}

func CalculateEstimatedPopulationFlatResult(growthRates cayvalues.CayValues[float64], baseYearPopulation cayvalues.CayValues[int]) (cayvalues.CayValues[FlatResult], error) {
	f := func(code string, ageRange string, year int, growthRate float64) (FlatResult, error) {
		baseYearPop, ok := baseYearPopulation.At(code, ageRange, 2023)
		if !ok {
			return FlatResult{}, fmt.Errorf("CalcuateGrowthRates did not find a base year population for %s %s %d", code, ageRange, 2023)
		}
		estimatedPop := int(math.Round(growthRate * float64(baseYearPop)))
		return FlatResult{Code: code, AgeRange: ageRange, Current: baseYearPop, Projected: estimatedPop, GrowthRate: growthRate}, nil
	}

	return cayvalues.Map(growthRates, f)
}
func FlatResultToResult(flatResult cayvalues.CayValues[FlatResult]) (map[string]map[string]Result, error) {
	r := map[string]map[string]Result{}
	f := func(code string, ageRange string, year int, fvalue FlatResult) error {
		_, ok := r[code]
		if !ok {
			r[code] = map[string]Result{ageRange: {Code: code, AgeRange: ageRange, Current: fvalue.Current, Projected: []int{fvalue.Projected}, GrowthRates: []float64{fvalue.GrowthRate}}}
		}
		_, ok = r[code][ageRange]
		if !ok {
			r[code][ageRange] = Result{Code: code, AgeRange: ageRange, Current: fvalue.Current, Projected: []int{fvalue.Projected}, GrowthRates: []float64{fvalue.GrowthRate}}
		}
		tmp := r[code][ageRange]
		if tmp.Current != fvalue.Current {
			return fmt.Errorf("error while iterating over cay values for keys %s %s value of current: %d differ from %d ", code, ageRange, tmp.Current, fvalue.Current)
		}
		projected := append(tmp.Projected, fvalue.Projected)
		grates := append(tmp.GrowthRates, fvalue.GrowthRate)
		r[code][ageRange] = Result{Code: code, AgeRange: ageRange, Current: fvalue.Current, Projected: projected, GrowthRates: grates}
		return nil
	}
	// dont need the result from this call - this is an accumulator and the result is in 'r'
	err := cayvalues.Iterate(flatResult, f)
	if err != nil {
		return r, err
	}
	return r, nil
}
func CalcEstimatedPopulationsFlatResult(growthRates cayvalues.CayValues[float64], baseYearPopulation cayvalues.CayValues[int]) ([]Result, error) {
	r1, err := CalculateEstimatedPopulationFlatResult(growthRates, baseYearPopulation)
	if err != nil {
		return []Result{}, err
	}
	r2, err2 := FlatResultToResult(r1)
	if err2 != nil {
		return []Result{}, err2
	}
	res := []Result{}
	for _, v1 := range r2 {
		for _, v2 := range v1 {
			res = append(res, v2)
		}
	}
	return res, nil
}
