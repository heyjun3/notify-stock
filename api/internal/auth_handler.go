package notifystock

import (
	"net/http"
	"net/url"
	"strings"
)

const googleOauthURL = "https://accounts.google.com/o/oauth2/v2/auth"

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

type AuthHandler struct {
	Sessions         *Sessions
	googleClient     *GoogleClient
	memberRepository *MemberRepository
}

func NewAuthHandler(sessions *Sessions, googleClient *GoogleClient, memberRepository *MemberRepository) *AuthHandler {
	return &AuthHandler{
		Sessions:         sessions,
		googleClient:     googleClient,
		memberRepository: memberRepository,
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

	token, err := h.googleClient.ExchangeToken(r.Context(), code)
	if err != nil {
		WriteErrorResponse(w, WrapError(err, ErrCodeExternalService, "Failed to exchange OAuth code"))
		return
	}

	userInfo, err := h.googleClient.GetUserInfo(r.Context(), token)
	if err != nil {
		WriteErrorResponse(w, WrapError(err, ErrCodeExternalService, "Failed to get user info"))
		return
	}

	member, err := NewGoogleMember(nil, userInfo.ID, userInfo.Email, userInfo.VerifiedEmail,
		userInfo.Name, userInfo.GivenName, userInfo.FamilyName, userInfo.Picture)
	if err != nil {
		WriteErrorResponse(w, WrapError(err, ErrCodeInternalServer, "Failed to create Google member"))
		return
	}
	member, err = h.memberRepository.GetOrCreateGoogleMember(r.Context(), member)
	if err != nil {
		WriteErrorResponse(w, WrapError(err, ErrCodeDatabase, "Failed to save or retrieve member"))
		return
	}
	logger.Info("User info response received", "response", userInfo)

	// セッションの有効化
	session.MemberID = member.ID
	session.IsActive = true
	err = h.Sessions.Store(r.Context(), session)
	if err != nil {
		WriteErrorResponse(w, WrapError(err, ErrCodeSession, "Failed to update session"))
		return
	}

	logger.Info("User successfully authenticated", "session_id", session.ID)
	http.Redirect(w, r, Cfg.FrontendURL, http.StatusFound)
}
