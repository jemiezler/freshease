package authoidc

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/oauth2"

	"freshease/backend/ent"
	"freshease/backend/ent/identity"
	"freshease/backend/ent/user"
)

type ProviderName string

const (
	ProviderGoogle ProviderName = "google"
	ProviderLINE   ProviderName = "line"
)

type providerClient struct {
	Name     ProviderName
	Provider *oidc.Provider
	Verifier *oidc.IDTokenVerifier
	Config   *oauth2.Config
}

type Service struct {
	db        *ent.Client
	jwtSecret []byte
	ttl       time.Duration
	clients   map[ProviderName]*providerClient
	baseURL   string
}

func NewService(ctx context.Context, db *ent.Client) (*Service, error) {
	base := mustEnv("OAUTH_BASE_URL")
	ttlMin := getenvIntDefault("JWT_ACCESS_TTL_MIN", 15)

	s := &Service{
		db:        db,
		jwtSecret: []byte(mustEnv("JWT_SECRET")),
		ttl:       time.Duration(ttlMin) * time.Minute,
		clients:   map[ProviderName]*providerClient{},
		baseURL:   base,
	}

	// Google
	gc, err := newClient(ctx,
		ProviderGoogle,
		mustEnv("OIDC_GOOGLE_ISSUER"),
		mustEnv("OIDC_GOOGLE_CLIENT_ID"),
		mustEnv("OIDC_GOOGLE_CLIENT_SECRET"),
		base+mustEnv("OIDC_GOOGLE_REDIRECT_PATH"),
	)
	if err != nil {
		return nil, fmt.Errorf("google oidc init: %w", err)
	}
	s.clients[ProviderGoogle] = gc

	// LINE
	lc, err := newClient(ctx,
		ProviderLINE,
		mustEnv("OIDC_LINE_ISSUER"),
		mustEnv("OIDC_LINE_CLIENT_ID"),
		mustEnv("OIDC_LINE_CLIENT_SECRET"),
		base+mustEnv("OIDC_LINE_REDIRECT_PATH"),
	)
	if err != nil {
		return nil, fmt.Errorf("line oidc init: %w", err)
	}
	s.clients[ProviderLINE] = lc

	return s, nil
}

func newClient(ctx context.Context, name ProviderName, issuer, id, secret, redirect string) (*providerClient, error) {
	prov, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, err
	}
	conf := &oauth2.Config{
		ClientID:     id,
		ClientSecret: secret,
		Endpoint:     prov.Endpoint(),
		RedirectURL:  redirect,
		Scopes:       []string{oidc.ScopeOpenID, "email", "profile"},
	}
	verifier := prov.Verifier(&oidc.Config{ClientID: id})
	return &providerClient{Name: name, Provider: prov, Verifier: verifier, Config: conf}, nil
}

func (s *Service) AuthCodeURL(p ProviderName, state, nonce, codeChallenge string) (string, error) {
	c, ok := s.clients[p]
	if !ok {
		return "", errors.New("unknown provider")
	}
	opts := []oauth2.AuthCodeOption{oauth2.SetAuthURLParam("nonce", nonce)}
	if codeChallenge != "" {
		opts = append(opts,
			oauth2.SetAuthURLParam("code_challenge", codeChallenge),
			oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		)
	}
	return c.Config.AuthCodeURL(state, opts...), nil
}

type oidcClaims struct {
	Sub     string `json:"sub"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func (s *Service) ExchangeAndLogin(ctx context.Context, p ProviderName, code, codeVerifier string) (string, error) {
	c, ok := s.clients[p]
	if !ok {
		return "", errors.New("unknown provider")
	}

	var tok *oauth2.Token
	var err error
	if codeVerifier != "" {
		tok, err = c.Config.Exchange(ctx, code, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	} else {
		tok, err = c.Config.Exchange(ctx, code)
	}
	if err != nil {
		return "", err
	}

	rawID, ok2 := tok.Extra("id_token").(string)
	if !ok2 {
		return "", errors.New("missing id_token")
	}

	idTok, err := c.Verifier.Verify(ctx, rawID)
	if err != nil {
		return "", err
	}

	var cl oidcClaims
	if err := idTok.Claims(&cl); err != nil {
		return "", err
	}

	uid, email, err := s.upsertIdentity(ctx, string(p), cl.Sub, cl.Email, cl.Name, cl.Picture, tok)
	if err != nil {
		return "", err
	}

	return s.issueJWT(uid, email)
}

func (s *Service) upsertIdentity(ctx context.Context, provider, sub, email, name, avatar string, tok *oauth2.Token) (uuid.UUID, string, error) {
	idn, err := s.db.Identity.Query().Where(identity.Provider(provider), identity.Subject(sub)).First(ctx)
	var u *ent.User

	if err == nil {
		u, err = s.db.User.Get(ctx, idn.UserID)
		if err != nil {
			return uuid.Nil, "", err
		}
		up := s.db.Identity.UpdateOneID(idn.ID).
			SetUpdatedAt(time.Now()).
			SetEmail(email).SetName(name).SetAvatar(avatar)
		if tok != nil {
			if at := tok.AccessToken; at != "" {
				up.SetAccessToken(at)
			}
			if rt := tok.RefreshToken; rt != "" {
				up.SetRefreshToken(rt)
			}
			if tok.Expiry.Unix() > 0 {
				up.SetExpiresAt(tok.Expiry)
			}
		}
		if _, err = up.Save(ctx); err != nil {
			return uuid.Nil, "", err
		}
	} else {
		// link by email if exists
		if email != "" {
			u, _ = s.db.User.Query().Where(user.Email(email)).First(ctx)
		}
		if u == nil {
			u, err = s.db.User.Create().
				SetEmail(email).SetName(name).SetAvatar(avatar).
				Save(ctx)
			if err != nil {
				return uuid.Nil, "", err
			}
		}
		cr := s.db.Identity.Create().
			SetUserID(u.ID).
			SetProvider(provider).
			SetSubject(sub).
			SetEmail(email).SetName(name).SetAvatar(avatar)
		if tok != nil {
			if at := tok.AccessToken; at != "" {
				cr.SetAccessToken(at)
			}
			if rt := tok.RefreshToken; rt != "" {
				cr.SetRefreshToken(rt)
			}
			if tok.Expiry.Unix() > 0 {
				cr.SetExpiresAt(tok.Expiry)
			}
		}
		if _, err = cr.Save(ctx); err != nil {
			return uuid.Nil, "", err
		}
	}
	return u.ID, u.Email, nil
}

func (s *Service) issueJWT(uid uuid.UUID, email string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   uid.String(),
		"email": email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(s.ttl).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(s.jwtSecret)
}

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		panic("missing env: " + k)
	}
	return v
}
func getenvIntDefault(k string, def int) int {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return i
}
