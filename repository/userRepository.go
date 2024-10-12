package repository

import (
	"ginorm/model"
)

type UserRepository struct {
	*Repository
}

func NewUserRepository(name ...string) *UserRepository {
	return &UserRepository{
		Repository: NewRepository(name...),
	}
}

func (r *UserRepository) GetUser(ID interface{}) (model.User, error) {
	var user model.User
	err := r.DB.Model(&model.User{}).Where("id = ?", ID).First(&user).Error
	return user, err
}

func (r *UserRepository) CreateUser(user *model.User) error {
	err := r.DB.Model(&model.User{}).Create(user).Error
	return err
}

func (r *UserRepository) UpdateUser(user *model.User) error {
	err := r.DB.Model(&model.User{}).Save(user).Error
	return err
}

func (r *UserRepository) DeleteUser(ID interface{}) error {
	err := r.DB.Model(&model.User{}).Delete(ID).Error
	return err
}

func (r *UserRepository) CountUserByNickname(nickname string) (int64, error) {
	count := int64(0)
	err := r.DB.Model(&model.User{}).Where("nickname = ?", nickname).Count(&count).Error
	return count, err
}

func (r *UserRepository) CountUserByUsername(username string) (int64, error) {
	count := int64(0)
	err := r.DB.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	return count, err
}

func (r *UserRepository) GetUserByUsername(username string) (model.User, error) {
	var user model.User
	err := r.DB.Model(&model.User{}).Where("username = ?", username).First(&user).Error
	return user, err
}
