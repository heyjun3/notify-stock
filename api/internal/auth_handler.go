package notifystock

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const googleOauthURL = "https://accounts.google.com/o/oauth2/v2/auth"
const googleAccessTokenURL = "https://oauth2.googleapis.com/token"
const googleUserInfoURL = "https://www.googleapis.com/oauth2/v1/userinfo"

var (
	scope = []string{
		"openid",
		"https://www.googleapis.com/auth/userinfo.profile",
		"https://www.googleapis.com/auth/userinfo.email",
	}
)

type Token struct {
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

type AuthHandler struct {
	Sessions *Sessions
}

func NewAuthHandler(sessions *Sessions) *AuthHandler {
	return &AuthHandler{
		Sessions: sessions,
	}
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// 既存セッションの確認
	session, err := h.Sessions.Get(r)
	if err == nil && session.IsActive {
		http.Redirect(w, r, Cfg.FrontendURL, http.StatusFound)
		return // already logged in
	}

	// 新しいセッションの作成
	newSession, err := h.Sessions.New(w, r.Context())
	if err != nil {
		WriteErrorResponse(w, WrapError(err, ErrCodeSession, "Failed to create session"))
		return
	}

	// OAuth URLの構築
	u, err := url.Parse(googleOauthURL)
	if err != nil {
		WriteErrorResponse(w, WrapError(err, ErrCodeInternalServer, "Failed to parse OAuth URL"))
		return
	}

	q := url.Values{}
	q.Add("scope", strings.Join(scope, " "))
	q.Add("access_type", "offline")
	q.Add("include_granted_scopes", "true")
	q.Add("response_type", "code")
	q.Add("state", newSession.State)
	q.Add("redirect_uri", Cfg.OauthRedirectURL)
	q.Add("client_id", Cfg.OauthClientID)

	u.RawQuery = q.Encode()
	logger.Info("OAuth redirect URL generated", "url", u.String())
	http.Redirect(w, r, u.String(), http.StatusFound)
}

func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	err := h.Sessions.Clear(w, r)
	if err != nil {
		logger.Info("No active session found for logout", "error", err)
		http.Redirect(w, r, Cfg.FrontendURL, http.StatusFound)
		return
	}

	logger.Info("User logged out successfully")
	http.Redirect(w, r, Cfg.FrontendURL, http.StatusFound)
}

func (h *AuthHandler) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	// セッションの取得
	session, err := h.Sessions.Get(r)
	if err != nil {
		WriteErrorResponse(w, WrapError(err, ErrCodeSession, "Failed to get session"))
		return
	}

	// OAuth認証コードの確認
	code := r.URL.Query().Get("code")
	if code == "" {
		WriteErrorResponse(w, NewValidationError("Authorization code is missing", "code parameter is required"))
		return
	}

	// CSRF対策: stateパラメータの確認
	state := r.URL.Query().Get("state")
	if state != session.State {
		WriteErrorResponse(w, NewUnauthorizedError("Invalid state parameter"))
		return
	}

	// OAuth トークン交換リクエストの準備
	tokenRequest := url.Values{}
	tokenRequest.Add("code", code)
	tokenRequest.Add("client_id", Cfg.OauthClientID)
	tokenRequest.Add("client_secret", Cfg.OauthClientSecret)
	tokenRequest.Add("redirect_uri", Cfg.OauthRedirectURL)
	tokenRequest.Add("grant_type", "authorization_code")

	logger.Debug("Exchanging OAuth code for token")

	// OAuth トークンエンドポイントの呼び出し
	tokenURL, err := url.Parse(googleAccessTokenURL)
	if err != nil {
		WriteErrorResponse(w, WrapError(err, ErrCodeInternalServer, "Failed to parse token URL"))
		return
	}

	req, err := http.NewRequestWithContext(r.Context(),
		http.MethodPost, tokenURL.String(),
		strings.NewReader(tokenRequest.Encode()))
	if err != nil {
		WriteErrorResponse(w, WrapError(err, ErrCodeInternalServer, "Failed to create token request"))
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		WriteErrorResponse(w, WrapError(err, ErrCodeExternalService, "Failed to exchange OAuth code"))
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		WriteErrorResponse(w, NewAppError(ErrCodeExternalService,
			fmt.Sprintf("OAuth token exchange failed with status %d", res.StatusCode), nil))
		return
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		WriteErrorResponse(w, WrapError(err, ErrCodeExternalService, "Failed to read token response"))
		return
	}
	var token Token
	if err := json.Unmarshal(body, &token); err != nil {
		WriteErrorResponse(w, WrapError(err, ErrCodeExternalService, "Failed to parse token response"))
		return
	}

	req, err = http.NewRequest(http.MethodGet, googleUserInfoURL, nil)
	if err != nil {
		WriteErrorResponse(w, WrapError(err, ErrCodeInternalServer, "Failed to create user info request"))
		return
	}
	req.Header.Set("Authorization", token.TokenType+" "+token.AccessToken)
	res, err = client.Do(req)
	if err != nil {
		WriteErrorResponse(w, WrapError(err, ErrCodeExternalService, "Failed to fetch user info"))
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		WriteErrorResponse(w, NewAppError(ErrCodeExternalService,
			fmt.Sprintf("Failed to fetch user info with status %d", res.StatusCode), nil))
		return
	}
	body, err = io.ReadAll(res.Body)
	if err != nil {
		WriteErrorResponse(w, WrapError(err, ErrCodeExternalService, "Failed to read user info response"))
		return
	}

	logger.Info("User info response received", "response", string(body))

	// セッションの有効化
	session.IsActive = true
	err = h.Sessions.Store(r.Context(), session)
	if err != nil {
		WriteErrorResponse(w, WrapError(err, ErrCodeSession, "Failed to update session"))
		return
	}

	logger.Info("User successfully authenticated", "session_id", session.ID)
	http.Redirect(w, r, Cfg.FrontendURL, http.StatusFound)
}
