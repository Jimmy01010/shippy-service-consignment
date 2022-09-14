package main

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/Jimmy01010/protocol/shippy-user"
	"golang.org/x/crypto/bcrypt"
)

//// CustomClaims is our custom metadata, which will be hashed
//// and sent as the second segment in our JWT
//type CustomClaims struct {
//	User *pb.User
//	// jwt.StandardClaims
//}

type authAble interface {
	Decode(token string) (*CustomClaims, error)
	Encode(user *pb.User) (string, error)
}

/*
Create(context.Context, *User, *Response) error
Get(context.Context, *User, *Response) error
GetAll(context.Context, *Request, *Response) error
Auth(context.Context, *User, *Token) error
*/
type handler struct {
	repository   Repository
	tokenService authAble
}

// Get 通过userID获取用户信息
func (s *handler) Get(ctx context.Context, req *pb.User, res *pb.Response) error {
	result, err := s.repository.Get(ctx, req.Id)
	if err != nil {
		return err
	}

	user := UnmarshalUser(result)
	res.User = user

	return nil
}

// Create 创建一个新用户
func (s *handler) Create(ctx context.Context, user *pb.User, res *pb.Response) error {
	// 保存哈希后的密码
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPass)

	if err := s.repository.Create(ctx, MarshalUser(user)); err != nil {
		return err
	}

	// Strip the password back out, so's we're not returning it
	user.Password = ""
	res.User = user

	return nil
}

// GetAll 获取所有用户
func (s *handler) GetAll(ctx context.Context, req *pb.Request, res *pb.Response) error {
	result, err := s.repository.GetAll(ctx)
	if err != nil {
		return err
	}

	res.Users = UnmarshalUserCollection(result)
	return nil
}

// Auth 用户认证, 用户通过密码登录时验证有效性, 并返回该次登录的token
func (s *handler) Auth(ctx context.Context, req *pb.User, res *pb.Token) error {
	user, err := s.repository.GetByEmail(ctx, req.Email)
	if err != nil {
		return fmt.Errorf("get user from db on auth failed: %s", err.Error())
	}

	// 将bcrypt散列密码与需要认证的明文密码进行比较
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return fmt.Errorf("user auth failed: %s", err.Error())
	}

	token, err := s.tokenService.Encode(req)
	if err != nil {
		return fmt.Errorf("encode auth failed: %s", err.Error())
	}

	res.Token = token

	return nil
}

// ValidateToken 验证token的合法性
func (s *handler) ValidateToken(ctx context.Context, req *pb.Token, res *pb.Token) error {
	// 解码并验证TOKEN
	claims, err := s.tokenService.Decode(req.Token)
	if err != nil {
		return err
	}

	// 在这里可以继续验证token的payload部分的有效性
	if claims.User.Email == "" {
		return errors.New("invalid user")
	}

	res.Valid = true
	return nil
}
