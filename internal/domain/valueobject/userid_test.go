package valueobject

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewUserID(t *testing.T) {
	t.Run("should create a valid UserID from UUID string", func(t *testing.T) {
		// Arrange
		idStr := uuid.New().String()

		// Act
		userID, err := NewUserID(idStr)

		// Assert
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if userID.String() != idStr {
			t.Errorf("expected %s, got %s", idStr, userID.String())
		}
	})

	t.Run("should return error for invalid UUID string", func(t *testing.T) {
		// Arrange
		invalidID := "invalid-uuid"

		// Act
		_, err := NewUserID(invalidID)

		// Assert
		if err == nil {
			t.Error("expected error for invalid UUID, got nil")
		}
		if err.Error() != "invalid user ID format" {
			t.Errorf("expected 'invalid user ID format', got %v", err)
		}
	})

	t.Run("should create a valid UserID from UUID", func(t *testing.T) {
		// Arrange
		id := uuid.New()

		// Act
		userID := UserIDFromUUID(id)

		// Assert
		if userID.String() != id.String() {
			t.Errorf("expected %s, got %s", id.String(), userID.String())
		}
	})

	t.Run("should generate a new random UserID", func(t *testing.T) {
		// Act
		userID1 := NewUserIDRandom()
		userID2 := NewUserIDRandom()

		// Assert
		if userID1.String() == userID2.String() {
			t.Error("expected different UserIDs, got the same")
		}

		// Verify it's a valid UUID
		_, err := uuid.Parse(userID1.String())
		if err != nil {
			t.Errorf("expected valid UUID, got error: %v", err)
		}
	})

	t.Run("should compare UserIDs correctly", func(t *testing.T) {
		// Arrange
		id := uuid.New()
		userID1 := UserIDFromUUID(id)
		userID2 := UserIDFromUUID(id)
		userID3 := NewUserIDRandom()

		// Act & Assert
		if !userID1.Equals(userID2) {
			t.Error("expected UserIDs to be equal")
		}
		if userID1.Equals(userID3) {
			t.Error("expected UserIDs to be different")
		}
	})
}