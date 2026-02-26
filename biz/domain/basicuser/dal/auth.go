package dal

import (
	"github.com/xh-polaris/synapse4b/biz/domain/basicuser/dal/query"
	"github.com/xh-polaris/synapse4b/biz/infra/contract/orm"
)

func NewAuthDAO(db *orm.DB) *AuthDAO {
	return &AuthDAO{query: query.Use(db)}
}

type AuthDAO struct {
	query *query.Query
}
