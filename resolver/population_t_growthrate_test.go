package resolver

import (
	"fmt"
	"forecast_model/mockdb"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test that forecast populations are the same as projected populations when the base year starting population
// is equal to its projected population.
// With Intermediates
func TestBaseCaseVerifyWithIntermediates(t *testing.T) {
	ctx := mockdb.Context{}
	code := "E06000002"
	ar := "20-24"
	baseYearNum := 2023

	// get a slice of LadPopulationProjections
	pp, er := GetProjectedPopulationsByCodes(ctx, []string{code}, 2020, 5, 20, 60, 10, true)
	assert.True(t, er == nil)

	// convert the slice of LadPopulationProjections into a structure more amenable to the rest of the processing
	mpv := CollateProjectedPopulationsByCode(pp)
	assert.True(t, IsValidMapOfPopVec(mpv))

	// get a slice of LadPopulationProjection for a single code and single age range.
	// thats the input to the basic growth rate calculation for population forecase
	projectedPopsForOneCodeOneAgeRange, ok := Index2LevelMap(mpv, code, ar)
	assert.True(t, ok)
	assert.True(t, IsValidPopVec(projectedPopsForOneCodeOneAgeRange))
	// next line proves intermediates
	assert.True(t, len(projectedPopsForOneCodeOneAgeRange) == 11)

	// calculation of growth rates requires the population projection for the same code, same ageRange
	// and baseYear (usually 2023). All growth rates are relative to the population in that base year.
	// Sometimes that data will be available in the slice projectedPopsForOneCodeOneAgeRange
	// but often not. So there is a function that will get it independently.
	// Notice it required ageRange specification
	baseYear, er := GetBaseYearProjectedPopulations(code, ar, baseYearNum, 5, 20, 60)
	assert.True(t, er == nil)

	// now calculate the growth rates, relative to baseYear for one code one ageRange
	rates := CalculateGrowthRatesBaseCase(projectedPopsForOneCodeOneAgeRange, baseYear)

	// Perform a population forecast. Such forecasts require starting population value for the base year.
	// In this case we are using the projected population for the baseYear; if all is working this should result
	// in our forecast populations being the same as the forecast populations.
	pops := CalculateEstimatedPopulationsForSomeYears(rates, baseYear.TotalPopulation)

	for ix, el := range projectedPopsForOneCodeOneAgeRange {
		fmt.Printf("usingGrowthRates : %d  %d %d \n", el.Year, el.TotalPopulation, pops[ix].Population)
		if el.TotalPopulation != pops[ix].Population {
			t.Errorf("Failed code:%s ix: %d year: %d pop1: %d pop2:%d", code, ix, el.Year, el.TotalPopulation, pops[ix].Population)
		}
	}
}

// Test that forecast populations are the same as projected populations when the base year starting population
// is equal to its projected population.
// This is a pretty obvious sanity test
// Without Intermediates
func TestBaseCaseVerifyWithOutIntermediates(t *testing.T) {
	ctx := mockdb.Context{}
	code := "E06000002"
	ar := "20-24"
	baseYearNum := 2023

	// get a slice of LadPopulationProjections
	pp, er := GetProjectedPopulationsByCodes(ctx, []string{code}, 2020, 5, 20, 60, 10, false)
	assert.True(t, er == nil)

	// convert the slice of LadPopulationProjections into a structure more amenable to the rest of the processing
	mpv := CollateProjectedPopulationsByCode(pp)
	assert.True(t, IsValidMapOfPopVec(mpv))

	// get a slice of LadPopulationProjection for a single code and single age range.
	// thats the input to the basic growth rate calculation for population forecase
	projectedPopsForOneCodeOneAgeRange, ok := Index2LevelMap(mpv, code, ar)
	assert.True(t, ok)
	assert.True(t, IsValidPopVec(projectedPopsForOneCodeOneAgeRange))
	// next line proves NO intermediates
	assert.True(t, len(projectedPopsForOneCodeOneAgeRange) == 2)

	// calculation of growth rates requires the population projection for the same code, same ageRange
	// and baseYear (usually 2023). All growth rates are relative to the population in that base year.
	// Sometimes that data will be available in the slice projectedPopsForOneCodeOneAgeRange
	// but often not. So there is a function that will get it independently.
	// Notice it required ageRange specification
	baseYear, er := GetBaseYearProjectedPopulations(code, ar, baseYearNum, 5, 20, 60)
	assert.True(t, er == nil)

	// now calculate the growth rates, relative to baseYear for one code one ageRange
	rates := CalculateGrowthRatesBaseCase(projectedPopsForOneCodeOneAgeRange, baseYear)

	// Perform a population forecast. Such forecasts require starting population value for the base year.
	// In this case we are using the projected population for the baseYear; if all is working this should result
	// in our forecast populations being the same as the forecast populations.
	pops := CalculateEstimatedPopulationsForSomeYears(rates, baseYear.TotalPopulation)

	// check the expected result
	for ix, el := range projectedPopsForOneCodeOneAgeRange {
		fmt.Printf("usingGrowthRates : %d  %d %d \n", el.Year, el.TotalPopulation, pops[ix].Population)
		if el.TotalPopulation != pops[ix].Population {
			t.Errorf("Failed code:%s ix: %d year: %d pop1: %d pop2:%d", code, ix, el.Year, el.TotalPopulation, pops[ix].Population)
		}
	}
}

// requireIntermediates false only one code
func Test01(t *testing.T) {
	projectedPops := []LadPopulationProjection{
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1000, Year: 2020},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 2000, Year: 2030},
	}
	baseYearPop := LadPopulationProjection{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1500, Year: 2023}

	gr := CalculateGrowthRatesBaseCase(projectedPops, baseYearPop)
	estimated_pop := CalculateEstimatedPopulationsForSomeYears(gr, 1500)

	assert.Equal(t, estimated_pop[0].Population, projectedPops[0].TotalPopulation, "2020 populations are the same")
	assert.Equal(t, estimated_pop[1].Population, projectedPops[1].TotalPopulation, "2030 populations are the same")
	fmt.Printf("%v\n", estimated_pop)
}
func TestBaseCaseExplicitData(t *testing.T) {
	// The term BaseCase is used when growth rates are being calculated and applied
	// for only a single code and a single ageRange.
	// In which case all entries in the slice []LadPopulationProjections has the same value
	// for Code, TYpe, AgeRange.
	// In the calcs below the variables of types []GrowthRate and []EstimatedPopulation
	// do not have fields for Code, Type, AgeRange but the common const values in
	// the variables of type []LadPopulationProjection
	//
	// this is a mockup of data from a database query where includeIntermediates == false
	projectedPopsNoIntermediates := []LadPopulationProjection{
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1000, Year: 2020},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 2000, Year: 2030},
	}
	// this verifies thatCode, Type, AgeRange fields all have the same value
	// that is a prerequisite for growth and estimated population calculations
	assert.True(t, IsValidPopVec(projectedPopsNoIntermediates))

	// this is a mockup of a database query where includeIntermediates == true
	projectedPopsWithIntermediates := []LadPopulationProjection{
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
	// this verifies thatCode, Type, AgeRange fields all have the same value
	// that is a prerequisite for growth and estimated population calculations
	assert.True(t, IsValidPopVec(projectedPopsWithIntermediates))

	baseYearPop := LadPopulationProjection{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1500, Year: 2023}
	gr := CalculateGrowthRatesBaseCase(projectedPopsNoIntermediates, baseYearPop)
	estimatedPopsNoIntermediates := CalculateEstimatedPopulationsForSomeYears(gr, baseYearPop.TotalPopulation)
	// without intermediaries the estimated population calculated using growth rates
	// are the same as the projected population values
	assert.Equal(t, estimatedPopsNoIntermediates[0].Population, projectedPopsNoIntermediates[0].TotalPopulation, "2020 populations are the same")
	assert.Equal(t, estimatedPopsNoIntermediates[1].Population, projectedPopsNoIntermediates[1].TotalPopulation, "2030 populations are the same")

	gr2 := CalculateGrowthRatesBaseCase(projectedPopsWithIntermediates, baseYearPop)
	estimatedPopsWithIntermediates := CalculateEstimatedPopulationsForSomeYears(gr2, baseYearPop.TotalPopulation)
	// with intermediaries the estimated populations calculated using growth rates
	// are the same as the projected population values
	assert.True(t, len(estimatedPopsWithIntermediates) == len(projectedPopsWithIntermediates))
	for ix := range estimatedPopsWithIntermediates {
		assert.Equal(t, estimatedPopsWithIntermediates[ix].Population, projectedPopsWithIntermediates[ix].TotalPopulation, "2020 populations are the same")
	}
	// finally the results for the years 2020 and 2030 are the same with or with out intermediaries
	// Warning: determined the matching indices manually
	assert.Equal(t, estimatedPopsWithIntermediates[2].Population, estimatedPopsNoIntermediates[0].Population, "2020 populations are the same")
	assert.Equal(t, estimatedPopsWithIntermediates[12].Population, estimatedPopsNoIntermediates[1].Population, "2030 populations are the same")
}
