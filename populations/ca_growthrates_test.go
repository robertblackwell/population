package populations

import (
	"fmt"
	"math"
	"popmodel/cayvalues"
	"popmodel/mockdb"
	"popmodel/repo"
	"testing"

	"github.com/stretchr/testify/assert"
)

type LadPopulation struct {
	LadCode    string  `json:"lad_code"`
	Share      float64 `json:"contribution"`
	Population int     `json:"population"`
	Year       int     `json:"year"`
}

func TestCaGrowthRates01(tt *testing.T) {

	contributions := []LadPopulation{
		LadPopulation{
			LadCode:    "E09000002",
			Share:      0.19567854260348208,
			Population: 42428,
			Year:       2022,
		},
		LadPopulation{
			LadCode:    "E09000025",
			Share:      0.0899875555262529,
			Population: 32757,
			Year:       2022,
		},
		LadPopulation{
			LadCode:    "E09000026",
			Share:      0.1055334939131567,
			Population: 32422,
			Year:       2022,
		},
	}
	ladPopMap := map[string]int{
		"E09000002": 42428,
		"E09000025": 32757,
		"E09000026": 32422,
	}
	// get a slice of LadPopulationProjections
	ctx := mockdb.Context{}

	totalPopulations, er := repo.GetProjectedPopulationsByCodes(ctx, []string{"E09000002", "E09000025", "E09000026"}, 2022, 100, 0, 99, 1, false)
	assert.True(tt, er == nil)

	projectedPopulations, er := repo.GetProjectedPopulationsByCodes(ctx, []string{"E09000002", "E09000025", "E09000026"}, 2020, 5, 20, 60, 10, true)
	assert.True(tt, er == nil)
	// geoProjectedPopulations, er := cayvalues.NewCayValuesFromArr(projectedPopulations)
	// assert.True(tt, er == nil)
	geoProjectedPopulationsInt, er := CollateProjectedPopulationsByCode(projectedPopulations)
	assert.True(tt, er == nil)
	baseYearProjectedPopulation, er := repo.GetProjectedPopulationsByCodes(ctx, []string{"E09000002", "E09000025", "E09000026"}, 2023, 5, 20, 60, 0, true)
	assert.True(tt, er == nil)
	// geoBaseYearProjectedPopulation, er := cayvalues.NewCayValuesFromArr(baseYearProjectedPopulation)
	// assert.True(tt, er == nil)
	geoBaseYearProjectedPopulationInt, er := CollateProjectedPopulationsByCode(baseYearProjectedPopulation)
	assert.True(tt, er == nil)

	growthRates, er := CalculateGrowthRates(geoProjectedPopulationsInt, geoBaseYearProjectedPopulationInt)
	assert.True(tt, er == nil)

	// apply the growth rates to the pojectedPopulations for the baseYear. Should give the projected populations
	estimatedPopulations, er := CalculateEstimatedPopulations(growthRates, geoBaseYearProjectedPopulationInt)
	assert.True(tt, er == nil)

	// verify the projectedPopulations as what we expected
	f := func(code string, ageRange string, year int, value int) error {
		v, _ := estimatedPopulations.At(code, ageRange, year)
		assert.True(tt, (value == v))
		fmt.Printf("XXX %s %s %d %d %d %v\n", code, ageRange, year, value, v, (value == v))
		return nil
	}
	cayvalues.Iterate(geoProjectedPopulationsInt, f)

	// make a map of LadCodes to weights - for example - weights['E09000002] == 0.19567854260348208
	weights := map[string]float64{}
	for _, v := range contributions {
		if _, ok := weights[v.LadCode]; !ok {
			weights[v.LadCode] = v.Share
		}
	}
	// apply the weights to the estimated populations to get the estimated contribution of each
	// lad to the catchment area still organized by code, ageRange and year
	weightedEstimatedPopulations := cayvalues.NewCayValues[int]()
	apw := func(code string, ageRange string, year int, value int) error {
		epop := value
		if w, ok := weights[code]; !ok {
			return fmt.Errorf("")
		} else {
			wpop := int(math.Round(float64(epop)*w + 1.0))
			err := weightedEstimatedPopulations.Add(code, ageRange, year, wpop)
			if err != nil {
				return err
			}
		}
		return nil
	}
	err := cayvalues.Iterate(estimatedPopulations, apw)
	assert.True(tt, err == nil)

	// sum the weightedEstimated population over lad so that there is one value for each (ageRange,year) combination.
	// do this by making a cayvalue.CayValue with a single code key
	oaPopulations := cayvalues.NewCayValues[int]()
	oaf := func(code string, ageRange string, year int, value int) error {
		v, ok := oaPopulations.At("oacode", ageRange, year)
		if !ok {
			oaPopulations.Set("oacode", ageRange, year, value)
		} else {
			oaPopulations.Set("oacode", ageRange, year, value+v)
		}
		return nil
	}
	err = cayvalues.Iterate(weightedEstimatedPopulations, oaf)
	assert.True(tt, err == nil)

	// as a santiy check - by comparing values in the weightedEstmatedPopulation object against values
	// in the estimatedPopulations object compute a weight/share for each ageRange, year combination.
	// See how close those weights are to the values in weights[code]
	deducedWeights := cayvalues.NewCayValues[float64]()
	dw := func(code string, ageRange string, year int, value int) error {
		wp, ok := weightedEstimatedPopulations.At(code, ageRange, year)
		if !ok {
			return fmt.Errorf("")
		}
		ep, ok := estimatedPopulations.At(code, ageRange, year)
		if !ok {
			return fmt.Errorf("")
		}
		w := float64(wp) / float64(ep)
		refW := weights[code]
		// fmt.Printf("\nFF  %f %f\n", w, refW)
		if math.Abs((w-refW)/w) > 0.1 {
			fmt.Printf("XXXXXXXXX   %f %f difference exceeds 15 percent\n", w, refW)
			assert.True(tt, false)
		}
		deducedWeights.Add(code, ageRange, year, w)
		return nil
	}
	err = cayvalues.Iterate(estimatedPopulations, dw)
	assert.True(tt, err == nil)
	// fmt.Printf("%v", estimatedPopulations)
	// fmt.Printf("%v", growthRates)
	// fmt.Printf("%v", geoBaseYearProjectedPopulationInt)
	// fmt.Printf("%v", geoBaseYearProjectedPopulation)
	// fmt.Printf("%v", geoProjectedPopulations)
	// fmt.Printf("%v", geoProjectedPopulationsInt)
	// fmt.Printf("%v", projectedPopulations)
	fmt.Printf("%v", totalPopulations)
	fmt.Printf("%v", &contributions)
	fmt.Printf("%v", &ladPopMap)
}

