module github.com/AMSt1010/microservises-course/week2/inventory

go 1.24.0

require (
	github.com/AMSt1010/microservises-course/week2/shared v0.0.0-00010101000000-000000000000
	github.com/brianvoe/gofakeit/v7 v7.17.1
	github.com/google/uuid v1.6.0
	github.com/samber/lo v1.49.1
	github.com/stretchr/testify v1.12.1
	google.golang.org/grpc v1.70.0
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/stretchr/objx v0.5.3 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.34.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250115164207-1a7da9e5054f // indirect
)

replace github.com/AMSt1010/microservises-course/week2/shared => ../shared
