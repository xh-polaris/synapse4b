package basicuser

import (
	"context"
	"fmt"
	"strconv"
	"time"

	model "github.com/xh-polaris/synapse4b/biz/api/model/basicuser"
	"github.com/xh-polaris/synapse4b/biz/application/base/token"
	"github.com/xh-polaris/synapse4b/biz/application/basicuser/internal"
	application "github.com/xh-polaris/synapse4b/biz/application/internal"
	"github.com/xh-polaris/synapse4b/biz/conf"
	"github.com/xh-polaris/synapse4b/biz/domain/basicuser/entity"
	basicuser "github.com/xh-polaris/synapse4b/biz/domain/basicuser/service"
	"github.com/xh-polaris/synapse4b/biz/infra/contract/email"
	"github.com/xh-polaris/synapse4b/biz/infra/contract/risk"
	"github.com/xh-polaris/synapse4b/biz/infra/contract/sms"
	ctxcache "github.com/xh-polaris/synapse4b/biz/pkg/ctxcache/ctx_cache"
	"github.com/xh-polaris/synapse4b/biz/pkg/errorx"
	"github.com/xh-polaris/synapse4b/biz/pkg/lang/util"
	"github.com/xh-polaris/synapse4b/biz/pkg/logs"
	"github.com/xh-polaris/synapse4b/biz/types/cst"
	"github.com/xh-polaris/synapse4b/biz/types/errno"
)

var BasicUserSVC = &BasicUserService{}

type BasicUserService struct {
	sms       sms.Provider
	email     email.Provider
	DomainSVC basicuser.BasicUser
}

// RegisterNewBasicUser 注册一个新用户
// 暂时只允许手机号注册新用户
func (s *BasicUserService) RegisterNewBasicUser(ctx context.Context, req *model.BasicUserRegisterReq) (resp *model.BasicUserRegisterResp, err error) {
	if err = conf.ValidApp(req.GetApp()); err != nil {
		return nil, err
	}
	var authType string
	switch req.AuthType {
	case cst.AuthTypePhoneVerify:
		if err = s.validPhoneVerify(ctx, req.App.Name, req.AuthId, req.Verify); err != nil {
			return nil, err
		}
		ok, err := s.DomainSVC.PhoneExist(ctx, req.AuthId)
		if err != nil {
			return nil, err
		}
		if ok {
			return nil, errorx.New(errno.PhoneHasExisted, errorx.KV("phone", req.AuthId))
		}
		authType = cst.TokenAuthType
	//case cst.AuthTypeCodePassword:
	//	if req.ExtraAuthId == nil {
	//		return nil, errorx.New(errno.MissingParameter, errorx.KV("parameter", "学号"))
	//	}
	//	ok, err := s.DomainSVC.CodeExist(ctx, req.AuthId, *req.ExtraAuthId)
	//	if err != nil {
	//		return nil, err
	//	}
	//	if ok {
	//		return nil, errorx.New(errno.CodeHasExisted, errorx.KV("code", *req.ExtraAuthId))
	//	}
	//	authType = cst.TokenAuthType
	default:
		return nil, errorx.New(errno.UnSupportAuthType, errorx.KV("type", req.AuthType))
	}

	if req.GetPassword() == "" {
		return nil, errorx.New(errno.MustPassword)
	}

	var extraAuthId string
	if req.ExtraAuthId != nil {
		extraAuthId = *req.ExtraAuthId
	}
	var u *entity.BasicUser
	if u, err = s.DomainSVC.Register(ctx, req.AuthType, req.AuthId, extraAuthId, *req.Password); err != nil {
		return nil, err
	}

	info := &token.Info{
		BasicUserId: u.ID,
		UnitId:      util.UnPtr(u.UnitID),
		Code:        util.UnPtr(u.Code),
		Phone:       util.UnPtr(u.Phone),
		Email:       util.UnPtr(u.Email),
		LoginTime:   time.Now().Unix(),
		AuthType:    authType,
		Extra:       nil,
	}
	jwt, err := token.SignJWT(conf.GetConfig().Token, info)
	if err != nil {
		return nil, err
	}

	resp = &model.BasicUserRegisterResp{
		Resp:      application.Success(),
		Token:     jwt,
		BasicUser: internal.BasicUserPO2VO(u),
	}
	return
}

// 判断是否达到风控限制
func (s *BasicUserService) validRisk(ctx context.Context, key string) (err error) {
	// 判断是否到上限
	limit, _, err := risk.CheckUpperLimit(ctx, key, conf.GetConfig().Token.MaxInPeriod)
	if err != nil {
		return err
	}
	if limit { // 达到上限, 不允许校验
		return errorx.New(errno.TooOftenLoginError, errorx.KV("period", strconv.Itoa(conf.GetConfig().SMS.Period/60)))
	}
	return nil
}

