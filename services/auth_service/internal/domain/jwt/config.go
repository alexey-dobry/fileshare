package jwt

type Config struct {
	AccessSecret  string `validate:"required" yaml:"access-secret"`
	RefreshSecret string `validate:"required" yaml:"refresh-secret"`
	TTL           TTL    `validate:"required" yaml:"ttl"`
}

type TTL struct {
	AccessTTL  string `validate:"required" yaml:"access-ttl"`
	RefreshTTL string `validate:"required" yaml:"refresh-ttl"`
}
