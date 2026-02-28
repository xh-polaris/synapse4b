package system

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/xh-polaris/synapse4b/biz/infra/contract/email"
	"math/rand"
	"strconv"
	"time"

	"github.com/xh-polaris/synapse4b/biz/api/model/system"
	"github.com/xh-polaris/synapse4b/biz/application/base/token"
	"github.com/xh-polaris/synapse4b/biz/application/internal"
	"github.com/xh-polaris/synapse4b/biz/conf"
	"github.com/xh-polaris/synapse4b/biz/infra/contract/cache"
	"github.com/xh-polaris/synapse4b/biz/infra/contract/sms"
	ctxcache "github.com/xh-polaris/synapse4b/biz/pkg/ctxcache/ctx_cache"
	"github.com/xh-polaris/synapse4b/biz/pkg/errorx"
	"github.com/xh-polaris/synapse4b/biz/pkg/logs"
	"github.com/xh-polaris/synapse4b/biz/types/cst"
	"github.com/xh-polaris/synapse4b/biz/types/errno"
)

var SystemSVC = &SystemService{}

type SystemService struct {
	sms   sms.Provider
	email email.Provider
	cache cache.Cmdable
}

func (s *SystemService) Send(ctx context.Context, req *system.SendVerifyCodeReq) (*system.SendVerifyCodeResp, error) {
	switch req.AuthType {
	case cst.AuthTypePhoneVerify:
		param := &sms.SMSParam{Code: genCode(), Expire: time.Duration(req.Expire) * time.Second}
		if err := s.sms.Send(ctx, req.App.Name, req.Cause, req.AuthId, param); err != nil {
			return nil, errorx.WrapByCode(err, errno.ErrSendPhoneVerify)
		}
	case cst.AuthTypeEmailVerify:
		param := &email.EmailParam{Code: genCode(), Expire: time.Duration(req.Expire) * time.Second}
		if err := s.email.Send(ctx, req.App.Name, req.Cause, req.AuthId, param); err != nil {
			return nil, errorx.WrapByCode(err, errno.ErrSendPhoneVerify)
		}
	default:
		return nil, errorx.New(errno.ErrInvalidAuthType, errorx.KV("type", req.AuthType))
	}

	return &system.SendVerifyCodeResp{Resp: internal.Success()}, nil
}

// 生成n位随机验证码
func genCode() string {
	return strconv.Itoa(rand.Intn(999999-100000) + 100000)
}

func (s *SystemService) Check(ctx context.Context, req *system.CheckVerifyCodeReq) (*system.CheckVerifyCodeResp, error) {
	var check bool
	var err error

	switch req.AuthType {
	case cst.AuthTypePhoneVerify:
		check, err = s.sms.Check(ctx, req.App.Name, req.Cause, req.AuthId, req.Verify)
		if err != nil {
			return nil, err
		}
	case cst.AuthTypeEmailVerify:
		check, err = s.email.Check(ctx, req.App.Name, req.Cause, req.AuthId, req.Verify)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errorx.New(errno.ErrInvalidAuthType, errorx.KV("type", req.AuthType))
	}

	return &system.CheckVerifyCodeResp{Resp: internal.Success(), Verify: check}, nil
}

// SignTicket 签发token
func (s *SystemService) SignTicket(ctx context.Context, req *system.SignTicketReq) (_ *system.SignTicketResp, err error) {
	// 校验app
	if err = conf.ValidApp(req.GetApp()); err != nil {
		return nil, err
	}
	info, _ := ctxcache.Get[*token.Info](ctx, cst.TokenInfo)

	// ticket逻辑为 app:basic_user_id:token:current_time的MD5值
	rawTicket := req.GetApp().GetName() + ":" + info.BasicUserId + ":" + info.RawToken + ":" + strconv.FormatInt(time.Now().Unix(), 10)
	ticket := fmt.Sprintf("%x", md5.Sum([]byte(rawTicket)))

	if err = s.cache.Set(ctx, req.GetApp().GetName()+":"+ticket, info.RawToken, time.Second*5).Err(); err != nil {
		logs.Errorf("set ticket failed, app: %s, ticket: %s, err: %v", req.GetApp().GetName(), ticket, err)
		return nil, errorx.WrapByCode(err, errno.ErrSignTicket)
	}
	return &system.SignTicketResp{Resp: internal.Success(), Ticket: ticket}, nil
}
func (s *SystemService) ExchangeTicket(ctx context.Context, req *system.ExchangeTicketReq) (*system.ExchangeTicketResp, error) {
	var (
		ok  bool
		err error
		jwt string
	)
	// 校验app与密钥
	if err = conf.ValidApp(req.GetApp()); err != nil {
		return nil, err
	}
	if err, ok = conf.VerifyTicketKey(req.GetApp(), req.GetTicketKey()); err != nil {
		logs.Errorf("verify ticket key failed, app: %s, ticket_key: %s, err: %v", req.GetApp().GetName(), req.GetTicketKey(), err)
		return nil, errorx.WrapByCode(err, errno.ErrExchangeTicket)
	} else if !ok {
		return nil, errorx.New(errno.ErrExchangeTicketByInvalidKey)
	}
	// 获取token
	resp := s.cache.Get(ctx, req.GetApp().GetName()+":"+req.GetTicket())
	if err = resp.Err(); err != nil {
		logs.Errorf("get token failed, app: %s, ticket: %s, err: %v", req.GetApp().GetName(), req.GetTicket(), err)
		return nil, errorx.WrapByCode(err, errno.ErrExchangeTicket)
	}
	if jwt = resp.Val(); jwt == "" {
		logs.Errorf("get token failed caused by empty, app: %s, ticket: %s", req.GetApp().GetName(), req.GetTicket())
		return nil, errorx.New(errno.ErrExchangeTicket)
	}
	return &system.ExchangeTicketResp{Resp: internal.Success(), Token: jwt}, nil
}
