package pg

import (
	"github.com/alexey-dobry/fileshare/services/auth_service/internal/domain/model"
)

func (ur *Repository) Add(userCredentials model.Credentials) error {
	return ur.db.Create(&userCredentials).Error
}

func (ur *Repository) GetOneByMail(email string) (model.Credentials, error) {
	user := model.Credentials{}

	result := ur.db.Select("uuid", "email", "password_hash", "role").Where("email = ?", email).First(&user)
	if result.Error != nil {
		return model.Credentials{}, result.Error
	}
	return user, nil
}

func (ur *Repository) GetOneByID(id string) (model.Credentials, error) {
	user := model.Credentials{}

	result := ur.db.Select("uuid", "email", "password_hash", "role").Where("uuid = ?", id).First(&user)
	if result.Error != nil {
		return model.Credentials{}, result.Error
	}
	return user, nil
}

func (ur *Repository) Delete(email string) error {
	result := ur.db.Where("email = ?", email).Delete(&model.Credentials{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (ur *Repository) UpdatePassword(id string, newHash string) error {
	result := ur.db.Model(&model.Credentials{}).Where("uuid = ?", id).Update("password_hash", newHash)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
