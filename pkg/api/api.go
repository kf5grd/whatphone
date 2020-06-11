package whatphone // import "samhofi.us/x/whatphone/pkg/api"

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	baseurl = "https://api.everyoneapi.com/v1/phone/"
)

// New returns a new API object
func New(accountsid string, authtoken string) *API {
	return &API{
		AccountSID: accountsid,
		AuthToken:  authtoken,
	}
}

// Lookup performs a phone number lookup and returns the Result
func (a *API) Lookup(phonenumber string, opts ...Option) (*Result, error) {
	f := new(fields)
	for _, opt := range opts {
		opt(f)
	}

	var data string
	if len(*f) > 0 {
		data = fmt.Sprintf("?data=%s", strings.Join(*f, ","))
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, baseurl+phonenumber+data, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(a.AccountSID, a.AuthToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s", resp.Status)
	}

	var ret Result
	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, err
	}

	return &ret, nil
}
