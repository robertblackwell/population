package repo

import (
	"fmt"
	"maps"
	"slices"

	"popmodel/mockdb"

	"popmodel/ageranges"
)

// copied from graphql project
type Population struct {
	AgeRange        string `db:"age_range"`
	Sex             string `db:"sex"`
	TotalPopulation int    `db:"total_population"`
	Year            int    `db:"year"`
}

// copied from graphql project
type LadPopulationProjection struct {
	Code            string `db:"code"`
	Type            string `db:"type"`
	AgeRange        string `db:"age_range"`
	TotalPopulation int    `db:"total_population"`
	Year            int    `db:"year"`
}

func (lap LadPopulationProjection) GetCode() string {
	return lap.Code
}
func (lap LadPopulationProjection) GetAgeRange() string {
	return lap.AgeRange
}
func (lap LadPopulationProjection) GetYear() int {
	return lap.Year
}

func GetParentCodes(ctx mockdb.Context, ladCode string) ([]string, error) {
	return []string{ladCode}, nil
}
func GetCurrentPopulationsByCodes(ctx mockdb.Context, parents []string, rangeSize int, minAge int, maxAge int, currentPopulationYear int) ([]LadPopulationProjection, error) {
	return []LadPopulationProjection{}, nil
}
func GetPopulationByCodes(ctx mockdb.Context, codes []string, startYear, rangeSize, minAge, maxAge, futureOffset int, includeIntermediates bool) ([]LadPopulationProjection, error) {
	return GetProjectedPopulationsByCodes(ctx, codes, startYear, rangeSize, minAge, maxAge, futureOffset, includeIntermediates)
}

// This is a mock version of a function written for test purposes.
//
// In the live system the analogous functions quesries the populations_projection_v2 table to get:Gets the projected population total for the specified age range, from the projected_populations_v2 table for all years between 2018 and 2035.
// amalgamate (ie sum over ages within the range)
// amalgamate over codes - this is a mock function only one code is allowed and is ignored
// startYear, rangesize, futureOffset, includeIntermediaries are ignored.
// startyear and future offset define the years for which we want a forecast - one forecast for startYear and another for startYear+futureOffset
// minAge, maxAge and rangeSize define the age ranges to be used in grouping results. If minAge is 20, maxAge is 60 and rangeSize is 5 we want the following groups:
// 20-24, 25-29, 30-34, 35-39, 40-44, 45-49, 50-54, 55-59, 60-64
func GetProjectedPopulationsByCodes(ctx mockdb.Context, codes []string, startYear, rangeSize, minAge, maxAge, futureOffset int, includeIntermediates bool) ([]LadPopulationProjection, error) {
	target := mockdb.LoadMockDb()
	map_result := make(map[string]map[string]map[int]LadPopulationProjection, 0)
	flat_result := make([]LadPopulationProjection, 0)
	ageRanges, err := ageranges.CreateAgeRanges(minAge, maxAge, rangeSize)
	if err != nil {
		return flat_result, err
	}
	fmt.Printf("%v", target)
	for kode, v1 := range target {
		if slices.Contains(codes, kode) {
			ageKeys := slices.Sorted(maps.Keys(v1))
			for _, age := range ageKeys {
				v2 := v1[age]
				r, err := ageranges.AgeRangesContainAge(ageRanges, age)
				if err == nil {
					ageRangeString := ageranges.AgeRangeToString(r)
					dateKeys := slices.Sorted(maps.Keys(v2))
					for _, dateK := range dateKeys {
						jrec := target[kode][age][dateK]
						map_result = Xadd(map_result, kode, ageRangeString, int(ageranges.YearFromDate(dateK)), jrec)
					}
				}
			}
		}
	}
	for _, v1 := range map_result {
		ageRangeKeys := slices.Sorted(maps.Keys(v1))
		for _, arKey := range ageRangeKeys {
			v2 := v1[arKey]
			yearKeys := slices.Sorted(maps.Keys(v2))
			for _, yearKey := range yearKeys {
				if (includeIntermediates && yearKey >= startYear && yearKey <= (startYear+futureOffset)) ||
					(yearKey == startYear || yearKey == (startYear+futureOffset)) {
					flat_result = append(flat_result, v2[yearKey])
				}
			}
		}
	}
	fmt.Printf("I am here again")
	return flat_result, nil
}
func GetBaseYearProjectedPopulations(code string, ageRange string, baseYear int, rangeSize int, minAge int, maxAge int) (LadPopulationProjection, error) {
	ctx := mockdb.Context{}
	pp, er := GetProjectedPopulationsByCodes(ctx, []string{code}, baseYear, rangeSize, minAge, maxAge, 0, false)
	if er != nil {
		return LadPopulationProjection{}, er
	}
	pv := CollateProjectedPopulationsByCode(pp)
	x1, ok := pv[code]
	if !ok {
		return LadPopulationProjection{}, fmt.Errorf("index by code %s failed", code)
	}
	x2, ok := x1[ageRange]
	if !ok {
		return LadPopulationProjection{}, fmt.Errorf("index by code %s and ageRange code %s failed", code, ageRange)
	}
	return x2[0], nil
}

