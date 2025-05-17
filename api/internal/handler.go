package notifystock

import (
	"net/http"
	"net/url"
	"strings"
	"io"
)

const googleOauthURL = "https://accounts.google.com/o/oauth2/v2/auth"
const googleAccessTokenURL = "https://oauth2.googleapis.com/token"

var scope = []string{
	"openid",
	"https://www.googleapis.com/auth/userinfo.profile",
	"https://www.googleapis.com/auth/userinfo.email",
}

// TODO: random state per request
var state = "state"

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(googleOauthURL)
	if err != nil {
		http.Error(w, "Failed to parse URL", http.StatusInternalServerError)
		return
	}
	q := url.Values{}
	q.Add("scope", strings.Join(scope, " "))
	q.Add("access_type", "offline")
	q.Add("include_granted_scopes", "true")
	q.Add("response_type", "code")
	q.Add("state", state)
	q.Add("redirect_uri", Cfg.OauthRedirectURL)
	q.Add("client_id", Cfg.OauthClientID)

	u.RawQuery = q.Encode()
	logger.Info("LoginHandler", "url", u.String())
	http.Redirect(w, r, u.String(), http.StatusFound)
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}
	s := r.URL.Query().Get("state")
	if s != state {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}

	f := url.Values{}
	f.Add("code", code)
	f.Add("client_id", Cfg.OauthClientID)
	f.Add("client_secret", Cfg.OauthClientSecret)
	f.Add("redirect_uri", Cfg.OauthRedirectURL)
	f.Add("grant_type", "authorization_code")
	form := f.Encode()
	logger.Info("CallbackHandler", "form", form)

	u, err := url.Parse(googleAccessTokenURL)
	if err != nil {
		http.Error(w, "Failed to parse URL", http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest(
		http.MethodPost, u.String(),
		strings.NewReader(form),
	)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		http.Error(w, "Failed to get access token", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}
	logger.Info("CallbackHandler", "body", string(body))

	http.Redirect(w, r, "http://localhost:8080", http.StatusFound)
}
