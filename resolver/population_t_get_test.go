package resolver

import (
	"fmt"
	"forecast_model/mockdb"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test that forecast populations are the same as projected populations when the base year starting population
// is equal to its projected population
func TestGetPop01(t *testing.T) {
	ctx := mockdb.Context{}
	codes := []string{"E06000002"}
	x, err := GetProjectedPopulationsByCodes(ctx, codes, 2020, 5, 20, 60, 10, true)
	fmt.Printf("%v", err)
	fmt.Printf("%v", x)
	y := CollateProjectedPopulationsByCode(x)
	fmt.Printf("%v", y)
	_, e1 := y["E06000002"]
	z2, e2 := y["E06000002"]["20-24"]
	l := len(z2)
	assert.True(t, e1)
	assert.True(t, e2)
	assert.True(t, l == 11)
}
func TestGetPop02(t *testing.T) {
	ctx := mockdb.Context{}
	codes := []string{"E06000002"}
	x, err := GetProjectedPopulationsByCodes(ctx, codes, 2020, 5, 20, 60, 10, false)
	fmt.Printf("%v", err)
	fmt.Printf("%v", x)
	y := CollateProjectedPopulationsByCode(x)
	fmt.Printf("%v", y)
	_, e1 := y["E06000002"]
	z2, e2 := y["E06000002"]["20-24"]
	l := len(z2)
	assert.True(t, e1)
	assert.True(t, e2)
	assert.True(t, l == 2)
}
