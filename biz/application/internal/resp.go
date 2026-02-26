package internal

import "github.com/xh-polaris/synapse4b/biz/api/model/base"

func Success() *base.Response {
	return &base.Response{Code: 0, Msg: ""}
}
