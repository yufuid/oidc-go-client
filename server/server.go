package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"net/url"
	"strings"
	"yufuid.com/oidc-go-client/config"
	"yufuid.com/oidc-go-client/utils"
)

var (
	oidcConfig   config.OIDCConfig
	clientConfig config.ClientConfig
)

func InitServer(config config.ClientConfig) {
	clientConfig = config;
	oidcConfig = utils.GenOIDCConfig(config.ClientId, config.ClientSecret, config.WellKnownUrl);
	e := echo.New()
	e.GET("/", center);
	e.GET("/login", login);
	e.GET("/logout", logout);
	e.GET("/dashboard", dashboard);
	e.GET("/callback", callback);
	e.Logger.Fatal(e.Start(config.Port));
}

func center(context echo.Context) error {
	if _, err := context.Cookie("data"); err != nil {
		return context.Redirect(http.StatusFound, "/login");
	}
	return context.Redirect(http.StatusFound, "/dashboard");
}

func login(context echo.Context) error {
	authUrl := fmt.Sprintf("%s?client_id=%s&response_type=%s&redirect_uri=%s&scope=%s",
		oidcConfig.AuthorizationEndpoint,
		url.QueryEscape(oidcConfig.ClientId),
		url.QueryEscape("code"),
		url.QueryEscape(clientConfig.RedirectUrl),
		url.QueryEscape("openid offline_access"));
	return context.Redirect(http.StatusFound, authUrl);
}

func logout(context echo.Context) error {
	cookie := http.Cookie{};
	cookie.Name = "data";
	cookie.Value = "";
	cookie.Path = "/"
	context.SetCookie(&cookie);
	return context.String(http.StatusOK, "you are logged out now.");
}

func dashboard(context echo.Context) error {
	data, err := context.Cookie("data");
	userInfo := map[string]interface{}{};
	if err != nil {
		return context.JSON(http.StatusForbidden, "no permission,not login in !");
	}
	result, _ := base64.StdEncoding.DecodeString(data.Value);
	json.Unmarshal(result, &userInfo);
	return context.JSON(http.StatusOK, userInfo);
}

func callback(context echo.Context) error {
	code := context.QueryParam("code");
	if code == "" {
		return context.String(http.StatusBadRequest, "no code found in request￿");
	}

	client := utils.NewHTTPClient();
	token, err := getToken(client, code);
	if err != nil {
		return err;
	}
	userInfo, err := getUserInfo(client, token);
	if err != nil {
		return err;
	}
	userInfoStr, _ := json.Marshal(userInfo);
	data := base64.StdEncoding.EncodeToString([]byte(userInfoStr));
	cookie := http.Cookie{};
	cookie.Name = "data";
	cookie.Value = data;
	cookie.Path = "/"
	context.SetCookie(&cookie);
	return context.Redirect(http.StatusFound, "/dashboard");
}

func getToken(client *http.Client, code string) (config.TokenResponse, error) {
	basicString := base64.StdEncoding.EncodeToString([]byte(oidcConfig.ClientId + ":" + oidcConfig.ClientSecret));
	payload := url.Values{
		"grant_type":   []string{"authorization_code"},
		"redirect_uri": []string{clientConfig.RedirectUrl},
		"code":         []string{code},
	}

	req, _ := http.NewRequest("POST", oidcConfig.TokenEndpoint, strings.NewReader(payload.Encode()))
	req.Header.Set("Authorization", "Basic "+basicString);
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	token := config.TokenResponse{};
	if err != nil {
		return token, err
	}
	defer resp.Body.Close();
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return token, err;
	}
	return token, nil;
}

func getUserInfo(client *http.Client, token config.TokenResponse) (interface{}, error) {
	req, _ := http.NewRequest("GET", oidcConfig.UserInfoEndpoint, nil)
	req.Header.Set("Authorization", "Bearer "+token.AccessToken);
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close();
	userInfo := new(interface{});
	e := json.NewDecoder(resp.Body).Decode(&userInfo);
	return userInfo, e;
}
