package pg

import (
	"github.com/alexey-dobry/fileshare/services/auth_service/internal/model"
)

func (ur *Repository) Add(userCredentials model.Credentials) error {
	return ur.db.Create(&userCredentials).Error
}

func (ur *Repository) GetOneByMail(email string) (model.Credentials, error) {
	user := model.Credentials{}

	result := ur.db.Select("username", "hash_password", "first_name", "last_name", "is_admin").Where("email = ?", email).First(&user)
	if result.Error != nil {
		return model.Credentials{}, result.Error
	}
	return user, nil
}

func (ur *Repository) GetOneByID(ID uint) (model.Credentials, error) {
	user := model.Credentials{}

	result := ur.db.Select("username", "first_name", "last_name", "is_admin").Where("id = ?", ID).First(&user)
	if result.Error != nil {
		return model.Credentials{}, result.Error
	}
	return user, nil
}

func (ur *Repository) UpdatePassword(ID uint, newHash string) error {
	result := ur.db.Update("password_hash", newHash).Where("id = ?", ID)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
