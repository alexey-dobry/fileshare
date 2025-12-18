package grpc

type Config struct {
	PublicPort  string `validate:"required" yaml:"public_port"`
	GatewayPort string `validate:"required" yaml:"gateway_port"`
}
