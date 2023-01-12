package graph

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/ralpioxxcs/go-onedrive-cli/model"

	"github.com/gorilla/mux"
	"github.com/pkg/browser"
	"github.com/spf13/viper"
)

var (
	authCodeValue string
	wg            sync.WaitGroup
)

// Login performing login action using refresh token and return access, refresh token
func Login(refreshToken string) (string, string) {
	// [Code Flow Authentication]
	// https://docs.microsoft.com/ko-kr/azure/active-directory/develop/v2-oauth2-auth-code-flow
	// 1. Authenticate to get code
	// 2. Get access token using code
	// 3. Call API using access token
	// ---------------------------------------------------------------------------------------

	// 1. Authenticate
	port := fmt.Sprintf(":%s", viper.GetString("auth.redirect_port"))

	router := mux.NewRouter()
	router.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.Header().Add("Content-Type", "application/json")
			h.ServeHTTP(rw, r)
		})
	})

	// auth page에서 redirect url로 post 요청되는 form에서 "code" 쿼리값 파싱하여 저장
	//router.HandleFunc("/authcode", parseAuthCode).Methods("GET")
	router.HandleFunc("/"+viper.GetString("auth.redirect_path"), func(w http.ResponseWriter, r *http.Request) {
		authCodeValue = r.URL.Query().Get("code")
		wg.Done()
	}).Methods("GET")

	wg.Add(1)
	log.Printf("Listening on http://localhost%s\n", port)
	go http.ListenAndServe(port, router)
	browser.OpenURL(generateAuthURL())
	wg.Wait()

	// 2. Get access token
	accessToken, refreshToken, err := getAccessToken(refreshToken)
	if err != nil {
		log.Fatalf("failed to get access token (err : %v)\n", err)
	}

	return accessToken, refreshToken
}

// generateAuthURL generate microsoft authentication url with clients parameters
func generateAuthURL() string {
	// reference : https://docs.microsoft.com/ko-kr/azure/active-directory/develop/v2-oauth2-auth-code-flow#request-an-authorization-code
	// required query components
	// - tenant
	// - client_id
	// - scope
	// - response_type
	// - redirect_uri

	// 인증 코드 요청 페이지 호출
	const responseType = "code"
	const queryParamsFormat = "%s?client_id=%s&scope=%s&response_type=%s&redirect_uri=%s"

	// must be encoded URL from string
	encodedScopeURI := url.QueryEscape(viper.GetString("auth.scope"))
	encodedRedirectURI := url.QueryEscape(makeRedirectURI(viper.GetString("auth.redirectPort"), viper.GetString("auth.redirectPath")))

	url := fmt.Sprintf(queryParamsFormat,
		makeAuthURI(viper.GetString("auth.tenant")), viper.GetString("auth.client_id"), encodedScopeURI, responseType, encodedRedirectURI)

	log.Println("authorization URL: ", url)

	return url
}

// getAccessToken refresh tokens using refresh token and return access, refresh token
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
	formValue.Set("client_id", viper.GetString("auth.client_id"))
	formValue.Set("code", authCodeValue)
	formValue.Set("redirect_uri", makeRedirectURI(viper.GetString("auth.redirectPort"), viper.GetString("auth.redirectPath")))
	formValue.Set("grant_type", grantType)
	if refreshToken != "" {
		formValue.Set("refresh_token", refreshToken)
	}
	formValue.Set("client_secret", viper.GetString("auth.client_secret"))

	resp, err := client.PostForm(makeTokenURI(viper.GetString("auth.tenant")), formValue)
	if err != nil {
		log.Fatalln("error on POST ", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		var errResp model.ErrorResponse
		err := json.Unmarshal(body, &errResp)
		if err != nil {
			log.Fatalln("error on unmarshall", err)
		}
		marshalled, _ := json.Marshal(errResp)
		log.Fatalf("error response : %v\n", string(marshalled))
	}

	var authResp model.AuthtokenResponse
	err = json.Unmarshal(body, &authResp)
	if err != nil {
		log.Fatalln("error on unmarshall", err)
	}

	marshalled, _ := json.Marshal(authResp)
	log.Printf("response body : %v\n", string(marshalled))

	return authResp.AccessToken, authResp.RefreshToken, nil
}

func makeAuthURI(tenant string) string {
	const authURIFormat = "https://login.microsoftonline.com/%s/oauth2/v2.0/authorize"
	return fmt.Sprintf(authURIFormat, tenant)
}

func makeTokenURI(tenant string) string {
	const tokenURIFormat = "https://login.microsoftonline.com/%s/oauth2/v2.0/token"
	return fmt.Sprintf(tokenURIFormat, tenant)
}

func makeRedirectURI(port, path string) string {
	const redirectURLFormat = "http://localhost:%s/%s"
	return fmt.Sprintf(redirectURLFormat, port, path)
}
