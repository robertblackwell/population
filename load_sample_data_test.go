package model

import (
	"fmt"
	"testing"
)

func TestJsonLoad(t *testing.T) {
	inputData := PopulationForecastInputData_LoadFromFile("./sample_forecast_data.json")
	fmt.Printf("%v\n", inputData)
	gr := inputData.AnnualizedGrowthRatesByAge(45)
	fmt.Printf("GrowthRates: %v\n", gr)
	pop := inputData.ForecastPopulationByAgeAndNumberofYears(45, 20, 1200)
	fmt.Printf("Population Forecast %v \n", pop)
	pop2 := inputData.ForecastPopulationByAgeToHorizon(45, 1200, []int{2028, 2035})
	fmt.Printf("Population Forecast 2 %v \n", pop2)
}
