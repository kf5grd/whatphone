package everyoneapi // import "samhofi.us/x/everyoneapi"

import (
	"testing"
)

func TestNew(t *testing.T) {
	accountsid := "acountsid1234"
	authtoken := "authtoken1234"
	api := New(accountsid, authtoken)
	if api.AccountSID != accountsid {
		t.Errorf("AccountSID not properly set. Got: %s, Want: %s", api.AccountSID, accountsid)
	}
	if api.AuthToken != authtoken {
		t.Errorf("AuthToken not properly set. Got: %s, Want: %s", api.AuthToken, authtoken)
	}
}

func TestSingleField(t *testing.T) {
	name := "Michael Seaver"
	expandedname := ExpandedName{First: "Michael", Last: "Seaver"}

	api := New("test", "test")
	res, err := api.Lookup("+15551234567", WithName())
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if *res.Data.Name != name {
		t.Errorf("Error: Unexpected name field value. Got: %v, Want: %v", res.Data.Name, &name)
	}
	if res.Data.Address != nil {
		t.Errorf("Error: Unexpected address field value. Got: %v, Want: %v", res.Data.Address, nil)
	}
	if res.Data.Carrier != nil {
		t.Errorf("Error: Unexpected carrier field value. Got: %v, Want: %v", res.Data.Carrier, nil)
	}
	if res.Data.CarrierO != nil {
		t.Errorf("Error: Unexpected carrier_o field value. Got: %v, Want: %v", res.Data.CarrierO, nil)
	}
	if res.Data.Cnam != nil {
		t.Errorf("Error: Unexpected cnam field value. Got: %v, Want: %v", res.Data.Cnam, nil)
	}
	if *res.Data.ExpandedName != expandedname {
		t.Errorf("Error: Unexpected expanded_name field value. Got: %v, Want: %v", res.Data.ExpandedName, &expandedname)
	}
	if res.Data.Gender != nil {
		t.Errorf("Error: Unexpected gender field value. Got: %v, Want: %v", res.Data.Gender, nil)
	}
	if res.Data.Image != nil {
		t.Errorf("Error: Unexpected image field value. Got: %v, Want: %v", res.Data.Image, nil)
	}
	if res.Data.LineProvider != nil {
		t.Errorf("Error: Unexpected line_provider field value. Got: %v, Want: %v", res.Data.LineProvider, nil)
	}
	if res.Data.Location != nil {
		t.Errorf("Error: Unexpected location field value. Got: %v, Want: %v", res.Data.Location, nil)
	}
	if res.Data.Profile != nil {
		t.Errorf("Error: Unexpected profile field value. Got: %v, Want: %v", res.Data.Profile, nil)
	}
}

func TestNoField(t *testing.T) {
	api := New("test", "test")
	res, err := api.Lookup("+15551234567")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if res.Data.Name == nil {
		t.Errorf("Error: name field is nil but should not be")
	}
	if res.Data.Address == nil {
		t.Errorf("Error: address field is nil but should not be")
	}
	if res.Data.Carrier == nil {
		t.Errorf("Error: carrier field is nil but should not be")
	}
	if res.Data.CarrierO == nil {
		t.Errorf("Error: carrier_o field is nil but should not be")
	}
	if res.Data.Cnam == nil {
		t.Errorf("Error: cnam field is nil but should not be")
	}
	if res.Data.ExpandedName == nil {
		t.Errorf("Error: expanded_name field is nil but should not be")
	}
	if res.Data.Gender == nil {
		t.Errorf("Error: gender field is nil but should not be")
	}
	if res.Data.Image == nil {
		t.Errorf("Error: image field is nil but should not be")
	}
	if res.Data.LineProvider == nil {
		t.Errorf("Error: line_provider field is nil but should not be")
	}
	if res.Data.Location == nil {
		t.Errorf("Error: location field is nil but should not be")
	}
	if res.Data.Profile == nil {
		t.Errorf("Error: pected profile field is nil but should not be")
	}
}
