package internal

import (
	"github.com/xh-polaris/synapse4b/biz/api/model/basicuser"
	"github.com/xh-polaris/synapse4b/biz/domain/basicuser/entity"
)

func BasicUserPO2VO(u *entity.BasicUser) *basicuser.BasicUser {
	return &basicuser.BasicUser{
		BasicUserId: u.ID,
	}
}
