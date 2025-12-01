package pg

import (
	"github.com/alexey-dobry/fileshare/services/auth_service/internal/domain/model"
	"github.com/google/uuid"
)

func (ur *Repository) Add(userCredentials model.Credentials) error {
	return ur.db.Create(&userCredentials).Error
}

func (ur *Repository) GetOneByMail(email string) (model.Credentials, error) {
	user := model.Credentials{}

	result := ur.db.Select("uuid", "email", "hash_password").Where("email = ?", email).First(&user)
	if result.Error != nil {
		return model.Credentials{}, result.Error
	}
	return user, nil
}

func (ur *Repository) GetOneByID(ID uuid.UUID) (model.Credentials, error) {
	user := model.Credentials{}

	result := ur.db.Select("uuid", "email", "hash_password").Where("uuid = ?", ID).First(&user)
	if result.Error != nil {
		return model.Credentials{}, result.Error
	}
	return user, nil
}

func (ur *Repository) UpdatePassword(ID uuid.UUID, newHash string) error {
	result := ur.db.Update("password_hash", newHash).Where("uuid = ?", ID)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
