module shopping/notification

go 1.12

require (
	github.com/golang/protobuf v1.3.1
	github.com/micro/go-grpc v1.0.1
	github.com/micro/go-log v0.1.0
	github.com/micro/go-micro v1.1.0
	github.com/micro/go-plugins v1.1.0
	product v0.0.0
	user v0.0.0
)

replace (
	product => ../product
	user => ../user
)

replace github.com/golang/lint => golang.org/x/lint v0.0.0-20190409202823-959b441ac422

replace github.com/testcontainers/testcontainer-go => github.com/testcontainers/testcontainers-go v0.0.3
