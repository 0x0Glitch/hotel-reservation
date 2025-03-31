package middleware

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error{
	fmt.Println("--JWT AUTH")
	token := c.Get("X-Api-Token")
	
	if err := parseToken(token); err!=nil{
	return err
}
	return nil
}



func parseToken(tokenStr string) error{
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _,ok := token.Method.(*jwt.SigningMethodHMAC); !ok{
			fmt.Println("invalid signing methods",token.Header)
			return nil, fmt.Errorf("Unauthorized")
		}
		secret := os.Getenv("JWT_SECRET")
		fmt.Println("never print secret",secret)
		return []byte(secret), nil
})
if err != nil{
	fmt.Println("failed to parse JWT token",err)
	return fmt.Errorf("unauthorized")
}
 

if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid{
	fmt.Println(claims)
} 
	return fmt.Errorf("unauthorized")

}