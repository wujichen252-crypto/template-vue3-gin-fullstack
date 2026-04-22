package jwt

import (
	"testing"
	"time"
)

func TestJWT_GenerateAndParseToken(t *testing.T) {
	secret := "test-secret-key-for-testing"
	jwtMgr := NewJWT(secret, time.Hour, time.Hour*24)

	userID := uint(123)

	token, err := jwtMgr.GenerateToken(userID)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	if token == "" {
		t.Fatal("GenerateToken returned empty token")
	}

	claims, err := jwtMgr.ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken failed: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("UserID mismatch: got %d, want %d", claims.UserID, userID)
	}
}

func TestJWT_GenerateRefreshToken(t *testing.T) {
	secret := "test-secret-key-for-testing"
	jwtMgr := NewJWT(secret, time.Hour, time.Hour*24)

	userID := uint(456)

	token, err := jwtMgr.GenerateRefreshToken(userID)
	if err != nil {
		t.Fatalf("GenerateRefreshToken failed: %v", err)
	}

	if token == "" {
		t.Fatal("GenerateRefreshToken returned empty token")
	}

	claims, err := jwtMgr.ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken failed: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("UserID mismatch: got %d, want %d", claims.UserID, userID)
	}
}

func TestJWT_InvalidToken(t *testing.T) {
	secret := "test-secret-key-for-testing"
	jwtMgr := NewJWT(secret, time.Hour, time.Hour*24)

	_, err := jwtMgr.ParseToken("invalid-token")
	if err == nil {
		t.Error("ParseToken should fail for invalid token")
	}
}

func TestJWT_WrongSecret(t *testing.T) {
	jwtMgr1 := NewJWT("secret-1", time.Hour, time.Hour*24)
	jwtMgr2 := NewJWT("secret-2", time.Hour, time.Hour*24)

	token, err := jwtMgr1.GenerateToken(123)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	_, err = jwtMgr2.ParseToken(token)
	if err == nil {
		t.Error("ParseToken should fail with wrong secret")
	}
}