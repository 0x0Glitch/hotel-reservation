package api

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/0x0Glitch/hotel-reservation/db"
	"github.com/0x0Glitch/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct{
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler{
	return &AuthHandler{
		userStore: userStore,
	}
}

type AuthParams struct{
	Email    string    `json:"email"`
	Password string    `json:"password"`
}


type AuthResponse struct{
	User *types.User `json:"User"`
	Token string	 `json:"token"`

}
func (h *AuthHandler) HandleAuthentication(c *fiber.Ctx) error{
	var params AuthParams
	if err := c.BodyParser(&params); err != nil{
		return err
	}
	user, err := h.userStore.GetUserByEmail(c.Context(),params.Email)
	if err != nil{
		if errors.Is(err,mongo.ErrNoDocuments){
			return fmt.Errorf("invalid credentials")
		}
		return err
	}
	if !types.IsValidPassword(user.EncryptedPassword,params.Password){
		return fmt.Errorf("Invalid credentials")
	}
	token := createTokenFromUser(user)
	resp := AuthResponse{
		User:user,
		Token: token,
	}
	return c.JSON(resp)
	

	return nil
}
func createTokenFromUser(user *types.User) string{
	now := time.Now()
	expires := now.Add(time.Hour*4).Unix()
	claims := jwt.MapClaims{
		"id":        user.ID,
		"email":     user.Email,
		"expires": expires,

	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	secret := os.Getenv("JWT_SECRET")
	tokenStr,err := token.SignedString([]byte(secret))
	fmt.Println(secret)
if err != nil{
	fmt.Println("Failed to sign token with secret",err)
}
return tokenStr
}