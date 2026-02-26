namespace go system

include "./system.thrift"

// 系统基础服务
service SystemService {
  system.SendVerifyCodeResp SendVerifyCode(1:system.SendVerifyCodeReq req) (api.post = "/system/send_verify_code")
  system.CheckVerifyCodeResp CheckVerifyCode(1:system.CheckVerifyCodeReq req) (api.post = "/system/check_verify_code")
  // 签发ticket
  system.SignTicketResp SignTicket(1:system.SignTicketReq req) (api.post = "/system/sign_ticket")
  // 兑换ticket
  system.ExchangeTicketResp ExchangeTicket(1:system.ExchangeTicketReq req) (api.post = "/system/exchange_ticket")
}