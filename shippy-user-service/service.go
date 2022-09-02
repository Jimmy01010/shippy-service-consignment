package main

import pb "github.com/Jimmy01010/protocol/shippy-user"

type TokenService struct {
	repo Repository
}

// Decode todo: implement
func (TokenService) Decode(token string) (*CustomClaims, error) {
	return nil, nil
}

// Encode todo: implement
func (TokenService) Encode(user *pb.User) (string, error) {
	return "", nil
}
