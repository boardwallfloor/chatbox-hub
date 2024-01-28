// in this package i tries to use getter when accessing protobuf structs, it's dumb and i'm too lazy to change it
package sqlite

import (
	pbAuth "boardwallfloor/auth_module/pb/auth/v1"
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"log"
	"math/big"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type sqliteDB struct {
	pbAuth.UnimplementedAuthServiceServer
	Db *sql.DB
}

func NewSqliteDB(sqlConnString string) (sqliteDB, error) {
	sqlConn, err := sql.Open("sqlite3", sqlConnString)
	if err != nil {
		log.Println(err)
		return sqliteDB{}, err
	}

	sqlitedb := sqliteDB{Db: sqlConn}
	return sqlitedb, nil
}

func encryptPassword(pass string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return string(hashedPass), nil
}

func generateSessionId() (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ses := make([]byte, 16)
	for i := 0; i < 16; i++ {
		randInt, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ses[i] = letters[randInt.Int64()]
	}
	return string(ses), nil
}

func (sl *sqliteDB) SignUp(ctx context.Context, in *pbAuth.SignUpRequest) (*pbAuth.SignUpResponse, error) {
	hashPass, err := encryptPassword(in.GetPassword())
	if err != nil {
		log.Println(err)
		return &pbAuth.SignUpResponse{}, err
	}
	res, err := sl.Db.Exec("insert or ignore into user (username,password,email)values(?,?,?)", in.GetUsername(), hashPass, in.GetEmail())
	if err != nil {
		return &pbAuth.SignUpResponse{}, err
	}

	rowsChanged, err := res.RowsAffected()
	if err != nil {
		return &pbAuth.SignUpResponse{}, err
	}

	if rowsChanged == 0 {
		return &pbAuth.SignUpResponse{Desc: "username is already taken", Status: false, UserId: 0}, nil
	}

	var userId int
	err = sl.Db.QueryRow("select user_id from user where username = ?", in.GetUsername()).Scan(&userId)
	if err != nil {
		return &pbAuth.SignUpResponse{}, err
	}
	return &pbAuth.SignUpResponse{UserId: int32(userId), Desc: "Sign up successfull", Status: true}, nil
}

// sign in always assume that associated session_id is not in session table in db, and sign in are required and called only there's no session_id
// either because expired session_id or first sign in
func (sl *sqliteDB) SignIn(ctx context.Context, in *pbAuth.SignInRequest) (*pbAuth.SignInResponse, error) {
	var userId int
	var hashPass string
	err := sl.Db.QueryRow("select user_id, password from user where username = ?", in.GetUsername()).Scan(&userId, &hashPass)
	if err == sql.ErrNoRows {
		return &pbAuth.SignInResponse{}, fmt.Errorf("user credential is not found")
	}
	if err != nil {
		return &pbAuth.SignInResponse{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(in.GetPassword()))
	if err != nil {
		return &pbAuth.SignInResponse{}, fmt.Errorf("user credential is not found")
	}

	sessionId, err := generateSessionId()
	if err != nil {
		return &pbAuth.SignInResponse{}, fmt.Errorf("session_id generation failed: %s", err)
	}
	_, err = sl.Db.Exec("insert into session (session_id,user_id,is_long_term) values (?, ?,?)", sessionId, userId, false)
	if err != nil {
		return &pbAuth.SignInResponse{}, err
	}

	return &pbAuth.SignInResponse{SessionId: sessionId, Desc: "Login successfull", Status: true}, nil
}

func (sl *sqliteDB) CheckLoggedIn(ctx context.Context, in *pbAuth.CheckLoggedInRequest) (*pbAuth.CheckLoggedInResponse, error) {
	var session_id string
	err := sl.Db.QueryRow("select session_id from session where session_id = ?", in.GetSessionId()).Scan(&session_id)
	if err == sql.ErrNoRows {
		return &pbAuth.CheckLoggedInResponse{IsLoggedIn: false}, nil
	}
	if err != nil {
		return &pbAuth.CheckLoggedInResponse{}, err
	}
	return &pbAuth.CheckLoggedInResponse{IsLoggedIn: true}, nil
}

func (sl *sqliteDB) Logout(ctx context.Context, in *pbAuth.LogoutRequest) (*pbAuth.LogoutResponse, error) {
	res, err := sl.Db.Exec("delete from session where session_id = ?", in.GetSession())
	if err != nil {
		return &pbAuth.LogoutResponse{}, err
	}

	rowsChanged, err := res.RowsAffected()
	if err != nil {
		return &pbAuth.LogoutResponse{}, err
	}

	if rowsChanged > 0 {
		return &pbAuth.LogoutResponse{Success: true}, nil
	} else {
		return &pbAuth.LogoutResponse{Success: false}, nil
	}
}
