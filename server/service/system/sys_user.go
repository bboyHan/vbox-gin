package system

import (
	"errors"
	"fmt"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/organization/model"
	utils2 "github.com/flipped-aurora/gin-vue-admin/server/plugin/organization/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/captcha"
	"github.com/songzhibin97/gkit/tools/rand_string"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type UserService struct{}

// OpenAccRegister
// @function: OpenAccRegister
// @description: 开户，注册子用户（子角色存在才能建立）
// @param: u model.SysUser
// @return: userInter system.SysUser, err error
func (userService *UserService) OpenAccRegister(u systemReq.OpenAccRegister) (userInter system.SysUser, err error) {
	var user system.SysUser
	if !errors.Is(global.GVA_DB.Where("username = ?", u.Username).First(&user).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
		return userInter, errors.New("用户名已注册")
	}

	if u.NewPassword != u.ConfirmPassword {
		return userInter, errors.New("两次输入密码不一致")
	}

	if utils.IsAlphaNumericUnderscore(u.Username) {
		return userInter, errors.New("用户名只能包含字母、数字、下划线")
	}

	// 只允许建 当前角色的子角色账户(目前只支持 子角色为单一角色，多角色不支持)
	roleID := uint(1001)

	u.Password = utils.BcryptHash(u.NewPassword)

	authorities := []system.SysAuthority{
		{
			AuthorityId: roleID,
		},
	}

	create := system.SysUser{
		UUID:        uuid.Must(uuid.NewV4()),
		Username:    u.Username,
		Password:    u.Password,
		Nickname:    u.Username,
		Enable:      u.Enable,
		EnableAuth:  u.EnableAuth,
		AuthorityId: roleID,
		Authorities: authorities,
	}

	create.AuthCaptcha, err = captcha.AuthQrCode(u.Username)
	if err != nil {
		return userInter, errors.New("当前用户设置防爆码异常，请联系管理员")
	}

	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		TxErr := tx.Create(&create).Error
		if TxErr != nil {
			return TxErr
		}

		if u.OrgName == "" {
			u.OrgName = create.Username + "团队"
		}

		var orgId uint
		if u.OrgName != "" {
			var org model.Organization
			if !errors.Is(tx.Where("name = ?", u.OrgName).First(&org).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
				return errors.New("团队名已注册")
			}

			orgCreat := model.Organization{
				Name:     u.OrgName,
				ParentID: 1,
			}
			if err := tx.Create(&orgCreat).Error; err != nil {
				return err
			}

			orgId = orgCreat.ID
		}

		roleID, err = utils.GetRoleID(u.CreateBy)
		if roleID == 888 || roleID == 1000 {
		} else {
			return errors.New("该账号无直充权限")
		}

		eventId := utils.Prefix(global.WalletEventRechargePrefix, rand_string.RandomInt(8))
		// 直充
		rowSelf := &vbox.UserWallet{
			Uid:       create.ID,
			Username:  create.Username,
			Recharge:  u.Recharge,
			Type:      global.WalletRechargeType,
			EventId:   eventId,
			Remark:    fmt.Sprintf(global.WalletEventRecharge, u.Recharge),
			CreatedBy: create.ID,
		}
		if err := tx.Model(&vbox.UserWallet{}).Create(rowSelf).Error; err != nil {
			return err
		}

		var Users []model.OrgUser
		var CUsers []model.OrgUser
		//if !errors.Is(tx.Where("name = ?", u.OrgName).First(&org).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册

		err := tx.Find(&Users, "organization_id = ?", orgId).Error
		if err != nil {
			return err
		}
		var UserIdMap = make(map[uint]bool)
		for i := range Users {
			UserIdMap[Users[i].SysUserID] = true
		}

		if !UserIdMap[create.ID] {
			CUsers = append(CUsers, model.OrgUser{SysUserID: create.ID, OrganizationID: orgId})
		}
		err = tx.Create(&CUsers).Error

		// 产品入队
		var Products []vbox.OrgProduct
		var CProducts []vbox.OrgProduct
		err = tx.Find(&Products, "organization_id = ?", orgId).Error
		if err != nil {
			return err
		}
		var ProductMap = make(map[uint]bool)
		for i := range Products {
			ProductMap[Products[i].ChannelProductID] = true
		}

		for i := range u.CodeProdIDS {
			if !ProductMap[u.CodeProdIDS[i]] {
				CProducts = append(CProducts, vbox.OrgProduct{ChannelProductID: u.CodeProdIDS[i], OrganizationID: orgId})
			}
		}
		err = tx.Create(&CProducts).Error

		// 返回 nil 提交事务
		return nil
	})
	return create, err
}

