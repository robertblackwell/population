package mockdb

import (
	"fmt"
	"testing"
)

func TestLoad(t *testing.T) {
	db := LoadMockDb()
	fmt.Printf("%v\n", db)

}
