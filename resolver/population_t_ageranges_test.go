package resolver

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAgeRanges01(t *testing.T) {
	x, e := CreateAgeRange(2, 5)
	assert.True(t, (e == nil))
	fmt.Printf("%v", x)
}
func TestAgeRanges02(t *testing.T) {
	x, e := CreateAgeRange(2, 2)
	assert.True(t, (e == nil))
	fmt.Printf("%v", x)
}
func TestAgeRanges03(t *testing.T) {
	x, e := CreateAgeRange(2, 1)
	assert.True(t, (e != nil))
	fmt.Printf("%v", x)
}
func TestAgeRanges04(t *testing.T) {
	x, e := CreateAgeRanges(20, 60, 5)
	assert.True(t, e == nil)
	fmt.Printf("%v", x)
}
