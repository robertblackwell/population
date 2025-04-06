package cayvalues

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type anotherType struct {
	TotalPopulation int
	Year            int
}

type vType struct {
	Code            string
	Type            string
	AgeRange        string
	Year            int
	TotalPopulation int
}

func (v vType) GetAgeRange() string {
	return v.AgeRange
}
func (v vType) GetCode() string {
	return v.Code
}
func (v vType) GetYear() int {
	return v.Year
}

func TestIV(t *testing.T) {
	iv := NewCayValues[int]()
	assert.True(t, iv.Add("ZZZ", "three", 2019, 101) == nil)
	assert.True(t, iv.Add("XXX", "two", 2023, 99) == nil)
	assert.True(t, iv.Add("XXX", "two", 2024, 999) == nil)
	assert.True(t, iv.Add("YYY", "two", 2023, 99) == nil)
	assert.False(t, iv.Add("YYY", "two", 2023, 199) == nil)
	vv, _ := iv.At("YYY", "two", 2023)
	fmt.Printf("%v %v", iv, vv)
}
func TestByTransform(t *testing.T) {
	vArr := []anotherType{
		{TotalPopulation: 821},
		{TotalPopulation: 933},
		{TotalPopulation: 100}, // <<====== index 2
		{TotalPopulation: 160},
		{TotalPopulation: 137},
		{TotalPopulation: 150}, // <======== index 5
		{TotalPopulation: 150},
		{TotalPopulation: 154},
		{TotalPopulation: 159},
		{TotalPopulation: 149},
		{TotalPopulation: 148},
		{TotalPopulation: 150},
		{TotalPopulation: 200}, //<<====== index 12
		{TotalPopulation: 190},
		{TotalPopulation: 189},
		{TotalPopulation: 188},
		{TotalPopulation: 187},
		{TotalPopulation: 186},
	}
	y := 2018
	f := func(t anotherType) (string, string, int, int) {
		tmp := y
		y = y + 1
		return "XXX", "99-202", tmp, t.TotalPopulation
	}
	cayValues, err := NewCayValuesByTransform(vArr, f)
	assert.True(t, err == nil)
	a, ok := cayValues.At("LAD1", "0-4", 2023)
	fmt.Printf("%v %v %v", cayValues, a, ok)
}
func TestMap(t *testing.T) {
	vArr := []vType{
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 821, Year: 2018},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 933, Year: 2019},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1000, Year: 2020}, // <<====== index 2
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1602, Year: 2021},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1378, Year: 2022},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1500, Year: 2023}, // <======== index 5
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1501, Year: 2024},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1540, Year: 2025},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1590, Year: 2026},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1490, Year: 2027},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1480, Year: 2028},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1501, Year: 2029},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 2000, Year: 2030}, //<<====== index 12
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1900, Year: 2031},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1890, Year: 2032},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1880, Year: 2033},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1870, Year: 2034},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1860, Year: 2035},
	}
	f := func(code string, ageRange string, year int, value vType) (float64, error) {
		return float64(value.TotalPopulation*10) + .1, nil
	}
	cayValues, err := NewCayValuesFromArr(vArr)
	assert.True(t, err == nil)
	r, err := Map(cayValues, f)
	assert.True(t, err == nil)
	fmt.Printf("%v %v ", r, err)

}
func TestFromCayAble(t *testing.T) {
	vArr := []vType{
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 821, Year: 2018},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 933, Year: 2019},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1000, Year: 2020}, // <<====== index 2
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1602, Year: 2021},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1378, Year: 2022},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1500, Year: 2023}, // <======== index 5
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1501, Year: 2024},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1540, Year: 2025},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1590, Year: 2026},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1490, Year: 2027},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1480, Year: 2028},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1501, Year: 2029},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 2000, Year: 2030}, //<<====== index 12
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1900, Year: 2031},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1890, Year: 2032},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1880, Year: 2033},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1870, Year: 2034},
		{Code: "LAD1", Type: "type", AgeRange: "0-4", TotalPopulation: 1860, Year: 2035},
	}
	cayValues, err := NewCayValuesFromArr(vArr)
	assert.True(t, err == nil)
	a, ok := cayValues.At("LAD1", "0-4", 2023)
	fmt.Printf("%v %v %v", cayValues, a, ok)

}
