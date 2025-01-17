package ipfs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/your-project/auth"
)

type TatumClient struct {
	apiKey  string
	baseURL string
	client  *http.Client
	logger  *logrus.Logger
	auth    *auth.Authenticator
}

type UploadResponse struct {
	IPFSID string `json:"ipfsHash"`
}

func NewTatumClient(apiKey string, logger *logrus.Logger) *TatumClient {
	return &TatumClient{
		apiKey:  apiKey,
		baseURL: "https://api.tatum.io/v3/ipfs",
		client: &http.Client{
			Timeout: time.Second * 30,
		},
		logger: logger,
	}
}

// UploadGraph uploads a trust graph to IPFS via Tatum
func (t *TatumClient) UploadGraph(ctx context.Context, data interface{}) (string, error) {
	if err := t.auth.ValidateContext(ctx); err != nil {
		return "", fmt.Errorf("authentication failed: %w", err)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal graph data: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", t.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", t.apiKey)

	t.logger.WithField("size", len(jsonData)).Info("Uploading graph to IPFS")

	resp, err := t.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to upload: status %d: %s", resp.StatusCode, string(body))
	}

	var uploadResp UploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	t.logger.WithField("ipfs_id", uploadResp.IPFSID).Info("Successfully uploaded graph to IPFS")
	return uploadResp.IPFSID, nil
}

// GetGraph retrieves a trust graph from IPFS via Tatum
func (t *TatumClient) GetGraph(ctx context.Context, ipfsID string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", t.baseURL, ipfsID)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("x-api-key", t.apiKey)

	t.logger.WithField("ipfs_id", ipfsID).Info("Retrieving graph from IPFS")

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get graph: status %d: %s", resp.StatusCode, string(body))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return data, nil
} 