func Xadd(m map[string]map[string]map[int]LadPopulationProjection, code string, ageRangeStr string, year int, jr mockdb.JsonRecord) map[string]map[string]map[int]LadPopulationProjection {
	p := LadPopulationProjection{Code: code, Type: jr.Type, AgeRange: ageRangeStr, Year: year, TotalPopulation: jr.Value}
	_, ok := m[code]
	if !ok {
		m[code] = map[string]map[int]LadPopulationProjection{ageRangeStr: map[int]LadPopulationProjection{year: p}}
		return m
	}
	_, ok = m[code][ageRangeStr]
	if !ok {
		m[code][ageRangeStr] = map[int]LadPopulationProjection{year: p}
		return m
	}
	_, ok = m[code][ageRangeStr][year]
	if !ok {
		m[code][ageRangeStr][year] = LadPopulationProjection{Code: code, Type: jr.Type, AgeRange: ageRangeStr, Year: year, TotalPopulation: jr.Value}
	} else {
		p2 := m[code][ageRangeStr][year]
		p2.TotalPopulation = p2.TotalPopulation + p.TotalPopulation
		m[code][ageRangeStr][year] = p2
	}
	return m
}
func CollateProjectedPopulationsByCode(pops []LadPopulationProjection) map[string]map[string][]LadPopulationProjection {

	result := map[string]map[string][]LadPopulationProjection{}
	for _, v := range pops {
		result = y_append(result, v)
	}
	return result
}
func CollateBaseYearProjectedPopulationsByCode(pops []LadPopulationProjection) (map[string]map[string]LadPopulationProjection, error) {
	result := map[string]map[string]LadPopulationProjection{}
	for _, v := range pops {
		result, err := z_append(result, v)
		if err != nil {
			return result, err
		}
	}
	return result, nil
}
func y_append(r map[string]map[string][]LadPopulationProjection, p LadPopulationProjection) map[string]map[string][]LadPopulationProjection {
	_, ok := r[p.Code]
	if !ok {
		r[p.Code] = map[string][]LadPopulationProjection{p.AgeRange: {p}}
		return r
	}
	_, ok = r[p.Code][p.AgeRange]
	if !ok {
		r[p.Code][p.AgeRange] = []LadPopulationProjection{p}
		return r
	}
	r[p.Code][p.AgeRange] = append(r[p.Code][p.AgeRange], p)
	return r
}
func z_append(r map[string]map[string]LadPopulationProjection, p LadPopulationProjection) (map[string]map[string]LadPopulationProjection, error) {
	_, ok := r[p.Code]
	if !ok {
		r[p.Code] = map[string]LadPopulationProjection{p.AgeRange: p}
		return r, nil
	}
	_, ok = r[p.Code][p.AgeRange]
	if !ok {
		r[p.Code][p.AgeRange] = p
		return r, nil
	}
	return r, fmt.Errorf("already has entry for Code: %s  ageRange: %s", p.Code, p.AgeRange)
}
