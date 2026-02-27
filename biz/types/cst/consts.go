package cst

// AuthType
const (
	AuthTypePhone         = "phone"
	AuthTypeEmail         = "email"
	AuthTypeCode          = "code"
	AuthTypePhoneVerify   = "phone-verify"
	AuthTypePhonePassword = "phone-password"
	AuthTypeCodePassword  = "code-password"
	AuthTypeEmailPassword = "email-password"
	AuthTypeEmailVerify   = "email-verify"
)

// Token 中存储的信息
const (
	TokenInfo        = "token_info"
	TokenBasicUserID = "basic_user_id"
	TokenUnitID      = "unit_id"
	TokenCode        = "code"
	TokenPhone       = "phone"
	TokenEmail       = "email"
	TokenLoginTime   = "login_time"
	TokenAuthType    = "auth_type"
	TokenExtra       = "extra"
)
