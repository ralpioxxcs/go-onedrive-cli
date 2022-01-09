package graph

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/pkg/browser"
)

const (
	loginServerPort = 6789
	authURI         = "https://login.microsoftonline.com/%s/oauth2/v2.0/authorize"
	tokenURI        = "https://login.microsoftonline.com/%s/oauth2/v2.0/token"
	scope           = "files.readwrite offline_access" // white space is replace to be "%20"
	redirectURI     = "http://localhost:6789/authcode"
)

// [Code Flow Authentication]
// https://docs.microsoft.com/ko-kr/azure/active-directory/develop/v2-oauth2-auth-code-flow
// 1. Authenticate to get code
// 2. Get access token using code
// 3. Call API using access token
func Login(refreshToken string) (string, string) {
	// 1. Authenticate
	port := fmt.Sprintf(":%d", loginServerPort)

	router := mux.NewRouter()
	router.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.Header().Add("Content-Type", "application/json")
			h.ServeHTTP(rw, r)
		})
	})

	// parsing post request form value ("code") from redirected URL
	router.HandleFunc("/authcode", parseAuthCode).Methods("GET")

	wg.Add(1)
	log.Printf("Listening on http://localhost%s\n", port)
	go http.ListenAndServe(port, router)
	browser.OpenURL(generateAuthURL())
	wg.Wait()

	// 2. Get access token
	accessToken, refreshToken, err := getAccessToken(refreshToken)
	if err != nil {
		log.Fatalf("failed to get access token (%v)\n", err)
	}

	return accessToken, refreshToken
}

var (
	tenant        string
	clientID      string
	clientSecret  string
	authCodeValue string
	wg            sync.WaitGroup
)

const (
	ACCESS_TOKEN = 1 + iota
	REFRESH_TOKEN
)

// response of get access token request
// (https://docs.microsoft.com/ko-kr/azure/active-directory/develop/v2-oauth2-auth-code-flow#successful-response-2)
type authtokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiressIn   int    `json:"expires_in"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
	IdToken      string `json:"id_token"`
}

type errorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorCode        int    `json:"erorr_codes"`
	TimeStamp        string `json:"timestamp"`
	TraceID          string `json:"trace_id"`
	CorrelationID    string `json:"correlation_id"`
}

func generateAuthURL() string {
	// reference : https://docs.microsoft.com/ko-kr/azure/active-directory/develop/v2-oauth2-auth-code-flow#request-an-authorization-code

	// [required query components]
	//  - tenant, client_id, scope, response_type, redirect_uri
	const responseType = "code"
	authURI := fmt.Sprintf(authURI, tenant)

	// must be encoded URL from string
	encodedScopeURI := url.QueryEscape(scope)
	encodedRedirectURI := url.QueryEscape(redirectURI)

	url := fmt.Sprintf("%s?client_id=%s&scope=%s&response_type=%s&redirect_uri=%s",
		authURI, clientID, encodedScopeURI, responseType, encodedRedirectURI)

	log.Println(url)

	return url
}

func parseAuthCode(rw http.ResponseWriter, r *http.Request) {
	authCodeValue = r.URL.Query().Get("code")

	wg.Done()
}

func getAccessToken(refreshToken string) (string, string, error) {
	// ref : https://docs.microsoft.com/ko-kr/azure/active-directory/develop/v2-oauth2-auth-code-flow#request-an-access-token-with-a-client_secret

	// [required query components]
	//  - tenant, client_id, code, redirect_uri, grant_type, client_secret
	//  - refresh_token (when refresh "true")
	var grantType string
	if refreshToken != "" {
		grantType = "refresh_token"
	} else {
		grantType = "authorization_code"
	}

	client := http.DefaultClient
	formValue := url.Values{}
	formValue.Set("client_id", clientID)
	formValue.Set("code", authCodeValue)
	formValue.Set("redirect_uri", redirectURI)
	formValue.Set("grant_type", grantType)
	if refreshToken != "" {
		formValue.Set("refresh_token", refreshToken)
	}
	formValue.Set("client_secret", clientSecret)

	resp, err := client.PostForm(fmt.Sprintf(tokenURI, tenant), formValue)
	if err != nil {
		log.Fatalln("error on POST ", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		var errResp errorResponse
		err := json.Unmarshal(body, &errResp)
		if err != nil {
			log.Fatalln("error on unmarshall", err)
		}
		marshalled, _ := json.Marshal(errResp)
		log.Fatalf("error response : %v\n", string(marshalled))
	}

	var authResp authtokenResponse
	err = json.Unmarshal(body, &authResp)
	if err != nil {
		log.Fatalln("error on unmarshall", err)
	}

	marshalled, _ := json.Marshal(authResp)
	log.Printf("response body : %v\n", string(marshalled))

	return authResp.AccessToken, authResp.RefreshToken, nil
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("failed to load .env", err)
	}
	tenant = os.Getenv("TENANT")
	clientID = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
}
