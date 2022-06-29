package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const PINNING_URL = "https://api.pinata.cloud/pinning/pinByHash"

type PinningClientRespErr struct {
	err error
	msg string
}

func newPinningClientRespErr(err error, msg string) *PinningClientRespErr {
	return &PinningClientRespErr{
		err: err,
		msg: msg,
	}
}

func (e *PinningClientRespErr) Error() string {
	s := e.msg
	if e.err != nil {
		s = fmt.Sprintf("%s %s", e.err.Error(), s)
	}

	return s
}

type PinningClient struct {
	httpClient *http.Client
	JWT        string
}

func newPinningClient(jwt string) *PinningClient {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	return &PinningClient{
		httpClient: httpClient,
		JWT:        jwt,
	}
}

func (p *PinningClient) Pin(cid string) error {
	payload := strings.NewReader(fmt.Sprintf(`{
		"hashToPin": "%s"
	}`, cid))

	req, err := http.NewRequest("POST", PINNING_URL, payload)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", p.JWT))
	req.Header.Add("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return newPinningClientRespErr(err, fmt.Sprintf("response: %s", body))
	}

	if errMsg, ok := data["error"]; ok {
		return newPinningClientRespErr(nil, fmt.Sprintf("%s", errMsg))
	}

	return nil
}
