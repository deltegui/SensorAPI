package domain

import (
	"golang.org/x/crypto/bcrypt"
)

type LoginUserRequest struct {
	UserName     string
	UserPassword string
}

type LoginUserCase UseCase

func NewLoginUserCase(userRepo UserRepo) LoginUserCase {
	return func(req UseCaseRequest) (UseCaseResponse, error) {
		loginReq := req.(LoginUserRequest)
		user, err := userRepo.GetUserByName(loginReq.UserName)
		if err != nil {
			return nil, UserNotFoundErr
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.UserPassword)); err != nil {
			return nil, NotValidPasswordErr
		}
		return user, nil
	}
}
