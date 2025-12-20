package grpc

type Config struct {
	PublicPort   string `validate:"required" yaml:"public_port"`
	InternalPort string `validate:"required" yaml:"internal_port"`
	GatewayPort  string `validate:"required" yaml:"gateway_port"`
	JWTSecret    string `validate:"required" yaml:"jwt_key"`
}
