package response

type ChartData struct {
	LegendData []string     `json:"legendData"`
	XAxisData  []string     `json:"xAxisData"`
	SeriesData []SeriesItem `json:"seriesData"`
}

type SeriesItem struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Stack  string `json:"stack"`
	Smooth bool   `json:"smooth"`
	Data   []int  `json:"data"`
}
