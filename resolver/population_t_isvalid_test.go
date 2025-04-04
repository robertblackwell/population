package resolver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidPopVec(t *testing.T) {
	projectedPops1 := []LadPopulationProjection{
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1000, Year: 2020},
		{Code: "LAD2", Type: "type", AgeRange: "0-4", TotalPopulation: 2000, Year: 2030},
	}
	assert.False(t, IsValidPopVec(projectedPops1))
	projectedPops2 := []LadPopulationProjection{
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1000, Year: 2020},
		{Code: "LAD1", Type: "type2", AgeRange: "0-4", TotalPopulation: 2000, Year: 2030},
	}
	assert.False(t, IsValidPopVec(projectedPops2))
	projectedPops3 := []LadPopulationProjection{
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1000, Year: 2020},
		{Code: "LAD1", Type: "type2", AgeRange: "1-4", TotalPopulation: 2000, Year: 2030},
	}
	assert.False(t, IsValidPopVec(projectedPops3))
	projectedPops4 := []LadPopulationProjection{
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1000, Year: 2020},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 2000, Year: 2030},
	}
	assert.True(t, IsValidPopVec(projectedPops4))
}
func TestValidMapPopVec(t *testing.T) {
	pm1 := map[string][]LadPopulationProjection{"LAD3": {
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1000, Year: 2020},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 2000, Year: 2030},
	}}
	assert.False(t, IsValidMapOfPopVec(pm1))
	pm2 := map[string][]LadPopulationProjection{"LAD1": {
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1000, Year: 2020},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 2000, Year: 2030},
	}}
	assert.True(t, IsValidMapOfPopVec(pm2))
}
