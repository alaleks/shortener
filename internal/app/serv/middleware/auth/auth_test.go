package auth_test

import (
	"testing"

	"github.com/alaleks/shortener/internal/app/serv/middleware/auth"
)

func TestSigning(t *testing.T) {
	a := auth.TurnOn(nil, []byte("SECRET_KEY"))
	userID := 1
	sign := a.CreateSigning(uint(userID))
	userIDFromSign, _ := a.ReadSigning(sign)

	if userID != int(userIDFromSign) {
		t.Errorf("invalid userID %d. but should: %d", userIDFromSign, userID)
	}
}

func BenchmarkCreateSigning(b *testing.B) {
	b.StopTimer()
	a := auth.TurnOn(nil, []byte("SECRET_KEY"))
	userID := 1
	b.StartTimer()

	b.Run("before optimize", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = a.CreateSigningOld(uint(userID))
		}
	})

	b.Run("after optimize", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = a.CreateSigning(uint(userID))
		}
	})
}

func BenchmarkReadSigning(b *testing.B) {
	b.StopTimer()
	a := auth.TurnOn(nil, []byte("SECRET_KEY"))
	userID := 1
	sign := a.CreateSigningOld(uint(userID))
	sign2 := a.CreateSigning(uint(userID))
	b.StartTimer()

	b.Run("before optimize", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = a.ReadSigningOld(sign)
		}
	})

	b.Run("after optimize", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = a.ReadSigning(sign2)
		}
	})
}
