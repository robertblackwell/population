package populations

import (
	"fmt"
	"testing"

	"popmodel/populations/cayvalues"
	"popmodel/repo"

	"github.com/stretchr/testify/assert"
)

// // Test that forecast populations are the same as projected populations when the base year starting population
// // is equal to its projected population.
// // With Intermediates
// func TestBaseCaseVerifyWithIntermediates(t *testing.T) {
// 	ctx := mockdb.Context{}
// 	code := "E06000002"
// 	ar := "20-24"
// 	baseYearNum := 2023

// 	// get a slice of LadPopulationProjections
// 	pp, er := GetProjectedPopulationsByCodes(ctx, []string{code}, 2020, 5, 20, 60, 10, true)
// 	assert.True(t, er == nil)

// 	// convert the slice of LadPopulationProjections into a structure more amenable to the rest of the processing
// 	mpv := CollateProjectedPopulationsByCode(pp)
// 	assert.True(t, IsValidMapOfPopVec(mpv))

// 	// get a slice of LadPopulationProjection for a single code and single age range.
// 	// thats the input to the basic growth rate calculation for population forecase
// 	projectedPopsForOneCodeOneAgeRange, ok := Index2LevelMap(mpv, code, ar)
// 	assert.True(t, ok)
// 	assert.True(t, IsValidPopVec(projectedPopsForOneCodeOneAgeRange))
// 	// next line proves intermediates
// 	assert.True(t, len(projectedPopsForOneCodeOneAgeRange) == 11)

// 	// calculation of growth rates requires the population projection for the same code, same ageRange
// 	// and baseYear (usually 2023). All growth rates are relative to the population in that base year.
// 	// Sometimes that data will be available in the slice projectedPopsForOneCodeOneAgeRange
// 	// but often not. So there is a function that will get it independently.
// 	// Notice it required ageRange specification
// 	baseYear, er := GetBaseYearProjectedPopulations(code, ar, baseYearNum, 5, 20, 60)
// 	assert.True(t, er == nil)

// 	// now calculate the growth rates, relative to baseYear for one code one ageRange
// 	rates := CalculateGrowthRatesBaseCase(projectedPopsForOneCodeOneAgeRange, baseYear)

// 	// Perform a population forecast. Such forecasts require starting population value for the base year.
// 	// In this case we are using the projected population for the baseYear; if all is working this should result
// 	// in our forecast populations being the same as the forecast populations.
// 	pops := CalculateEstimatedPopulationsBaseCase(rates, baseYear.TotalPopulation)

// 	for ix, el := range projectedPopsForOneCodeOneAgeRange {
// 		fmt.Printf("usingGrowthRates : %d  %d %d \n", el.Year, el.TotalPopulation, pops[ix].Population)
// 		if el.TotalPopulation != pops[ix].Population {
// 			t.Errorf("Failed code:%s ix: %d year: %d pop1: %d pop2:%d", code, ix, el.Year, el.TotalPopulation, pops[ix].Population)
// 		}
// 	}
// }

// // Test that forecast populations are the same as projected populations when the base year starting population
// // is equal to its projected population.
// // This is a pretty obvious sanity test
// // Without Intermediates
// func TestBaseCaseVerifyWithOutIntermediates(t *testing.T) {
// 	ctx := mockdb.Context{}
// 	code := "E06000002"
// 	ar := "20-24"
// 	baseYearNum := 2023

// 	// get a slice of LadPopulationProjections
// 	pp, er := GetProjectedPopulationsByCodes(ctx, []string{code}, 2020, 5, 20, 60, 10, false)
// 	assert.True(t, er == nil)

// 	// convert the slice of LadPopulationProjections into a structure more amenable to the rest of the processing
// 	mpv := CollateProjectedPopulationsByCode(pp)
// 	assert.True(t, IsValidMapOfPopVec(mpv))

// 	// get a slice of LadPopulationProjection for a single code and single age range.
// 	// thats the input to the basic growth rate calculation for population forecase
// 	projectedPopsForOneCodeOneAgeRange, ok := Index2LevelMap(mpv, code, ar)
// 	assert.True(t, ok)
// 	assert.True(t, IsValidPopVec(projectedPopsForOneCodeOneAgeRange))
// 	// next line proves NO intermediates
// 	assert.True(t, len(projectedPopsForOneCodeOneAgeRange) == 2)

