package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type NoSQLClient struct {
	BaseURL string
	APIKey  string
}

func (c *NoSQLClient) CreateTable(name string) error {
	url := fmt.Sprintf("%s/createTable", c.BaseURL)
	data := map[string]string{"name": name}
	return c.post(url, data)
}

func (c *NoSQLClient) DeleteTable(name string) error {
	url := fmt.Sprintf("%s/deleteTable", c.BaseURL)
	data := map[string]string{"name": name}
	return c.post(url, data)
}

func (c *NoSQLClient) Set(table, key, value string) error {
	url := fmt.Sprintf("%s/set", c.BaseURL)
	data := map[string]string{"table_name": table, "key": key, "value": value}
	return c.post(url, data)
}

func (c *NoSQLClient) Get(table, key string) (string, error) {
	url := fmt.Sprintf("%s/get?table_name=%s&key=%s", c.BaseURL, table, key)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("X-API-Key", c.APIKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)
	return result["value"], nil
}

func (c *NoSQLClient) Delete(table, key string) error {
	url := fmt.Sprintf("%s/delete", c.BaseURL)
	data := map[string]string{"table_name": table, "key": key}
	return c.post(url, data)
}

func (c *NoSQLClient) post(url string, data map[string]string) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.APIKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d; body: %s", resp.StatusCode, string(body))
	}

	return nil
}
