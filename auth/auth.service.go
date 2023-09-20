package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/space-w-alker/chat-room-server/model/generic"
	"github.com/space-w-alker/chat-room-server/model/user"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type CustomClaim struct {
	Value string
	jwt.RegisteredClaims
}

type SignInDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=50"`
}
func SignUp(dto *user.CreateUserDTO) (jwtToken string, err error) {
	dto.Password, err = HashPassword(dto.Password)
	if err != nil {
		return "", err
	}
	_, err = user.Create(dto)
	if err != nil {
		return "", err
	}
	newUser, err := user.GetWhere(&user.GetOrUpdateUserDTO{Email: dto.Email}, &generic.PaginationArgs{})
	if err != nil {
		return "", err
	}
	jwtToken, err = SignJWT(newUser[0].Id)
	return jwtToken, err
}

func SignIn(email string, password string) (jwtToken string, err error) {
	u, err := user.GetWhere(&user.GetOrUpdateUserDTO{Email: email}, &generic.PaginationArgs{})
	if err != nil {
		return "", err
	}
	if (len(u) == 0){
		return "", fmt.Errorf("user not found")
	}
	if CheckPasswordHash(password, u[0].Password){
		jwtToken, err = SignJWT(u[0].Id)
		return jwtToken, err;
	}else {
		return "", fmt.Errorf("wrong password")
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func SignJWT(value string) (token string, err error) {
	claim := CustomClaim{Value: value, RegisteredClaims: jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour * 365)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    "test",
		Subject:   "somebody",
		ID:        "1",
		Audience:  []string{"somebody_else"},
	}}
	t := jwt.NewWithClaims(jwt.SigningMethodHS512, claim)
	secret := viper.GetString("JWT_SECRET")
	token, err = t.SignedString([]byte(secret))
	return token, err
}

func VerifyJWT(tokenString string) (value string, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return viper.GetStringSlice("JWT_SECRET"), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return fmt.Sprintf("%v", claims["value"]), nil
	} else {
		return "", fmt.Errorf("unknown error")
	}
}

func GetUserByToken(tokenString string) (u user.User, err error)  {
	uid, err := VerifyJWT(tokenString)
	if err != nil {
		return u, err
	}
	u, err = user.GetById(&uid)
	if err != nil {
		return u, err
	}
	return u, nil
}