// 	// calculation of growth rates requires the population projection for the same code, same ageRange
// 	// and baseYear (usually 2023). All growth rates are relative to the population in that base year.
// 	// Sometimes that data will be available in the slice projectedPopsForOneCodeOneAgeRange
// 	// but often not. So there is a function that will get it independently.
// 	// Notice it required ageRange specification
// 	baseYear, er := GetBaseYearProjectedPopulations(code, ar, baseYearNum, 5, 20, 60)
// 	assert.True(t, er == nil)

// 	// now calculate the growth rates, relative to baseYear for one code one ageRange
// 	rates := CalculateGrowthRatesBaseCase(projectedPopsForOneCodeOneAgeRange, baseYear)

// 	// Perform a population forecast. Such forecasts require starting population value for the base year.
// 	// In this case we are using the projected population for the baseYear; if all is working this should result
// 	// in our forecast populations being the same as the forecast populations.
// 	pops := CalculateEstimatedPopulationsBaseCase(rates, baseYear.TotalPopulation)

// 	// check the expected result
// 	for ix, el := range projectedPopsForOneCodeOneAgeRange {
// 		fmt.Printf("usingGrowthRates : %d  %d %d \n", el.Year, el.TotalPopulation, pops[ix].Population)
// 		if el.TotalPopulation != pops[ix].Population {
// 			t.Errorf("Failed code:%s ix: %d year: %d pop1: %d pop2:%d", code, ix, el.Year, el.TotalPopulation, pops[ix].Population)
// 		}
// 	}
// }

// requireIntermediates false only one code

