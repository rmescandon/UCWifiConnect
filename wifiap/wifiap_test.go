// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2017 Canonical Ltd
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License version 3 as
// published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package wifiap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

// Testing Show()
type mockTransportShow struct{}

func (mock *mockTransportShow) Do(req *http.Request) (*http.Response, error) {

	url := req.URL.String()
	if url != "http://unix/v1/configuration" {
		return nil, fmt.Errorf("Not valid request URL: %v", url)
	}

	if req.Method != "GET" {
		return nil, fmt.Errorf("Method is not valid. Expected GET, got %v", req.Method)
	}

	rawBody := `{"result":{
		"debug":false, 
		"dhcp.lease-time": "12h", 
		"dhcp.range-start": "10.0.60.2", 
		"dhcp.range-stop": "10.0.60.199", 
		"disabled": true, 
		"share.disabled": false, 
		"share-network-interface": "tun0", 
		"wifi-address": "10.0.60.1", 
		"wifi.channel": "6", 
		"wifi.hostapd-driver": "nl80211", 
		"wifi.interface": "wlan0", 
		"wifi.interface-mode": "direct", 
		"wifi.netmask": "255.255.255.0", 
		"wifi.operation-mode": "g", 
		"wifi.security": "wpa2", 
		"wifi.security-passphrase": "passphrase123", 
		"wifi.ssid": "AP"},"status":"OK","status-code":200,"type":"sync"}`

	response := http.Response{
		StatusCode: 200,
		Status:     "200 OK",
		Body:       ioutil.NopCloser(strings.NewReader(rawBody)),
	}

	return &response, nil
}

func TestShow(t *testing.T) {
	client := NewClient(&mockTransportShow{})
	response, err := client.Show()
	if err != nil {
		t.Errorf("Failed to show current config: %v", err)
	}

	if len(response) != 17 {
		t.Error("Got different number of response objects in map than expected")
	}

	if response["debug"] != false {
		t.Errorf("'debug' value is not valid")
	}

	if response["dhcp.lease-time"] != "12h" {
		t.Errorf("'dhcp.lease-time' value is not valid")
	}

	if response["dhcp.range-start"] != "10.0.60.2" {
		t.Errorf("'dhcp.range-start' value is not valid")
	}

	if response["dhcp.range-stop"] != "10.0.60.199" {
		t.Errorf("'dhcp.range-stop' value is not valid")
	}

	if response["disabled"] != true {
		t.Errorf("'disabled' value is not valid")
	}

	if response["share.disabled"] != false {
		t.Errorf("'share.disabled' value is not valid")
	}

	if response["share-network-interface"] != "tun0" {
		t.Errorf("'share-network-interface' value is not valid")
	}

	if response["wifi-address"] != "10.0.60.1" {
		t.Errorf("'wifi-address' value is not valid")
	}

	if response["wifi.channel"] != "6" {
		t.Errorf("'wifi.channel' value is not valid")
	}

	if response["wifi.hostapd-driver"] != "nl80211" {
		t.Errorf("'wifi.hostapd-driver' value is not valid")
	}

	if response["wifi.interface"] != "wlan0" {
		t.Errorf("'wifi.interface' value is not valid")
	}

	if response["wifi.interface-mode"] != "direct" {
		t.Errorf("'wifi-interface-mode' value is not valid")
	}

	if response["wifi.netmask"] != "255.255.255.0" {
		t.Errorf("'wifi.netmask' value is not valid")
	}

	if response["wifi.operation-mode"] != "g" {
		t.Errorf("'wifi.operation-mode' value is not valid")
	}

	if response["wifi.security"] != "wpa2" {
		t.Errorf("'wifi.security' value is not valid")
	}

	if response["wifi.security-passphrase"] != "passphrase123" {
		t.Errorf("'wifi.security-passphrase' value is not valid")
	}

	if response["wifi.ssid"] != "AP" {
		t.Errorf("'wifi.ssid' value is not valid")
	}
}

