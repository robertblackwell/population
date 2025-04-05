package resolver

import (
	"fmt"
	"forecast_model/mockdb"
	"testing"
)

//	func TestGetProjectedPopulation(t *testing.T) {
//	    ctx := mockdb.Context{}
//	    m, ok := GetProjectedPopulationByCodes(ctx, []string{"E06000002", "E06000003"}, 2020, 5, 20, 60, 3, true)
//	    assert.True(t, ok == nil)
//	    tmp := CollateProjectedPopulationsByCode(m)
//	    assert.Equal(t, len(tmp["E06000002"]), (2035 - 2018 + 1))
//	    iv := IsValidMapOfPopVec(tmp)
//	    assert.True(t, iv)
//	    assert.True(t, IsValidMapOfPopVec(m))
//	    fmt.Printf("%v\n", tmp)
//	}
//
//	func TestGetAllProjectedPopulation(t *testing.T) {
//	    ctx := mockdb.Context{}
//	    m, ok := GetAllProjectedPopulationsByCodes(ctx, []string{"E06000002", "E06000003"}, 40, 42)
//	    // x, err := mockdb.GetProjectedPopulationByCodes(mockdb.Context{}, []string{"XX", "YY"}, 2020, 2, 40, 42, 3, true)
//	    if ok != nil {
//	        return
//	    }
//	    assert.Equal(t, len(m["E06000002"]), (2035 - 2018 + 1))
//	    iv := IsValidMapOfPopVec(m)
//	    assert.True(t, iv)
//	    assert.True(t, IsValidMapOfPopVec(m))
//	    // fmt.Printf("%v\n", m)
//	}
func TestXadd(t *testing.T) {
	m := map[string]map[string]map[int]LadPopulationProjection{}
	j1 := mockdb.JsonRecord{Code: "XXX", Type: "type", Age: 45, Date: "2018-01-01", Value: 1234}
	j2 := mockdb.JsonRecord{Code: "XXX", Type: "type", Age: 46, Date: "2018-01-01", Value: 1234}
	m = Xadd(m, j1.Code, "45-49", 2018, j1)
	m = Xadd(m, j2.Code, "45-49", 2018, j2)
	j11 := mockdb.JsonRecord{Code: "XXX", Type: "type", Age: 55, Date: "2018-01-01", Value: 1234}
	j12 := mockdb.JsonRecord{Code: "XXX", Type: "type", Age: 56, Date: "2018-01-01", Value: 1234}
	m = Xadd(m, j11.Code, "55-59", 2018, j11)
	m = Xadd(m, j12.Code, "55-59", 2018, j12)
	fmt.Printf("%v", m)
}
