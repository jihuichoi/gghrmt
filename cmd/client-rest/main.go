package main

import (
	"flag"
	"time"
	"net/http"
	"strings"
	"fmt"
	"log"
	"io/ioutil"
	"encoding/json"
)

func main() {
	//get configuration
	address := flag.String("server", "http://localhost:8080", "HTTP gateway url, e.g. http://localhost:8080")
	flag.Parse()

	t := time.Now().In(time.UTC)
	pfx := t.Format(time.RFC3339Nano)

	var body string

	// Call Create
	resp, err := http.Post(*address+"/v1/todo", "applicatin/json", strings.NewReader(fmt.Sprintf(`
	{
		"api": "v1",
		"toDo": {
			"title": "title (%s)",
			"description": "description (%s)",
			"reminder": "%s"
		}
	}
`, pfx, pfx, pfx)))
	if err != nil {
		log.Fatalf("failed to call Create method: %v", err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("failed to read Create response body: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("Create response: code=%d, Body=%s\n\n", resp.StatusCode, body)

	// parse ID of created ToDo
	var created struct {
		API string `json:"api"`
		ID string `json:"id"`
	}
	err = json.Unmarshal(bodyBytes, &created)
	if err != nil {
		log.Fatalf("failed to unamrshal JSON response of Create method: %v", err)
		fmt.Println("error", err)
	}

	// Call Read
	resp, err = http.Get(fmt.Sprintf("%s%s/%s", *address, "/v1/todo", created.ID))
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		body = fmt.Sprintf("failed to Read response body: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("Read response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

	// Call Update
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s%s/%s", *address, "/v1/todo", created.ID),
		strings.NewReader(fmt.Sprintf(`
		{
			"api":"v1",
			"toDo": {
					"title":"title (%s) + updated",
					"description": "description (%s) + updated",
					"reminder": "%s"
			}
		}
		`, pfx, pfx, pfx)))
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("failed to call Updated method: %v", err)
	}
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("failed to read Update response body: %v", err)
	} else {
		body = string(bodyBytes)
	}

	// Call ReadAll
	resp, err = http.Get(*address + "/v1/todo/all")
	if err != nil {
		log.Fatalf("failed to call ReadAll method: %v", err)
	}
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("failed to read ReadAll response body: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("ReadAll response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

	// Call Delete
	req, err = http.NewRequest("DELETE", fmt.Sprintf("%s%s/%s", *address, "/v1/todo", created.ID), nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("failed to call Delete method: %v", err)
	}
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("failed to read Delete response body: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("Delete response: Code=%d, Body=%s\n\n", resp.StatusCode, body)
}