package auth

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/enki/daemon/internal/config"
	"golang.org/x/oauth2"
)

// Provider wraps the OIDC provider and OAuth2 config
type Provider struct {
	oidcProvider *oidc.Provider
	oauth2Config *oauth2.Config
	verifier     *oidc.IDTokenVerifier
}

// NewProvider creates a new OIDC provider from configuration
func NewProvider(ctx context.Context, cfg *config.OIDCConfig) (*Provider, error) {
	provider, err := oidc.NewProvider(ctx, cfg.IssuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider: %w", err)
	}

	scopes := cfg.Scopes
	if len(scopes) == 0 {
		scopes = []string{oidc.ScopeOpenID, "profile", "email"}
	}

	oauth2Cfg := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       scopes,
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: cfg.ClientID,
	})

	return &Provider{
		oidcProvider: provider,
		oauth2Config: oauth2Cfg,
		verifier:     verifier,
	}, nil
}

// AuthCodeURL returns the URL to redirect users for authentication
func (p *Provider) AuthCodeURL(state string) string {
	return p.oauth2Config.AuthCodeURL(state)
}

// Exchange exchanges the authorization code for tokens
func (p *Provider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return p.oauth2Config.Exchange(ctx, code)
}

// VerifyIDToken verifies the ID token and returns the claims
func (p *Provider) VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("no id_token in token response")
	}

	return p.verifier.Verify(ctx, rawIDToken)
}

// UserInfo fetches user info from the OIDC provider
func (p *Provider) UserInfo(ctx context.Context, token *oauth2.Token) (*oidc.UserInfo, error) {
	return p.oidcProvider.UserInfo(ctx, p.oauth2Config.TokenSource(ctx, token))
}

// OAuth2Config returns the OAuth2 config for token operations
func (p *Provider) OAuth2Config() *oauth2.Config {
	return p.oauth2Config
}
