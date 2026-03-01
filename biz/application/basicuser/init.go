package basicuser

import (
	"context"

	"github.com/xh-polaris/synapse4b/biz/domain/basicuser/repo"
	"github.com/xh-polaris/synapse4b/biz/domain/basicuser/service"
	"github.com/xh-polaris/synapse4b/biz/infra/contract/email"
	"github.com/xh-polaris/synapse4b/biz/infra/contract/id"
	"github.com/xh-polaris/synapse4b/biz/infra/contract/orm"
	"github.com/xh-polaris/synapse4b/biz/infra/contract/sms"
)

func InitService(ctx context.Context, sms sms.Provider, email email.Provider, db *orm.DB, idGen id.IDGenerator) *BasicUserService {
	BasicUserSVC.sms = sms
	BasicUserSVC.email = email
	BasicUserSVC.DomainSVC = service.NewBasicUserDomain(ctx, &service.Component{
		BasicUserRepo: repo.NewBasicUserRepo(db),
		UnitRepo:      repo.NewUnitRepo(db),
		IdGen:         idGen,
	})
	return BasicUserSVC
}
