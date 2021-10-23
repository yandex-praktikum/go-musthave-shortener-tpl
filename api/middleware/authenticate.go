package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/service/auth"
)

const AuthCookieName = "USER-ID"

type AuthContextKeyType struct{}

// Authenticator middleware authenticates a request
// based on the signed cookie containing a user ID.
// In case authentication has failed, it signs up a new user.
func Authenticator(idService auth.IDService) func(http.Handler) http.Handler {
	ra := requestAuth{idService}

	return func(next http.Handler) http.Handler {
		serveHTTP := func(w http.ResponseWriter, r *http.Request) {
			userID := ra.extractUserID(r)
			if userID == nil {
				log.Printf("Signing up new user")

				var errSignUp error
				userID, errSignUp = ra.signUp(w)
				if errSignUp != nil {
					log.Printf("Cannot authenticate: %s", errSignUp.Error())
					http.Error(w, "Cannot authenticate", http.StatusUnauthorized)
					return
				}
			}

			ctxWithUserID := context.WithValue(r.Context(), AuthContextKeyType{}, *userID)

			next.ServeHTTP(w, r.WithContext(ctxWithUserID))
		}

		return http.HandlerFunc(serveHTTP)
	}
}

type requestAuth struct {
	IDService auth.IDService
}

func (a *requestAuth) extractUserID(r *http.Request) *int64 {
	cookie, errGetCookie := r.Cookie(AuthCookieName)
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
	hmac := parts[1]

	userID, errParseID := strconv.ParseInt(userIDStr, 10, 64)
	if errParseID != nil {
		log.Printf("Cannot parse user ID [%s]", userIDStr)
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

func (a *requestAuth) signUp(w http.ResponseWriter) (*int64, error) {
	user, errSignUp := a.IDService.SignUp()
	if errSignUp != nil {
		return nil, fmt.Errorf("cannot sign up: %w", errSignUp)
	}

	signedUserID, errSign := auth.SignUserID(*user)
	if errSign != nil {
		return nil, fmt.Errorf("cannot sign user id: %w", errSign)
	}

	v := fmt.Sprintf("%d|%s", signedUserID.ID, signedUserID.HMAC)
	c := http.Cookie{
		Name:  AuthCookieName,
		Value: v,
	}
	http.SetCookie(w, &c)

	return &signedUserID.ID, nil
}
