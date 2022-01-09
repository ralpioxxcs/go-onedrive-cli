package graph

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const endPoint = "https://graph.microsoft.com/v1.0/me/drive/root/children"

type listResponse struct {
	Value []item `json:"value"`
}

type item struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

func List(accessToken string) []string {
	client := http.DefaultClient

	// log.Println(accessToken)
	req, _ := http.NewRequest("GET", endPoint, nil)
	bearerToken := "Bearer " + accessToken
	req.Header.Add("Authorization", bearerToken)

	resp, err := client.Do(req)

	if err != nil {
		log.Println("error on POST ", err)
		return nil
	}
	defer resp.Body.Close()

	// log.Printf("status : %v\n", resp.Status)
	if resp.StatusCode != 200 {
		fmt.Println(resp.Status)
		return nil
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var unmarshalledResponse listResponse
	json.Unmarshal(body, &unmarshalledResponse)

	items := []string{}
	for _, item := range unmarshalledResponse.Value {
		items = append(items, item.Name)
	}

	return items
}
