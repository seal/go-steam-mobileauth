package mobileauth

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

var _isTimeAligned bool
var _steamTimeDifference int64

func GetSteamTime() int64 {
	if !_isTimeAligned {
		AlignTime()
	}
	return time.Now().Unix() + _steamTimeDifference
}

// Aligns system time with the Steam server time. Not super advanced;
// probably not taking some things into account that it should.
// Necessary to generate up-to-date codes. In general, this will have an error
// of less than a second, assuming Steam is operational.
func AlignTime() error {
	now := time.Now().Unix()
	client := new(http.Client)
	resp, err := client.Post(UrlTwoFactorService+"/QueryTime/v0001", "application/x-www-form-urlencoded", bytes.NewBuffer([]byte("steamid=0")))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	r := timeQueryResponse{}
	if err = json.Unmarshal(respBody, &r); err != nil {
		return err
	}
	if r.Response == nil {
		return errors.New("steam returned empty time query response")
	}
	_steamTimeDifference = r.Response.ServerTime - now
	_isTimeAligned = true
	return nil
}

type timeQueryResponse struct {
	Response *timeQueryResult `json:"response"`
}

type timeQueryResult struct {
	ServerTime int64 `json:"server_time,string"`
}
