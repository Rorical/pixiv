package pixiv

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"github.com/dghubble/sling"
)

const (
	clientID         = "MOBrBDS8blbauoSck0ZfDbtuzpyT"
	clientSecret     = "lsACyCD94FhDUtGTXi3QzcFE2uU1hqtDaKeqrdwj"
	clientHashSecret = "28c1fdd170a5204386cb1313c7077b34f83e4aaf4aa829ce78c231e05b0bae2c"
	authHosts        = "https://oauth.secure.pixiv.net/"
)

//BasePixiv struct
type BasePixiv struct {
	AccessToken      string
	UserID           string
	RefreshToken     string
	sling            *sling.Sling
	TokenDeadline    time.Time
	AccessChangeHook func(string)
}

func genClientHash(clientTime string) string {
	h := md5.New()
	io.WriteString(h, clientTime)
	io.WriteString(h, clientHashSecret)
	return hex.EncodeToString(h.Sum(nil))
}

//BasePixivAPI is basic tools which will be called by others
func BasePixivAPI() *BasePixiv {
	s := sling.New().Base(authHosts).Set("User-Agent", "PixivAndroidApp/5.0.64 (Android 6.0)")
	return &BasePixiv{
		sling:         s,
		UserID:        "",
		AccessToken:   "default",
		RefreshToken:  "",
		TokenDeadline: time.Now(),
	}
}

func (api *BasePixiv) refreshXTime() {
	clientTime := time.Now().Format(time.RFC3339)
	api.sling.Set("X-Client-Time", clientTime).Set("X-Client-Hash", genClientHash(clientTime))
}

//SetAuth requires AccessToken and RefreshToken to auth
func (api *BasePixiv) SetAuth(AccessToken string, RefreshToken string) {
	api.RefreshToken = RefreshToken
	if AccessToken != "" {
		api.AccessToken = AccessToken
	} else {
		api.Login("", "")
	}
	api.TokenDeadline = time.Now().Add(time.Duration(3600) * time.Second)
}

type AuthParams struct {
	GetSecureURL int    `url:"get_secure_url,omitempty"`
	ClientID     string `url:"client_id,omitempty"`
	ClientSecret string `url:"client_secret,omitempty"`
	GrantType    string `url:"grant_type,omitempty"`
	Username     string `url:"username,omitempty"`
	Password     string `url:"password,omitempty"`
	RefreshToken string `url:"refresh_token,omitempty"`
}
type LoginResponse struct {
	Response *AuthInfo `json:"response"`
}
type LoginError struct {
	HasError bool              `json:"has_error"`
	Errors   map[string]Perror `json:"errors"`
}
type Perror struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
type AuthInfo struct {
	AccessToken  string   `json:"access_token"`
	ExpiresIn    int      `json:"expires_in"`
	TokenType    string   `json:"token_type"`
	Scope        string   `json:"scope"`
	RefreshToken string   `json:"refresh_token"`
	User         *Account `json:"user"`
	DeviceToken  string   `json:"device_token"`
}

//Account struct
type Account struct {
	ID               string               `json:"id"`
	Name             string               `json:"name"`
	Account          string               `json:"account"`
	MailAddress      string               `json:"mail_address"`
	IsPremium        bool                 `json:"is_premium"`
	XRestrict        int                  `json:"x_restrict"`
	IsMailAuthorized bool                 `json:"is_mail_authorized"`
	ProfileImage     AccountProfileImages `json:"profile_image_urls"`
}

//AccountProfileImages struct
type AccountProfileImages struct {
	Px16  string `json:"px_16x16"`
	Px50  string `json:"px_50x50"`
	Px170 string `json:"px_170x170"`
}

//Auth requires authParams
func (api *BasePixiv) Auth(params *AuthParams) (*AuthInfo, error) {
	api.refreshXTime()
	res := &LoginResponse{
		Response: &AuthInfo{
			User: &Account{},
		},
	}
	loginErr := &LoginError{
		Errors: map[string]Perror{},
	}
	_, err := api.sling.New().Post("auth/token").BodyForm(params).Receive(res, loginErr)
	if err != nil {
		return nil, err
	}
	if loginErr.HasError {
		for k, v := range loginErr.Errors {
			return nil, fmt.Errorf("Login %s error: %s", k, v.Message)
		}
	}
	api.TokenDeadline = time.Now().Add(time.Duration(res.Response.ExpiresIn) * time.Second)
	api.AccessToken = res.Response.AccessToken
	//fmt.Println(api.AccessToken)
	if api.AccessChangeHook != nil {
		api.AccessChangeHook(api.AccessToken)

	}
	api.RefreshToken = res.Response.RefreshToken
	api.UserID = res.Response.User.ID
	return res.Response, nil
}

//HookAccessToken call a function when the AccessToken changed
func (api *BasePixiv) HookAccessToken(token func(string)) {
	api.AccessChangeHook = token
}

//Login can login
func (api *BasePixiv) Login(username, password string) (*AuthInfo, error) {
	params := &AuthParams{
		GetSecureURL: 1,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
	if username != "" && password != "" {
		params.GrantType = "password"
		params.Username = username
		params.Password = password
	} else {
		if time.Now().Before(api.TokenDeadline) {
			return nil, nil
		}
		params.GrantType = "refresh_token"
		params.RefreshToken = api.RefreshToken
	}
	a, err := api.Auth(params)
	if err != nil {
		return nil, err
	}
	return a, nil
}