// SelfRegister
// @function: SelfRegister
// @description: 注册子用户（子角色存在才能建立）
// @param: u model.SysUser
// @return: userInter system.SysUser, err error
func (userService *UserService) SelfRegister(u systemReq.SelfRegister) (userInter system.SysUser, err error) {
	var user system.SysUser
	if !errors.Is(global.GVA_DB.Where("username = ?", u.Username).First(&user).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
		return userInter, errors.New("用户名已注册")
	}

	if u.NewPassword != u.ConfirmPassword {
		return userInter, errors.New("两次输入密码不一致")
	}

	if !utils.IsAlphaNumericUnderscore(u.Username) {
		return userInter, errors.New("用户名只能包含字母、数字、下划线")
	}

	// 只允许建 当前角色的子角色账户(目前只支持 子角色为单一角色，多角色不支持)
	roleID, err := utils.GetSubRoleID(u.CreateBy)
	if err != nil || roleID == 0 {
		return userInter, errors.New("当前用户无权创建成员用户，请联系管理员")
	}

	u.Password = utils.BcryptHash(u.NewPassword)

	authorities := []system.SysAuthority{
		{
			AuthorityId: roleID,
		},
	}

	create := system.SysUser{
		UUID:        uuid.Must(uuid.NewV4()),
		Username:    u.Username,
		Password:    u.Password,
		Nickname:    u.Username,
		Enable:      u.Enable,
		EnableAuth:  u.EnableAuth,
		AuthorityId: roleID,
		Authorities: authorities,
	}

	if u.EnableAuth == 1 {
		create.AuthCaptcha, err = captcha.AuthQrCode(u.Username)
		if err != nil {
			return userInter, errors.New("当前用户设置防爆码异常，请联系管理员")
		}
	}

	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		TxErr := tx.Create(&create).Error
		if TxErr != nil {
			return TxErr
		}

		orgIds := utils2.GetSelfOrg(u.CreateBy)

		for _, orgId := range orgIds {
			var Users []model.OrgUser
			var CUsers []model.OrgUser
			err := global.GVA_DB.Find(&Users, "organization_id = ?", orgId).Error
			if err != nil {
				return err
			}
			var UserIdMap = make(map[uint]bool)
			for i := range Users {
				UserIdMap[Users[i].SysUserID] = true
			}

			if !UserIdMap[create.ID] {
				CUsers = append(CUsers, model.OrgUser{SysUserID: create.ID, OrganizationID: orgId})
			}
			err = global.GVA_DB.Create(&CUsers).Error
		}

		// 返回 nil 提交事务
		return nil
	})
	return create, err
}

