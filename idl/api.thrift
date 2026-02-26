include "./basic_user/basic_user_svc.thrift"
include "./system/system_svc.thrift"
include "./thirdparty/thirdparty_svc.thrift"

namespace go synapse4b

service SystemService extends system_svc.SystemService {}
service BasicUserService extends basic_user_svc.BasicUserService {}
service ThirdPartyService extends thirdparty_svc.ThirdPartyService {}