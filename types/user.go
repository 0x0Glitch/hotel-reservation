package types

import (
	"fmt"

	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)
const (
	bcryptCost = 12
	miniFirstNameLen = 2
	miniLastNameLen= 2
	miniPasswordLen = 7
)

type CreateUserParams struct{
	FirstName   string `json:"firstName"`
	LastName 	string `json:"lastName"`
	Email 		string `json:"email"`
	Password 	string `json:"password"`

}

type User struct {
	ID 				  primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName   	  string `bson:"firstName " json:"firstName"`
	LastName 		  string `bson:"lastName " json:"lastName"`
	Email 			  string `bson:"email" json:"email"`
	EncryptedPassword string `bson:"EncryptedPassword" json:"-"`

}

func NewUserFromParams(params CreateUserParams) (*User,error){
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil{
		return nil,err
	}
	return &User{FirstName: params.FirstName,
		LastName: params.LastName,
		Email: params.Email,
		EncryptedPassword: string(encpw),
	},nil
}
func (params CreateUserParams) Validate() []string{
	errors := []string{}
	if len(params.FirstName)<miniFirstNameLen{
		errors = append(errors,fmt.Sprintf("firstName length should be at least %d characters",miniFirstNameLen))
	}
	if len(params.LastName)<miniFirstNameLen{
		 errors = append(errors,fmt.Sprintf("LastNamelength should be at least %d characters",miniLastNameLen))
	}

	if len(params.Password)<miniPasswordLen{
		errors=append(errors,fmt.Sprintf("minimum password length should be at least %d characters",miniPasswordLen))
	}
	if !isEmailValid(params.Email){
		errors = append(errors,fmt.Sprintf("Email is invalid"))
	}
	return errors
}
func isEmailValid(e string) bool{
		var emailRegex = regexp.MustCompile(`(?i)^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
		return emailRegex.MatchString(e)
	}