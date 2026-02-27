package email

import (
	"context"
	"fmt"
	"strconv"

	"github.com/xh-polaris/synapse4b/biz/conf"
	"github.com/xh-polaris/synapse4b/biz/infra/contract/cache"
	"github.com/xh-polaris/synapse4b/biz/infra/contract/email"
	"github.com/xh-polaris/synapse4b/biz/infra/contract/risk"
	"github.com/xh-polaris/synapse4b/biz/pkg/errorx"
	"github.com/xh-polaris/synapse4b/biz/pkg/logs"
	"github.com/xh-polaris/synapse4b/biz/types/errno"
)

type SafeEmailProvider struct {
	email.Provider
	Cache cache.Cmdable
}

func NewSafeEmailProvider(provider email.Provider, cacheCli cache.Cmdable) (*SafeEmailProvider, error) {
	return &SafeEmailProvider{
		Provider: provider,
		Cache:    cacheCli,
	}, nil
}

// Send 发送验证码
func (s *SafeEmailProvider) Send(ctx context.Context, app, cause, email string, param *email.EmailParam) error {
	// 判断是否到上限
	key := fmt.Sprintf("risk:sendVerifyCode:%s:%s:%s", app, cause, email)
	limit, _, err := risk.CheckUpperLimit(ctx, key, conf.GetConfig().Email.MaxInPeriod)
	if err != nil {
		return err
	}
	if limit { // 达到上限, 不允许发送
		return errorx.New(errno.ErrSendUpperLimit, errorx.KV("period", strconv.Itoa(conf.GetConfig().Email.Period)))
	}
	// 发送验证码
	err = s.Provider.Send(ctx, app, cause, email, param)
	if err != nil {
		return err
	}
	// 记录操作
	if err = risk.AddOnce(ctx, key, conf.GetConfig().Email.Period); err != nil {
		logs.Errorf("record send verify err:%s", err)
	}
	return nil
}

// Check 校验验证码
func (s *SafeEmailProvider) Check(ctx context.Context, app, cause, email, code string) (bool, error) {
	return s.Provider.Check(ctx, app, cause, email, code)
}
