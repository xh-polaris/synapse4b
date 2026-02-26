namespace go system

include "base.thrift"

/*
 系统基础服务API
 */
struct SendVerifyCodeReq {
    1: string authType,
    2: string authId,
    3: i64 expire,  // 过期时间, 秒级时间戳, 最长24小时
    4: string cause, // 发验证码的原因, 需提前配置
    255: base.App app
}

struct SendVerifyCodeResp {
    1: base.Response resp
}

struct CheckVerifyCodeReq {
    1: string authType,
    2: string authId,
    3: string cause, // 发验证码的原因, 需提前配置
    4: string verify,
    255: base.App app
}
struct CheckVerifyCodeResp{
    1: base.Response resp
    2: bool verify
}

// 签发ticket
struct SignTicketReq {
    255: base.App app
}

struct SignTicketResp {
    1: base.Response resp
    2: string ticket      // 签发的token
}

// 换取token
struct ExchangeTicketReq {
    1: string ticket      // 签发的token
    2: string ticketKey   // app对应的兑换密钥
    255: base.App app
}

struct ExchangeTicketResp {
    1: base.Response resp
    2: string token      // 签发的token
}