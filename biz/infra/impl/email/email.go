package email

import (
	"context"
	"fmt"

	"github.com/xh-polaris/synapse4b/biz/conf"
	"github.com/xh-polaris/synapse4b/biz/infra/contract/cache"
	"github.com/xh-polaris/synapse4b/biz/infra/contract/email"
	"github.com/xh-polaris/synapse4b/biz/infra/impl/email/tencent"
)

const (
	Tencent = "tencent"
)

func New(ctx context.Context, cacheCli cache.Cmdable) (provider email.Provider, err error) {
	c := conf.GetConfig().SMS
	ch := email.NewEmailCache(ctx, cacheCli)
	switch c.Provider {
	case Tencent:
		provider, err = tencent.New(ctx, ch, c.Account, c.Token)
	default:
		return nil, fmt.Errorf("no such Email provider: %s", c.Provider)
	}
	if err != nil {
		return nil, err
	}
	return NewSafeEmailProvider(provider, cacheCli)
}
