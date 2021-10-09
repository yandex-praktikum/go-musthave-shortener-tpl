package middleware

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/service/auth"
)

const COOKIE_NAME = "USER-ID"

type Authenticator struct {
	IDService auth.IDService
}

func (a *Authenticator) Authenticate(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		cookie, errGetCookie := r.Cookie(COOKIE_NAME)
		if errGetCookie != nil && errGetCookie != http.ErrNoCookie {
			log.Printf("Cannot get authentication cookie: %s", errGetCookie.Error())
			http.Error(w, "Cannot get authentication cookie", http.StatusUnauthorized)
			return
		}

		if errGetCookie == http.ErrNoCookie {
			a.SignUp(w)
		} else if errValidate := a.Validate(cookie.Value); errValidate != nil {
			a.SignUp(w)
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (a *Authenticator) SignUp(w http.ResponseWriter) {
	user, errSignUp := a.IDService.SignUp()
	if errSignUp != nil {
		log.Printf("Cannot sign up: %s", errSignUp.Error())
		http.Error(w, "Cannot sign up", http.StatusInternalServerError)
		return
	}

	userID := auth.SignUserID(*user)
	v := fmt.Sprintf("%d|%x", userID.ID, userID.HMAC)
	c := http.Cookie{
		Name:  COOKIE_NAME,
		Value: v,
	}
	http.SetCookie(w, &c)
	fmt.Print(userID)
}

func (a *Authenticator) Validate(v string) error {
	parts := strings.Split(v, "|")
	if len(parts) != 2 {
		msg := fmt.Sprintf("cannot parse signed user ID [%s]", v)
		return errors.New(msg)
	}

	userIDStr := parts[0]
	hmacStr := parts[1]

	userID, errParseID := strconv.Atoi(userIDStr)
	if errParseID != nil {
		msg := fmt.Sprintf("cannot parse user ID [%s]", userIDStr)
		return errors.New(msg)
	}

	hmac, errParseHmac := hex.DecodeString(hmacStr)
	if errParseHmac != nil {
		msg := fmt.Sprintf("cannot parse signature [%s]", hmacStr)
		return errors.New(msg)
	}

	sgn := model.SignedUserID{
		ID:   userID,
		HMAC: hmac,
	}

	return a.IDService.Validate(sgn)
}
