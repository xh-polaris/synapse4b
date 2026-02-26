package thirdparty

import (
	"context"

	"github.com/xh-polaris/synapse/biz/application/base/token"
	"github.com/xh-polaris/synapse/biz/application/thirdparty/internal"
	"github.com/xh-polaris/synapse/biz/conf"
	"github.com/xh-polaris/synapse/biz/domain/thirdparty/entity"
	"github.com/xh-polaris/synapse/biz/pkg/errorx"
	"github.com/xh-polaris/synapse/biz/types/errno"
	model "github.com/xh-polaris/synapse4b/biz/api/model/thirdparty"
	application "github.com/xh-polaris/synapse4b/biz/application/internal"
	thirdparty "github.com/xh-polaris/synapse4b/biz/domain/thirdparty/service"
)

const (
	bupt = "bupt" // 北邮
)

var ThirdPartySVC = &ThirdPartyService{}

type ThirdPartyService struct {
	DomainSVC thirdparty.ThirdParty
}

// Login 第三方用户登录, 不存在则注册
func (s *ThirdPartyService) Login(ctx context.Context, req *model.ThirdPartyLoginReq) (_ *model.ThirdPartyLoginResp, err error) {
	var u *entity.ThirdPartyUser
	switch req.GetThirdparty() {
	case bupt:
		u, err = s.DomainSVC.BUPTLogin(ctx, req.Ticket)
	default:
		return nil, errorx.New(errno.UnSupportThirdParty)
	}
	if err != nil {
		return nil, err
	}

	info := &token.Info{BasicUserId: u.ID, UserRole: ""}
	jwt, err := token.SignJWT(conf.GetConfig().Token, info)
	if err != nil {
		return nil, err
	}
	resp := &model.ThirdPartyLoginResp{
		Resp:      application.Success(),
		Token:     jwt,
		BasicUser: internal.ThirdPartyUser2BasicUserVO(u),
	}
	return resp, nil
}
