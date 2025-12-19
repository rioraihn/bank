package valueobject

import (
	"errors"
	"github.com/google/uuid"
)

type UserID struct {
	value uuid.UUID
}

func NewUserID(idStr string) (UserID, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return UserID{}, errors.New("invalid user ID format")
	}

	return UserID{value: id}, nil
}

func UserIDFromUUID(id uuid.UUID) UserID {
	return UserID{value: id}
}

func NewUserIDRandom() UserID {
	return UserID{value: uuid.New()}
}

func (u UserID) String() string {
	return u.value.String()
}

func (u UserID) Equals(other UserID) bool {
	return u.value == other.value
}
