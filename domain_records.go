package porkbungo

import (
	"context"
	"errors"
	"fmt"
	"log"
)

type RecordType string

const (
	A     RecordType = "A"
	MX    RecordType = "MX"
	CNAME RecordType = "CNAME"
	ALIAS RecordType = "ALIAS"
	TXT   RecordType = "TXT"
	NS    RecordType = "NS"
	AAAA  RecordType = "AAAA"
	SRV   RecordType = "SRV"
	TLSA  RecordType = "TLSA"
	CAA   RecordType = "CAA"
)

type RecordOptions struct {
	*Keys
	// ID of DNS record created
	ID	string	`json:"-"`
	// Second level domain
	Domain	 string	`json:"-"`
	// The subdomain for the record being created or updated, not including the domain itself.
	// Leave blank to create/update record on root. Use * to create/update wildcard
	Name string `json:"name,omitempty"`
	// Type of record being created/updated
	Type RecordType `json:"type,omitempty"`
	// Answer content for the record.
	Content  string `json:"content,omitempty"`
	// The time to live for the record. The minimum allowed and default
	// is 600 seconds
	TTL      int   `json:"ttl,omitempty"`
	// Priority of the record, for the record types that support it
	Priority int   `json:"prio,omitempty"`
}

type Record struct {
	ID			string		`json:"id"`
	Name		string		`json:"name"`
	Type		RecordType	`json:"type"`
	Content		string		`json:"content"`
	TTL			string		`json:"ttl"`
	Priority	string		`json:"prio"`
	Notes		string		`json:"notes"`
}

// Create a DNS record
func (c *Client) CreateDNSRecord(ctx context.Context, opts RecordOptions) (int,error) {
	result := struct {
		Status	string	`json:"status"`
		Id		int		`json:"id"`
	}{}

	req := c.resty.NewRequest().SetContext(ctx).SetResult(&result)
	u := fmt.Sprintf("/dns/create/%s",opts.Domain)
	
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
		return 0,err
	}

	return result.Id,nil
}

// Edit a DNS record using the Second level domain and the ID of the dns record
func (c *Client) EditRecordByID(ctx context.Context, opts RecordOptions) error {
	req := c.resty.NewRequest().SetContext(ctx)
	u := fmt.Sprintf("/dns/edit/%s/%s",opts.Domain, opts.ID)

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

// Edit all records for the domain that match a specific subdomain and type
func (c *Client) EditRecordsByNameAndType(ctx context.Context, opts RecordOptions) error {
	req := c.resty.NewRequest().SetContext(ctx)
	u := fmt.Sprintf("/dns/editByNameType/%s/%s/%s", opts.Domain,opts.Type,opts.Name)

	if opts.Keys == nil {
		opts.Keys = c.keys
	}

	// Remove values to prevent them from being serialized
	opts.Type = ""
	opts.Name = ""
	
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

func (c *Client) DeleteRecordByID(ctx context.Context, domain string, id int, keys *Keys) error {
	req := c.resty.NewRequest().SetContext(ctx)
	u := fmt.Sprintf("/dns/delete/%s/%v",domain,id)

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

func (c *Client) DeleteRecordsByNameAndType(ctx context.Context, opts RecordOptions) error {
	req := c.resty.NewRequest().SetContext(ctx)
	u := fmt.Sprintf("/dns/deleteByNameType/%s/%s",opts.Domain,opts.Type)

	if opts.Name != "" {
		u += fmt.Sprintf("/%s",opts.Name)
		opts.Name = ""
	}

	if opts.Keys == nil {
		opts.Keys = c.keys
	}

	// Remove values to prevent them from being serialized
	opts.Type = ""

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

func (c *Client) GetAllRecords(ctx context.Context, domain string, keys *Keys) ([]Record,error) {
	result := struct {
		Status	string	 `json:"status"`
		Records []Record `json:"records"`
	}{}

	req := c.resty.NewRequest().SetContext(ctx).SetResult(&result)
	u := fmt.Sprintf("/dns/retrieve/%s",domain)

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

	return result.Records,nil
}

func (c *Client) GetRecordByID(ctx context.Context, domain string, id int, keys *Keys) (*Record,error) {
	result := struct {
		Status	string	 `json:"status"`
		Records []Record `json:"records"`
	}{}

	req := c.resty.NewRequest().SetContext(ctx).SetResult(&result)
	u := fmt.Sprintf("/dns/retrieve/%s/%v",domain,id)

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

	if len(result.Records) <= 0 {
		return nil,errors.New("record not found")
	}
	
	return &result.Records[0],nil
}

func (c *Client) GetRecordsByNameAndType(ctx context.Context, opts RecordOptions) ([]Record,error) {
	result := struct {
		Status	string	 `json:"status"`
		Records []Record `json:"records"`
	}{}

	req := c.resty.NewRequest().SetContext(ctx).SetResult(&result)
	u := fmt.Sprintf("/dns/retrieveByNameType/%s/%s",opts.Domain,opts.Type)

	if opts.Name != "" {
		u += fmt.Sprintf("/%s",opts.Name)
		opts.Name = ""
	}

	if opts.Keys == nil {
		opts.Keys = c.keys
	}
	req.SetBody(opts.Keys)

	resp,err := req.Post(u)
	if err != nil || resp.StatusCode() != 200 {
		if err == nil {
			e,_ := resp.Error().(*APIError)
			err = errors.New(e.Message)
		}
		return nil,err
	}

	log.Println()

	return result.Records,nil
}