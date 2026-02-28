package tencent

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/cloudwego/hertz/pkg/common/json"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tencenterrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencentemail "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"
	"github.com/xh-polaris/synapse4b/biz/conf"
	"github.com/xh-polaris/synapse4b/biz/infra/contract/cache"
	"github.com/xh-polaris/synapse4b/biz/infra/contract/email"
	"github.com/xh-polaris/synapse4b/biz/pkg/logs"
)

const (
	// TriggerType 邮件触发类型 0:非触发类，默认类型，营销类邮件、非即时类邮件等选择此类型
	// 1:触发类，验证码等即时发送类邮件，若邮件超过一定大小，系统会自动选择非触发类型通道
	TriggerType    = uint64(1)
	NonTriggerType = uint64(0)
)

type tencentEmail struct {
	cache  *email.Cache
	client *tencentemail.Client
}

func New(ctx context.Context, cache *email.Cache, account, token string) (email.Provider, error) {
	s, err := getTencentEmailProvider(ctx, cache, account, token)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func getTencentEmailProvider(_ context.Context, cache *email.Cache, secretId, secretKey string) (email.Provider, error) {
	credential := common.NewCredential(secretId, secretKey)
	client, err := tencentemail.NewClient(credential, "ap-guangzhou", profile.NewClientProfile())
	if err != nil {
		return nil, err
	}
	return &tencentEmail{cache: cache, client: client}, nil
}

func (t *tencentEmail) Send(ctx context.Context, app, cause, email string, param *email.EmailParam) (err error) {
	// 发送邮件
	if _, err = t.send(ctx, app, cause, email, param); err != nil {
		return err
	}
	expire := param.Expire

	if err = t.cache.Store(ctx, app, cause, email, param.Code, expire); err != nil {
		return err
	}
	return nil
}

func (t *tencentEmail) send(_ context.Context, app, cause, email string, param *email.EmailParam) (map[string]any, error) {
	// 参数设置
	req := tencentemail.NewSendEmailRequest()

	req.FromEmailAddress = common.StringPtr(conf.GetConfig().Email.Account)
	req.Destination = common.StringPtrs([]string{email})
	req.Subject = common.StringPtr(conf.GetConfig().Email.SubjectConf[app][cause])
	req.Template = &tencentemail.Template{
		TemplateID: common.Uint64Ptr(conf.GetConfig().Email.TemplateConf[app][cause]),
		TemplateData: common.StringPtr(
			fmt.Sprintf("{\"code\":\"%s\",\"expire\":\"%s\"}",
				param.Code,
				strconv.Itoa(int(param.Expire.Minutes()))),
		),
	}
	req.TriggerType = common.Uint64Ptr(TriggerType)

	// 发送响应
	resp, err := t.client.SendEmail(req)

	// SDK 错误
	var tencentCloudSDKError *tencenterrors.TencentCloudSDKError
	if errors.As(err, &tencentCloudSDKError) {
		logs.Errorf("An Tencent API error has returned: %s", err)
	}
	if err != nil {
		return nil, err
	}
	// 正常响应
	r, respStr := map[string]any{}, resp.ToJsonString()
	if err = json.Unmarshal([]byte(respStr), &r); err != nil {
		logs.Infof("Tencent email return resp %s but unmarsharl failed %s", respStr, err)
	}
	return r, nil
}
func (t *tencentEmail) Check(ctx context.Context, app, cause, email, code string) (bool, error) {
	if code == "xh-polaris" && conf.GetConfig().State == "test" {
		return true, nil
	}
	ori, err := t.cache.Load(ctx, app, cause, email)
	if errors.Is(err, cache.Nil) {
		return false, nil
	}
	return ori == code, err
}
