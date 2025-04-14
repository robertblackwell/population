package cayvalues

import (
	"fmt"
	"maps"
	"slices"
)

//
// The CayAble data structure is intended to facilitate the manipulation of datasets where values can be uniquely indexed
// by (Code string, ageRange string, Year int).
//
// The primary use of this type is manipulation of population data where:
// -	Code identified a geographic region,
// -	the ageRange is a string version of something like people between the age of 40 and 44 represented as "40-44"
// -    the Year is a number like 2023 representing the year 2023
//

type CayValues[T any] struct {
	values map[string]map[string]map[int]T
}
type CayValuesx[T any] map[string]map[string]T

type CayAble interface {
	GetCode() string
	GetAgeRange() string
	GetYear() int
}

// Make a new empty CayValue with values of type T
func NewCayValues[T any]() CayValues[T] {
	r := CayValues[T]{values: make(map[string]map[string]map[int]T, 0)}
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
			return cv, err
		}
	}
	if err := cv.CheckKeys(); err != nil {
		return NewCayValues[T](), err
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
			return tmp, err
		}
	}
	return tmp, nil
}

// adds a new value with its index to the CayValue instance.
//
// will return error if the index of the value to be added already exists
func (iv *CayValues[T]) Add(code string, ageRange string, year int, value T) error {
	// _, ok := iv.values[code]
	if _, ok := iv.values[code]; !ok {
		iv.values[code] = map[string]map[int]T{ageRange: {year: value}}
		return nil
	}
	if _, ok := iv.values[code][ageRange]; !ok {
		iv.values[code][ageRange] = map[int]T{year: value}
		return nil
	}
	if _, ok := iv.values[code][ageRange][year]; !ok {
		iv.values[code][ageRange][year] = value
		return nil
	}
	return fmt.Errorf("keys %s %s %d already in use", code, ageRange, year)
}

// get a value by index. Will return bool true if found and bool false if not found
func (cayv *CayValues[T]) At(code string, ageRange string, year int) (T, bool) {
	var x T
	if _, ok := cayv.values[code]; !ok {
		return x, false
	}
	if _, ok := cayv.values[code][ageRange]; !ok {
		return x, false
	}
	if _, ok := cayv.values[code][ageRange][year]; !ok {
		return x, false
	}
	r := (cayv.values[code][ageRange][year])
	return r, true
}

// Construct a new CayValue object by applying a function to each element of the input CayValues.
//
// This call will fail and return an error if the function f returns an error
func Map[T any, V any](iv CayValues[T], f func(code string, ageRange string, year int, t T) (V, error)) (CayValues[V], error) {
	result := NewCayValues[V]()
	for k1, v1 := range iv.values {
		for k2, v2 := range v1 {
			for k3, v3 := range v2 {
				r, err := f(k1, k2, k3, v3)
				if err != nil {
					return NewCayValues[V](), err
				}
				result.Add(k1, k2, k3, r)
			}
		}
	}
	if err := result.CheckKeys(); err != nil {
		return result, err
	}
	return result, nil
}

// iterates over a CayValues in depth first order. The func f is called on each leaf node
func Iterate[T any](iv CayValues[T], f func(code string, ageRange string, year int, t T) error) error {
	for k1, v1 := range iv.values {
		for k2, v2 := range v1 {
			for k3, v3 := range v2 {
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
	keys1 := slices.Collect(maps.Keys(cayv.values))
	keys2 := slices.Collect(maps.Keys(cayv.values[keys1[0]]))
	keys3 := slices.Collect(maps.Keys(cayv.values[keys1[0]][keys2[0]]))

	for k1, v1 := range cayv.values {
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
