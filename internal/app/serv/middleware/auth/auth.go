package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/alaleks/shortener/internal/app/storage"
)

var ErrInvalidSign = errors.New("this signing is invalid")

const (
	cookieName     = "Authorization"
	lifeTimeCookie = 2592000
)

type Auth struct {
	store     *storage.Store
	secretKey []byte
}

func TurnOn(store *storage.Store, secretKey []byte) Auth {
	return Auth{store: store, secretKey: secretKey}
}

func (a *Auth) createSigning(uid uint) string {
	mac := hmac.New(sha256.New, a.secretKey)
	mac.Write([]byte(strconv.Itoa(int(uid))))
	signature := mac.Sum(nil)
	signature = append(signature, []byte(strconv.Itoa(int(uid)))...)

	return base64.URLEncoding.EncodeToString(signature)
}

func (a *Auth) readSigning(cookieVal string) (uint, error) {
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

func (a *Auth) Authorization(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		authCookie, err := getCookie(req, cookieName)
		var userID uint

		if authCookie == nil || err != nil {
			userID = a.store.Store.Create()

			setCookie(writer, req, a.createSigning(userID))
			req.URL.User = url.User(strconv.Itoa(int(userID)))
			handler.ServeHTTP(writer, req)

			return
		}

		userID, err = a.readSigning(authCookie.Value)
		if err != nil {
			userID = a.store.Store.Create()

			setCookie(writer, req, a.createSigning(userID))
			req.URL.User = url.User(strconv.Itoa(int(userID)))
			handler.ServeHTTP(writer, req)

			return
		}

		req.URL.User = url.User(strconv.Itoa(int(userID)))
		handler.ServeHTTP(writer, req)
	})
}

func setCookie(writer http.ResponseWriter, req *http.Request, sign string) {
	cookie := http.Cookie{
		Name:     "Authorization",
		Value:    sign,
		Path:     "/",
		MaxAge:   lifeTimeCookie,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	if req.TLS == nil {
		cookie.Secure = false
	}

	http.SetCookie(writer, &cookie)
}

func getCookie(req *http.Request, name string) (*http.Cookie, error) {
	cookie, err := req.Cookie(name)
	if err != nil {
		err = fmt.Errorf("error getting authorization cookie: %w", err)
	}

	return cookie, err
}
