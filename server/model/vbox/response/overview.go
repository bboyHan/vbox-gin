package response

type DataOverView struct {
	X interface{} `json:"x" form:"x" url:"x"`
	Y interface{} `json:"y" form:"y" url:"y"`
}

type DataRateOverView struct {
	X1 int `json:"x1" form:"x1" url:"x1"`
	X2 int `json:"x2" form:"x2" url:"x2"`
	X3 int `json:"x3" form:"x3" url:"x3"`
	X4 int `json:"x4" form:"x4" url:"x4"`
}

type DataSOverView struct {
	X string `json:"x" form:"x" url:"x"`
	Y string `json:"y" form:"y" url:"y"`
}

type DataSExtOverView struct {
	X1 string `json:"x1" form:"x1" url:"x1"`
	X2 string `json:"x2" form:"x2" url:"x2"`
	X3 string `json:"x3" form:"x3" url:"x3"`
	X4 string `json:"x4" form:"x4" url:"x4"`
}