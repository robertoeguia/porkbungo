package porkbungo

import (
	"context"
	"errors"
)

type DomainPricing struct {
	Status  string `json:"status"`
	Pricing map[string] struct {
		Registration	string	`json:"registration"`
		Renewal			string`json:"renewal"`
		Transfer 		string`json:"transfer"`
	} `json:"pricing"`
}

func (c *Client) GetDomainPricing(ctx *context.Context) (*DomainPricing, error) {
	request := c.resty.NewRequest()
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