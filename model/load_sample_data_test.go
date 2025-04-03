package model

import (
	"fmt"
	"testing"
)

func TestJsonLoad(t *testing.T) {
	inputData := PopulationForecastInputData_LoadFromFile("./prop-db.json", 40, 50)
	fmt.Printf("%v\n", inputData)
	gr := inputData.AnnualizedGrowthRatesByAge(45)
	fmt.Printf("GrowthRates: %v\n", gr)
	pop := inputData.ForecastPopulationByAgeAndNumberofYears(45, 20, 1200)
	fmt.Printf("Population Forecast %v \n", pop)
	pop2 := inputData.ForecastPopulationByAgeToHorizon(45, 1200, []int{2028, 2035})
	fmt.Printf("Population Forecast 2 %v \n", pop2)
}

// func GetPopulationProjectionsByCodes(ctx Context, parents []String, startYear int, rangeSize int, minAge int, maxAge int)
// {
// 	inputData := PopulationForecastInputData_LoadFromFile("./prop-db.json");
// 	for
// 	for age := minAge; age <= maxAge; age++ {
// 		d := inputData[]
// 	}

// }
