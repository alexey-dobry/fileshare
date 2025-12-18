package minio

type Config struct {
	Port      string `validate:"required" yaml:"port"`
	Host      string `validate:"required" yaml:"host"`
	AccessKey string `validate:"required" yaml:"access_key"`
	SecretKey string `validate:"required" yaml:"secret_key"`
	Bucket    string `validate:"required" yaml:"bucket"`
}
