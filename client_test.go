package porkbungo

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

var keys = &Keys{
	APIKey: "TestAPIKey",
	SecretAPIKey: "TestSecretAPIKey",
}

func TestClient_New(t *testing.T) {
	expectedHost := "https://porkbun.com/api/json/v3"

	client := NewClient()
	
	if client.useIPv4 != false {
		t.Error(cmp.Diff(client.useIPv4,false))
	}

	if client.keys != nil {
		t.Error(cmp.Diff(client.keys,nil))
	}

	if client.resty.BaseURL != expectedHost {
		t.Fatal(cmp.Diff(client.resty.BaseURL,expectedHost))
	}
}

func TestClient_NewWithAPIKeys(t *testing.T) {

	expectedHost := "https://porkbun.com/api/json/v3"

	client := NewClientWithAPIKeys(keys)
	
	if client.keys != keys {
		t.Errorf(cmp.Diff(client.keys,keys))
	}

	if client.useIPv4 != false {
		t.Error(cmp.Diff(client.useIPv4,false))
	}

	if client.resty.BaseURL != expectedHost {
		t.Fatal(cmp.Diff(client.resty.BaseURL,expectedHost))
	}
}

func TestClient_SetAPIKeys(t *testing.T) {
	client := NewClient()

	client.SetAPIKeys(keys)

	if client.keys != keys {
		t.Errorf(cmp.Diff(client.keys,keys))
	}
}

func TestClient_UseIPv4(t *testing.T) {
	client := NewClient()

	ipv4host := "https://api-ipv4.porkbun.com/api/json/v3"
	apihost := "https://porkbun.com/api/json/v3"

	client.UseIPv4(true)
	if client.useIPv4 != true {
		t.Fatal(cmp.Diff(client.useIPv4,true))
	}

	if client.resty.BaseURL != ipv4host {
		t.Fatal(cmp.Diff(client.resty.BaseURL,ipv4host))
	}

	client.UseIPv4(false)
	if client.useIPv4 != false {
		t.Fatal(cmp.Diff(client.useIPv4,false))
	}

	if client.resty.BaseURL != apihost {
		t.Fatal(cmp.Diff(client.resty.BaseURL,apihost))
	}
}