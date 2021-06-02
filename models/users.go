package models

import (
	"errors"

	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrNotFound = errors.New("models: resource not found")
)

func NewUserService(connInfo string) (*UserService, error) {
	db, err := gorm.Open(postgres.Open(connInfo), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &UserService{db: db}, nil
}

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;uniqueIndex"`
}

type UserService struct {
	db *gorm.DB
}

func (us *UserService) ById(id uint) (*User, error) {
	var u User
	err := us.db.Where("id = ?", id).First(&u).Error

	switch err {
	case nil:
		return &u, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (us *UserService) Close() error {
	sqlDB, err := us.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (us *UserService) DestructiveReset() {
	us.db.Migrator().DropTable(&User{})
	us.db.AutoMigrate(&User{})
}
