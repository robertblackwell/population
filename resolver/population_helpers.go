package resolver

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"sort"
	"strconv"
	"strings"
)

// Years are the years for which population forcasts have been made and from which we will build
// a population forecasting model. Currently they are 2018 .. 2035
type Year int

const FirstYear = 2018
const LastYear = 2035

func IsValidYear(y int) bool {
	return y >= 2018 && y <= 2035
}
func YearFromInt(y int) Year {
	if IsValidYear(y) {
		return Year(y)
	}
	panic(fmt.Sprintf("YearFromInt y: %d", y))
}

type Age int

func IsValidAge(a int) bool {
	return a >= 0 && a <= 100
}

func int2Age(a int) (Age, error) {
	if IsValidAge(a) {
		return Age(a), nil
	}
	return Age(0), errors.New("invalid age")
}

type AgeRange = struct {
	Start Age
	End   Age
}

func createAgeRange(start int, end int) (AgeRange, error) {
	if start <= end {
		s, err1 := int2Age(start)
		e, err2 := int2Age(end)
		if err1 == nil && err2 == nil {
			return AgeRange{Start: s, End: e}, nil
		}
	}
	return AgeRange{}, errors.New("invalid ages or range")
}
func AgeRangeToString(ar AgeRange) string {
	return fmt.Sprintf("%d-%d", int(ar.Start), int(ar.End))
}

// Takes a string of the form 2024-01-01 and extracts the year value as an int
// Checks that the year number is in the range 2018 .. 2023
//
// panics on error
func YearFromDate(dateStr string) Year {
	bits := strings.Split(dateStr, "-")
	if len(bits) != 3 {
		fmt.Printf("YearFromDate %s\n", dateStr)
		panic("YearFromDate failed")
	}
	y := bits[0]
	if ynum, ok := strconv.Atoi(y); ok == nil {
		if ynum >= 2018 && ynum <= 2035 {
			return Year(ynum)
		} else {
			fmt.Printf("YearFromDate %d\n", ynum)
			panic("YearFromDate out of range failed")
		}
	} else {
		panic("YearFromDate Atoi failed")
	}
}

// Generic - Extracts the string keys from a value of type map[string]T
// and sorts the keys
func SortedMapKeys[T any](m map[string]T) []string {
	keys := slices.Collect(maps.Keys(m))
	sort.Strings(keys)
	return keys
}
