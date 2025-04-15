package populations

import (
	"popmodel/populations/cayvalues"
	"popmodel/repo"
)

func CollateProjectedPopulationsByCode(pops []repo.LadPopulationProjection) (cayvalues.CayValues[int], error) {
	f := func(lap LadPopulationProjecttion) (string, string, int, int) {
		return lap.Code, lap.AgeRange, lap.Year, lap.TotalPopulation
	}
	return cayvalues.NewCayValuesByTransform(pops, f)
}
func CollateBaseYearProjectedPopulationsByCode(pops []repo.LadPopulationProjection) (cayvalues.CayValues[int], error) {
	f := func(lap LadPopulationProjecttion) (string, string, int, int) {
		return lap.Code, lap.AgeRange, lap.Year, lap.TotalPopulation
	}
	return cayvalues.NewCayValuesByTransform(pops, f)
}
