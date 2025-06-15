package notifystock

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type GoogleToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

type GoogleClient struct {
	Client       *http.Client
	clientSecret string
	clientID     string
	redirectURI  string
}

type GoogleClientOption struct {
	ClientID    string
	Secret      string
	RedirectURI string
}

func NewGoogleClient(client http.Client, option GoogleClientOption) *GoogleClient {
	return &GoogleClient{
		Client:       &client,
		clientSecret: option.Secret,
		clientID:     option.ClientID,
		redirectURI:  option.RedirectURI,
	}
}

func (g *GoogleClient) ExchangeToken(ctx context.Context, code string) (*GoogleToken, error) {
	values := url.Values{}
	values.Add("code", code)
	values.Add("client_id", g.clientID)
	values.Add("client_secret", g.clientSecret)
	values.Add("redirect_uri", g.redirectURI)
	values.Add("grant_type", "authorization_code")

	u, err := url.Parse("https://oauth2.googleapis.com/token")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := g.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to exchange token: %s", res.Status)
	}
	defer res.Body.Close()
	var token GoogleToken
	if err := json.NewDecoder(res.Body).Decode(&token); err != nil {
		return nil, err
	}

	return &token, nil
}

func (g *GoogleClient) GetUserInfo(ctx context.Context, token *GoogleToken) (*GoogleUserInfo, error) {
	req, err := http.NewRequestWithContext(
		ctx, http.MethodGet, "https://www.googleapis.com/oauth2/v1/userinfo", nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", token.TokenType+" "+token.AccessToken)

	res, err := g.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: %s", res.Status)
	}
	defer res.Body.Close()

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(res.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}
