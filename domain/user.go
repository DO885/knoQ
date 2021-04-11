package domain

import "github.com/gofrs/uuid"

type User struct {
	ID          uuid.UUID
	Name        string
	DisplayName string
	Icon        string
	Privileged  bool
	State       int
}

type UserRepository interface {
	// LoginUser OAuthによってユーザーを得る
	LoginUser(query, state, codeVerifier string) error
	GetUser(userID uuid.UUID, info *ConInfo) (*User, error)
	GetAllUsers(*ConInfo) ([]*User, error)
	//ReplaceToken(userID uuid.UUID, token string, info *ConInfo) error
	//GetToken(info *ConInfo) (string, error)
	ReNewMyiCalSecret(info *ConInfo) (string, error)
	GetMyiCalSecret(info *ConInfo) (string, error)

	IsPrevilege(info *ConInfo) bool
}
