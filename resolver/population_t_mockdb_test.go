package resolver

import (
	"fmt"
	"forecast_model/mockdb"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProjectedPopulation(t *testing.T) {
	ctx := mockdb.Context{}
	m, ok := GetProjectedPopulationByCodes(ctx, []string{"E06000002", "E06000003"}, 2020, 2, 40, 42, 3, true)
	if ok != nil {
		return
	}
	assert.True(t, IsValidMapOfPopVec(z))
	fmt.Printf("%v\n", m)
}
func TestGetAllProjectedPopulation(t *testing.T) {
	ctx := mockdb.Context{}
	m, ok := GetAllProjectedPopulationsByCodes(ctx, []string{"E06000002", "E06000003"}, 40, 42)
	// x, err := mockdb.GetProjectedPopulationByCodes(mockdb.Context{}, []string{"XX", "YY"}, 2020, 2, 40, 42, 3, true)
	if ok != nil {
		return
	}
	assert.True(t, IsValidMapOfPopVec(m))
	fmt.Printf("%v\n", m)
}
