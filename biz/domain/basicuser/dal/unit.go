package dal

import (
	"context"
	"errors"

	"github.com/xh-polaris/synapse4b/biz/domain/basicuser/dal/model"
	"github.com/xh-polaris/synapse4b/biz/domain/basicuser/dal/query"
	"github.com/xh-polaris/synapse4b/biz/infra/contract/id"
	"github.com/xh-polaris/synapse4b/biz/infra/contract/orm"
	"gorm.io/gorm"
)

func NewUnitDAO(db *orm.DB) *UnitDAO {
	return &UnitDAO{
		query: query.Use(db),
	}
}

type UnitDAO struct {
	query *query.Query
}

func (d *UnitDAO) FindByID(ctx context.Context, idStr string) (*model.Unit, error) {
	sid, err := id.FromHex(idStr)
	if err != nil {
		return nil, err
	}
	unit, err := d.query.WithContext(ctx).Unit.Where(d.query.Unit.ID.Eq(sid)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return unit, err
}
