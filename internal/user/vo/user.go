package vo

import (
	"dashboard-server/internal/user/po"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type LoginReq struct {
	AccountName string `json:"account_name"`
	Password    string `json:"password"`
}

type LoginResp struct {
	Message string           `json:"message"`
	Users   map[string]*User `json:"users"`
}

type User struct {
	Username   string `json:"username"`
	ProfileURL string `json:"url"`
}

type RegisterUsersReq struct {
	AccountName string `json:"account_name"`
	Username1   string `json:"username_1"`
	ProfileURL1 string `json:"profile_url_1"`
	Username2   string `json:"username_2"`
	ProfileURL2 string `json:"profile_url_2"`
	Username3   string `json:"username_3"`
	ProfileURL3 string `json:"profile_url_3"`
}

type RegisterUsersResp struct {
	Message          string `json:"message"`
	Username1Created bool   `json:"username_1_created"`
	Username2Created bool   `json:"username_2_created"`
	Username3Created bool   `json:"username_3_created"`
}

type CreateUserReq struct {
	AccountName string `json:"account_name"`
	Username    string `json:"username"`
}

func Login(accountName string, password string) (*LoginResp, error) {
	var (
		account *po.Account
		poUsers []*po.User
		err     error
	)

	account, err = po.GetAccount(accountName)
	if account == nil && err == nil {
		//account does not exist, create account
		inputPasswordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		if err != nil {
			return nil, errors.Wrap(err, "unable to generate password hash")
		}

		account, err = po.CreateAccount(accountName, inputPasswordHash)
		if err != nil {
			//account does not exist and failed to create account
			return nil, errors.Wrap(err, "unable to create account")
		}

		return &LoginResp{Message: "Created New Account", Users: map[string]*User{}}, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "get account error")
	}

	if err := bcrypt.CompareHashAndPassword(account.PasswordHash,
		[]byte(password)); err != nil {
		return nil, errors.Wrap(err, "password incorrect")
	}

	poUsers, err = po.GetUsersFromAccount(accountName)
	if err != nil {
		return nil, errors.Wrap(err, "get users error")
	}

	users := make(map[string]*User)
	for _, user := range poUsers {
		users[user.Username] = &User{user.Username, user.ProfileURL}
	}

	return &LoginResp{Message: "Login Successful", Users: users}, nil
}

func RegisterUsersIfNotExist(req RegisterUsersReq) (*RegisterUsersResp, error) {
	var (
		resp RegisterUsersResp
		user *po.User
		err  error
	)

	user, err = po.GetUserInAccount(req.AccountName, req.Username1)
	if user == nil && err == nil {
		err = po.CreateUserFromAccount(req.AccountName, req.Username1, req.ProfileURL1)
		if err != nil {
			return nil, errors.Wrapf(err, "error creating %v+", req.Username1)
		}
		resp.Username1Created = true
	}

	user, err = po.GetUserInAccount(req.AccountName, req.Username2)
	if user == nil && err == nil {
		err = po.CreateUserFromAccount(req.AccountName, req.Username2, req.ProfileURL2)
		if err != nil {
			return nil, errors.Wrapf(err, "error creating %v+", req.Username2)
		}
		resp.Username2Created = true
	}

	user, err = po.GetUserInAccount(req.AccountName, req.Username3)
	if user == nil && err == nil {
		err = po.CreateUserFromAccount(req.AccountName, req.Username3, req.ProfileURL3)
		if err != nil {
			return nil, errors.Wrapf(err, "error creating %v+", req.Username3)
		}
		resp.Username3Created = true
	}

	return &resp, nil
}
