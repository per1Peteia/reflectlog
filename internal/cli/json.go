package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func postJSON(url string, payload interface{}) (*http.Response, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling JSON: %s\n", err)
	}

	client := &http.Client{Timeout: TIMEOUT}
	res, err := client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("Error posting JSON: %v\n", err)
	}

	if res.StatusCode != http.StatusOK {
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("Error reading response body: %v\n", err)
		}
		fmt.Fprintf(os.Stderr, "Server Error (%d): %s\n", res.StatusCode, string(data))
	}

	return res, nil
}
