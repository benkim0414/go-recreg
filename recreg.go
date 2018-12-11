// Package recreg provides a client for using the REC Registry API.
package recreg

import (
	"encoding/json"
	"net/http"
	"net/url"
	"path"
	"time"
)

const (
	defaultBaseURL = "http://rec-registry.gov.au/rec-registry/app/api/public-register/"
	userAgent      = "go-recreg"
)

// Response represents a response from the REC Registry API.
type Response struct {
	Status  string   `json:"status"`
	Actions []Action `json:"result,omitempty"`
	Err     string   `json:"errorMessage,omitempty"`
}

// Action represents a certificate action.
type Action struct {
	// The type of the action.
	Type string `json:"actionType"`
	// The timestamp of when the action completed.
	CompleteTime time.Time `json:"completedTime"`
	// The certificate ranges involved in the action.
	Ranges []Range `json:"certificateRanges"`
}

// Range represents the certificate range involved in the action.
type Range struct {
	// The type of certificates.
	CertificateType string `json:"certificateType"`
	// The unique identifier for the registered person.
	PersonID int `json:"registeredPersonNumber"`
	// The unique identifier for accreditation of the installation.
	AccreditationCode string `json:"accreditationCode"`
	// The year when the certificate is created.
	Year int `json:"generationYear"`
	// The state where the systems are installed.
	State string `json:"generationState"`
	// The starting serial number for the certificate range.
	Start int `json:"startSerialNumber"`
	// The ending serial number for the certificate range.
	End int `json:"endSerialNumber"`
	// The fuel source or type of the system.
	FuelSource string `json:"fuelSource"`
	// The owner of the certificate after the action happened.
	Owner string `json:"ownerAccount"`
	// The unique identifier for the owner account.
	OwnerID int `json:"ownerAccountID"`
	// The status of the certificate after the certificate action.
	Status string `json:"status"`
}

// Client is a client manages communication with the REC Registry API.
type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	// Base URL for API requests.
	BaseURL *url.URL
	// User agent used when communicating with the REC Registry API.
	UserAgent string
}

// NewClient returns a new REC Registry API client.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL, _ := url.Parse(defaultBaseURL)
	return &Client{client: httpClient, BaseURL: baseURL, UserAgent: userAgent}
}

const (
	// FirstDate is the first date the REC Registry has data for.
	FirstDate = "2001-05-18"
	// ISO8601Date is the international standard for the representation of dates.
	ISO8601Date = "2006-01-02"
)

// ListActions returns certificate actions for a particular day.
func (c *Client) ListActions(date time.Time) ([]Action, error) {
	v := url.Values{}
	v.Set("date", date.Format(ISO8601Date))
	p := path.Join(c.BaseURL.Path, "certificate-actions")
	ref := &url.URL{Path: p, RawQuery: v.Encode()}
	u := c.BaseURL.ResolveReference(ref)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.UserAgent)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var r Response
	err = json.NewDecoder(resp.Body).Decode(&r)
	return r.Actions, err
}
