package resolver

import (
	"fmt"
	"forecast_model/mockdb"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test that forecast populations are the same as projected populations when the base year starting population
// is equal to its projected population
func TestSingleCodeVerify(t *testing.T) {
	ctx := mockdb.Context{}
	code := "E06000002"
	projected_pops, ok := GetAllProjectedPopulationsByCodes(ctx, []string{"E06000002"}, 40, 42)
	if ok != nil {
		t.Errorf("GetAllProjectedPopulationByCode failed")
	}
	assert.True(t, IsValidMapOfPopVec(projected_pops))

	rates := CalculateGrowthRatesAllYears(projected_pops[code], 2023)
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
		t.Errorf("get all projected pops failed")
	}
	baseYearPops, e := GetProjectedPopulationsByCodeForBaseYear(ctx, []string{"E06000002"}, 2023, 40, 42)
	if e != nil {
		t.Errorf("get baseYear projected pops failed")
	}
	fmt.Printf("%v\n", baseYearPops)
	rates := CalculateGrowthRatesAllYearsMultipleCodes(projected_pops, 2023)
	allYears := make([]int, 0)
	for i := 2018; i <= 2035; i++ {
		allYears = append(allYears, i)
	}
	tmp := baseYearPops[code].TotalPopulation
	baseYearPopulation := map[string]int{code: tmp}
	pops, ok := CalculateEstimatedPopulationsForSomeYearsMultipleCodes(rates, baseYearPopulation, allYears)
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
func Test01(t *testing.T) {
	projectedPops := []LadPopulationProjection{
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1000, Year: 2020},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 2000, Year: 2030},
	}
	baseYearPop := LadPopulationProjection{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1500, Year: 2023}

	gr := CalculateGrowthRatesRelativeToBaseYear(projectedPops, baseYearPop)
	estimated_pop := CalculateEstimatedPopulationsForSomeYears(gr, 1500, []int{2020, 2030})

	assert.Equal(t, estimated_pop[0].Population, projectedPops[0].TotalPopulation, "2020 populations are the same")
	assert.Equal(t, estimated_pop[1].Population, projectedPops[1].TotalPopulation, "2030 populations are the same")
	fmt.Printf("%v\n", estimated_pop)

}
