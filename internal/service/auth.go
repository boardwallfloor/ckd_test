package service

import (
	"boardwallfloor/ckd/internal/db"
	utils "boardwallfloor/ckd/internal/util"
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthClaims struct {
	UserID int32  `json:"user_id"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

type AuthService interface {
	Login(ctx context.Context, email, password string) (string, error)
	Register(ctx context.Context, params db.CreateUserParams) (db.User, error)
}

type authService struct {
	queries *db.Queries
}

func NewAuthService(queries *db.Queries) AuthService {
	return &authService{queries: queries}
}

func getJwtKey() []byte {
	key := os.Getenv("JWT_SECRET_KEY")
	if key == "" {
		key = "my_super_secret_key"
	}
	return []byte(key)
}

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	fetchedUser, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return "", utils.ServiceError{
			Message:     "Wrong email or password",
			ServiceName: "Auth",
			ErrorMsg:    err,
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(fetchedUser.Password), []byte(password))
	if err != nil {
		return "", utils.ServiceError{
			Message:     "Wrong email or password",
			ServiceName: "Auth",
			ErrorMsg:    err,
		}
	}

	claims := AuthClaims{
		UserID: fetchedUser.ID,
		Name:   fetchedUser.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			Issuer:    "Auth Service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(getJwtKey())
	if err != nil {
		return "", utils.ServiceError{
			Message:     "Failure to sign JWT",
			ServiceName: "Auth",
			ErrorMsg:    err,
		}
	}

	return tokenStr, nil
}

func (s *authService) Register(ctx context.Context, params db.CreateUserParams) (db.User, error) {
	_, err := s.queries.GetUserByEmail(ctx, params.Email)
	if err == nil {
		return db.User{}, utils.ServiceError{
			Message:     "Email is already registered",
			ServiceName: "Auth",
			ErrorMsg:    errors.New("email already in use"),
		}
	}
	if err != sql.ErrNoRows {
		return db.User{}, utils.ServiceError{
			Message:     "Failed to check user",
			ServiceName: "Auth",
			ErrorMsg:    err,
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return db.User{}, utils.ServiceError{
			Message:     "Failed to hash password",
			ServiceName: "Auth",
			ErrorMsg:    err,
		}
	}
	params.Password = string(hashedPassword)
	params.Token = sql.NullString{
		Valid:  false,
		String: "",
	}
	log.Println(params)

	newUser, err := s.queries.CreateUser(ctx, params)
	if err != nil {
		return db.User{}, utils.ServiceError{
			Message:     "Failed to create user. Email may already be in use.",
			ServiceName: "Auth",
			ErrorMsg:    err,
		}
	}
	newUser.Password = ""
	return newUser, nil
}

func ValidateToken(tokenString string) (*AuthClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return getJwtKey(), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
