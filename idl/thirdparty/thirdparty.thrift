namespace go thirdparty

include "base.thrift"

struct ThirdPartyBasicUser {
    1: string basicUserId
}


struct ThirdPartyLoginReq {
    1: string thirdparty,
    2: string ticket,
}

struct ThirdPartyLoginResp {
    1:   base.Response           resp,
    2:   string                  token,
    255: ThirdPartyBasicUser     basicUser,
}