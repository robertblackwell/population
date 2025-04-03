package model

type JsonRecord struct {
	Value  int    `json:"value"`
	Code   string `json:"code"`
	Type   string `json:"type"`
	Date   string `json:"date"`
	Age    int    `json:"age"`
	Gender string `json:"gender"`
}

type JsonList struct {
	Records []JsonRecord `json:"records"`
}
