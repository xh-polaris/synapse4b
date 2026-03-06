package system

import (
	"context"
	"github.com/xh-polaris/synapse4b/biz/infra/contract/email"

	"github.com/xh-polaris/synapse4b/biz/infra/contract/cache"
	"github.com/xh-polaris/synapse4b/biz/infra/contract/sms"
)

func InitService(ctx context.Context, sms sms.Provider, cache cache.Cmdable, email email.Provider) *SystemService {
	SystemSVC.sms = sms
	SystemSVC.cache = cache
	SystemSVC.email = email
	return SystemSVC
}
