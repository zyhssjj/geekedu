module geekedu/logic-server

go 1.25.3

require (
	geekedu/common v0.0.0-00010101000000-000000000000
	geekedu/proto v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.44.0
	google.golang.org/grpc v1.78.0
	gorm.io/driver/mysql v1.5.2
	gorm.io/gorm v1.25.3
)

require (
	github.com/aliyun/aliyun-oss-go-sdk v3.0.2+incompatible // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	golang.org/x/time v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251029180050-ab9386a59fda // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace geekedu/common => ../common

replace geekedu/proto => ../proto
