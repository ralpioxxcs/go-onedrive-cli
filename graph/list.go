package graph

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	allItems    = "https://graph.microsoft.com/v1.0/me/drive/root/children"
	recentFiles = "https://graph.microsoft.com/v1.0/me/drive/recent"
)

var (
	ErrRequest = errors.New("failed to request")
)

type listResponse struct {
	Items []Items `json:"value"`
}

type Items struct {
	// inherited from baseItem
	Name           string `json:"name"`
	ID             string `json:"id"`
	CreateDateTime string `json:"createdDateTime"`
	ETag           string `json:"eTag"`
	DownloadURL    string `json:"@microsoft.graph.downloadUrl"`
}

func List(accessToken string) ([]Items, error) {
	client := http.DefaultClient

	req, _ := http.NewRequest("GET", allItems, nil)
	bearerToken := "Bearer " + accessToken
	req.Header.Add("Authorization", bearerToken)

	resp, err := client.Do(req)
	if err != nil {
		log.Println("error on POST ", err)
		return []Items{}, ErrRequest
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println(resp.Status)
		return []Items{}, ErrRequest
	}

	body, _ := io.ReadAll(resp.Body)
	unmarshalledResponse := listResponse{}
	if err := json.Unmarshal(body, &unmarshalledResponse); err != nil {
		log.Fatal(err)
	}

	return unmarshalledResponse.Items, nil
}
