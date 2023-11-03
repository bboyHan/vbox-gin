package model

type Geo struct {
	ID         uint   `gorm:"primary_key" json:"id" form:"id"`
	Name       string `gorm:"column:name" json:"name" form:"name"`
	Level      int    `gorm:"column:level" json:"level" form:"level"`
	Code       string `gorm:"column:code" json:"code" form:"code"`
	Geocode    string `gorm:"column:geocode" json:"geocode" form:"geocode"`
	Latitude   string `gorm:"column:latitude" json:"latitude" form:"latitude"`
	Longitude  string `gorm:"column:longitude" json:"longitude" form:"longitude"`
	Sort       int    `gorm:"column:sort" json:"sort" form:"id"`
	Children   []Geo  `gorm:"-" json:"children" form:"children"`
	ParentCode string `gorm:"column:parentCode" json:"parentCode" form:"parentCode"`
}
