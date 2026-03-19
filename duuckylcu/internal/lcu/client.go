package lcu

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"time"
)

type Client struct {
	baseURL  string
	password string
	client   *http.Client
}

func NewClient() (*Client, error) {
	lockfile, err := ReadLockfile()
	if err != nil {
		return nil, fmt.Errorf("failed to read flockfile: %w", err)
	}

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	return &Client{
		baseURL:  fmt.Sprintf("https://127.0.0.1:%s", lockfile.Port),
		password: lockfile.Password,
		client:   httpClient,
	}, nil
}

func (c *Client) Get(endpoint string) ([]byte, error) {
	req, err := http.NewRequest("GET", c.baseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}

	auth := base64.StdEncoding.EncodeToString([]byte("riot:" + c.password))
	req.Header.Set("Authorization", "Basic "+auth)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LCU returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *Client) GetFriendList() ([]Friend, error) {
	data, err := c.Get("/lol-chat/v1/friends")
	if err != nil {
		return nil, err
	}

	var friends []Friend
	if err := json.Unmarshal(data, &friends); err != nil {
		return nil, fmt.Errorf("failed to parse friends: %w", err)
	}

	allowed := []string{"Arthur#ツツツ", "Il L7 Il#L77", "Naba#naba1", "dream#heart"}
	trackedFriends := slices.DeleteFunc(friends, func(f Friend) bool {
		name := f.GameName + "#" + f.GameTag
		if slices.Contains(allowed, name) {
			return false
		}
		return true
	})

	return trackedFriends, nil
}

func IsARAM(gameStatus *GameStatus) bool {
	return gameStatus.GameMode == "ARAM" || gameStatus.MapID == "12"
}
