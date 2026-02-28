
namespace go basicuser

include "./basic_user.thrift"

// 基础用户相关接口
service BasicUserService {
  basic_user.BasicUserRegisterResp BasicUserRegister(1:basic_user.BasicUserRegisterReq req) (api.post = "/basic_user/register")
  basic_user.BasicUserLoginResp BasicUserLogin(1:basic_user.BasicUserLoginReq req)  (api.post = "/basic_user/login")
  basic_user.BasicUserResetPasswordResp BasicUserResetPassword(1:basic_user.BasicUserResetPasswordReq req) (api.post = "/basic_user/reset_password")
  basic_user.BasicUserCreateResp CreateBasicUser(1:basic_user.BasicUserCreateReq req) (api.post = "/basic_user/create")
}