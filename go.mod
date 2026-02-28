module github.com/xh-polaris/synapse4b

go 1.24.6

require (
	github.com/apache/thrift v0.22.0
	github.com/bytedance/sonic v1.14.0
	github.com/cloudwego/hertz v0.10.2
	github.com/golang-jwt/jwt/v4 v4.5.2
	github.com/hertz-contrib/monitor-prometheus v0.1.3
	github.com/redis/go-redis/v9 v9.12.1
	github.com/zeromicro/go-zero v1.9.0
	go.mongodb.org/mongo-driver v1.17.4
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.63.0
	go.opentelemetry.io/contrib/propagators/b3 v1.38.0
	go.opentelemetry.io/otel v1.38.0
	golang.org/x/crypto v0.41.0
	gorm.io/datatypes v1.2.4
	gorm.io/driver/mysql v1.6.0
	gorm.io/gen v0.3.27
	gorm.io/gorm v1.31.0
	gorm.io/plugin/dbresolver v1.6.2
)

require (
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common v1.3.50 // indirect
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses v1.3.49 // indirect
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms v1.1.49 // indirect
)

replace github.com/apache/thrift => github.com/apache/thrift v0.13.0
