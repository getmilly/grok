package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	auth0 "github.com/auth0-community/go-auth0"
	jose "gopkg.in/square/go-jose.v2"
)

//AuthService ...
type AuthService interface {
	Authorize(req *http.Request) (Claims, error)
}

type auth struct {
	JwksURI  string
	Issuer   string
	Audience []string
}

//Claims wraps user data
type Claims map[string]interface{}

var (
	//ErrClaimNotFound is returned when any claim is found.
	ErrClaimNotFound = errors.New("Claim not found")
)

//GetKey return the claim value if exists
func (claims Claims) GetKey(key string) (interface{}, error) {
	value, ok := claims[key]

	if ok {
		return value, nil
	}

	return nil, ErrClaimNotFound
}

//NewAuthService creates a new pointer to an auth handler.
func NewAuthService(jwks, issuer string, audience []string) AuthService {
	return &auth{
		JwksURI:  jwks,
		Issuer:   issuer,
		Audience: audience,
	}
}

func (service auth) Authorize(req *http.Request) (Claims, error) {
	client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: service.JwksURI}, nil)
	configuration := auth0.NewConfiguration(client, service.Audience, service.Issuer, jose.RS256)
	validator := auth0.NewValidator(configuration, nil)

	jwt, err := validator.ValidateRequest(req)

	if err != nil {
		return nil, err
	}

	var claims map[string]interface{}
	validator.Claims(req, jwt, &claims)

	return claims, nil
}

//Authentication authenticates a request against an auth service.
func Authentication(service AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := service.Authorize(c.Request)

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		for k, v := range claims {
			c.Set(strings.ToLower(k), v)
		}

		c.Next()
	}
}
