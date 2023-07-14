package tests

import (
	"fmt"
	"testing"
	"time"
)

import (
	"bytes"
	"encoding/json"
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

func TestPerformance(t *testing.T) {
	client := &NoSQLClient{
		BaseURL: "http://localhost:8000",
		APIKey:  "your_api_key",
	}

	tableName := "test_table"
	client.CreateTable(tableName)
	defer client.DeleteTable(tableName)

	t.Run("SetPerformance", func(t *testing.T) { testSetPerformance(t, client, tableName) })
	t.Run("GetPerformance", func(t *testing.T) { testGetPerformance(t, client, tableName) })
	t.Run("DeletePerformance", func(t *testing.T) { testDeletePerformance(t, client, tableName) })
}

func testSetPerformance(t *testing.T, client *NoSQLClient, tableName string) {
	start := time.Now()
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key_%d", i)
		value := fmt.Sprintf("value_%d", i)
		client.Set(tableName, key, value)
	}
	duration := time.Since(start)
	t.Logf("Set performance: %v", duration)
}

func testGetPerformance(t *testing.T, client *NoSQLClient, tableName string) {
	start := time.Now()
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key_%d", i)
		client.Get(tableName, key)
	}
	duration := time.Since(start)
	t.Logf("Get performance: %v", duration)
}

func testDeletePerformance(t *testing.T, client *NoSQLClient, tableName string) {
	start := time.Now()
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key_%d", i)
		client.Delete(tableName, key)
	}
	duration := time.Since(start)
	t.Logf("Delete performance: %v", duration)
}