// Testing Enable()
func validateHeaders(m map[string]string, req *http.Request) error {
	buf, _ := ioutil.ReadAll(req.Body)
	var headers map[string]interface{}
	if err := json.Unmarshal(buf, &headers); err != nil {
		return fmt.Errorf("Error reading request headers: %v", err)
	}

	n := len(m)
	if len(headers) != n {
		return fmt.Errorf("Expected %v headers", n)
	}

	for key, value := range m {
		if headers[key] != value {
			return fmt.Errorf("Header '%v' has not valid value", key)
		}
	}

	return nil
}

type mockTransportEnable struct{}

func (mock *mockTransportEnable) Do(req *http.Request) (*http.Response, error) {

	url := req.URL.String()

	switch url {
	case "http://unix/v1/configuration":
		if req.Method != "POST" {
			return nil, fmt.Errorf("Method is not valid. Expected POST, got %v", req.Method)
		}

		err := validateHeaders(map[string]string{"disabled": "false"}, req)
		if err != nil {
			return nil, err
		}

		rawBody := `{"result":{},"status":"OK","status-code":200,"type":"sync"}`

		response := http.Response{
			StatusCode: 200,
			Status:     "200 OK",
			Body:       ioutil.NopCloser(strings.NewReader(rawBody)),
		}

		return &response, nil

	case "http://unix/v1/status":
		if req.Method != "GET" {
			return nil, fmt.Errorf("Method is not valid. Expected GET, got %v", req.Method)
		}

		rawBody := `{"result":{"ap.active": true},"status":"OK","status-code":200,"type":"sync"}`

		response := http.Response{
			StatusCode: 200,
			Status:     "200 OK",
			Body:       ioutil.NopCloser(strings.NewReader(rawBody)),
		}

		return &response, nil

	default:
		return nil, fmt.Errorf("Not valid request URL: %v", url)
	}
}

func TestEnable(t *testing.T) {
	client := NewClient(&mockTransportEnable{})
	err := client.Enable()
	if err != nil {
		t.Errorf("Failed to enable ap: %v\n", err)
	}
}

// Testing Disable()
type mockTransportDisable struct{}

func (mock *mockTransportDisable) Do(req *http.Request) (*http.Response, error) {

	url := req.URL.String()

	switch url {
	case "http://unix/v1/configuration":

		if req.Method != "POST" {
			return nil, fmt.Errorf("Method is not valid. Expected POST, got %v", req.Method)
		}

		err := validateHeaders(map[string]string{"disabled": "true"}, req)
		if err != nil {
			return nil, err
		}

		rawBody := `{"result":{},"status":"OK","status-code":200,"type":"sync"}`

		response := http.Response{
			StatusCode: 200,
			Status:     "200 OK",
			Body:       ioutil.NopCloser(strings.NewReader(rawBody)),
		}

		return &response, nil

	case "http://unix/v1/status":

		if req.Method != "GET" {
			return nil, fmt.Errorf("Method is not valid. Expected GET, got %v", req.Method)
		}

		rawBody := `{"result":{"ap.active": false},"status":"OK","status-code":200,"type":"sync"}`

		response := http.Response{
			StatusCode: 200,
			Status:     "200 OK",
			Body:       ioutil.NopCloser(strings.NewReader(rawBody)),
		}

		return &response, nil

	default:
		return nil, fmt.Errorf("Not valid request URL: %v", url)
	}
}

func TestDisable(t *testing.T) {
	client := NewClient(&mockTransportDisable{})
	err := client.Disable()
	if err != nil {
		t.Errorf("Failed to disable ap: %v\n", err)
	}
}

type mockTransportEnabledTrue struct{}

func (mock *mockTransportEnabledTrue) Do(req *http.Request) (*http.Response, error) {

	url := req.URL.String()
	if url != "http://unix/v1/configuration" {
		return nil, fmt.Errorf("Not valid request URL: %v", url)
	}

	if req.Method != "GET" {
		return nil, fmt.Errorf("Method is not valid. Expected GET, got %v", req.Method)
	}

	rawBody := `{"result":{
		"disabled": false},"status":"OK","status-code":200,"type":"sync"}`

	response := http.Response{
		StatusCode: 200,
		Status:     "200 OK",
		Body:       ioutil.NopCloser(strings.NewReader(rawBody)),
	}
	return &response, nil
}

