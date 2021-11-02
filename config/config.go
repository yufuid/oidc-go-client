package config

type ClientConfig struct {
	Type         string
	ClientId     string
	ClientSecret string
	WellKnownUrl string
	RedirectUrl  string
	Port         string
}

type OIDCConfig struct {
	ClientId              string
	ClientSecret          string
	Issuer                string
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
	UserInfoEndpoint      string `json:"userinfo_endpoint"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token,omitempty"`
}

type UserInfoResponse struct {
	Sub          string `json:"sub,omitempty"`
	Email        string `json:"token_type"`
	Name         int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token,omitempty"`
}

type Server struct {
	Port string
}
