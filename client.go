package porkbungo

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

const (
    // Hostname to force the use of IPv4
    IPv4ApiHost = "api-ipv4.porkbun.com"
    // Porkbun API hostname
    APIHost = "porkbun.com"
    // Porkbun API base path
    BasePath = "/api/json"
    // Porkbun API version
    APIVersion = "v3"
    // connect to API with http(s)
    APIProto = "https"
)

type Keys struct {
    // API Key usually starts with 'pk1_'
    APIKey          string      `json:"apikey"`
    // Secret key, usually starts with 'sk1_'
    SecretAPIKey    string      `json:"secretapikey"`
}

type APIError struct {
    // Status indicating the command was not successful
    Status  string  `json:"status"`
    // Error message
    Message string  `json:"message"`
}

// Client is a wrapper around the Resty client
type Client struct {
	resty       *resty.Client
    keys        *Keys
    useIPv4     bool
}

// NewClient factory to create a new Client struct
func NewClient() (client Client) {
    client.resty = resty.New()
    client.useIPv4 = false

    u := fmt.Sprintf("%s://%s%s/%s",
                        APIProto,
                        APIHost,
                        BasePath,
                        APIVersion)
    client.resty.SetBaseURL(u).SetError(&APIError{})

    return
}

// NewClient factory to create a new client with the
// the specified API keys
func NewClientWithAPIKeys(keys *Keys) (client Client) {
    client = NewClient()

    client.keys = keys

    return
}

// SetAPIKeys sets the API keys for to be used
// with all the requests from this client. The API keys
// can be overridden for individual requests using the
// requests options
func (c *Client) SetAPIKeys(keys *Keys) {
    c.keys = keys
}

// Force client to use IPv4
func (c *Client) UseIPv4(value bool) {
    c.useIPv4 = value

    c.updateHostUrl()
}

func (c *Client) updateHostUrl() {
    useIPv4 := c.useIPv4

    var apiHost string;
    if useIPv4 {
        apiHost = IPv4ApiHost
    } else {
        apiHost = APIHost
    }

    u := fmt.Sprintf("%s://%s%s/%s",
                    APIProto,
                    apiHost,
                    BasePath,
                    APIVersion)
    c.resty.SetBaseURL(u)
}
