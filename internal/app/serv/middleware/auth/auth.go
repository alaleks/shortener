// The auth package implements Cookie-based authorization.
package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/alaleks/shortener/internal/app/storage"
)

// ErrInvalidSign is an indicator that the signature is invalid.
var ErrInvalidSign = errors.New("this signing is invalid")

const (
	cookieName     = "Authorization"
	lifeTimeCookie = 2592000
)

// Auth stores storing a link to storage and a secret key as an array of bytes.
type Auth struct {
	store     *storage.Store
	secretKey []byte
}

// TurnOn enables on-site authorization.
func TurnOn(store *storage.Store, secretKey []byte) Auth {
	return Auth{store: store, secretKey: secretKey}
}

// CreateSigningOld (Deprecated) creates a signature for the cookie.
func (a *Auth) CreateSigningOld(uid uint) string {
	mac := hmac.New(sha256.New, a.secretKey)
	mac.Write([]byte(strconv.Itoa(int(uid))))
	signature := mac.Sum(nil)
	signature = append(signature, []byte(strconv.Itoa(int(uid)))...)

	return base64.URLEncoding.EncodeToString(signature)
}

// CreateSigning creates a signature for the cookie.
//
// Encryption is performed according to the sha256 algorithm.
// The secret key and user ID are used for encryption.
func (a *Auth) CreateSigning(uid uint) string {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(uid))

	mac := hmac.New(sha256.New, a.secretKey)
	mac.Write(b)
	signature := mac.Sum(nil)
	signature = append(signature, b...)

	return base64.URLEncoding.EncodeToString(signature)
}

// ReadSigningOld (Deprecated) decrypts cookie value and returns user ID and error value.
func (a *Auth) ReadSigningOld(cookieVal string) (uint, error) {
	signedVal, err := base64.URLEncoding.DecodeString(cookieVal)
	if err != nil {
		return 0, fmt.Errorf("cookie decoding error: %w", err)
	}

	signature := string(signedVal)[:sha256.Size]
	uid, err := strconv.Atoi(string(signedVal)[sha256.Size:])
	if err != nil {
		return 0, fmt.Errorf("UID conversion error: %w", err)
	}

	mac := hmac.New(sha256.New, a.secretKey)
	mac.Write([]byte(strconv.Itoa(uid)))
	expectedSignature := mac.Sum(nil)

	if !hmac.Equal([]byte(signature), expectedSignature) {
		return 0, ErrInvalidSign
	}

	return uint(uid), nil
}

// ReadSigning decrypts cookie value and returns user ID and error value.
func (a *Auth) ReadSigning(cookieVal string) (uint, error) {
	signedVal, err := base64.URLEncoding.DecodeString(cookieVal)
	if err != nil {
		return 0, fmt.Errorf("cookie decoding error: %w", err)
	}

	signature := signedVal[:sha256.Size]
	mac := hmac.New(sha256.New, a.secretKey)
	mac.Write(signedVal[sha256.Size:])
	expectedSignature := mac.Sum(nil)

	if !hmac.Equal(signature, expectedSignature) {
		return 0, ErrInvalidSign
	}

	return uint(int64(binary.LittleEndian.Uint64(signedVal[sha256.Size:]))), nil
}

// Authorization method middleware, which performs an authorization check
// or creates a new user if the authorization cookie value is empty or invalid.
func (a *Auth) Authorization(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		authCookie, err := getCookie(req, cookieName)
		var userID uint

		if authCookie == nil || err != nil {
			userID = a.store.St.Create()
			http.SetCookie(writer, setCookie(a.CreateSigning(userID), req.TLS != nil))
			req.URL.User = url.User(strconv.Itoa(int(userID)))
			handler.ServeHTTP(writer, req)

			return
		}

		userID, err = a.ReadSigning(authCookie.Value)
		if err != nil {
			userID = a.store.St.Create()
			http.SetCookie(writer, setCookie(a.CreateSigning(userID), req.TLS != nil))
			req.URL.User = url.User(strconv.Itoa(int(userID)))
			handler.ServeHTTP(writer, req)

			return
		}

		req.URL.User = url.User(strconv.Itoa(int(userID)))
		handler.ServeHTTP(writer, req)
	})
}

func setCookie(sign string, ssl bool) *http.Cookie {
	cookie := http.Cookie{
		Name:     cookieName,
		Value:    sign,
		Path:     "/",
		MaxAge:   lifeTimeCookie,
		HttpOnly: true,
		Secure:   ssl,
		SameSite: http.SameSiteLaxMode,
	}

	return &cookie
}

func getCookie(req *http.Request, name string) (*http.Cookie, error) {
	cookie, err := req.Cookie(name)
	if err != nil {
		return cookie, fmt.Errorf("error getting authorization cookie: %w", err)
	}

	return cookie, nil
}
