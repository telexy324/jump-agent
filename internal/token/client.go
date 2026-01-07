package token

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"jump-agent/internal/model"
)

const api = "https://jump.example.com/api/agent/consume-token"

func Consume(token string) (*model.ConnInfo, error) {
	body, _ := json.Marshal(map[string]string{
		"token": token,
	})

	req, _ := http.NewRequest("POST", api, bytes.NewReader(body))
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

	var conn model.ConnInfo
	if err := json.NewDecoder(resp.Body).Decode(&conn); err != nil {
		return nil, err
	}

	return &conn, nil
}
