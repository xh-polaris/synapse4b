package service

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/xh-polaris/synapse4b/biz/conf"
	"github.com/xh-polaris/synapse4b/biz/domain/basicuser/dal/model"
	"github.com/xh-polaris/synapse4b/biz/domain/basicuser/entity"
	"github.com/xh-polaris/synapse4b/biz/domain/basicuser/repo"
	"github.com/xh-polaris/synapse4b/biz/infra/contract/id"
	"github.com/xh-polaris/synapse4b/biz/infra/contract/risk"
	"github.com/xh-polaris/synapse4b/biz/pkg/errorx"
	"github.com/xh-polaris/synapse4b/biz/pkg/lang/crypt"
	"github.com/xh-polaris/synapse4b/biz/pkg/lang/util"
	"github.com/xh-polaris/synapse4b/biz/pkg/logs"
	"github.com/xh-polaris/synapse4b/biz/types/cst"
	"github.com/xh-polaris/synapse4b/biz/types/errno"
)

type Component struct {
	BasicUserRepo repo.BasicUserRepo
	AuthRepo      repo.AuthRepo
	SchoolRepo    repo.SchoolRepo
	IdGen         id.IDGenerator
}

func NewBasicUserDomain(ctx context.Context, c *Component) BasicUser {
	return &userImpl{Component: c}
}

type userImpl struct {
	*Component
}

func (i *userImpl) CreateBasicUser(ctx context.Context, unitId, code, phone, email, password string, encryptType int64) (user *entity.BasicUser, err error) {
	// 查询是否存在完全匹配的账号, 若是则绑定到该账号上, 忽视初始密码
	u, err := i.BasicUserRepo.FindCompletely(ctx, unitId, code, phone, email)
	if err != nil {
		return nil, err
	} else if u != nil { // 存在完全匹配
		if user, err = basicUserModel2Entity(u); err != nil {
			return nil, err
		}
		return user, nil
	}
	// 查询是否存在部分匹配的账号, 若是则返回部分匹配, 绑定失败
	u, err = i.BasicUserRepo.FindPartly(ctx, unitId, code, phone, email)
	if err != nil {
		return nil, err
	} else if u != nil { // 存在部分匹配
		return nil, errorx.New(errno.ErrPartlyCreate)
	}

	// 否则, 创建新用户, 使用初始密码
	var pass string
	if password == "" {
		return nil, errorx.New(errno.MustPassword)
	}
	if encryptType == 0 { // 使用bcrypt加密
		pass, err = crypt.Hash(password)
	} else { // 传入的是md5密文
		pass = password
	}
	nu := &model.BasicUser{ID: i.IdGen.GenID(ctx), Password: util.Of(pass), Encrypt: uint8(encryptType)}

	// 学号
	if unitId != "" && code != "" {
		uid, err := id.FromHex(unitId)
		if err != nil {
			return nil, err
		}
		nu.UnitID = util.Of(uid)
		nu.Code = util.Of(code)
	} else if unitId != "" || code != "" { // 不完整
		return nil, errorx.New(errno.MissingParameter, errorx.KV("parameter", "unitId or code"))
	}
	// 手机号
	if phone != "" {
		nu.Phone = util.Of(phone)
	}
	// 邮箱
	if email != "" {
		nu.Email = util.Of(email)
	}

	nu, err = i.BasicUserRepo.Create(ctx, nu)
	if err != nil {
		return nil, errorx.New(errno.ErrRegister)
	}
	return basicUserModel2Entity(nu)
}

func (i *userImpl) UserIDExist(ctx context.Context, userId string) (user *entity.BasicUser, is bool, err error) {
	u, err := i.BasicUserRepo.FindByID(ctx, userId)
	if err != nil || u == nil {
		return nil, false, err
	}
	if user, err = basicUserModel2Entity(u); err != nil {
		return nil, false, err
	}
	return user, true, err
}

func (i *userImpl) LoginByPhone(ctx context.Context, requirePassword bool, phone, verify string) (*entity.BasicUser, error) {
	u, err := i.BasicUserRepo.FindByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}
	if u == nil { // 未注册过
		return nil, errorx.New(errno.PhoneNotExisted)
	}
	if requirePassword {
		if err = loginLimiter(ctx, verify, u.Password, phone); err != nil {
			return nil, err
		}
	}
	return basicUserModel2Entity(u)
}

func (i *userImpl) LoginByEmail(ctx context.Context, requirePassword bool, email, verify string) (*entity.BasicUser, error) {
	u, err := i.BasicUserRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if u == nil { // 未注册过
		return nil, errorx.New(errno.EmailNotExisted)
	}
	if requirePassword {
		if err = loginLimiter(ctx, verify, u.Password, email); err != nil {
			return nil, err
		}
	}
	return basicUserModel2Entity(u)
}