type mockTransportEnabledFalse struct{}

func (mock *mockTransportEnabledFalse) Do(req *http.Request) (*http.Response, error) {

	url := req.URL.String()
	if url != "http://unix/v1/configuration" {
		return nil, fmt.Errorf("Not valid request URL: %v", url)
	}

	if req.Method != "GET" {
		return nil, fmt.Errorf("Method is not valid. Expected GET, got %v", req.Method)
	}

	rawBody := `{"result":{
		"disabled": true},"status":"OK","status-code":200,"type":"sync"}`

	response := http.Response{
		StatusCode: 200,
		Status:     "200 OK",
		Body:       ioutil.NopCloser(strings.NewReader(rawBody)),
	}
	return &response, nil
}

func TestEnabled(t *testing.T) {
	client := NewClient(&mockTransportEnabledTrue{})
	res, err := client.Enabled()
	if err != nil {
		t.Errorf("Enabled() (expecting true response) has error: %v", err)
	}
	if !res {
		t.Errorf("Enabled() incorrectly returned false")
	}
	client = NewClient(&mockTransportEnabledFalse{})
	res, err = client.Enabled()
	if err != nil {
		t.Errorf("Enabled() (expecting false response) has error: %v", err)
	}
	if res {
		t.Errorf("Enabled() incorrectly returned true")
	}
}

// Testing SetSsid(ssid)
type mockTransportSetSsid struct{}

func (mock *mockTransportSetSsid) Do(req *http.Request) (*http.Response, error) {

	url := req.URL.String()
	if url != "http://unix/v1/configuration" {
		return nil, fmt.Errorf("Not valid request URL: %v", url)
	}

	if req.Method != "POST" {
		return nil, fmt.Errorf("Method is not valid. Expected POST, got %v", req.Method)
	}

	err := validateHeaders(map[string]string{"wifi.ssid": "MySsid"}, req)
	if err != nil {
		return nil, err
	}

	rawBody := `{"result":{},"status":"OK","status-code":200,"type":"sync"}`

	response := http.Response{
		StatusCode: 200,
		Status:     "200 OK",
		Body:       ioutil.NopCloser(strings.NewReader(rawBody)),
	}

	return &response, nil
}

func TestSetSsid(t *testing.T) {
	client := NewClient(&mockTransportSetSsid{})
	err := client.SetSsid("MySsid")
	if err != nil {
		t.Errorf("Failed to set ssid: %v\n", err)
	}
}

// Testing SetPassphrase(passphrase)
type mockTransportSetPassphrase struct{}

func (mock *mockTransportSetPassphrase) Do(req *http.Request) (*http.Response, error) {

	url := req.URL.String()
	if url != "http://unix/v1/configuration" {
		return nil, fmt.Errorf("Not valid request URL: %v", url)
	}

	if req.Method != "POST" {
		return nil, fmt.Errorf("Method is not valid. Expected POST, got %v", req.Method)
	}

	err := validateHeaders(map[string]string{"wifi.security": "wpa2", "wifi.security-passphrase": "passphrase123"}, req)
	if err != nil {
		return nil, err
	}

	rawBody := `{"result":{},"status":"OK","status-code":200,"type":"sync"}`

	response := http.Response{
		StatusCode: 200,
		Status:     "200 OK",
		Body:       ioutil.NopCloser(strings.NewReader(rawBody)),
	}

	return &response, nil
}

func TestSetPassphrase(t *testing.T) {
	client := NewClient(&mockTransportSetPassphrase{})
	err := client.SetPassphrase("passphrase123")
	if err != nil {
		t.Errorf("Failed to set passphrase: %v\n", err)
	}
}
