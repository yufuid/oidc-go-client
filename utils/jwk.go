package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"yufuid.com/oidc-go-client/config"
)

func GenOIDCConfig(clientId string, clientSecret string, wellKnownUrl string) config.OIDCConfig {
	config := config.OIDCConfig{};
	client := NewHTTPClient();
	req, _ := http.NewRequest("GET", wellKnownUrl, nil)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err);
	}
	defer resp.Body.Close()

	config.ClientId = clientId;
	config.ClientSecret = clientSecret;
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		log.Printf("failed to load .well-known: %v\n", err)
		return config;
	}
	return config;
}
