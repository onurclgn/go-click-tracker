package models

type Click struct {
	DateTime string `json:"date_time"`
	X        int16  `json:"x"`
	Y        int16  `json:"y"`
}
