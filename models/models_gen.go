package models

type PopulationAgeRange struct {
	AgeRange    string    `json:"ageRange"`
	Values      []int     `json:"values"`
	GrowthRates []float64 `json:"growthRates"`
}

type PopulationGeography struct {
	Code      string               `json:"code"`
	AgeRanges []PopulationAgeRange `json:"ageRanges"`
}

type PopulationProjections struct {
	Geographies []PopulationGeography `json:"geographies"`
	Years       []int                 `json:"years"`
}
