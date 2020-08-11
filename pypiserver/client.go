package pypiserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Client struct {
	address string
	apiKey  string
}

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func NewClient(address, api_key string) *Client {
	return &Client{
		address: address,
		apiKey:  api_key,
	}
}

func (c *Client) Users() (*[]User, error) {
	resp, err := c.sendRequest("GET", "/users", []byte{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body: %v", err)
	}
	var userList []User
	if err = json.Unmarshal(body, &userList); err != nil {
		return nil, fmt.Errorf("Failed to decode response: %v", err)
	}
	return &userList, nil
}

func (c *Client) User(username string) (*User, error) {
	resp, err := c.sendRequest("GET", "/users/"+username, []byte{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body: %v", err)
	}
	var user User
	if err = json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("Failed to decode response: %v, %s", err, string(body))
	}
	return &user, nil
}

func (c *Client) CreateUser(user *User) error {
	body, err := json.Marshal(user)
	if err != nil {
		return err
	}
	_, err = c.sendRequest("POST", "/users", body)
	return err
}

func (c *Client) UpdateUser(username string, user *User) error {
	body, err := json.Marshal(user)
	if err != nil {
		return err
	}
	_, err = c.sendRequest("PUT", "/users/"+username, body)
	return err
}

func (c *Client) DeleteUser(username string) error {
	_, err := c.sendRequest("DELETE", "/users/"+username, []byte{})
	return err
}

func (c *Client) sendRequest(method, endpoint string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, c.address+endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("Failed to create request: %v", err)
	}
	req.Header.Add("Authorization", "Token "+c.apiKey)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	log.Printf("[INFO] [Pypi client] %s %s (%v)", method, c.address+endpoint, resp.StatusCode)

	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		return nil, fmt.Errorf("Wrong authorization key")
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 && resp.StatusCode != 204 {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("Failed to read response body: %v", err)
		}
		return nil, fmt.Errorf("Failed request: (%d) %s", resp.StatusCode, string(body))
	}

	return resp, nil
}
