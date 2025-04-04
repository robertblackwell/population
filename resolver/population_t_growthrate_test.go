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

// requireIntermediates false only one code
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
func Test02(t *testing.T) {
	projectedPops1 := []LadPopulationProjection{
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1000, Year: 2020},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 2000, Year: 2030},
	}
	projectedPops2 := []LadPopulationProjection{
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
	baseYearPop := LadPopulationProjection{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1500, Year: 2023}

	gr := CalculateGrowthRatesRelativeToBaseYear(projectedPops1, baseYearPop)
	estimated_pop := CalculateEstimatedPopulationsForSomeYears(gr, baseYearPop.TotalPopulation, []int{2020, 2030})

	// without intermediaries we recompute the projected population values
	assert.Equal(t, estimated_pop[0].Population, projectedPops1[0].TotalPopulation, "2020 populations are the same")
	assert.Equal(t, estimated_pop[1].Population, projectedPops1[1].TotalPopulation, "2030 populations are the same")

	gr2 := CalculateGrowthRatesRelativeToBaseYear(projectedPops2, baseYearPop)
	estimated_pop2 := CalculateEstimatedPopulationsForSomeYears(gr2, baseYearPop.TotalPopulation, []int{2020, 2030})

	// with intermediaries also recompute the projected population values
	assert.Equal(t, estimated_pop2[0].Population, projectedPops2[2].TotalPopulation, "2020 populations are the same")
	assert.Equal(t, estimated_pop2[1].Population, projectedPops2[12].TotalPopulation, "2030 populations are the same")
	// both calculations get the same result
	assert.Equal(t, estimated_pop2[0].Population, estimated_pop[0].Population, "2020 populations are the same")
	assert.Equal(t, estimated_pop2[1].Population, estimated_pop[1].Population, "2030 populations are the same")

	gr3 := CalculateGrowthRatesRelativeToBaseYear(projectedPops2, baseYearPop)
	estimated_pop3 := CalculateEstimatedPopulationsForSomeYears(gr3, baseYearPop.TotalPopulation, []int{2018, 2019, 2020, 2021, 2022, 2023, 2024, 2025, 2026, 2027, 2028, 2029, 2030, 2031, 2032, 2033, 2034, 2035})
	assert.Equal(t, estimated_pop3[2].Population, estimated_pop[0].Population, "2020 populations are the same")
	assert.Equal(t, estimated_pop3[12].Population, estimated_pop[1].Population, "2030 populations are the same")

	fmt.Printf("%v\n", estimated_pop)
}
