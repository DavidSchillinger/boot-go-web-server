package auth

import (
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestCreateJWT(t *testing.T) {
	id := uuid.New()

	tok, err := CreateJWT(id, "AllYourBase", time.Minute*5)
	if err != nil {
		t.Errorf("expected err to be nil, got %v", err)
		t.Fail()
	}

	if _, err := ValidateJWT(tok, "AllYourBase"); err != nil {
		t.Errorf("expected err to be nil, got %v", err)
		t.Fail()
	}
}
