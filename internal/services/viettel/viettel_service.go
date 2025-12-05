// internal/services/viettel_service.go
package viettel

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
   "strings"
	_ "YSKH_DMS/internal/models" // nếu cần lưu log, user, v.v.
)

// Struct nhận token từ Viettel
type ViettelTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// Struct để cache token (thread-safe)
type tokenCache struct {
	Token     string
	ExpiresAt time.Time
	mu        sync.RWMutex
}

var (
	viettelTokenCache = &tokenCache{}
	client            = &http.Client{Timeout: 30 * time.Second}
)

// Lấy token (có cache tự động)
func GetViettelAccessToken() (string, error) {
	// Kiểm tra cache trước
	viettelTokenCache.mu.RLock()
	if viettelTokenCache.Token != "" && time.Now().Before(viettelTokenCache.ExpiresAt) {
		defer viettelTokenCache.mu.RUnlock()
		return viettelTokenCache.Token, nil
	}
	viettelTokenCache.mu.RUnlock()

	// Nếu hết hạn → gọi lại API lấy token mới
	data := url.Values{
		"client_id":     {"0199cc1f-d566-7188-b1f2-85bbc3bdcdf3"},
		"client_secret": {"1ZJcVj0Qa7YxhiarBb7dmJXyM3hLJUZn"},
		"grant_type":    {"client_credentials"},
	}

	resp, err := http.PostForm("https://app.vietteldms.com/openapi/identity/token", data)
	if err != nil {
		return "", fmt.Errorf("lỗi kết nối Viettel DMS: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("lấy token thất bại %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp ViettelTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("parse token thất bại: %v", err)
	}

	// Cache token (hết hạn trước 60 giây để an toàn)
	expireTime := time.Now().Add(time.Duration(tokenResp.ExpiresIn-60) * time.Second)
	viettelTokenCache.mu.Lock()
	viettelTokenCache.Token = tokenResp.AccessToken
	viettelTokenCache.ExpiresAt = expireTime
	viettelTokenCache.mu.Unlock()

	return tokenResp.AccessToken, nil
}
func BuildViettelDistURL(urlApi string,queryApi any) (string, error) {
	queryJSON, _ := json.Marshal(queryApi)

	params := url.Values{}
	params.Add("query", string(queryJSON))
	return urlApi+ "?" + params.Encode(), nil
}
// Gọi GET API Viettel có Bearer token
func ViettelGet(urlendpoint string) ([]byte, error) {
	token, err := GetViettelAccessToken()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET",urlendpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// Gọi POST API Viettel có Bearer token
func ViettelPost(endpoint string, payload any) ([]byte, error) {
	token, err := GetViettelAccessToken()
	if err != nil {
		return nil, err
	}

	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", "https://app.vietteldms.com"+endpoint, strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}