func (s *BasicUserService) validPhoneVerify(ctx context.Context, app, phone, code string) error {
	key := fmt.Sprintf("risk:login:passport:%s", phone)
	if err := s.validRisk(ctx, key); err != nil {
		return err
	}
	ok, err := s.sms.Check(ctx, app, "passport", phone, code)
	if err != nil {
		return err
	}
	if !ok {
		if err = risk.AddOnce(ctx, key, conf.GetConfig().Token.Period); err != nil {
			logs.Errorf("record send verify err:%s", err)
		}
		return errorx.New(errno.ErrVerifyCode)
	}
	return err
}

func (s *BasicUserService) validEmailVerify(ctx context.Context, app, email, code string) error {
	key := fmt.Sprintf("risk:login:passport:%s", email)
	if err := s.validRisk(ctx, key); err != nil {
		return err
	}
	ok, err := s.email.Check(ctx, app, "passport", email, code)
	if err != nil {
		return err
	}
	if !ok {
		if err = risk.AddOnce(ctx, key, conf.GetConfig().Token.Period); err != nil {
			logs.Errorf("record send verify err:%s", err)
		}
		return errorx.New(errno.ErrVerifyCode)
	}
	return err
}

func (s *BasicUserService) Login(ctx context.Context, req *model.BasicUserLoginReq) (resp *model.BasicUserLoginResp, err error) {
	if err = conf.ValidApp(req.GetApp()); err != nil {
		return nil, err
	}

	var ok = true
	var u *entity.BasicUser
	var authType string
	switch req.AuthType {
	case cst.AuthTypePhoneVerify: // 手机号验证码登录
		err = s.validPhoneVerify(ctx, req.App.Name, req.AuthId, req.Verify)
		if err != nil {
			return nil, err
		}
		ok, err = s.DomainSVC.PhoneExist(ctx, req.AuthId)
		if err != nil {
			return nil, err
		}
		if ok {
			u, err = s.DomainSVC.LoginByPhone(ctx, false, req.AuthId, req.Verify)
		} else {
			u, err = s.DomainSVC.Register(ctx, req.AuthType, req.AuthId, "", "")
		}
		authType = cst.AuthTypePhone
	case cst.AuthTypePhonePassword: // 手机号密码登录
		u, err = s.DomainSVC.LoginByPhone(ctx, true, req.AuthId, req.Verify)
		authType = cst.AuthTypePhone
	case cst.AuthTypeCodePassword: // 学号密码登录
		if req.ExtraAuthId == nil {
			return nil, errorx.New(errno.MissingParameter, errorx.KV("parameter", "学号"))
		}
		u, err = s.DomainSVC.LoginByCode(ctx, req.AuthId, *req.ExtraAuthId, req.Verify)
		authType = cst.AuthTypeCode
	case cst.AuthTypeEmailPassword: // 邮箱密码登录
		u, err = s.DomainSVC.LoginByEmail(ctx, true, req.AuthId, req.Verify)
		authType = cst.AuthTypeEmail
	case cst.AuthTypeEmailVerify: // 邮箱验证码登录
		err = s.validEmailVerify(ctx, req.App.Name, req.AuthId, req.Verify)
		if err != nil {
			return nil, err
		}
		ok, err = s.DomainSVC.EmailExist(ctx, req.AuthId)
		if err != nil {
			return nil, err
		}
		if ok {
			u, err = s.DomainSVC.LoginByEmail(ctx, false, req.AuthId, req.Verify)
		} else {
			u, err = s.DomainSVC.Register(ctx, req.AuthType, req.AuthId, "", "")
		}
		authType = cst.AuthTypeEmail
	default:
		return nil, errorx.New(errno.UnSupportAuthType, errorx.KV("type", req.AuthType))
	}
	if err != nil {
		return nil, err
	}

	info := &token.Info{
		BasicUserId: u.ID,
		UnitId:      util.UnPtr(u.UnitID),
		Code:        util.UnPtr(u.Code),
		Phone:       util.UnPtr(u.Phone),
		Email:       util.UnPtr(u.Email),
		LoginTime:   time.Now().Unix(),
		AuthType:    authType,
		Extra:       nil,
	}
	jwt, err := token.SignJWT(conf.GetConfig().Token, info)
	if err != nil {
		return nil, err
	}

	resp = &model.BasicUserLoginResp{
		Resp:      application.Success(),
		Token:     jwt,
		BasicUser: internal.BasicUserPO2VO(u),
		New:       !ok,
	}
	return
}

func (s *BasicUserService) ResetPassword(ctx context.Context, req *model.BasicUserResetPasswordReq) (resp *model.BasicUserResetPasswordResp, err error) {
	if err = conf.ValidApp(req.GetApp()); err != nil {
		return nil, err
	}
	info, _ := ctxcache.Get[*token.Info](ctx, cst.TokenInfo)

	if info == nil {
		return nil, errorx.New(errno.InvalidToken)
	}
	if err = s.DomainSVC.ResetPassword(ctx, info.BasicUserId, req.NewPassword); err != nil {
		return nil, errorx.New(errno.ErrResetPassword)
	}

	return &model.BasicUserResetPasswordResp{
		Resp: application.Success(),
	}, nil
}
