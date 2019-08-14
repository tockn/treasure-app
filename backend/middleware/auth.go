package middleware

import (
	"errors"
	"log"
	"net/http"

	"github.com/voyagegroup/treasure-app/service"

	"firebase.google.com/go/auth"
	"github.com/voyagegroup/treasure-app/httputil"
	"github.com/voyagegroup/treasure-app/model"
)

const (
	bearer = "Bearer"
)

type Auth struct {
	authService service.Auth
	userService service.User
}

func NewAuth(as service.Auth, us service.User) *Auth {
	return &Auth{
		authService: as,
		userService: us,
	}
}

func (auth *Auth) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		idToken, err := getTokenFromHeader(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		token, err := auth.authService.VerifyIDToken(ctx, idToken)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Failed to verify token", http.StatusForbidden)
			return
		}
		userRecord, err := auth.authService.GetUser(ctx, token.UID)

		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Failed to get userRecord", http.StatusInternalServerError)
			return
		}

		firebaseUser := toFirebaseUser(userRecord)
		_, syncErr := auth.userService.Sync(ctx, &firebaseUser)
		if syncErr != nil {
			log.Print(syncErr.Error())
			http.Error(w, "Failed to sync user", http.StatusInternalServerError)
			return
		}

		user, err := auth.userService.Get(ctx, firebaseUser.FirebaseUID)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Failed to get user", http.StatusInternalServerError)
			return
		}

		ctx = httputil.SetUserToContext(ctx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getTokenFromHeader(req *http.Request) (string, error) {
	header := req.Header.Get("Authorization")
	if header == "" {
		return "", errors.New("authorization header not found")
	}

	l := len(bearer)
	if len(header) > l+1 && header[:l] == bearer {
		return header[l+1:], nil
	}

	return "", errors.New("authorization header format must be 'Bearer {token}'")
}

func toFirebaseUser(u *auth.UserRecord) model.FirebaseUser {
	return model.FirebaseUser{
		FirebaseUID: u.UID,
		Email:       u.Email,
		PhotoURL:    u.PhotoURL,
		DisplayName: u.DisplayName,
	}
}