func (i *userImpl) LoginByCode(ctx context.Context, schoolId, code, verify string) (*entity.BasicUser, error) {
	u, err := i.BasicUserRepo.FindByCode(ctx, schoolId, code)
	if err != nil {
		return nil, err
	}
	if u == nil { // 未注册过
		return nil, errorx.New(errno.CodeNotExisted, errorx.KV("code", code))
	}
	if err = loginLimiter(ctx, verify, u.Password, schoolId, code); err != nil {
		return nil, err
	}
	return basicUserModel2Entity(u)
}

func (i *userImpl) PhoneExist(ctx context.Context, phone string) (is bool, err error) {
	mu, err := i.BasicUserRepo.FindByPhone(ctx, phone)
	if err != nil {
		return false, err
	}
	if mu != nil {
		return true, nil
	}
	return false, nil
}

func (i *userImpl) CodeExist(ctx context.Context, schoolId, code string) (is bool, err error) {
	mus, err := i.BasicUserRepo.FindManyBySchoolID(ctx, schoolId)
	if err != nil {
		return false, err
	}
	for _, mu := range mus {
		if mu.Code != nil && *mu.Code == code {
			return true, nil
		}
	}
	return false, nil
}

func (i *userImpl) EmailExist(ctx context.Context, email string) (is bool, err error) {
	mu, err := i.BasicUserRepo.FindByEmail(ctx, email)
	if err != nil {
		return false, err
	}
	if mu != nil {
		return true, nil
	}
	return false, nil
}

func (i *userImpl) Register(ctx context.Context, authType, authId, extraAuthId, password string) (u *entity.BasicUser, err error) {
	var hashed string
	if password != "" {
		hashed, err = crypt.Hash(password)
		if err != nil {
			return nil, err
		}
	}

	nu := &model.BasicUser{
		ID:       i.IdGen.GenID(ctx),
		Password: &hashed,
	}

	switch authType {
	case cst.AuthTypePhoneVerify:
		nu.Phone = &authId
	//case cst.AuthTypeCodePassword:
	//	aid, err := id.FromHex(authId)
	//	if err != nil {
	//		return nil, err
	//	}
	//	nu.UnitID = &aid
	//	nu.Code = &extraAuthId
	case cst.AuthTypeEmailVerify:
		nu.Email = &authId
	default:
		return nil, errorx.New(errno.UnSupportAuthType, errorx.KV("type", authType))
	}
	nu, err = i.BasicUserRepo.Create(ctx, nu)
	if err != nil {
		return nil, errorx.New(errno.ErrRegister)
	}
	return basicUserModel2Entity(nu)
}

func (i *userImpl) ResetPassword(ctx context.Context, basicUserId string, password string) error {
	if password == "" {
		return errorx.New(errno.MustPassword)
	}
	hashed, err := crypt.Hash(password)
	if err != nil {
		return err
	}
	return i.BasicUserRepo.ResetPassword(ctx, basicUserId, hashed)
}

func loginLimiter(ctx context.Context, password string, hashed *string, parts ...string) error {
	key := "risk:login:passport:" + strings.Join(parts, ",")
	limit, _, err := risk.CheckUpperLimit(ctx, key, conf.GetConfig().Token.MaxInPeriod)
	if err != nil {
		return err
	}
	if limit { // 达到上限, 不允许校验
		return errorx.New(errno.TooOftenLoginError, errorx.KV("period", strconv.Itoa(conf.GetConfig().SMS.Period/60)))
	}
	if hashed == nil || *hashed == "" {
		return errorx.New(errno.NoPassword)
	}
	if !crypt.Check(password, *hashed) {
		if err = risk.AddOnce(ctx, key, conf.GetConfig().Token.Period); err != nil {
			logs.Errorf("record send verify err:%s", err)
		}
		return errorx.New(errno.ErrPassword)
	}
	return nil
}

func basicUserModel2Entity(u *model.BasicUser) (*entity.BasicUser, error) {
	eu := &entity.BasicUser{
		ID:        u.ID.Hex(),
		Code:      u.Code,
		Phone:     u.Phone,
		Email:     u.Email,
		Password:  u.Password,
		Name:      u.Name,
		Gender:    u.Gender,
		CreatedAt: time.UnixMilli(u.CreatedAt),
		UpdatedAt: time.UnixMilli(u.UpdatedAt),
	}
	if len(u.Extra) != 0 {
		extra := map[string]any{}
		err := sonic.Unmarshal([]byte(u.Extra.String()), &extra)
		if err != nil {
			return nil, err
		}
		eu.Extra = &extra
	}
	return eu, nil
}
