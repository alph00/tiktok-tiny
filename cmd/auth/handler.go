package main

import (
	"context"
	auth "github.com/alph00/tiktok-tiny/kitex_gen/auth"
)

// AuthServiceImpl implements the last service interface defined in the IDL.
type AuthServiceImpl struct{}

// Register implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) Register(ctx context.Context, req *auth.UserRegisterRequest) (resp *auth.UserRegisterResponse, err error) {
	// TODO: Your code here...
	return
}

// Login implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) Login(ctx context.Context, req *auth.UserLoginRequest) (resp *auth.UserLoginResponse, err error) {
	// TODO: Your code here...
	return
}
