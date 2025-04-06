package populations

import (
	"popmodel/cayvalues"
	"popmodel/repo"
)

// // copied from graphql project
// type Population struct {
// 	AgeRange        string `db:"age_range"`
// 	Sex             string `db:"sex"`
// 	TotalPopulation int    `db:"total_population"`
// 	Year            int    `db:"year"`
// }

// // copied from graphql project
// type LadPopulationProjection struct {
// 	Code            string `db:"code"`
// 	Type            string `db:"type"`
// 	AgeRange        string `db:"age_range"`
// 	TotalPopulation int    `db:"total_population"`
// 	Year            int    `db:"year"`
// }

// func (lap LadPopulationProjection) GetCode() string {
// 	return lap.Code
// }
// func (lap LadPopulationProjection) GetAgeRange() string {
// 	return lap.AgeRange
// }
// func (lap LadPopulationProjection) GetYear() int {
// 	return lap.Year
// }

type Population = repo.Population
type LadPopulationProjecttion = repo.LadPopulationProjection

type GrowthRate struct {
	Year             int
	RateFromBaseYear float64
}

type EstimatedPopulation struct {
	Year       int
	Population int
}

type PopVec = []repo.LadPopulationProjection

type GG = cayvalues.CayValues[repo.LadPopulationProjection]
