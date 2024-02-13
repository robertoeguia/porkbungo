package porkbungo

import (
	"context"
	"errors"
	"fmt"
)

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

type BooleanType string
const (
	Yes	BooleanType = "yes"
	No	BooleanType = "no"
)

type ForwardType string
const (
	Temporary	ForwardType = "temporary"
	Permanent	ForwardType = "permanent"
)

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

type URLForward	struct {
	Id			string		`json:"id"`
	Subdomain	string	`json:"subdomain"`
	Location	string	`json:"location"`
	Type		string	`json:"type"`
	IncludePath	string	`json:"includePath"`
	Wildcard	string	`json:"wildcard"`
}

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

func (c *Client) UpdateNameServers(ctx context.Context, opts UpdateNameServerOptions) error {
	req := c.resty.NewRequest().SetContext(ctx)

	result := struct {
		Status	string	`json:"status"`
	}{}

	req.SetResult(&result)

	u := fmt.Sprintf("/domain/updateNs/%v", opts.Domain)

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

func (c *Client) CreateURLForward(ctx context.Context, opts URLForwardOptions) error {

	req := c.resty.NewRequest().SetContext(ctx)
	u := fmt.Sprintf("/domain/addUrlForward/%v", opts.Domain)

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

func (c *Client) GetUrlForwards(ctx context.Context, domain string, keys *Keys) ([]URLForward,error) {
	result := struct {
		Status		string			`json:"status"`
		Forwards	[]URLForward	`json:"forwards"`
	}{}

	req := c.resty.NewRequest().SetContext(ctx).SetResult(&result)
	u := fmt.Sprintf("/domain/getUrlForwarding/%v", domain)

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