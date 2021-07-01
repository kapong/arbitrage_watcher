package networking

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var client = &http.Client{Timeout: 15 * time.Second}

// GetContent - get JSON content from API and response via given interface
func GetContent(uri string) ([]byte, error) {
	if uri == "" {
		return nil, errors.New("URI cannot be nil")
	}

	var data []byte
	response, err := client.Get(uri)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// PostContent update the content to API
func PostContent(uri string, jsonStr string) ([]byte, error) {
	if uri == "" {
		return nil, errors.New("URI cannot be nil")
	}

	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return data, nil
}
