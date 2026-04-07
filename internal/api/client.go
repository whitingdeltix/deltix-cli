package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/whitingdeltix/deltix-cli/internal/config"
)

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		baseURL: cfg.APIURL,
		token:   cfg.Token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) do(method, path string, body interface{}, result interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, c.baseURL+path, bodyReader)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}

	return nil
}

// Auth
func (c *Client) Login(email, password string) (*TokenResponse, error) {
	var resp TokenResponse
	err := c.do("POST", "/auth/login", LoginRequest{Email: email, Password: password}, &resp)
	return &resp, err
}

func (c *Client) GetMe() (*UserResponse, error) {
	var resp UserResponse
	err := c.do("GET", "/auth/me", nil, &resp)
	return &resp, err
}

// Apps
func (c *Client) ListApps() ([]App, error) {
	var apps []App
	err := c.do("GET", "/apps", nil, &apps)
	return apps, err
}

func (c *Client) GetApp(id string) (*App, error) {
	var app App
	err := c.do("GET", "/apps/"+id, nil, &app)
	return &app, err
}

// Tasks
func (c *Client) ListTasks(appID string) ([]Task, error) {
	var tasks []Task
	err := c.do("GET", "/apps/"+appID+"/tasks", nil, &tasks)
	return tasks, err
}

// Runs
func (c *Client) TriggerRun(appID string, req TriggerRunRequest) (*Run, error) {
	var run Run
	err := c.do("POST", "/apps/"+appID+"/runs", req, &run)
	return &run, err
}

func (c *Client) GetRun(runID string) (*Run, error) {
	var run Run
	err := c.do("GET", "/runs/"+runID, nil, &run)
	return &run, err
}

func (c *Client) GetRunResults(runID string) ([]TaskResult, error) {
	var results []TaskResult
	err := c.do("GET", "/runs/"+runID+"/results", nil, &results)
	return results, err
}

// Specs / Playbooks
func (c *Client) ListSpecs(appID string) ([]Spec, error) {
	var specs []Spec
	err := c.do("GET", "/apps/"+appID+"/specs", nil, &specs)
	return specs, err
}

func (c *Client) TriggerPlayback(specID string) (*PlaybackResponse, error) {
	var resp PlaybackResponse
	err := c.do("POST", "/specs/"+specID+"/playback", struct{}{}, &resp)
	return &resp, err
}

func (c *Client) GetPlaybackResult(playbackRunID string) (*PlaybackResult, error) {
	var resp PlaybackResult
	err := c.do("GET", "/playback/"+playbackRunID, nil, &resp)
	return &resp, err
}

// SSE stream URL (for external use)
func (c *Client) StreamURL(runID string) string {
	return c.baseURL + "/runs/" + runID + "/stream"
}
