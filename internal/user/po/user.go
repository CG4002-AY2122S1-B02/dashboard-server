package po

import (
	"dashboard-server/dbutils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Account struct {
	AccountName  string `gorm:"primaryKey"`
	PasswordHash []byte
}

type User struct {
	AccountName string `gorm:"primaryKey"`
	Username    string `gorm:"primaryKey"`
	ProfileURL  string
}

func CreateAccount(accountName string, passwordHash []byte) (*Account, error) {
	account := &Account{AccountName: accountName,
		PasswordHash: passwordHash}
	if err := dbutils.GetDB().Create(account).Error; err != nil {
		return nil, errors.Wrap(err, "db create account error")
	}

	return account, nil
}

func GetAccount(accountName string) (*Account, error) {
	var account Account
	err := dbutils.GetDB().Where(&Account{
		AccountName: accountName},
	).First(&account).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		//no account name found. Create account
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "db get account error")
	}

	return &account, nil
}

func GetUsersFromAccount(accountName string) ([]*User, error) {
	var users []*User
	err := dbutils.GetDB().Where(&User{
		AccountName: accountName},
	).Find(&users).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		//no users found
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "db get users from account error")
	}

	return users, nil
}

func GetUserInAccount(accountName string, username string) (*User, error) {
	var user User
	err := dbutils.GetDB().Where(&User{
		AccountName: accountName,
		Username:    username,
	},
	).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		//no users found
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "db get user error")
	}

	return &user, nil
}

func CreateUserFromAccount(accountName string, username string, profileURL string) error {
	err := dbutils.GetDB().Create(User{
		accountName,
		username,
		profileURL,
	}).Error
	if err != nil {
		return errors.Wrap(err, "db create user error")
	}

	return nil
}
