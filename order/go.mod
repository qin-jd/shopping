module shopping/order

go 1.12

require (
	github.com/bwmarrin/snowflake v0.0.0-20180412010544-68117e6bbede
	github.com/golang/protobuf v1.3.1
	github.com/jinzhu/gorm v1.9.2
	github.com/micro/go-grpc v1.0.0
	github.com/micro/go-log v0.1.0
	github.com/micro/go-micro v1.0.0
	shopping/product v0.0.0
)

replace shopping/product => ../product
