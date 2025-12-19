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
   _ "strings"
	_ "YSKH_DMS/internal/models" // nếu cần lưu log, user, v.v.
	// "github.com/davecgh/go-spew/spew"
	
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
	client            = &http.Client{Timeout: 3600 * time.Second}
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

	// Cache token (hết hạn trước 10 phút để an toàn)
	expireTime := time.Now().Add(time.Duration(tokenResp.ExpiresIn-600) * time.Second)
	viettelTokenCache.mu.Lock()
	viettelTokenCache.Token = tokenResp.AccessToken
	viettelTokenCache.ExpiresAt = expireTime
	viettelTokenCache.mu.Unlock()

	return tokenResp.AccessToken, nil
}
func BuildViettelDistURL(urlApi string,queryApi any,page int,size int) (string, error) {
		queryJSON, err := json.Marshal(queryApi)
		if err != nil {
			return "", err // hoặc xử lý lỗi phù hợp
		}

		// Chuyển []byte thành string
		queryStr := string(queryJSON)

		baseURL, err := url.Parse(urlApi)
		if err != nil {
			return "", err
		}

		queryParams := baseURL.Query()
		queryParams.Add("query", queryStr) 
		queryParams.Add("page",fmt.Sprintf("%d", page)) 
		queryParams.Add("size",fmt.Sprintf("%d",size )) 
		baseURL.RawQuery = queryParams.Encode()

		fullURL := baseURL.String()
		
		return fullURL, nil
	
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

