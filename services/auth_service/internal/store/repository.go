package store

import (
	"time"

	"github.com/alexey-dobry/fileshare/services/auth_service/internal/domain/model"
)

type UserRepository interface {
	Add(model.Credentials) error

	GetOneByMail(email string) (model.Credentials, error)
	GetOneByID(ID string) (model.Credentials, error)

	UpdatePassword(ID string, newHash string) error

	Delete(email string) error
}

type TokenBlacklistRepository interface {
	BlacklistAccessToken(jti string, expiresIn time.Duration) error
	IsAccessTokenBlacklisted(jti string) (bool, error)

	StoreLogoutSession(jti string, expiresIn time.Duration) error
	IsSessionLoggedOut(jti string) (bool, error)
}
