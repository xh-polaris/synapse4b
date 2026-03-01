namespace go basicuser

include "base.thrift"

struct BasicUser {
    1: string basicUserId
    2: optional string unitId,      // 学校ID
    3: optional string code,        // 学号
    4: optional string phone,       // 手机号
    5: optional string email,       // 邮箱
    6: optional string name,        // 姓名
    7: optional i32    gender,      // 性别
}

struct Unit {
    1: string Id,
    2: optional string name,
    3: i64 createTime,
    4: i64 updateTime,
}

/*
 基础用户相关API
 */

// 基础用户注册
struct BasicUserRegisterReq {
    1:   string          authType,    // 认证类型
    2:   string          authId,      // 认证id
    3:   optional string extraAuthId, // 扩展认证id
    4:   string          verify,      // 认证凭证
    5:   optional string password,    // 是否初始化密码
    255: base.App        app,
}

struct BasicUserRegisterResp {
    1:   base.Response resp,
    2:   string        token,
    255: BasicUser     basicUser,
}

// 基础用户登录
struct BasicUserLoginReq {
    1: string authType,
    2: string authId,
    3: optional string extraAuthId, // 扩展认证id
    4: string verify,
    255: base.App        app,
}

struct BasicUserLoginResp {
    1:   base.Response resp,
    2:   string token,
    3:   bool   new,
    255: BasicUser basicUser,
}

// 修改密码
struct BasicUserResetPasswordReq {
    1: string newPassword,
    2: optional string resetKey, // 重置码
    3: optional string basicUserId,
    255: base.App        app,
}

struct BasicUserResetPasswordResp {
    1:   base.Response resp,
}

// 创建基础用户
struct BasicUserCreateReq {
    1: optional string unitId,      // 学校ID
    2: optional string code,        // 学号
    3: optional string phone,       // 手机号
    4: optional string email,       // 邮箱
    5: optional string password,    // 密码
    6: optional i64 encryptType, // 密码加密类型
    254: optional string createKey, // 创建密钥
    255: base.App        app,
}

// 创建基础用户响应
struct BasicUserCreateResp {
    1: base.Response resp,
    255: BasicUser basicUser,
}

// 获取Unit
struct GetUnitReq {
    1: string unitId,
    255: base.App        app,
}

struct GetUnitResp {
    1: base.Response resp,
    255: Unit unit,
}

// 查询Unit
struct QueryUnitReq {
    1: optional string name,
    255: base.App        app,
}
struct QueryUnitResp {
    1: base.Response resp,
    255: Unit unit,
}

// 创建Unit
struct CreateUnitReq {
    1: string     name,
    254: optional string createKey, // 创建密钥
    255: base.App app,
}
struct CreateUnitResp {
    1: base.Response resp,
    255: Unit unit,
}
