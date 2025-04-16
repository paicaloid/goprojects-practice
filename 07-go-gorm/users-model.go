package main

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string
}

func CreateUser(db *gorm.DB, user *User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password), bcrypt.DefaultCost,
	)

	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	result := db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func LoginUser(db *gorm.DB, user *User) (string, error) {
	resUser := new(User)
	result := db.Where("email = ?", user.Email).First(resUser)
	if result.Error != nil {
		return "", result.Error
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(resUser.Password), []byte(user.Password),
	); err != nil {
		return "", err
	}

	// Create token
	tempSecret := "TestSecret"
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = resUser.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(tempSecret))
	if err != nil {
		return "", err
	}

	return t, nil
}
