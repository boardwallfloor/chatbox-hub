package inmem

import (
	pbAuth "boardwallfloor/auth_module/pb/auth/v1"
	"context"
)

type TestUser struct {
	id       int
	username string
	hashpass string
	email    string
}

type InMemStorage struct {
	userList []TestUser
	session  map[string]string
	pbAuth.UnimplementedAuthServiceServer
}

func (im *InMemStorage) SignUp(ctx context.Context, in *pbAuth.SignUpRequest) (*pbAuth.SignUpResponse, error) {
	for _, v := range im.userList {
		if v.username == in.GetUsername() {
			return &pbAuth.SignUpResponse{Status: false, Desc: "Username are already registered"}, nil
		}
	}
	regUser := TestUser{
		id:       len(im.userList),
		username: in.GetUsername(),
		hashpass: in.GetPassword(),
		email:    in.GetEmail(),
	}
	im.userList = append(im.userList, regUser)
	return &pbAuth.SignUpResponse{Status: true, Desc: "Registration successfull", UserId: int32(regUser.id)}, nil
}

func (im *InMemStorage) SignIn(ctx context.Context, in *pbAuth.SignInRequest) (*pbAuth.SignInResponse, error) {
	for _, v := range im.userList {
		if v.username == in.GetUsername() && v.hashpass == in.GetPassword() {
			return &pbAuth.SignInResponse{Status: true, Desc: "Login successfull", UserId: int32(v.id)}, nil
		}
	}
	return &pbAuth.SignInResponse{Status: false, Desc: "Login unsuccessfull"}, nil
}
