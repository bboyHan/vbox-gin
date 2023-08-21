package student

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/student"
	studentReq "github.com/flipped-aurora/gin-vue-admin/server/model/student/request"
	"gorm.io/gorm"
)

type StudentService struct {
}

// CreateStudent 创建Student记录
// Author [piexlmax](https://github.com/piexlmax)
func (stdService *StudentService) CreateStudent(std *student.Student) (err error) {
	err = global.GVA_DB.Create(std).Error
	return err
}

// DeleteStudent 删除Student记录
// Author [piexlmax](https://github.com/piexlmax)
func (stdService *StudentService) DeleteStudent(std student.Student) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&student.Student{}).Where("id = ?", std.ID).Update("deleted_by", std.DeletedBy).Error; err != nil {
			return err
		}
		if err = tx.Delete(&std).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteStudentByIds 批量删除Student记录
// Author [piexlmax](https://github.com/piexlmax)
func (stdService *StudentService) DeleteStudentByIds(ids request.IdsReq, deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&student.Student{}).Where("id in ?", ids.Ids).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", ids.Ids).Delete(&student.Student{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateStudent 更新Student记录
// Author [piexlmax](https://github.com/piexlmax)
func (stdService *StudentService) UpdateStudent(std student.Student) (err error) {
	err = global.GVA_DB.Save(&std).Error
	return err
}

// GetStudent 根据id获取Student记录
// Author [piexlmax](https://github.com/piexlmax)
func (stdService *StudentService) GetStudent(id uint) (std student.Student, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&std).Error
	return
}

// GetStudentInfoList 分页获取Student记录
// Author [piexlmax](https://github.com/piexlmax)
func (stdService *StudentService) GetStudentInfoList(info studentReq.StudentSearch) (list []student.Student, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&student.Student{})
	var stds []student.Student
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Find(&stds).Error
	return stds, total, err
}
