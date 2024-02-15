package porkbungo

import (
	"context"
	"errors"
	"fmt"
)

// Domain Pricing represents the response of the
// domain pricing API
type DomainPricing struct {
	Status  string `json:"status"`
	Pricing map[string] struct {
		Registration	string	`json:"registration"`
		Renewal			string	`json:"renewal"`
		Transfer 		string	`json:"transfer"`
	}	`json:"pricing"`
}

// UpdateNSOptions fields
type UpdateNameServerOptions struct {
	*Keys
	Nameservers	[]string	`json:"ns"`
	Domain		string		`json:"-"`
}

// Type to represent yes or no
type BooleanType string
const (
	Yes	BooleanType = "yes"
	No	BooleanType = "no"
)

// Constants are the redirect types that can be assigned
// to a URL Forward
type ForwardType string
const (
	Temporary	ForwardType = "temporary"
	Permanent	ForwardType = "permanent"
)

// Options for creating a URL Forward
type URLForwardOptions struct {
	*Keys
	// [Optional] A subdomain that you would like to add email forwarding for.
	SubDomain	string	`json:"subdomain,omitempty"`
	// Destination for URL forwarding
	Location	string	`json:"location"`
	//	Type of forward. 
	//	Valid types: [permament,temporary]
	Type		ForwardType	`json:"type"`
	//	Include or don't include URI path in redirection.
	// Valid: [yes,no]
	IncludePath	BooleanType	`json:"includePath"`
	// Forward all subdomains of this domain
	// Valid: [yes,no]
	Wildcard	BooleanType	`json:"wildcard"`
	Domain		string	`json:"-"`
}

// URL forward returned from porkbun API
type URLForward	struct {
	// Subdomain id
	Id			string	`json:"id"`
	// Name of subdomain
	Subdomain	string	`json:"subdomain"`
	// URL forwarding destination
	Location	string	`json:"location"`
	// Type of forward temporary or permanent
	Type		ForwardType	`json:"type"`
	// Whether URI path is included in redirection
	IncludePath	BooleanType	`json:"includePath"`
	//Forward all subdomains of this domain/subdomain
	Wildcard	BooleanType	`json:"wildcard"`
}

// Domain represents a second level domain owned by the account
type Domain struct {
	Domain       string `json:"domain"`
	Status       string `json:"status"`
	TLD          string `json:"tld"`
	CreateDate   string `json:"createDate"`
	ExpireDate   string `json:"expireDate"`
	SecurityLock string `json:"securityLock"`
	WhoisPrivacy string `json:"whoisPrivacy"`
	AutoRenew    int    `json:"autoRenew,string"`
	NotLocal     int    `json:"notLocal"`
}

// Checks default domain pricing information for all supported TLDs.
func (c *Client) GetDomainPricing(ctx context.Context) (*DomainPricing, error) {
	request := c.resty.NewRequest().SetContext(ctx)
	response, err := request.SetResult(&DomainPricing{}).Get("/pricing/get")

	if err != nil || response.StatusCode() != 200{
		if err == nil {
			e,_ := response.Error().(*APIError)
			err = errors.New(e.Message)
		}
		return nil,err
	}

	
	if val,ok := response.Result().(*DomainPricing); ok {
		return val,nil
	}

	return nil,errors.New("error data received is not in correct format")
}

// Update the authoritative nameservers for the specified domain
func (c *Client) UpdateNameServers(ctx context.Context, opts UpdateNameServerOptions) error {
	req := c.resty.NewRequest().SetContext(ctx)

	result := struct {
		Status	string	`json:"status"`
	}{}

	req.SetResult(&result)

	u := fmt.Sprintf("/domain/updateNs/%s", opts.Domain)

	if opts.Keys == nil {
		opts.Keys = c.keys
	}
	
	req.SetBody(opts)

	resp,err := req.Post(u)
	if err != nil || resp.StatusCode() != 200 {
		if err == nil {
			e,_ := resp.Error().(*APIError)
			err = errors.New(e.Message)
		}
		return err
	}

	return nil
}

