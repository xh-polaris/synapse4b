namespace go basicuser

include "./thirdparty.thrift"

service ThirdPartyService {
    thirdparty.ThirdPartyLoginResp ThirdPartyLogin(1: thirdparty.ThirdPartyLoginReq req) (api.post = "/thirdparty/login")
}