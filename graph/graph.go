package graph

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/gorilla/mux"
	"github.com/pkg/browser"
)

// [Code Flow Authentication]
// https://docs.microsoft.com/ko-kr/azure/active-directory/develop/v2-oauth2-auth-code-flow
// 1. Authenticate to get code
// 2. Get access token using code
// 3. Call API using access token
func Login() string {
	// 1. Authenticate
	port := fmt.Sprintf(":%d", loginServerPort)

	router := mux.NewRouter()
	router.Use(jsonContentTypeMiddleware, loggerMiddleWare)
	router.HandleFunc("/authcode", parseAuthCode).Methods("GET")

	wg.Add(1)
	fmt.Printf("Listening on http://localhost%s\n", port)
	go http.ListenAndServe(port, router)

	browser.OpenURL(generateAuthURL())
	wg.Wait()

	// 2. Access token
	token, err := getAccessToken()
	if err != nil {
		panic(err)
	}
	accessToken = token
	//log.Printf("access token : %v\n", token)

	return accessToken
}

const (
	loginServerPort = 6789
	tenant          = "f740e02f-c3e0-4cff-b4e6-6df192596902"
	clientID        = "77601891-94f3-45fc-9e67-1d629f50ba9c"
	clientSecret    = "o8d7Q~vIFMe450-D8eTEIvwJ1m4S4pVqxkhR6"
	scope           = "files.readwrite"
	redirectURI     = "http://localhost:6789/authcode"
)

var (
	authCodeValue string
	accessToken   string
	wg            sync.WaitGroup
)

type authtokenResponse struct {
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func generateAccessTokenURL() string {
	tokenURI := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenant)
	return tokenURI
}

func generateAuthURL() string {
	authURI := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/authorize", tenant)
	responseType := "code"
	encodedRedirectURI := url.QueryEscape(redirectURI)

	url := fmt.Sprintf("%s?client_id=%s&scope=%s&response_type=%s&redirect_uri=%s",
		authURI, clientID, scope, responseType, encodedRedirectURI)

	//log.Printf("request URL : %v\n", url)

	return url
}

func parseAuthCode(rw http.ResponseWriter, r *http.Request) {
	codeKey := r.URL.Query().Get("code")
	log.Printf("code : %v\n", codeKey)

	authCodeValue = codeKey

	log.Printf("code : %v\n", authCodeValue)

	wg.Done()
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func loggerMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)
		next.ServeHTTP(rw, r)
	})
}

func getAccessToken() (string, error) {
	client := http.DefaultClient

	formValue := url.Values{}
	formValue.Set("client_id", clientID)
	formValue.Set("redirect_uri", redirectURI)
	formValue.Set("client_secret", clientSecret)
	formValue.Set("code", authCodeValue)
	formValue.Set("grant_type", "authorization_code")

	resp, err := client.PostForm(generateAccessTokenURL(), formValue)
	if err != nil {
		log.Println("error on POST ", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println(resp.Status)
		return "", nil
	}
	body, _ := ioutil.ReadAll(resp.Body)

	var unmarshalledResponse authtokenResponse
	err = json.Unmarshal(body, &unmarshalledResponse)
	if err != nil {
		log.Println("error on unmarshall", err)
		return "", err
	}
	log.Printf("response body : %v\n", unmarshalledResponse)
	return unmarshalledResponse.AccessToken, nil
}

func getList(path string) (string, error) {
	client := http.DefaultClient

	url := "https://graph.microsoft.com/v1.0/me/drive/root/children"

	req, _ := http.NewRequest("GET", url, nil)
	bearerToken := "Bearer " + accessToken
	req.Header.Add("Authorization", bearerToken)

	resp, err := client.Do(req)

	if err != nil {
		log.Println("error on POST ", err)
		return "", err
	}
	defer resp.Body.Close()

	log.Printf("status : %v\n", resp.Status)
	if resp.StatusCode != 200 {
		fmt.Println(resp.Status)
		return "", nil
	}
	//_, _ := ioutil.ReadAll(resp.Body)
	//log.Printf("body : %v\n", body)

	return "", nil
}
