package middleware

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/service/auth"
)

const AUTH_COOKIE_NAME = "USER-ID"

type AUTH_CONTEXT_KEY_TYPE struct{}

type Authenticator struct {
	IDService auth.IDService
}

func (a *Authenticator) Authenticate(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		userID := a.ExtractUserID(r)
		if userID == nil {
			var errSignUp error
			userID, errSignUp = a.SignUp(w)
			if errSignUp != nil {
				log.Printf("Cannot authenticate: %s", errSignUp.Error())
				http.Error(w, "Cannot authenticate", http.StatusUnauthorized)
				return
			}
		}

		ctxWithUserID := context.WithValue(r.Context(), AUTH_CONTEXT_KEY_TYPE{}, *userID)

		next.ServeHTTP(w, r.WithContext(ctxWithUserID))
	}

	return http.HandlerFunc(fn)
}

func (a *Authenticator) ExtractUserID(r *http.Request) *int {
	cookie, errGetCookie := r.Cookie(AUTH_COOKIE_NAME)
	if errGetCookie != nil {
		log.Printf("Cannot get authentication cookie: %s", errGetCookie.Error())
		return nil
	}

	parts := strings.Split(cookie.Value, "|")
	if len(parts) != 2 {
		log.Printf("Cannot parse signed user ID [%s]", cookie.Value)
		return nil
	}

	userIDStr := parts[0]
	hmacStr := parts[1]

	userID, errParseID := strconv.Atoi(userIDStr)
	if errParseID != nil {
		log.Printf("Cannot parse user ID [%s]", userIDStr)
		return nil
	}

	hmac, errParseHmac := hex.DecodeString(hmacStr)
	if errParseHmac != nil {
		log.Printf("Cannot parse signature [%s]", hmacStr)
		return nil
	}

	sgn := model.SignedUserID{
		ID:   userID,
		HMAC: hmac,
	}

	if invalid := a.IDService.Validate(sgn); invalid != nil {
		log.Printf("Signature is invalid: %s", invalid.Error())
		return nil
	}

	return &sgn.ID
}

func (a *Authenticator) SignUp(w http.ResponseWriter) (*int, error) {
	user, errSignUp := a.IDService.SignUp()
	if errSignUp != nil {
		return nil, fmt.Errorf("cannot sign up: %w", errSignUp)
	}

	signedUserID := auth.SignUserID(*user)
	v := fmt.Sprintf("%d|%x", signedUserID.ID, signedUserID.HMAC)
	c := http.Cookie{
		Name:  AUTH_COOKIE_NAME,
		Value: v,
	}
	http.SetCookie(w, &c)

	return &signedUserID.ID, nil
}