func TestCaGrowthRates02(tt *testing.T) {

	contributions := []LadPopulation{
		LadPopulation{
			LadCode:    "E09000002",
			Share:      0.19567854260348208,
			Population: 42428,
			Year:       2022,
		},
		LadPopulation{
			LadCode:    "E09000025",
			Share:      0.0899875555262529,
			Population: 32757,
			Year:       2022,
		},
		LadPopulation{
			LadCode:    "E09000026",
			Share:      0.1055334939131567,
			Population: 32422,
			Year:       2022,
		},
	}
	ladPopMap := map[string]int{
		"E09000002": 42428,
		"E09000025": 32757,
		"E09000026": 32422,
	}
	// get a slice of LadPopulationProjections
	ctx := mockdb.Context{}

	totalPopulations, er := repo.GetProjectedPopulationsByCodes(ctx, []string{"E09000002", "E09000025", "E09000026"}, 2022, 100, 0, 99, 1, false)
	assert.True(tt, er == nil)

	projectedPopulations, er := repo.GetProjectedPopulationsByCodes(ctx, []string{"E09000002", "E09000025", "E09000026"}, 2020, 5, 20, 60, 10, true)
	assert.True(tt, er == nil)
	// geoProjectedPopulations, er := cayvalues.NewCayValuesFromArr(projectedPopulations)
	// assert.True(tt, er == nil)
	geoProjectedPopulationsInt, er := CollateProjectedPopulationsByCode(projectedPopulations)
	assert.True(tt, er == nil)
	baseYearProjectedPopulation, er := repo.GetProjectedPopulationsByCodes(ctx, []string{"E09000002", "E09000025", "E09000026"}, 2023, 5, 20, 60, 0, true)
	assert.True(tt, er == nil)
	// geoBaseYearProjectedPopulation, er := cayvalues.NewCayValuesFromArr(baseYearProjectedPopulation)
	// assert.True(tt, er == nil)
	geoBaseYearProjectedPopulationInt, er := CollateProjectedPopulationsByCode(baseYearProjectedPopulation)
	assert.True(tt, er == nil)

	baseYearCurrentPopulations, er := repo.GetCurrentPopulationsByCodes(ctx, []string{"E09000002", "E09000025", "E09000026"}, 2023, 5, 20, 60, 0, true)
	assert.True(tt, er == nil)
	geoBaseYearCurrentPopulations, er := CollateProjectedPopulationsByCode(baseYearCurrentPopulations)
	assert.True(tt, er == nil)
	fmt.Printf("%v", &geoBaseYearCurrentPopulations)

	growthRates, er := CalculateGrowthRates(geoProjectedPopulationsInt, geoBaseYearProjectedPopulationInt)
	assert.True(tt, er == nil)

	// apply the growth rates to the pojectedPopulations for the baseYear. Should give the projected populations
	estimatedPopulations, er := CalculateEstimatedPopulations(growthRates, geoBaseYearProjectedPopulationInt)
	assert.True(tt, er == nil)

	// verify the projectedPopulations as what we expected
	f := func(code string, ageRange string, year int, value int) error {
		v, _ := estimatedPopulations.At(code, ageRange, year)
		assert.True(tt, (value == v))
		fmt.Printf("XXX %s %s %d %d %d %v\n", code, ageRange, year, value, v, (value == v))
		return nil
	}
	cayvalues.Iterate(geoProjectedPopulationsInt, f)

	// make a map of LadCodes to weights - for example - weights['E09000002] == 0.19567854260348208
	weights := map[string]float64{}
	for _, v := range contributions {
		if _, ok := weights[v.LadCode]; !ok {
			weights[v.LadCode] = v.Share
		}
	}
	// apply the weights to the estimated populations to get the estimated contribution of each
	// lad to the catchment area still organized by code, ageRange and year
	weightedEstimatedPopulations := cayvalues.NewCayValues[int]()
	apw := func(code string, ageRange string, year int, value int) error {
		epop := value
		if w, ok := weights[code]; !ok {
			return fmt.Errorf("")
		} else {
			wpop := int(math.Round(float64(epop)*w + 1.0))
			err := weightedEstimatedPopulations.Add(code, ageRange, year, wpop)
			if err != nil {
				return err
			}
		}
		return nil
	}
	err := cayvalues.Iterate(estimatedPopulations, apw)
	assert.True(tt, err == nil)

	// sum the weightedEstimated population over lad so that there is one value for each (ageRange,year) combination.
	// do this by making a cayvalue.CayValue with a single code key
	oaPopulations := cayvalues.NewCayValues[int]()
	oaf := func(code string, ageRange string, year int, value int) error {
		v, ok := oaPopulations.At("oacode", ageRange, year)
		if !ok {
			oaPopulations.Set("oacode", ageRange, year, value)
		} else {
			oaPopulations.Set("oacode", ageRange, year, value+v)
		}
		return nil
	}
	err = cayvalues.Iterate(weightedEstimatedPopulations, oaf)
	assert.True(tt, err == nil)

	// as a santiy check - by comparing values in the weightedEstmatedPopulation object against values
	// in the estimatedPopulations object compute a weight/share for each ageRange, year combination.
	// See how close those weights are to the values in weights[code]
	deducedWeights := cayvalues.NewCayValues[float64]()
	dw := func(code string, ageRange string, year int, value int) error {
		wp, ok := weightedEstimatedPopulations.At(code, ageRange, year)
		if !ok {
			return fmt.Errorf("")
		}
		ep, ok := estimatedPopulations.At(code, ageRange, year)
		if !ok {
			return fmt.Errorf("")
		}
		w := float64(wp) / float64(ep)
		refW := weights[code]
		// fmt.Printf("\nFF  %f %f\n", w, refW)
		if math.Abs((w-refW)/w) > 0.1 {
			fmt.Printf("XXXXXXXXX   %f %f difference exceeds 15 percent\n", w, refW)
			assert.True(tt, false)
		}
		deducedWeights.Add(code, ageRange, year, w)
		return nil
	}
	err = cayvalues.Iterate(estimatedPopulations, dw)
	assert.True(tt, err == nil)

	oa2, er := oaApplyGrowthRatesAndWeights("oacode", contributions, growthRates, geoBaseYearProjectedPopulationInt)
	assert.True(tt, er == nil)

	fn := func(code string, ageRange string, year int, value int) error {
		v1, ok1 := oa2.At(code, ageRange, year)
		v2, ok2 := oaPopulations.At(code, ageRange, year)
		assert.True(tt, ok1)
		assert.True(tt, ok2)
		assert.True(tt, v1 == v2)
		return nil
	}
	cayvalues.Iterate(oaPopulations, fn)
	// fmt.Printf("%v", estimatedPopulations)
	// fmt.Printf("%v", growthRates)
	// fmt.Printf("%v", geoBaseYearProjectedPopulationInt)
	// fmt.Printf("%v", geoBaseYearProjectedPopulation)
	// fmt.Printf("%v", geoProjectedPopulations)
	// fmt.Printf("%v", geoProjectedPopulationsInt)
	// fmt.Printf("%v", projectedPopulations)
	fmt.Printf("%v", totalPopulations)
	fmt.Printf("%v", &contributions)
	fmt.Printf("%v", &ladPopMap)

}

