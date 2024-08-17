package auth

import (
	"context"
	_ "embed"
	"errors"
	"net/http"
	"time"

	"github.com/fchimpan/simple-server/clock"
	"github.com/fchimpan/simple-server/entity"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

//go:embed cert/secret.pem
var rawPrivateKey []byte

//go:embed cert/public.pem
var rawPublicKey []byte

type JWTer struct {
	PrivateKey, PublicKey jwk.Key
	Store                 Store
	Clocker               clock.Clocker
}

//go:generate go run github.com/matryer/moq -out moq_test.go . Store
type Store interface {
	Save(ctx context.Context, key string, userID entity.UserID) error
	Load(ctx context.Context, key string) (entity.UserID, error)
}

func NewJWTer(store Store, clocker clock.Clocker) (*JWTer, error) {
	privateKey, err := parse(rawPrivateKey)
	if err != nil {
		return nil, err
	}
	publicKey, err := parse(rawPublicKey)
	if err != nil {
		return nil, err
	}
	return &JWTer{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Store:      store,
		Clocker:    clocker,
	}, nil
}

func parse(rawKey []byte) (jwk.Key, error) {
	return jwk.ParseKey(rawKey, jwk.WithPEM(true))
}

const (
	RoleKey     = "role"
	UserNameKey = "username"
)

func (j *JWTer) GenerateToken(ctx context.Context, user entity.User) ([]byte, error) {
	tok, err := jwt.NewBuilder().
		JwtID(uuid.New().String()).
		Issuer(`github.com/fchimpan/simple-server`).
		Subject("access_token").
		IssuedAt(j.Clocker.Now()).
		Expiration(j.Clocker.Now().Add(60*time.Minute)).
		Claim(RoleKey, user.Role).
		Claim(UserNameKey, user.Name).
		Build()
	if err != nil {
		return nil, errors.New("failed to build token: " + err.Error())
	}
	if err := j.Store.Save(ctx, tok.JwtID(), user.ID); err != nil {
		return nil, errors.New("failed to save token: " + err.Error())
	}

	signed, err := jwt.Sign(tok, jwt.WithKey(jwa.RS256, j.PrivateKey))
	if err != nil {
		return nil, errors.New("failed to sign token: " + err.Error())
	}

	return signed, nil
}

func (j *JWTer) GetToken(ctx context.Context, r *http.Request) (jwt.Token, error) {
	token, err := jwt.ParseRequest(r,
		jwt.WithKey(jwa.RS256, j.PublicKey),
		jwt.WithValidate(false),
	)
	if err != nil {
		return nil, err
	}
	if err := jwt.Validate(token, jwt.WithClock(j.Clocker)); err != nil {
		return nil, err
	}
	if _, err := j.Store.Load(ctx, token.JwtID()); err != nil {
		return nil, err
	}
	return token, nil
}

func (j *JWTer) FillContext(r *http.Request) (*http.Request, error) {
	token, err := j.GetToken(r.Context(), r)
	if err != nil {
		return nil, err
	}
	uid, err := j.Store.Load(r.Context(), token.JwtID())
	if err != nil {
		return nil, err
	}
	ctx := SetUserID(r.Context(), uid)

	ctx = SetRole(ctx, token)
	clone := r.Clone(ctx)
	return clone, nil
}

type userIDKey struct{}
type roleKey struct{}

func SetUserID(ctx context.Context, userID entity.UserID) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}

func GetUserID(ctx context.Context) (entity.UserID, bool) {
	v, ok := ctx.Value(userIDKey{}).(entity.UserID)
	return v, ok
}

func SetRole(ctx context.Context, tok jwt.Token) context.Context {
	get, ok := tok.Get(RoleKey)
	if !ok {
		return context.WithValue(ctx, roleKey{}, "")
	}
	return context.WithValue(ctx, roleKey{}, get)
}

func GetRole(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(roleKey{}).(string)
	return role, ok
}

func IsAdmin(ctx context.Context) bool {
	role, ok := GetRole(ctx)
	if !ok {
		return false
	}
	return role == "admin"
}
