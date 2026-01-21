package token

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"jump-agent/internal/model"
)

type Response struct {
	Data    model.ConnInfo `json:"data"`
	Message string         `json:"message"`
}

// const api = "https://jump.example.com/api/agent/consume-token"
const addr = "http://127.0.0.1:8080/api/jumpServer/getServer"

func Consume(token string) (*model.ConnInfo, error) {
	body, _ := json.Marshal(map[string]string{
		"token": token,
	})

	req, _ := http.NewRequest("POST", addr, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("token invalid or expired")
	}

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}