// Register
// @function: Register
// @description: 用户注册
// @param: u model.SysUser
// @return: userInter system.SysUser, err error
func (userService *UserService) Register(u system.SysUser) (userInter system.SysUser, err error) {
	var user system.SysUser
	if !errors.Is(global.GVA_DB.Where("username = ?", u.Username).First(&user).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
		return userInter, errors.New("用户名已注册")
	}
	// 否则 附加uuid 密码hash加密 注册
	u.Password = utils.BcryptHash(u.Password)
	u.UUID = uuid.Must(uuid.NewV4())
	err = global.GVA_DB.Create(&u).Error
	return u, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@author: [SliverHorn](https://github.com/SliverHorn)
//@function: Login
//@description: 用户登录
//@param: u *model.SysUser
//@return: err error, userInter *model.SysUser

func (userService *UserService) Login(u *system.SysUser) (userInter *system.SysUser, err error) {
	if nil == global.GVA_DB {
		return nil, fmt.Errorf("db not init")
	}

	var user system.SysUser
	err = global.GVA_DB.Debug().Table("sys_users AS u").Joins("join sys_user_authority as ua on u.id = ua.sys_user_id").
		Select("u.*,a.authority_id").Preload("Authorities").Preload("Authority").
		Joins("join sys_authorities as a on a.authority_id = ua.sys_authority_authority_id").First(&user, "username = ?", u.Username).Error
	if err == nil {
		if ok := utils.BcryptCheck(u.Password, user.Password); !ok {
			return nil, errors.New("密码错误")
		}
		if user.EnableAuth == 1 { //开启走自己的防爆码
			var secret string
			secret, _ = captcha.GetSecret(user.AuthCaptcha)
			if ok := captcha.ValidateCode(secret, u.AuthCaptcha); !ok {
				return nil, errors.New("双因子认证码错误")
			}
		} else { //未开启走系统自定义防爆码
			var capAuth string
			err = global.GVA_DB.Model(&vbox.Proxy{}).Select("url").
				Where("chan = ?", "auth_captcha").Where("type = ? and status = ?", 1, 1).
				Find(&capAuth).Error
			if u.AuthCaptcha != capAuth || err != nil {
				return nil, errors.New("双因子认证码错误")
			}
		}
		MenuServiceApp.UserAuthorityDefaultRouter(&user)
	}
	return &user, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: ChangePassword
//@description: 修改用户密码
//@param: u *model.SysUser, newPassword string
//@return: userInter *model.SysUser,err error

func (userService *UserService) ChangePassword(u *system.SysUser, newPassword string) (userInter *system.SysUser, err error) {
	var user system.SysUser
	if err = global.GVA_DB.Where("id = ?", u.ID).First(&user).Error; err != nil {
		return nil, err
	}
	if ok := utils.BcryptCheck(u.Password, user.Password); !ok {
		return nil, errors.New("原密码错误")
	}
	user.Password = utils.BcryptHash(newPassword)
	err = global.GVA_DB.Save(&user).Error
	return &user, err
}

// ResetAuthCaptcha
// @function: ResetAuthCaptcha
// @description: 重置防爆码
// @param: u *model.SysUser, newPassword string
// @return: userInter *model.SysUser,err error
func (userService *UserService) ResetAuthCaptcha(req systemReq.ChangeAuthCaptchaReq) (userInter *system.SysUser, err error) {
	var user system.SysUser
	// 1. 核对当前登录账户的密码过验
	if err = global.GVA_DB.Where("id = ?", req.ID).First(&user).Error; err != nil {
		return nil, err
	}
	if ok := utils.BcryptCheck(req.Password, user.Password); !ok {
		return nil, errors.New("密码错误，不允许重置双因子认证")
	}

	// 2. 核对是否为同组织下的用户，不允许跨组织修改别人账户
	uIds := utils2.GetDeepUserIDs(req.ID)
	exist := utils.Contains(uIds, req.ToUid)
	if !exist {
		return nil, errors.New("不允许修改非同一团队的账户信息")
	}

	if req.Type == 1 { // 重置子账户防爆
		var otherUser system.SysUser
		if err = global.GVA_DB.Where("id = ?", req.ToUid).First(&otherUser).Error; err != nil {
			return nil, err
		}
		otherUser.AuthCaptcha, err = captcha.AuthQrCode(otherUser.Username)
		otherUser.EnableAuth = 1
		err = global.GVA_DB.Save(&otherUser).Error
		return &otherUser, err
	} else { // 重置自己账户的防爆
		user.AuthCaptcha, err = captcha.AuthQrCode(user.Username)
		user.EnableAuth = 1
		err = global.GVA_DB.Save(&user).Error
	}

	return &user, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetUserInfoList
//@description: 分页获取数据
//@param: info request.PageInfo
//@return: err error, list interface{}, total int64

func (userService *UserService) GetUserInfoList(info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&system.SysUser{})
	var userList []system.SysUser
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Preload("Authorities").Preload("Authority").Find(&userList).Error
	return userList, total, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: SetUserAuthority
//@description: 设置一个用户的权限
//@param: uuid uuid.UUID, authorityId string
//@return: err error

func (userService *UserService) SetUserAuthority(id uint, authorityId uint) (err error) {
	assignErr := global.GVA_DB.Where("sys_user_id = ? AND sys_authority_authority_id = ?", id, authorityId).First(&system.SysUserAuthority{}).Error
	if errors.Is(assignErr, gorm.ErrRecordNotFound) {
		return errors.New("该用户无此角色")
	}
	err = global.GVA_DB.Where("id = ?", id).First(&system.SysUser{}).Update("authority_id", authorityId).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: SetUserAuthorities
//@description: 设置一个用户的权限
//@param: id uint, authorityIds []string
//@return: err error

func (userService *UserService) SetUserAuthorities(id uint, authorityIds []uint) (err error) {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		TxErr := tx.Delete(&[]system.SysUserAuthority{}, "sys_user_id = ?", id).Error
		if TxErr != nil {
			return TxErr
		}
		var useAuthority []system.SysUserAuthority
		for _, v := range authorityIds {
			useAuthority = append(useAuthority, system.SysUserAuthority{
				SysUserId: id, SysAuthorityAuthorityId: v,
			})
		}
		TxErr = tx.Create(&useAuthority).Error
		if TxErr != nil {
			return TxErr
		}
		TxErr = tx.Where("id = ?", id).First(&system.SysUser{}).Update("authority_id", authorityIds[0]).Error
		if TxErr != nil {
			return TxErr
		}
		// 返回 nil 提交事务
		return nil
	})
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteUser
//@description: 删除用户
//@param: id float64
//@return: err error

func (userService *UserService) DeleteUser(id int) (err error) {
	var user system.SysUser
	err = global.GVA_DB.Where("id = ?", id).Delete(&user).Error
	if err != nil {
		return err
	}
	err = global.GVA_DB.Delete(&[]system.SysUserAuthority{}, "sys_user_id = ?", id).Error
	err = global.GVA_DB.Delete(&[]model.OrgUser{}, "sys_user_id = ?", id).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: SetUserInfo
//@description: 设置用户信息
//@param: reqUser model.SysUser
//@return: err error, user model.SysUser

func (userService *UserService) SetUserInfo(req system.SysUser) error {
	return global.GVA_DB.Model(&system.SysUser{}).
		Select("updated_at", "nick_name", "header_img", "phone", "email", "sideMode", "enable").
		Where("id=?", req.ID).
		Updates(map[string]interface{}{
			"updated_at": time.Now(),
			"nick_name":  req.Nickname,
			"header_img": req.HeaderImg,
			"phone":      req.Phone,
			"email":      req.Email,
			"side_mode":  req.SideMode,
			"enable":     req.Enable,
		}).Error
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: SetUserInfo
//@description: 设置用户信息
//@param: reqUser model.SysUser
//@return: err error, user model.SysUser

func (userService *UserService) SetSelfInfo(req system.SysUser) error {
	return global.GVA_DB.Model(&system.SysUser{}).
		Where("id=?", req.ID).
		Updates(req).Error
}

//@author: [piexlmax](https://github.com/piexlmax)
//@author: [SliverHorn](https://github.com/SliverHorn)
//@function: GetUserInfo
//@description: 获取用户信息
//@param: uuid uuid.UUID
//@return: err error, user system.SysUser

func (userService *UserService) GetUserInfo(uuid uuid.UUID) (user system.SysUser, err error) {
	var reqUser system.SysUser
	//err = global.GVA_DB.Debug().Preload("Authorities").Preload("Authority").First(&reqUser, "uuid = ?", uuid).Error
	err = global.GVA_DB.Table("sys_users AS u").Joins("join sys_user_authority as ua on u.id = ua.sys_user_id").
		Select("u.*,a.authority_id").Preload("Authorities").Preload("Authority").
		Joins("join sys_authorities as a on a.authority_id = ua.sys_authority_authority_id").First(&reqUser, "uuid = ?", uuid).Error
	if err != nil {
		return reqUser, err
	}
	MenuServiceApp.UserAuthorityDefaultRouter(&reqUser)
	return reqUser, err
}

//@author: [SliverHorn](https://github.com/SliverHorn)
//@function: FindUserById
//@description: 通过id获取用户信息
//@param: id int
//@return: err error, user *model.SysUser

func (userService *UserService) FindUserById(id int) (user *system.SysUser, err error) {
	var u system.SysUser
	err = global.GVA_DB.Where("`id` = ?", id).First(&u).Error
	return &u, err
}

//@author: [SliverHorn](https://github.com/SliverHorn)
//@function: FindUserByUuid
//@description: 通过uuid获取用户信息
//@param: uuid string
//@return: err error, user *model.SysUser

func (userService *UserService) FindUserByUuid(uuid string) (user *system.SysUser, err error) {
	var u system.SysUser
	if err = global.GVA_DB.Where("`uuid` = ?", uuid).First(&u).Error; err != nil {
		return &u, errors.New("用户不存在")
	}
	return &u, nil
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: resetPassword
//@description: 修改用户密码
//@param: ID uint
//@return: err error

func (userService *UserService) ResetPassword(ID uint) (err error) {
	err = global.GVA_DB.Model(&system.SysUser{}).Where("id = ?", ID).Update("password", utils.BcryptHash("123456")).Error
	return err
}
