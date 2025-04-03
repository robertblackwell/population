package resolver

import (
	"fmt"
	"forecast_model/mockdb"
	"testing"
)

func TestGetProjectedPopulation(t *testing.T) {
	ctx := mockdb.Context{}
	z, ok := GetProjectedPopulationByCodes(ctx, []string{"E06000002", "E06000003"}, 2020, 2, 40, 42, 3, true)
	if ok != nil {
		return
	}
	fmt.Printf("%v\n", z)
}
func TestGetAllProjectedPopulation(t *testing.T) {
	ctx := mockdb.Context{}
	z, ok := GetAllProjectedPopulationsByCodes(ctx, []string{"E06000002", "E06000003"}, 40, 42)
	// x, err := mockdb.GetProjectedPopulationByCodes(mockdb.Context{}, []string{"XX", "YY"}, 2020, 2, 40, 42, 3, true)
	if ok != nil {
		return
	}
	fmt.Printf("%v\n", z)
}

// Test that we can get back the projected_population values by setting the base years population to its projected value
// this is a good sanity check.
func TestSingleCodeVerify(t *testing.T) {
	ctx := mockdb.Context{}
	code := "E06000002"
	projected_pops, ok := GetAllProjectedPopulationsByCodes(ctx, []string{"E06000002"}, 40, 42)
	if ok != nil {
		t.Errorf("GetAllProjectedPopulationByCode failed")
	}
	rates := GrowthRatesAllYears(projected_pops[code], 2023)
	years := make([]int, 0)
	for i := 2018; i <= 2035; i++ {
		years = append(years, i)
	}
	basePop, ok := FindBaseYearProjectedPopulation(projected_pops[code], 2023)
	if ok != nil {
		t.Errorf("something went wrong")
	}

	pops := CalculateEstimatedPopulationsForSomeYears(rates, basePop, years)
	if ok != nil {
		t.Errorf("%v", ok)
	}

	for ix, el := range projected_pops[code] {
		fmt.Printf("usingGrowthRates : %d  %d %d \n", el.Year, el.TotalPopulation, pops[ix].Population)
		if el.TotalPopulation != pops[ix].Population {
			t.Errorf("Failed code:%s ix: %d year: %d pop1: %d pop2:%d", code, ix, el.Year, el.TotalPopulation, pops[ix].Population)
		}
	}
}

func TestGetBackTheProjections(t *testing.T) {
	ctx := mockdb.Context{}
	code := "E06000002"
	projected_pops, ok := GetAllProjectedPopulationsByCodes(ctx, []string{"E06000002"}, 40, 42)
	// x, err := mockdb.GetProjectedPopulationByCodes(mockdb.Context{}, []string{"XX", "YY"}, 2020, 2, 40, 42, 3, true)
	if ok != nil {
		return
	}
	rates := GrowthRatesAllYearsMultipleCodes(projected_pops, 2023)
	years := make([]int, 0)
	for i := 2018; i <= 2035; i++ {
		years = append(years, i)
	}
	baseYearPopulation := map[string]int{code: 3880}
	pops, ok := CalculateEstimatedPopulationsForSomeYearsMultipleCodes(rates, baseYearPopulation, years)
	if ok != nil {
		t.Errorf("%v", ok)
	}

	for ix, el := range projected_pops[code] {
		fmt.Printf("usingGrowthRates : %d  %d %d \n", el.Year, el.TotalPopulation, pops[code][ix].Population)
		if el.TotalPopulation != pops[code][ix].Population {
			t.Errorf("Failed code:%s ix: %d year: %d pop1: %d pop2:%d", code, ix, el.Year, el.TotalPopulation, pops[code][ix].Population)
		}
	}
}
