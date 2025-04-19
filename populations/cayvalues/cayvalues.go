package cayvalues

import (
	"cmp"
	"fmt"
	"maps"
	"slices"
	"sort"
	"strconv"
	"strings"
)

//
// The CayAble data structure is intended to facilitate the manipulation of datasets where values can be uniquely indexed
// by (Code string, ageRange string, Year int).
//
// The primary use of this type is manipulation of population data where:
// -	Code identified a geographic region,
// -	the ageRange is a string version of something like people between the age of 40 and 44 represented as "40-44"
// -    the Year is a number like 2023 representing the year 2023
// - the values can be records or primitives like int and float64

// type CayValuesX[T any] struct {
// 	values map[string]map[string]map[int]T
// }

type CayValues[T any] map[string]map[string]map[int]T

type CayAble interface {
	GetCode() string
	GetAgeRange() string
	GetYear() int
}

// Make a new empty CayValue with values of type T
func NewCayValues[T any]() CayValues[T] {
	r := make(map[string]map[string]map[int]T, 0)
	return r
}

// Create a new CayValue[T] instance and fill it from an array []T
// T must be CayAble.
// Will fail if the input data has multiple values with the same code/ageRange/year
//
// That is will fail if the uniqueness of the index cannot be assured
// of the CheckKeys function fails
func NewCayValuesFromArr[T CayAble](d []T) (CayValues[T], error) {
	cv := NewCayValues[T]()
	for _, v := range d {
		err := cv.Add(v.GetCode(), v.GetAgeRange(), v.GetYear(), v)
		if err != nil {
			return nil, err
		}
	}
	if err := cv.CheckKeys(); err != nil {
		return nil, err
	}
	if err := cv.CheckKeys(); err != nil {
		panic(err)
	}
	return cv, nil
}

// Create a CayValue[V] instance from an arr []T where a transform functions both provides some of the
// keys - code, ageRange, year - and also may modify the value.
//
// This is a good place to use closures see the tests for CayValues for an example
func NewCayValuesByTransform[T any, V any](ts []T, transform func(t T) (string, string, int, V)) (CayValues[V], error) {

	tmp := NewCayValues[V]()
	for _, tv := range ts {
		code, ageRange, year, value := transform(tv)
		err := tmp.Add(code, ageRange, year, value)
		if err != nil {
			return nil, err
		}
	}
	if err := tmp.CheckKeys(); err != nil {
		panic(err)
	}
	return tmp, nil
}

// adds a new value with its index to the CayValue instance.
//
// will return error if the index of the value to be added already exists
func (iv *CayValues[T]) Add(code string, ageRange string, year int, value T) error {
	// _, ok := iv.values[code]
	if _, ok := (*iv)[code]; !ok {
		(*iv)[code] = map[string]map[int]T{ageRange: {year: value}}
		return nil
	}
	if _, ok := (*iv)[code][ageRange]; !ok {
		(*iv)[code][ageRange] = map[int]T{year: value}
		return nil
	}
	if _, ok := (*iv)[code][ageRange][year]; !ok {
		(*iv)[code][ageRange][year] = value
		return nil
	}
	return fmt.Errorf("keys %s %s %d already in use", code, ageRange, year)
}

// get a value by index. Will return bool true if found and bool false if not found
func (cayv CayValues[T]) At(code string, ageRange string, year int) (T, bool) {
	var x T
	if _, ok := cayv[code]; !ok {
		return x, false
	}
	if _, ok := cayv[code][ageRange]; !ok {
		return x, false
	}
	if _, ok := cayv[code][ageRange][year]; !ok {
		return x, false
	}
	r := (cayv[code][ageRange][year])
	return r, true
}

