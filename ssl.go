package porkbungo

import (
	"context"
	"errors"
	"fmt"
)

type SSLBundleOptions struct {
	*Keys
	Domain string `json:"-"` // This field is not part of the request body
}

type SSLBundle struct {
	Status				string	`json:"status"`
	IntermCert			string	`json:"intermediatecertificate"`
	CertChain			string	`json:"certificatechain"`
	PrivateKey			string	`json:"privatekey"`
	PublicKey			string	`json:"publickey"`
}

func (c *Client) GetSSLBundle(ctx context.Context, opts SSLBundleOptions) (*SSLBundle,error)  {
	req := c.resty.NewRequest().SetContext(ctx).SetResult(&SSLBundle{})
	u := fmt.Sprintf("/ssl/retrieve/%s",opts.Domain)

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
		return nil,err
	}
	
	if _,ok := resp.Result().(*SSLBundle); !ok {
		return nil,errors.New("error data recieved is not in correct format")
	}

	return resp.Result().(*SSLBundle),nil
}