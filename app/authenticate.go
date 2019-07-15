package app

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v3"
)

// User represents user object for JWT authentication
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// JwtToken represents an authorization token
type JwtToken struct {
	Token string `json:"token"`
}

// Claims is an embedded type to add fields to JWT StandardClaims
type Claims struct {
	Username           string `json:"username"`
	Password           string `json:"password"`
	jwt.StandardClaims `json:"claims"`
}

var (
	secretKey = os.Getenv("BR_API_SECRET_KEY")
)

// GenerateToken validates API user creds and returns a JWT token string; endpoint - /api/v1/generateToken
func (s *Server) generateToken() http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user User
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&user)
		if err != nil {
			json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: err.Error()})
			return
		}
		if validateCredentials(user, s.Dbc.Db) {
			expirationTime := time.Now().Add(24 * time.Hour)
			// Create the JWT claims, which includes the username, password and expiration time
			claims := &Claims{
				Username: user.Username,
				Password: user.Password,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: expirationTime.Unix(), // JWT expiration time is expressed as unix milliseconds
				},
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
			tokenString, err := token.SignedString([]byte(secretKey))
			if err != nil {
				json.NewEncoder(w).Encode(Exception{Status: http.StatusInternalServerError, Message: "error signing token string"})
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name:    "token",
				Value:   tokenString,
				Expires: expirationTime,
			})
			json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
			return
		}
		json.NewEncoder(w).Encode(Exception{Status: http.StatusUnauthorized, Message: "error validating credentials"})
	})
}

func validateCredentials(user User, db *sql.DB) bool {
	// Execute user_role(user, pass) func in database and scan result into null.String; null == ""
	var role null.String
	err := db.QueryRow("SELECT basic_auth.user_role($1, $2)", user.Username, user.Password).Scan(&role)
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
