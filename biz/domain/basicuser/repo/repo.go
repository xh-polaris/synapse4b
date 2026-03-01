package repo

import (
	"context"

	"github.com/xh-polaris/synapse4b/biz/domain/basicuser/dal"
	"github.com/xh-polaris/synapse4b/biz/domain/basicuser/dal/model"
	"github.com/xh-polaris/synapse4b/biz/infra/contract/orm"
)

func NewBasicUserRepo(db *orm.DB) BasicUserRepo {
	return dal.NewBasicUserDAO(db)
}

func NewUnitRepo(db *orm.DB) UnitRepo {
	return dal.NewUnitDAO(db)
}

type BasicUserRepo interface {
	FindByID(ctx context.Context, id string) (*model.BasicUser, error)
	FindByPhone(ctx context.Context, phone string) (*model.BasicUser, error)
	FindByEmail(ctx context.Context, email string) (*model.BasicUser, error)
	FindManyByUnitID(ctx context.Context, unitId string) ([]*model.BasicUser, error)
	FindByCode(ctx context.Context, unitId, code string) (*model.BasicUser, error)
	FindCompletely(ctx context.Context, unitId, code, phone, email string) (*model.BasicUser, error)
	FindPartly(ctx context.Context, unitId, code, phone, email string) (*model.BasicUser, error)
	Create(ctx context.Context, nu *model.BasicUser) (*model.BasicUser, error)
	ResetPassword(ctx context.Context, basicUserId, password string) error
}

type AuthRepo interface {
}

type UnitRepo interface {
	FindByID(ctx context.Context, id string) (*model.Unit, error)
	FindByName(ctx context.Context, name string) (*model.Unit, error)
	Create(ctx context.Context, nu *model.Unit) (*model.Unit, error)
}
