package codegen

// defaults package
//go:generate go run ./cmd/defaultsconfig -json ./models/sdk-default-configuration.json -output=../../aws/defaults/defaults.go
//go:generate go run ./cmd/defaultsconfig -json ./models/testdata/sdk-default-configuration.json -output=../../aws/defaults/defaults_codegen_test.go -r v1TestResolver
//go:generate gofmt -s -w ../../aws/defaults/defaults.go ../../aws/defaults/defaults_codegen_test.go
// defaults mode constants
//go:generate go run ./cmd/defaultsmode -json ./models/sdk-default-configuration.json -output=../../aws/defaultsmode.go
//go:generate gofmt -s -w ../../aws/defaultsmode.go
