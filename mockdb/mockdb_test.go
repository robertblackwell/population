package mockdb

import (
	"fmt"
	"testing"
)

func TestLoad(t *testing.T) {
	db := LoadProjectedMockDb("population_projections_v2.json")
	fmt.Printf("%v\n", db)
	db2 := LoadCurrentMockDb("populations_v2.json")
	fmt.Printf("%v\n", db2)

}
