package global

import (
	"time"

	"gorm.io/gorm"
)

type GVA_MODEL struct {
	ID        uint           `gorm:"primarykey"`                                                        // 主键ID
	CreatedAt *time.Time     `json:"CreatedAt" form:"CreatedAt" gorm:"column:created_at;comment:创建时间;"` // 创建时间
	UpdatedAt *time.Time     `json:"UpdatedAt" form:"UpdatedAt" gorm:"column:updated_at;comment:更新时间;"` // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                                                    // 删除时间
}