func K[T any, U any](t T, u U) {
	fmt.Printf("%v %v", t, u)
}
func Test01(t *testing.T) {
	projectedPops := []repo.LadPopulationProjection{
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1000, Year: 2020},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 2000, Year: 2030},
	}
	cayProjectedPops, er := CollateProjectedPopulationsByCode(projectedPops)
	assert.True(t, er == nil)

	baseYearPop := []repo.LadPopulationProjection{
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1500, Year: 2023},
	}
	cayBaseYearPop, er := CollateProjectedPopulationsByCode(baseYearPop)
	assert.True(t, er == nil)

	gr, er := CalculateGrowthRates(cayProjectedPops, cayBaseYearPop)
	assert.True(t, er == nil)

	cayEstimatedPop, er := CalculateEstimatedPopulations(gr, cayBaseYearPop)
	assert.True(t, er == nil)

	cEsti2020, ok := cayEstimatedPop.At("LAD1", "0-4", 2020)
	assert.True(t, ok)
	cProj2020, ok := cayProjectedPops.At("LAD1", "0-4", 2020)
	assert.True(t, ok)
	assert.Equal(t, cEsti2020, cProj2020, "2020 populations are the same")

	cEsti2030, ok := cayEstimatedPop.At("LAD1", "0-4", 2030)
	assert.True(t, ok)
	cProj2030, ok := cayProjectedPops.At("LAD1", "0-4", 2030)
	assert.True(t, ok)
	assert.Equal(t, cEsti2030, cProj2030, "2030 populations are the same")
	fmt.Printf("%v\n", cayEstimatedPop)
}
func TestBaseCaseExplicitData(t *testing.T) {
	// The term BaseCase is used when growth rates are being calculated and applied
	// for only a single code and a single ageRange.
	// In which case all entries in the slice []LadPopulationProjections have the same value
	// for Code, Type, AgeRange.
	// In the calcs below the variables of types []GrowthRate and []EstimatedPopulation
	// do not have fields for Code, Type, AgeRange but the common const values in
	// the variables of type []LadPopulationProjection are assumed
	//
	// This is a mockup of data from a database query where includeIntermediates == false
	projectedPopsNoIntermediates := []repo.LadPopulationProjection{
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1000, Year: 2020},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 2000, Year: 2030},
	}
	cayProjectedPopsNoIntermediates, er := CollateProjectedPopulationsByCode(projectedPopsNoIntermediates)
	assert.True(t, er == nil)

	baseYearPop := []repo.LadPopulationProjection{
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1500, Year: 2023},
	}
	cayBaseYearPop, er := CollateProjectedPopulationsByCode(baseYearPop)
	assert.True(t, er == nil)

	gr, er := CalculateGrowthRates(cayProjectedPopsNoIntermediates, cayBaseYearPop)
	assert.True(t, er == nil)

	cayEstimatedPopsNoIntermediates, er := CalculateEstimatedPopulations(gr, cayBaseYearPop)
	assert.True(t, er == nil)

	// without intermediaries the estimated population calculated using growth rates
	// are the same as the projected population values
	cEstimatedNoInt2020, ok := cayEstimatedPopsNoIntermediates.At("LAD1", "0-4", 2020)
	assert.True(t, ok)
	cProjectedNoInt2020, ok := cayProjectedPopsNoIntermediates.At("LAD1", "0-4", 2020)
	assert.True(t, ok)

	assert.Equal(t, cEstimatedNoInt2020, cProjectedNoInt2020, "2020 populations are the same")

	cEstimatedNoInt2030, ok := cayEstimatedPopsNoIntermediates.At("LAD1", "0-4", 2030)
	assert.True(t, ok)
	cProjectedNoInt2030, ok := cayProjectedPopsNoIntermediates.At("LAD1", "0-4", 2030)
	assert.True(t, ok)
	assert.Equal(t, cEstimatedNoInt2030, cProjectedNoInt2030, "2030 populations are the same")

	// this is a mockup of a database query where includeIntermediates == true
	projectedPopsWithIntermediates := []repo.LadPopulationProjection{
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 821, Year: 2018},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 933, Year: 2019},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1000, Year: 2020}, // <<====== index 2
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1602, Year: 2021},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1378, Year: 2022},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1500, Year: 2023}, // <======== index 5
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1501, Year: 2024},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1540, Year: 2025},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1590, Year: 2026},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1490, Year: 2027},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1480, Year: 2028},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1501, Year: 2029},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 2000, Year: 2030}, //<<====== index 12
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1900, Year: 2031},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1890, Year: 2032},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1880, Year: 2033},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1870, Year: 2034},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1860, Year: 2035},
	}
	cayProjectedPopsWithIntermediates, er := CollateProjectedPopulationsByCode(projectedPopsWithIntermediates)
	assert.True(t, er == nil)

	baseYearPopWith := []repo.LadPopulationProjection{
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1500, Year: 2023},
	}
	cayBaseYearPopWith, er := CollateProjectedPopulationsByCode(baseYearPopWith)
	assert.True(t, er == nil)

	grWith, er := CalculateGrowthRates(cayProjectedPopsWithIntermediates, cayBaseYearPopWith)
	assert.True(t, er == nil)

	cayEstimatedPopsWithIntermediates, er := CalculateEstimatedPopulations(grWith, cayBaseYearPopWith)
	assert.True(t, er == nil)

	// with intermediaries the estimated population calculated using growth rates
	// are the same as the projected population values
	cE2020, ok := cayEstimatedPopsWithIntermediates.At("LAD1", "0-4", 2020)
	assert.True(t, ok)
	cP2020, ok := cayProjectedPopsWithIntermediates.At("LAD1", "0-4", 2020)
	assert.True(t, ok)
	assert.Equal(t, cE2020, cP2020, "2020 populations are the same")

	cE2030, ok := cayEstimatedPopsWithIntermediates.At("LAD1", "0-4", 2030)
	assert.True(t, ok)
	cP2030, ok := cayProjectedPopsWithIntermediates.At("LAD1", "0-4", 2030)
	assert.True(t, ok)
	assert.Equal(t, cE2030, cP2030, "2030 populations are the same")

	// with intermediaries the estimated populations calculated using growth rates
	// are the same as the projected population values
	f := func(code string, ageRange string, year int, value int) (int, error) {
		x, ok := cayProjectedPopsWithIntermediates.At(code, ageRange, year)
		assert.True(t, ok)
		assert.True(t, x == value)
		return 0, nil
	}
	_, er2 := cayvalues.Map(cayEstimatedPopsWithIntermediates, f)
	assert.True(t, er2 == nil)
	// assert.True(t, len(estimatedPopsWithIntermediates) == len(projectedPopsWithIntermediates))
	// for ix := range estimatedPopsWithIntermediates {
	// 	assert.Equal(t, estimatedPopsWithIntermediates[ix].Population, projectedPopsWithIntermediates[ix].TotalPopulation, "2020 populations are the same")
	// }
	// // finally the results for the years 2020 and 2030 are the same with or with out intermediaries
	// // Warning: determined the matching indices manually
	// assert.Equal(t, estimatedPopsWithIntermediates[2].Population, estimatedPopsNoIntermediates[0].Population, "2020 populations are the same")
	// assert.Equal(t, estimatedPopsWithIntermediates[12].Population, estimatedPopsNoIntermediates[1].Population, "2030 populations are the same")
}
