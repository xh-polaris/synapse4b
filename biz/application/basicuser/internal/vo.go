package internal

import (
	"github.com/xh-polaris/synapse4b/biz/api/model/basicuser"
	"github.com/xh-polaris/synapse4b/biz/domain/basicuser/entity"
	"github.com/xh-polaris/synapse4b/biz/pkg/lang/util"
)

func BasicUserPO2VO(u *entity.BasicUser) *basicuser.BasicUser {
	return &basicuser.BasicUser{
		BasicUserId: u.ID,
		UnitId:      u.UnitID,
		Code:        u.Code,
		Phone:       u.Phone,
		Email:       u.Email,
		Name:        util.Of(u.Name),
		Gender:      util.Of(int32(u.Gender)),
	}
}

// UnitPO2VO .
func UnitPO2VO(u *entity.Unit) *basicuser.Unit {
	return &basicuser.Unit{
		Id:         u.ID,
		Name:       util.Of(u.Name),
		CreateTime: u.CreatedAt.Unix(),
		UpdateTime: u.CreatedAt.Unix(),
	}
}
