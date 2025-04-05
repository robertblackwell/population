package main

import (
	"fmt"
	"forecast_model/mockdb"
	"forecast_model/resolver"

	_ "github.com/lib/pq"
)

func main() {
	ctx := mockdb.Context{}
	db := mockdb.LoadMockDb()
	fmt.Printf("%v\n", db)
	fmt.Printf("bye")
	z, ok := resolver.GetProjectedPopulationByCodes(ctx, []string{"E06000002", "E06000003"}, 2020, 2, 40, 42, 3, true)
	// x, err := mockdb.GetProjectedPopulationByCodes(mockdb.Context{}, []string{"XX", "YY"}, 2020, 2, 40, 42, 3, true)
	if ok != nil {
		return
	}
	fmt.Printf("%v\n", z)

	// fmt.Printf("withIntermediaries\n")
	// withIntermediaries()
	// fmt.Printf("usingGrowthRatesAll\n")
	// usingGrowthRatesAll()
	// fmt.Printf("usingGrowthRatesSome\n")
	// usingGrowthRatesSomeYears()

	// inputData := model.PopulationForecastInputData_LoadFromFile("./model/prop-db.json", 40, 44)
	// fmt.Printf("%v\n", inputData)
	// gr := inputData.AnnualizedGrowthRatesByAge(45)
	// fmt.Printf("GrowthRates: %v\n", gr)
	// pop := inputData.ForecastPopulationByAgeAndNumberofYears(45, 20, 1200)
	// fmt.Printf("Population Forecast %v \n", pop)
	// pop2 := inputData.ForecastPopulationByAgeToHorizon(45, 1200, []int{2028, 2035})
	// fmt.Printf("Population Forecast 2 %v \n", pop2)

}

// func withIntermediaries() {
//     startYear := 2020
//     rangeSize := 5
//     futureOffset := 10
//     minAge := 40
//     maxAge := 42
//     includeIntermediaries := true
//     projected_pops, err := mockdb.GetProjectedPopulationByCodes(mockdb.Context{}, []string{"XX", "YY"}, startYear, rangeSize, minAge, maxAge, futureOffset, includeIntermediaries)
//     if err != nil {
//         return
//     }
//     ep := CalculateEstimatedPopulationWithBaseYear(projected_pops, 2023, 3880)
//     for ix, el := range projected_pops {
//         fmt.Printf("%d  %d %d \n", el.Year, el.TotalPopulation, ep[ix])
//     }
// }
// func usingGrowthRatesAll() {
//     startYear := 2020
//     rangeSize := 5
//     futureOffset := 10
//     minAge := 40
//     maxAge := 42
//     includeIntermediaries := true
//     projected_pops, err := mockdb.GetProjectedPopulationByCodes(mockdb.Context{}, []string{"XX", "YY"}, startYear, rangeSize, minAge, maxAge, futureOffset, includeIntermediaries)
//     if err != nil {
//         return
//     }
//     rates := CalculateGrowthRates(projected_pops, 2023)
//     years := make([]int, 0)
//     for i := 2018; i <= 2035; i++ {
//         years = append(years, i)
//     }
//     pops := CalculateEstimatedPopulationsForSomeYears(rates, 3880, years)
//     for ix, el := range projected_pops {
//         fmt.Printf("usingGrowthRates : %d  %d %d \n", el.Year, el.TotalPopulation, pops[ix].Population)
//     }
// }
// func usingGrowthRatesSomeYears() {
//     startYear := 2020
//     rangeSize := 5
//     futureOffset := 10
//     minAge := 40
//     maxAge := 42
//     includeIntermediaries := true
//     projected_pops, err := mockdb.GetProjectedPopulationByCodes(mockdb.Context{}, []string{"XX", "YY"}, startYear, rangeSize, minAge, maxAge, futureOffset, includeIntermediaries)
//     if err != nil {
//         return
//     }
//     rates := CalculateGrowthRates(projected_pops, 2023)
//     years := make([]int, 0)

//     years = append(years, 2020)
//     years = append(years, 2030)

//     pops := CalculateEstimatedPopulationsForSomeYears(rates, 3880, years)
//     el := projected_pops[2]
//     ix := 0
//     popix := pops[ix].Population
//     fmt.Printf("usingGrowthRatesSome : %d  %d %d \n", el.Year, el.TotalPopulation, popix)
//     el = projected_pops[12]
//     ix = 1
//     popix = pops[ix].Population
//     fmt.Printf("usingGrowthRatesSome : %d  %d %d \n", el.Year, el.TotalPopulation, popix)

// }

// type GrowthRate = struct {
//     Year             int
//     RateFromBaseYear float64
// }

// func CalculateGrowthRates(projected_pops []mockdb.LadPopulationProjection, baseYear int) []GrowthRate {
//     base_year_index, err := findBaseYearIndexInResult(projected_pops, 2023)
//     if err != nil {
//         panic("could not find base year in result")
//     }
//     growthRates := make([]GrowthRate, 0)
//     for ix, pp := range projected_pops {
//         var p float64
//         if ix == base_year_index {
//             p = 1.0
//         } else if ix < base_year_index {
//             p = float64(pp.TotalPopulation) / float64(projected_pops[base_year_index].TotalPopulation)
//         } else {
//             p = float64(pp.TotalPopulation) / float64(projected_pops[base_year_index].TotalPopulation)
//         }
//         growthRates = append(growthRates, GrowthRate{Year: pp.Year, RateFromBaseYear: p})
//     }
//     return growthRates
// }

// type EstimatedPopulation = struct {
//     Year       int
//     Population int
// }

// func CalculateEstimatedPopulationsForSomeYears(growthRates []GrowthRate, baseYearPopulation int, requiredYear []int) []EstimatedPopulation {

//     result := make([]EstimatedPopulation, 0)
//     for _, gr := range growthRates {
//         if slices.Contains(requiredYear, gr.Year) {
//             p := int(math.Round(float64(baseYearPopulation) * gr.RateFromBaseYear))
//             result = append(result, EstimatedPopulation{Year: gr.Year, Population: p})
//         }
//     }
//     return result

// }

// func CalculateEstimatedPopulationWithBaseYear(projected_pops []mockdb.LadPopulationProjection, baseYear int, currentPopulationForBaseYear int) []int {
//     base_year_index, err := findBaseYearIndexInResult(projected_pops, 2023)
//     if err != nil {
//         panic("could not find base year in result")
//     }
//     estimated_populations := make([]int, 0)
//     for ix, pp := range projected_pops {
//         var p int
//         if ix == base_year_index {
//             p = currentPopulationForBaseYear
//         } else if ix < base_year_index {
//             x := float64(pp.TotalPopulation) / float64(projected_pops[base_year_index].TotalPopulation) * float64(currentPopulationForBaseYear)
//             p = int(math.Round(x))
//         } else {
//             x := float64(pp.TotalPopulation) / float64(projected_pops[base_year_index].TotalPopulation) * float64(currentPopulationForBaseYear)
//             p = int(math.Round(x))
//         }
//         estimated_populations = append(estimated_populations, p)
//     }
//     return estimated_populations
// }

// func findBaseYearIndexInResult(result []mockdb.LadPopulationProjection, base_year int) (int, error) {
//     for i, el := range result {
//         if el.Year == base_year {
//             return i, nil
//         }
//     }
//     return 0, errors.New("could not find 2023")
// }
