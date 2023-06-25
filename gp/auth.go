package gp

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/cristalhq/jwt/v3"
	"github.com/google/uuid"
)

type JwtClaims struct {
	jwt.RegisteredClaims
	Authenticated      bool `json:"authenticated"`
	TwoFaAuthenticated bool `json:"two_fa_authenticated"`
}

func CreateToken(sUser string, tokenDuration time.Duration, isAuthenticate, is2FAAuthenticated bool) (string, error) {

	// create a Signer (HMAC in this example)
	key := []byte(PConfig.JWTSecret)
	signer, err := jwt.NewSignerHS(jwt.HS256, key)
	if err != nil {
		return "", err
	}

	// create claims
	claims := JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.Must(uuid.NewRandom()).String(),
			Audience:  jwt.Audience{PConfig.ProjectName},
			Issuer:    "Makyo",
			Subject:   sUser,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenDuration).UTC()),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Duration(time.Second * -1)).UTC()),
		},
		Authenticated:      isAuthenticate,
		TwoFaAuthenticated: is2FAAuthenticated,
	}

	// create a Builder
	builder := jwt.NewBuilder(signer)

	// and build a Token
	token, err := builder.Build(claims)
	if err != nil {
		return "", err
	}

	// here is token as byte slice
	return token.String(), nil
}

// VerifyToken verifies a token
func VerifyToken(tokenString string, requireAuthentication, require2FA bool) (string, error) {
	// create a Verifier
	key := []byte(PConfig.JWTSecret)
	verifier, err := jwt.NewVerifierHS(jwt.HS256, key)
	if err != nil {
		return "", err
	}

	// parse a Token
	token, err := jwt.ParseString(tokenString)
	if err != nil {
		return "", err
	}

	// and verify it's signature
	err = verifier.Verify(token.Payload(), token.Signature())
	if err != nil {
		return "", err
	}

	// get standard claims
	var newClaims JwtClaims
	err = json.Unmarshal(token.RawClaims(), &newClaims)
	if err != nil {
		return "", err
	}

	// verify claims
	// We check Audience
	if !newClaims.IsForAudience(PConfig.ProjectName) {
		return "", errors.New("the token can't be used for this application")
	}

	// We check exp
	if !newClaims.IsValidAt(time.Now().UTC()) {
		return "", errors.New("the token is not valid anymore")
	}

	if requireAuthentication && !newClaims.Authenticated {
		return "", errors.New("the token is not valid, email auth is required")
	}

	if require2FA && !newClaims.TwoFaAuthenticated {
		return "", errors.New("the token is not valid, 2fa auth is required")
	}

	return newClaims.Subject, nil
}

func GetTokenClaims(r *http.Request) (*jwt.RegisteredClaims, error) {
	// We check the user token
	userToken := r.Header.Get("user-token")

	// parse a Token
	token, err := jwt.ParseString(userToken)
	if err != nil {
		return nil, err
	}

	// get standard claims
	var newClaims jwt.RegisteredClaims
	errClaims := json.Unmarshal(token.RawClaims(), &newClaims)
	if errClaims != nil {
		return nil, err
	}

	return &newClaims, nil
}

// APISiteAuth handles authentication for site-data endpoints
func APISiteAuth(r *http.Request) (string, error) {

	if PConfig.DebugMode {
		return "", nil
	}

	// We check the user token
	userToken := r.Header.Get("user-token")
	return VerifyToken(userToken, true, true)
}

// APISiteAuthFor2FA handles authentication for site-data join/2fa endpoint
func APISiteAuthFor2FA(r *http.Request) (string, error) {

	if PConfig.DebugMode {
		return "", nil
	}

	// We check the user token
	userToken := r.Header.Get("user-token")
	return VerifyToken(userToken, true, false)
}

// APISiteAuthNotAuth handles authentication for site-data not authenticated endpoints
func APISiteAuthNotAuth(r *http.Request) (string, error) {

	if PConfig.DebugMode {
		return "", nil
	}

	// We check the user token
	userToken := r.Header.Get("user-token")
	return VerifyToken(userToken, false, false)
}
