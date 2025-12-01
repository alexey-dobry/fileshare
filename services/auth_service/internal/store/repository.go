package store

import (
	"time"

	"github.com/alexey-dobry/fileshare/services/auth_service/internal/model"
)

type UserRepository interface {
	Add(model.Credentials) error

	GetOneByMail(email string) (model.Credentials, error)
	GetOneByID(ID uint) (model.Credentials, error)

	UpdatePassword(ID uint, newHash string) error
}

type TokenBlacklistRepository interface {
	BlacklistAccessToken(jti string, expiresIn time.Duration) error
	IsAccessTokenBlacklisted(jti string) (bool, error)

	StoreLogoutSession(jti string, expiresIn time.Duration) error
	IsSessionLoggedOut(jti string) (bool, error)
}
