package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lestrrat-go/jwx/jwk"
	"log"
	"net/http"
	"yufuid.com/oidc-go-client/config"
)

func GenOIDCConfig(clientId string, clientSecret string, wellKnownUrl string) config.OIDCConfig {
	oidcConfig := config.OIDCConfig{};
	client := NewHTTPClient();
	req, _ := http.NewRequest("GET", wellKnownUrl, nil)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err);
	}
	defer resp.Body.Close()

	oidcConfig.ClientId = clientId;
	oidcConfig.ClientSecret = clientSecret;
	if err := json.NewDecoder(resp.Body).Decode(&oidcConfig); err != nil {
		log.Printf("failed to load .well-known: %v\n", err)
		return oidcConfig;
	}

	//写入keys文件
	set, err := jwk.Fetch(context.Background(), oidcConfig.Keys)
	oidcConfig.KeySet = set;
	return oidcConfig;
}
