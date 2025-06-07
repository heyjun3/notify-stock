package notifystock

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const googleOauthURL = "https://accounts.google.com/o/oauth2/v2/auth"
const googleAccessTokenURL = "https://oauth2.googleapis.com/token"

var (
	scope = []string{
		"openid",
		"https://www.googleapis.com/auth/userinfo.profile",
		"https://www.googleapis.com/auth/userinfo.email",
	}
)

func LoginHandler(sessions *Sessions) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// 既存セッションの確認
		session, err := sessions.Get(r)
		if err == nil && session.IsActive {
			http.Redirect(w, r, Cfg.FrontendURL, http.StatusFound)
			return // already logged in
		}

		// 新しいセッションの作成
		newSession, err := sessions.New(w, r.Context())
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
}

func CallbackHandler(sessions *Sessions) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// セッションの取得
		session, err := sessions.Get(r)
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

		logger.Debug("OAuth token exchange successful", "response_length", len(body))

		// セッションの有効化
		session.IsActive = true
		err = sessions.Store(r.Context(), session)
		if err != nil {
			WriteErrorResponse(w, WrapError(err, ErrCodeSession, "Failed to update session"))
			return
		}

		logger.Info("User successfully authenticated", "session_id", session.ID)
		http.Redirect(w, r, Cfg.FrontendURL, http.StatusFound)
	}
}
