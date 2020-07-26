package app

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v3"
)

// JwtToken represents an authorization token
type JwtToken struct {
	Token string
}

// Claims is an embedded type to add fields to JWT StandardClaims
type Claims struct {
	Username           string
	jwt.StandardClaims
}

// User represents user object for JWT authentication
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var (
	secretKey = os.Getenv("SBD_API_SECRET_KEY")
)

// Authenticate is a middleware that wraps an http.Handler and checks/validates a Bearer Token cookie or header
func (s *Server) Authenticate(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var bearerToken string
		bearerToken, err := checkBearerTokenCookie(r) // check for bearer token cookie
		switch err {
		case nil: // if err is nil do nothing
		case http.ErrNoCookie: // if no cookie is present check for the token header
			bearerToken, err = checkBearerTokenHeader(r)
			if ok := checkWriteError(err, http.StatusUnauthorized, w); ok {
				return
			}
		default:
			if ok := checkWriteError(errors.Wrap(err, "no authorization token found"), http.StatusUnauthorized, w); ok {
				return
			}
		}

		ok := authenticateToken(bearerToken) // if token is found then attempt to authenticate
		if ok {
			next.ServeHTTP(w, r)
			return
		}
		json.NewEncoder(w).Encode(Exception{Status: http.StatusUnauthorized, Message: "invalid authorization token"})
		return
	})
}

func (s *Server) validateCredentials(user User) bool {
	// execute user_role(user, pass) func in database and scan result into null.String; null == ""
	var role null.String
	err := s.Dbc.Db.Get(&role, "SELECT basic_auth.user_role($1, $2)", user.Username, user.Password)
	if err != nil {
		log.Println(errors.Wrap(err, "error querying users table"))
		return false
	}
	// if value is "" then user credentials do not exist in database; reject
	if role.NullString.String != "" {
		return true
	}
	return false
}

func authenticateToken(bearerToken string) bool {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(bearerToken, claims, jwtParseKeyFunc)
	if err != nil {
		log.Println(errors.Wrap(err, "error parsing bearer token"))
		return false
	}
	if token.Valid {
		return true
	}
	return false
}

func checkBearerTokenCookie(r *http.Request) (string, error) {
	c, err := r.Cookie("token")
	if err != nil {
		return "", err
	}
	return c.Value, nil
}

func checkBearerTokenHeader(r *http.Request) (string, error) {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return "", errors.New("no authorization header")
	}
	bearerTokenSlice := strings.Split(authorizationHeader, " ")
	if len(bearerTokenSlice) != 2 {
		return "", errors.New("empty authorization token")
	}
	return bearerTokenSlice[1], nil
}

func jwtParseKeyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New("error with token signing method")
	}
	return []byte(secretKey), nil
}
