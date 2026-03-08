package dto

import (
	bmsauth "github.com/pafthang/bms/internal/auth"
	serviceauth "github.com/pafthang/bms/internal/services/auth"
)

type AuthRegisterRequest struct {
	Email    string `json:"email" validate:"required,format=email"`
	Password string `json:"password" validate:"required,minlength=8,maxlength=128"`
}

type AuthLoginRequest struct {
	Email    string `json:"email" validate:"required,format=email"`
	Password string `json:"password" validate:"required,minlength=8,maxlength=128"`
}

type AuthRefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required,minlength=16"`
}

type AuthPayload struct {
	User   *serviceauth.UserView `json:"user"`
	Tokens *bmsauth.TokenPair    `json:"tokens"`
}

type RefreshPayload struct {
	Tokens *bmsauth.TokenPair `json:"tokens"`
}
