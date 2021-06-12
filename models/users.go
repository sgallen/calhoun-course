package models

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const userPasswordPepper = "ajDYmXelRfAC06G3oEjjXT2/+BucicO4"

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
	Name         string
	Email        string `gorm:"not null;uniqueIndex"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
}

type UserService struct {
	db *gorm.DB
}

func (us *UserService) ById(id uint) (*User, error) {
	var user User
	tx := us.db.Where("id = ?", id)
	err := first(tx, &user)
	return &user, err
}

func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	tx := us.db.Where("email = ?", email)
	err := first(tx, &user)
	return &user, err
}

func first(tx *gorm.DB, user *User) error {
	err := tx.First(user).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

func (us *UserService) Create(user *User) error {
	passwordBytes := []byte(user.Password + userPasswordPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(
		passwordBytes,
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedBytes)
	user.Password = ""

	return us.db.Create(user).Error
}

func (us *UserService) Update(user *User) error {
	return us.db.Save(user).Error
}

func (us *UserService) Delete(pk uint) error {
	// if pk < 1 {
	// 	return ErrInvalidID
	// }
	// user := User{Model: gorm.Model{ID: pk}}
	// return us.db.Delete(&user).Error
	return us.db.Delete(&User{}, pk).Error
}

func (us *UserService) Close() error {
	sqlDB, err := us.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (us *UserService) DestructiveReset() error {
	if err := us.db.Migrator().DropTable(&User{}); err != nil {
		return err
	}
	return us.AutoMigrate()
}

func (us *UserService) AutoMigrate() error {
	return us.db.AutoMigrate(&User{})
}
