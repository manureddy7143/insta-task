package repository

import (
	"github.com/manureddy7143/GolangStarter/source/model"
	"github.com/manureddy7143/GolangStarter/utils/database"
)

// AuthRepository -
type AuthRepository struct{}

//FindUsers - Finds Users based on the filters from DB
func (articleRepository AuthRepository) FindUsers(filter map[string]interface{}) ([]model.Users, error) {
	var us []model.Users
	result := database.GetInstance().Where(filter).Find(&us)
	if result.Error != nil {
		return us, result.Error
	}

	return us, nil
}

//CreatingUsers - Creating Users based on the filters .

func (articleRepository AuthRepository) CreateUsers(user model.Users) (int, error) {
	//Insert into Users table
	result := database.GetInstance().Create(&user)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(user.Id), nil
}
