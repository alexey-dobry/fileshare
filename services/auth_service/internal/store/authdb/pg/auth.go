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

	result := ur.db.Select("uuid", "email", "hash_password", "role").Where("email = ?", email).First(&user)
	if result.Error != nil {
		return model.Credentials{}, result.Error
	}
	return user, nil
}

func (ur *Repository) GetOneByID(id uuid.UUID) (model.Credentials, error) {
	user := model.Credentials{}

	result := ur.db.Select("uuid", "email", "hash_password", "role").Where("uuid = ?", id).First(&user)
	if result.Error != nil {
		return model.Credentials{}, result.Error
	}
	return user, nil
}

func (ur *Repository) Delete(email string) error {
	result := ur.db.Where("email = ?", email).Delete(model.Credentials{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (ur *Repository) UpdatePassword(id uuid.UUID, newHash string) error {
	result := ur.db.Update("password_hash", newHash).Where("uuid = ?", id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