// get a value by index. Will return bool true if found and bool false if not found
func (cayv *CayValues[T]) Set(code string, ageRange string, year int, value T) {
	if _, ok := (*cayv)[code]; !ok {
		(*cayv)[code] = map[string]map[int]T{ageRange: {year: value}}
	} else if _, ok := (*cayv)[code][ageRange]; !ok {
		(*cayv)[code][ageRange] = map[int]T{year: value}
	} else if _, ok := (*cayv)[code][ageRange][year]; !ok {
		(*cayv)[code][ageRange][year] = value
	}
}

// Construct a new CayValue object by applying a function to each element of the input CayValues.
//
// This call will fail and return an error if the function f returns an error
func Map[T any, V any](iv CayValues[T], f func(code string, ageRange string, year int, t T) (V, error)) (CayValues[V], error) {
	if err := iv.CheckKeys(); err != nil {
		panic(err)
	}
	result := NewCayValues[V]()
	for k1, v1 := range iv {
		for k2, v2 := range v1 {
			for k3, v3 := range v2 {
				r, err := f(k1, k2, k3, v3)
				if err != nil {
					return nil, err
				}
				result.Add(k1, k2, k3, r)
			}
		}
	}
	if err := result.CheckKeys(); err != nil {
		return nil, err
	}
	return result, nil
}

// iterates over a CayValues in depth first order. The func f is called on each leaf node
func Iterate[T any](iv CayValues[T], f func(code string, ageRange string, year int, value T) error) error {
	if err := iv.CheckKeys(); err != nil {
		panic(err)
	}
	// a function for comparing age ranges for sorting. An age range is a string of the form "20-24"
	// the will not sort correctly as strings as "20-24" preceeds "5-9"
	comparef := func(a string, b string) int {
		afirst, err1 := strconv.Atoi(strings.Split(a, "-")[0])
		if err1 != nil {
			panic(fmt.Errorf("failed comparing ageRanges %s is invalid as an ageRange", a))
		}
		bfirst, err2 := strconv.Atoi(strings.Split(b, "-")[0])
		if err2 != nil {
			panic(fmt.Errorf("failed comparing ageRanges %s is invalid as an ageRange", b))
		}
		return cmp.Compare(afirst, bfirst)
	}
	keys1 := (slices.Collect((maps.Keys(iv))))
	sort.Strings(keys1)
	for _, k1 := range keys1 {
		v1 := iv[k1]

		keys2 := slices.Collect(maps.Keys(v1))
		slices.SortFunc(keys2, comparef)
		for _, k2 := range keys2 {
			v2 := v1[k2]

			keys3 := slices.Collect(maps.Keys(v2))
			sort.Ints(keys3)
			for _, k3 := range keys3 {
				v3 := v2[k3]
				err := f(k1, k2, k3, v3)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Check keys - each first level key must have the same set of 2nd level keys
//
//	each 2nd level key must have the same set of 3rd level keys
func (cayv CayValues[T]) CheckKeys() error {
	keys1 := slices.Collect(maps.Keys(cayv))
	keys2 := slices.Collect(maps.Keys(cayv[keys1[0]]))
	keys3 := slices.Collect(maps.Keys(cayv[keys1[0]][keys2[0]]))

	for k1, v1 := range cayv {
		if !slices.Contains(keys1, k1) {
			return fmt.Errorf("top level key fail key: %s is extraneous ", k1)
		}
		if len(keys2) != len(slices.Collect(maps.Keys(v1))) {
			return fmt.Errorf("second level key count failed key: %s ", k1)
		}
		for k2, v2 := range v1 {
			if !slices.Contains(keys2, k2) {
				return fmt.Errorf("2nd level key fail key: %s:%s is extraneous ", k1, k2)
			}
			if len(keys3) != len(slices.Collect(maps.Keys(v2))) {
				return fmt.Errorf("3rd level key count failed key: %s:%s ", k1, k2)
			}
			for k3, _ := range v2 {
				if !slices.Contains(keys3, k3) {
					return fmt.Errorf("3rd level key fail key: %s:%s:%d is extraneous ", k1, k2, k3)
				}

			}
		}
	}
	return nil
}
