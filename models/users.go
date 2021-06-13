package models

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"lenslocked.com/hash"
	"lenslocked.com/rand"
)

const userPasswordPepper = "ajDYmXelRfAC06G3oEjjXT2/+BucicO4"
const secretKey = "dPVstQi1BUGMr8Q0jYU9iO9n9VRIgoGe"

var (
	ErrNotFound          = errors.New("models: resource not found")
	ErrIncorrectPassword = errors.New("models: incorrect password")
	ErrInvalidPassword   = errors.New("models: invalid password")
)

func NewUserService(connInfo string) (*UserService, error) {
	db, err := gorm.Open(postgres.Open(connInfo), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &UserService{db: db, hmac: hash.NewHMAC(secretKey)}, nil
}

type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;uniqueIndex"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;uniqueIndex"`
}

type UserService struct {
	db   *gorm.DB
	hmac *hash.HMAC
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

func (us *UserService) ByRemember(remember string) (*User, error) {
	rememberHash := us.hmac.Hash(remember)
	var user User
	tx := us.db.Where("remember_hash = ?", rememberHash)
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

func pepperPassword(password string) []byte {
	return []byte(password + userPasswordPepper)
}

func validatePassword(password string) error {
	if len(password) >= 10 {
		return nil
	}
	return ErrInvalidPassword
}

func (us *UserService) Create(user *User) error {
	err := validatePassword(user.Password)
	if err != nil {
		return err
	}

	hashedBytes, err := bcrypt.GenerateFromPassword(
		pepperPassword(user.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedBytes)
	user.Password = ""

	if user.Remember == "" {
		remember, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = remember
	}
	user.RememberHash = us.hmac.Hash(user.Remember)

	return us.db.Create(user).Error
}

func (us *UserService) Authenticate(email string, password string) (*User, error) {
	user, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		pepperPassword(password),
	)

	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrIncorrectPassword
		default:
			return nil, err
		}
	}

	return user, nil
}

func (us *UserService) Update(user *User) error {
	if user.Remember != "" {
		user.RememberHash = us.hmac.Hash(user.Remember)
	}
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
