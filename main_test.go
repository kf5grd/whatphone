package main

import (
	"bytes"
	"errors"
	"testing"

	whatphone "samhofi.us/x/whatphone/pkg/api"
)

func testReadConfig() (*whatphone.API, error) {
	return &whatphone.API{
		AccountSID: "test",
		AuthToken:  "test",
	}, nil
}

func TestLookup(t *testing.T) {
	lookups := []struct {
		args     []string
		expected string
	}{
		{
			[]string{"whatphone", "lookup", "-n", "15551234567"},
			`Name: Michael Seaver
Note: THIS IS A SAMPLE, YOU WILL NOT BE CHARGED
Price Total: -0.0100
`,
		},
		{
			[]string{"whatphone", "lookup", "-na", "15551234567"},
			`Name: Michael Seaver
Address: 15 Robin Hood Lane
Location:
  City, State, Zip: Long Island, NY, 10003
  Lat, Long: 40.799787, -73.971421
Note: THIS IS A SAMPLE, YOU WILL NOT BE CHARGED
Price Total: -0.0900
`,
		},
		{
			[]string{"whatphone", "lookup", "--all", "15551234567"},
			`Name: Michael Seaver
Profile:
  Edu: Thomas Dewey High School
  Job: Custodian
  Relationship: April Lerman
CNAM: MICHAEL SEAVER
Gender: M
Image:
  Cover: //teloimg-pub.com.s3.amazonaws.com/cover.jpg
  Small: //teloimg-pub.com.s3.amazonaws.com/small.jpg
  Medium: //teloimg-pub.com.s3.amazonaws.com/med.jpg
  Large: //teloimg-pub.com.s3.amazonaws.com/large.jpg
Address: 15 Robin Hood Lane
Location:
  City, State, Zip: Long Island, NY, 10003
  Lat, Long: 40.799787, -73.971421
Line Provider:
  ID: 215
  Name: MysticVoice
  MMS E-mail: 5551234567@mms.mysticvoice.com
  SMS E-mail: 5551234567@sms.mysticvoice.com
Carrier:
  ID: 214
  Name: Growing Wireless Inc.
Original Carrier:
  ID: 213
  Name: Paine Mobile Inc.
Linetype: mobile
Note: THIS IS A SAMPLE, YOU WILL NOT BE CHARGED
Price Total: -0.1610
`,
		},
		{
			[]string{"whatphone", "lookup", "-pb", "+15551234567"},
			`Profile:
  Edu: Thomas Dewey High School
  Job: Custodian
  Relationship: April Lerman
Note: THIS IS A SAMPLE, YOU WILL NOT BE CHARGED
Price Total: -0.0050
  Name: 0.0000
  Profile: -0.0050
  CNAM: 0.0000
  Gender: 0.0000
  Image: 0.0000
  Address: 0.0000
  Location: 0.0000
  Line Provider: 0.0000
  Carrier: 0.0000
  Original Carrier: 0.0000
  Linetype: 0.0000
`,
		},
	}

	for _, lookup := range lookups {
		var stdout bytes.Buffer
		err := run(lookup.args, &stdout, newConfigReader(testReadConfig))
		if err != nil {
			t.Errorf("%v returned error: %v", lookup.args, err)
		}
		out := stdout.String()
		if out != lookup.expected {
			t.Errorf("%v returned unexpected output.\nExpected: %s\nGot: %s\n", lookup.args, lookup.expected, out)
		}
	}
}

func TestErrors(t *testing.T) {
	lookups := []struct {
		args     []string
		expected error
	}{
		{
			[]string{"whatphone", "lookup", "15551234567"},
			errors.New("no data points selected; use --all to request all data points"),
		},
		{
			[]string{"whatphone", "lookup"},
			errors.New("missing phone number"),
		},
	}

	for _, lookup := range lookups {
		var stdout bytes.Buffer
		err := run(lookup.args, &stdout, newConfigReader(testReadConfig))

		// we expect all of these to return an error
		if err == nil {
			t.Errorf("%v should have returned an error but didn't", lookup.args)
		}
		if err.Error() != lookup.expected.Error() {
			t.Errorf("%v returned unexpected error.\nExpected: %s\nGot: %s\n", lookup.args, lookup.expected.Error(), err.Error())
		}
	}
}
