package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

// OpenAccRegister vb开户
type OpenAccRegister struct {
	Username        string `json:"username" example:"用户名"`
	Password        string `json:"password" example:"密码"`
	NewPassword     string `json:"newPassword" example:"新密码"`
	ConfirmPassword string `json:"confirmPassword" example:"确认密码"`
	CodeProdIDS     []uint `json:"codeProdIDS,omitempty"`
	OrgName         string `json:"orgName" example:"orgName"`
	Recharge        int    `json:"recharge" example:"int 充值金额"`
	Enable          int    `json:"enable" example:"int 是否启用"`
	EnableAuth      int    `json:"enableAuth" example:"int 是否启用防爆验证码"`
	CreateBy        uint   `json:"createBy" example:"创建者"`
}

// SelfRegister Vbox注册子账号
type SelfRegister struct {
	Username        string `json:"username" example:"用户名"`
	Password        string `json:"password" example:"密码"`
	NewPassword     string `json:"newPassword" example:"新密码"`
	ConfirmPassword string `json:"confirmPassword" example:"确认密码"`
	Enable          int    `json:"enable" example:"int 是否启用"`
	EnableAuth      int    `json:"enableAuth" example:"int 是否启用防爆验证码"`
	CreateBy        uint   `json:"createBy" example:"创建者"`
}

// Register User register structure
type Register struct {
	Username     string `json:"username" example:"用户名"`
	Password     string `json:"password" example:"密码"`
	Nickname     string `json:"nickname" example:"昵称"`
	HeaderImg    string `json:"headerImg" example:"头像链接"`
	AuthorityId  uint   `json:"authorityId" swaggertype:"string" example:"int 角色id"`
	Enable       int    `json:"enable" swaggertype:"string" example:"int 是否启用"`
	AuthorityIds []uint `json:"authorityIds" swaggertype:"string" example:"[]uint 角色id"`
	Phone        string `json:"phone" example:"电话号码"`
	Email        string `json:"email" example:"电子邮箱"`
}

// Login User login structure
type Login struct {
	Username    string `json:"username"`    // 用户名
	Password    string `json:"password"`    // 密码
	Captcha     string `json:"captcha"`     // 验证码
	CaptchaId   string `json:"captchaId"`   // 验证码ID
	AuthCaptcha string `json:"authCaptcha"` // 防爆验证码
}

// ChangePasswordReq Modify password structure
type ChangePasswordReq struct {
	ID          uint   `json:"-"`           // 从 JWT 中提取 user id，避免越权
	Password    string `json:"password"`    // 密码
	NewPassword string `json:"newPassword"` // 新密码
}

// ChangeAuthCaptchaReq Modify auth captcha structure
type ChangeAuthCaptchaReq struct {
	ID       uint   `json:"-"`        // 从 JWT 中提取 user id，避免越权
	ToUid    uint   `json:"toUid"`    // 重置子账户
	Type     uint   `json:"type"`     // 1-重置子账户、 0-重置自己账户
	Password string `json:"password"` // 密码
}

// SetUserAuth Modify  user's auth structure
type SetUserAuth struct {
	AuthorityId uint `json:"authorityId"` // 角色ID
}

// SetUserAuthorities Modify  user's auth structure
type SetUserAuthorities struct {
	ID           uint
	AuthorityIds []uint `json:"authorityIds"` // 角色ID
}

type ChangeUserInfo struct {
	ID           uint                  `gorm:"primarykey"`                                                                           // 主键ID
	Nickname     string                `json:"nickname" gorm:"default:系统用户;comment:用户昵称"`                                            // 用户昵称
	Phone        string                `json:"phone"  gorm:"comment:用户手机号"`                                                          // 用户手机号
	AuthorityIds []uint                `json:"authorityIds" gorm:"-"`                                                                // 角色ID
	Email        string                `json:"email"  gorm:"comment:用户邮箱"`                                                           // 用户邮箱
	HeaderImg    string                `json:"headerImg" gorm:"default:https://qmplusimg.henrongyi.top/gva_header.jpg;comment:用户头像"` // 用户头像
	SideMode     string                `json:"sideMode"  gorm:"comment:用户侧边主题"`                                                      // 用户侧边主题
	Enable       int                   `json:"enable" gorm:"comment:冻结用户"`                                                           //冻结用户
	Authorities  []system.SysAuthority `json:"-" gorm:"many2many:sys_user_authority;"`
}
