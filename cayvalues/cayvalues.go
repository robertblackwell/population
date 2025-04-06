package cayvalues

import (
	"fmt"
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
// That is will fail if the uniqueness of the index cannot ge assured
func NewCayValuesFromArr[T CayAble](d []T) (CayValues[T], error) {
	cv := NewCayValues[T]()
	for _, v := range d {
		err := cv.Add(v.GetCode(), v.GetAgeRange(), v.GetYear(), v)
		if err != nil {
			return cv, err
		}
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
	return result, nil
}
