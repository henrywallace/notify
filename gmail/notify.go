package gmail

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

var (
	SecretsDir string
	To         string
	From       string

	Scopes = []string{gmail.GmailSendScope}
)

func init() {
	SecretsDir = mustEnv("NOTIFY_SECRETS")
	To = mustEnv("NOTIFY_TO")
	From = mustEnv("NOTIFY_FROM")
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	hash := fnv.New32a()
	hash.Write([]byte(config.ClientID))
	hash.Write([]byte(config.ClientSecret))
	hash.Write([]byte(strings.Join(config.Scopes, " ")))
	hash.Write([]byte(strings.Join(Scopes, " ")))
	basename := fmt.Sprintf("token-%d.json", hash.Sum32())
	tokFile := filepath.Join(SecretsDir, basename)
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func mustEnv(env string) string {
	val := os.Getenv(env)
	if val == "" {
		log.Fatalf("must set env var %s to non-empty value", env)
	}
	return val
}

func Run(subject, body string) {
	cred := filepath.Join(SecretsDir, "google-credentials.json")
	b, err := ioutil.ReadFile(cred)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(
		b,
		Scopes...,
	)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	call := srv.Users.Messages.Send("me", Message{
		From:    From,
		To:      To,
		Subject: subject,
		Body:    body,
	}.Format())
	msg, err := call.Do()
	if err != nil {
		log.Fatalf("failed to send mail: %v", err)
	}
	log.Infof("gmail: sent message: %v", msg.ServerResponse.HTTPStatusCode)
}

type Message struct {
	From    string
	To      string
	Subject string
	Body    string
}

func (m Message) Format() *gmail.Message {
	raw := fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n\r\n"+
			"%s",
		m.From,
		m.To,
		m.Subject,
		m.Body,
	)
	enc := base64.URLEncoding.EncodeToString([]byte(raw))
	return &gmail.Message{
		Raw: enc,
	}
}
