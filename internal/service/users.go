package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"project/initialize"
	dal "project/internal/dal"
	"project/internal/logic"
	model "project/internal/model"
	query "project/internal/query"
	common "project/pkg/common"
	utils "project/pkg/utils"

	"github.com/sirupsen/logrus"
	"gorm.io/gen/field"
)

type UsersService struct {
}

// GetTenant
// @AUTHOR:zxq
// @DATE: 2024-03-04 11:04
// @DESCRIPTIONS: 租户数:租户总数&昨日新增&本月新增&月历史数据
func (*UsersService) GetTenant(ctx context.Context) (model.GetTenantRes, error) {
	var (
		list []*model.GetBoardUserListMonth
		data model.GetTenantRes

		user = query.User
		db   = dal.UserQuery{}
	)
	// 总数据
	total, err := db.Count(ctx)
	if err != nil {
		logrus.Error(ctx, "[GetTenant]Users data failed:", err)
	}
	// 昨日数据
	yesterday, err := db.CountByWhere(ctx, user.CreatedAt.Gte(common.GetYesterdayBegin()))
	if err != nil {
		logrus.Error(ctx, "[GetTenant]Users data failed:", err)
	}
	// 月数据
	month, err := db.CountByWhere(ctx, user.CreatedAt.Gte(common.GetMonthStart()))
	if err != nil {
		logrus.Error(ctx, "[GetTenant]Users data failed:", err)
	}
	// 历史数据
	list = db.GroupByMonthCount(ctx, nil)

	if err != nil {
		logrus.Error(ctx, "[GetTenant]Users data failed:", err)
		return data, err
	}

	data = model.GetTenantRes{
		UserTotal:          total,
		UserAddedYesterday: yesterday,
		UserAddedMonth:     month,
		UserListMonth:      list,
	}
	return data, err
}

// GetTenantUserInfo
// @AUTHOR:zxq
// @DATE: 2024-03-04 11:04
// @DESCRIPTIONS: 租户用户下数据
func (*UsersService) GetTenantUserInfo(ctx context.Context, email string) (model.GetTenantRes, error) {
	var (
		err                     error
		total, yesterday, month int64
		list                    []*model.GetBoardUserListMonth
		data                    model.GetTenantRes

		user = query.User
		db   = dal.UserQuery{}
	)
	// 租户总数据
	total, err = db.CountByWhere(ctx, user.Email.Eq(email))
	if err != nil {
		logrus.Error(ctx, "[GetTenant]Users data failed:", err)
	}
	// 昨日数据
	yesterday, err = db.CountByWhere(ctx, user.CreatedAt.Gte(common.GetYesterdayBegin()), user.Email.Eq(email))
	if err != nil {
		logrus.Error(ctx, "[GetTenant]Users data failed:", err)
	}
	// 月数据
	month, err = db.CountByWhere(ctx, user.CreatedAt.Gte(common.GetMonthStart()), user.Email.Eq(email))
	if err != nil {
		logrus.Error(ctx, "[GetTenant]Users data failed:", err)
	}
	// 历史数据
	list = db.GroupByMonthCount(ctx, &email)

	if err != nil {
		logrus.Error(ctx, "[GetTenant]Users data failed:", err)
		return data, err
	}

	data = model.GetTenantRes{
		UserTotal:          total,
		UserAddedYesterday: yesterday,
		UserAddedMonth:     month,
		UserListMonth:      list,
	}
	return data, err
}

// GetTenantInfo
// @AUTHOR:zxq
// @DATE: 2024-03-04 11:04
// @DESCRIPTIONS: 租户个人信息
func (*UsersService) GetTenantInfo(ctx context.Context, email string) (*model.UsersRes, error) {
	var (
		info *model.UsersRes

		db   = dal.UserQuery{}
		user = query.User
	)
	// 总数据
	UserInfo, err := db.First(ctx, user.Email.Eq(email))
	if err != nil {
		logrus.Error(ctx, "[GetTenantInfo]Users info failed:", err)
		return info, err
	}
	info = dal.UserVo{}.PoToVo(UserInfo)

	return info, err
}

// UpdateTenantInfo
// @AUTHOR:zxq
// @DATE: 2024-03-04 11:04
// @DESCRIPTIONS: 更新租户个人信息
func (*UsersService) UpdateTenantInfo(ctx context.Context, userInfo *utils.UserClaims, param *model.UsersUpdateReq) error {
	var (
		db   = dal.UserQuery{}
		user = query.User
	)
	info, err := db.First(ctx, user.Email.Eq(userInfo.Email))
	if err != nil {
		logrus.Error(ctx, "[UpdateTenantInfo]Get Users info failed:", err)
		return err
	}
	var columns []field.Expr
	columns = append(columns, user.Name)
	info.Name = &param.Name
	if param.AdditionalInfo != nil {
		info.AdditionalInfo = param.AdditionalInfo
		columns = append(columns, user.AdditionalInfo)
	}
	if param.PhoneNumber != nil {
		var phonePrefix string
		if param.PhonePrefix != nil {
			phonePrefix = *param.PhonePrefix
		}
		info.PhoneNumber = fmt.Sprintf("%s %s", phonePrefix, *param.PhoneNumber)
		columns = append(columns, user.PhoneNumber)
	}
	if err = db.UpdateByEmail(ctx, info, columns...); err != nil {
		logrus.Error(ctx, "[UpdateTenantInfo]Update Users info failed:", err)
		return err
	}
	return err
}

// UpdateTenantInfoPassword
// @AUTHOR:zxq
// @DATE: 2024-03-05 13:04
// @DESCRIPTIONS: 更新租户个人密码
func (*UsersService) UpdateTenantInfoPassword(ctx context.Context, userInfo *utils.UserClaims, param *model.UsersUpdatePasswordReq) error {
	// test@test.cn不允许修改密码
	if userInfo.Email == "test@test.cn" {
		return errors.New("该用户不允许修改密码")
	}
	err := utils.ValidatePassword(param.Password)
	if err != nil {
		return err
	}
	var (
		db   = dal.UserQuery{}
		user = query.User
	)
	info, err := db.First(ctx, user.Email.Eq(userInfo.Email))
	if err != nil {
		logrus.Error(ctx, "[UpdateTenantInfoPassword]Get Users info failed:", err)
		return err
	}

	// 是否加密配置
	if logic.UserIsEncrypt(ctx) {
		password, err := initialize.DecryptPassword(param.Password)
		if err != nil {
			return fmt.Errorf("wrong decrypt password")
		}
		passwords := strings.TrimSuffix(string(password), param.Salt)
		param.Password = passwords
	}

	// 验证旧密码
	if !utils.BcryptCheck(param.OldPassword, info.Password) {
		return errors.New("OldPassword Failed,Please again~")
	}

	t := time.Now().UTC()
	info.UpdatedAt = &t
	info.PasswordLastUpdated = &t

	info.Password = utils.BcryptHash(param.Password)
	if err = db.UpdateByEmail(ctx, info, user.Password, user.UpdatedAt, user.PasswordLastUpdated); err != nil {
		logrus.Error(ctx, "[UpdateTenantInfoPassword]Update Users info failed:", err)
		return err
	}
	return err
}
