package sms

import (
	"context"
	"fmt"

	"github.com/xh-polaris/synapse4b/biz/conf"
	"github.com/xh-polaris/synapse4b/biz/infra/contract/cache"
	"github.com/xh-polaris/synapse4b/biz/infra/contract/sms"
	"github.com/xh-polaris/synapse4b/biz/infra/impl/sms/bluecloud"
	"github.com/xh-polaris/synapse4b/biz/infra/impl/sms/tencent"
)

const (
	BlueCloud = "blue-cloud"
	Tencent   = "tencent"
)

func New(ctx context.Context, cacheCli cache.Cmdable) (provider sms.Provider, err error) {
	c := conf.GetConfig().SMS
	ch := sms.NewSMSCache(ctx, cacheCli)
	switch c.Provider {
	case BlueCloud:
		provider, err = bluecloud.New(ctx, ch, c.Account, c.Token)
	case Tencent:
		provider, err = tencent.New(ctx, ch, c.Account, c.Token)
	default:
		return nil, fmt.Errorf("no such SMS provider: %s", c.Provider)
	}
	if err != nil {
		return nil, err
	}
	return NewSafeSMSProvider(provider, cacheCli)
}