// Get Authoritative nameservers listed at the registry for the specified domain
func (c *Client) GetNameServers(ctx context.Context, domain string, keys *Keys) ([]string,error) {

	result := struct {
		Status	string		`json:"status"`
		NS		[]string	`json:"ns"`
	}{}

	req := c.resty.NewRequest().SetContext(ctx).SetResult(&result)
	u := fmt.Sprintf("/domain/getNs/%s", domain)
	
	if keys == nil {
		keys = c.keys
	}

	req.SetBody(keys)

	resp,err := req.Post(u)
	if err != nil || resp.StatusCode() != 200 {
		if err == nil {
			e,_ := resp.Error().(*APIError)
			err = errors.New(e.Message)
		}
		return nil,err
	}

	return result.NS,nil
}

// Create a URL forward for the specified domain
func (c *Client) CreateURLForward(ctx context.Context, opts URLForwardOptions) error {

	req := c.resty.NewRequest().SetContext(ctx)
	u := fmt.Sprintf("/domain/addUrlForward/%s", opts.Domain)

	if opts.Keys == nil {
		opts.Keys = c.keys
	}

	req.SetBody(opts)

	resp,err := req.Post(u)
	if err != nil || resp.StatusCode() != 200 {
		if err == nil {
			e,_ := resp.Error().(*APIError)
			err = errors.New(e.Message)
		}
		return err
	}

	return nil
}
// Get URL forwarding for the specified domain
func (c *Client) GetUrlForwards(ctx context.Context, domain string, keys *Keys) ([]URLForward,error) {
	result := struct {
		Status		string			`json:"status"`
		Forwards	[]URLForward	`json:"forwards"`
	}{}

	req := c.resty.NewRequest().SetContext(ctx).SetResult(&result)
	u := fmt.Sprintf("/domain/getUrlForwarding/%s", domain)

	if keys == nil {
		keys = c.keys
	}

	req.SetBody(keys)

	resp,err := req.Post(u)
	if err != nil || resp.StatusCode() != 200 {
		if err == nil {
			e,_ := resp.Error().(*APIError)
			err = errors.New(e.Message)
		}
		return nil,err
	}

	return result.Forwards,nil
}

// Delete the specified URL forwarding for a domain
func (c *Client) DeleteURLForward(ctx context.Context, domain string, forwardId string, keys *Keys) error {
	req := c.resty.NewRequest().SetContext(ctx)
	u := fmt.Sprintf("/domain/deleteUrlForward/%s/%s", domain, forwardId)

	if keys == nil {
		keys = c.keys
	}

	req.SetBody(keys)

	resp,err := req.Post(u)
	if err != nil || resp.StatusCode() != 200 {
		if err == nil {
			e,_ := resp.Error().(*APIError)
			err = errors.New(e.Message)
		}
		return err
	}

	return nil
}

// Get all domain names for the current account
func (c *Client) GetAllDomains(ctx context.Context, keys *Keys) ([]Domain,error) {
	respResult := struct {
		Status	string	 `json:"status"`
		Domains	[]Domain `json:"domains"`
	}{}

	body := struct {
		Keys
		Start	int	`json:"start"`
	}{}

	if keys == nil {
		body.Keys = *c.keys
	} else {
		body.Keys = *keys
	}
	body.Start = 0

	u := "/domain/listAll"
	req := c.resty.NewRequest().SetContext(ctx).SetResult(&respResult)

	result := make([]Domain,0)
	
	req.SetBody(body)
	resp,err := req.Post(u)

	// TODO: Turn this into reusable function
	if err != nil || resp.StatusCode() != 200 {
		if err == nil {
			e,_ := resp.Error().(*APIError)
			err = errors.New(e.Message)
		}
		return result,err
	}

	for len(respResult.Domains) > 0 {
		result = append(result, respResult.Domains...)

		body.Start += 1000
		req.SetBody(body)

		resp,err = req.Post(u)
		if err != nil || resp.StatusCode() != 200 {
			if err == nil {
				e,_ := resp.Error().(*APIError)
				err = errors.New(e.Message)
			}
			return result,err
		}
	} 
	return result,nil
}