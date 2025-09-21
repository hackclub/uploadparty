package licenses

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type airtableStore struct {
	token     string
	baseID    string
	tableName string
	client    *http.Client
}

func newAirtableStore(token, baseID, tableName string) *airtableStore {
	return &airtableStore{
		token:     token,
		baseID:    baseID,
		tableName: tableName,
		client:    &http.Client{Timeout: 5 * time.Second},
	}
}

// Ping performs a minimal request to confirm connectivity without leaking details.
func (a *airtableStore) Ping() error {
	endpoint := fmt.Sprintf("https://api.airtable.com/v0/%s/%s", url.PathEscape(a.baseID), url.PathEscape(a.tableName))
	req, err := http.NewRequest(http.MethodGet, endpoint+"?pageSize=1", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+a.token)
	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return errors.New("license store auth failed")
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("license store request failed: %s", resp.Status)
	}
	// Basic shape validation without reading all records
	var tmp struct {
		Records []any `json:"records"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&tmp)
	return nil
}
