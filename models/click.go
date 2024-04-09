package models

type Click struct {
	DateTime string `json:"date_time"`
	X        int16  `json:"x"`
	Y        int16  `json:"y"`
}

func NewClick(dateTime string, x int16, y int16) *Click {
	return &Click{
		DateTime: dateTime,
		X:        x,
		Y:        y,
	}
}