func oaApplyGrowthRatesAndWeights(oaCode string, contributions []LadPopulation, growthRates cayvalues.CayValues[float64], baseYearCurrentPopulations cayvalues.CayValues[int]) (cayvalues.CayValues[int], error) {
	// apply the growth rates to the pojectedPopulations for the baseYear. Should give the projected populations
	estimatedPopulations, err := CalculateEstimatedPopulations(growthRates, baseYearCurrentPopulations)
	if err != nil {
		return nil, err
	}
	// make a map of LadCodes to weights - for example - weights['E09000002] == 0.19567854260348208
	weights := map[string]float64{}
	for _, v := range contributions {
		if _, ok := weights[v.LadCode]; !ok {
			weights[v.LadCode] = v.Share
		}
	}
	// apply the weights to the estimated populations to get the estimated contribution of each
	// lad to the catchment area still organized by code, ageRange and year
	weightedEstimatedPopulations := cayvalues.NewCayValues[int]()
	apw := func(code string, ageRange string, year int, value int) error {
		epop := value
		if w, ok := weights[code]; !ok {
			return fmt.Errorf("")
		} else {
			wpop := int(math.Round(float64(epop)*w + 1.0))
			err := weightedEstimatedPopulations.Add(code, ageRange, year, wpop)
			if err != nil {
				return err
			}
		}
		return nil
	}
	err = cayvalues.Iterate(estimatedPopulations, apw)
	if err != nil {
		return nil, err
	}

	// sum the weightedEstimated population over lad so that there is one value for each (ageRange,year) combination.
	// do this by making a cayvalue.CayValue with a single code key
	oaPopulations := cayvalues.NewCayValues[int]()
	oaf := func(code string, ageRange string, year int, value int) error {
		v, ok := oaPopulations.At(oaCode, ageRange, year)
		if !ok {
			oaPopulations.Set(oaCode, ageRange, year, value)
		} else {
			oaPopulations.Set(oaCode, ageRange, year, value+v)
		}
		return nil
	}
	err = cayvalues.Iterate(weightedEstimatedPopulations, oaf)
	if err != nil {
		return nil, err
	}
	return oaPopulations, nil
}
