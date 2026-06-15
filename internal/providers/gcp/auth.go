package gcp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/pkg/browser"

	"github.com/nitinchouhan1/cloudctl/internal/schemas"
	"github.com/nitinchouhan1/cloudctl/internal/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func Login() error {

	listener, err := net.Listen("tcp", "localhost:8085")
	if err != nil {
		return err
	}

	redirectURI := "http://localhost:8085/callback"
	oauthConfig := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/cloud-platform",
			"https://www.googleapis.com/auth/userinfo.email",
		},

		Endpoint:    google.Endpoint,
		RedirectURL: redirectURI,
	}

	codeChan := make(chan string)

	mux := http.NewServeMux()

	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {

		code := r.URL.Query().Get("code")

		fmt.Fprintln(
			w,
			"Authentication successful. You may close this tab.",
		)

		codeChan <- code
	})

	server := &http.Server{
		Handler: mux,
	}

	go server.Serve(listener)

	authURL := oauthConfig.AuthCodeURL(
		"cloudctl",
		oauth2.AccessTypeOffline,
	)

	fmt.Println("Opening browser...")

	if err := browser.OpenURL(authURL); err != nil {
		return err
	}

	var code string

	select {
	case code = <-codeChan:

	case <-time.After(5 * time.Minute):
		return fmt.Errorf("authentication timed out")
	}

	token, err := oauthConfig.Exchange(
		context.Background(),
		code,
	)

	if err != nil {
		return err
	}

	cfg, err := utils.LoadConfig()
	if err != nil {
		return err
	}

	cfg.CurrentProvider = "gcp"

	if cfg.Providers == nil {
		cfg.Providers = make(map[string]schemas.Provider)
	}

	cfg.Providers["gcp"] = schemas.Provider{
		RefreshToken: token.RefreshToken,
	}

	if err := utils.SaveConfig(cfg); err != nil {
		return err
	}

	fmt.Println("✓ Logged into GCP")

	return nil
}
