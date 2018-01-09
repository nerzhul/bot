package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type scalewayActionResponse struct {
	Task struct {
		Description string `json:"description"`
		HrefFrom    string `json:"href_from"`
		ID          string `json:"id"`
		Progress    string `json:"progress"`
		Status      string `json:"status"`
	} `json:"task"`
}

func (r *commandRouter) handlerStartBuilder(args string, user string, channel string) *string {
	if len(gconfig.Scaleway.Token) == 0 {
		log.Errorf("Scaleway token is empty")
		return nil
	}

	if len(gconfig.Scaleway.BuildServerID) == 0 {
		log.Errorf("Scaleway build-server-id is empty")
		return nil
	}

	url := fmt.Sprintf("%s/servers/%s/action", gconfig.Scaleway.URL, gconfig.Scaleway.BuildServerID)
	startAction := fmt.Sprintf(`{"action": "poweron"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(startAction)))
	if err != nil {
		log.Errorf("HTTP request error: %v", err)
		return nil
	}

	// Add token
	req.Header.Set("X-Auth-Token", gconfig.Scaleway.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("Failed to read body on start builder error.")
		} else {
			log.Errorf("Failed to start builder (code %d). Response was: %s",
				resp.StatusCode, body)
		}
		return nil
	}

	sar := &scalewayActionResponse{}
	if err := json.NewDecoder(resp.Body).Decode(sar); err != nil {
		log.Errorf("Failed to decode scaleway response when starting builder.")
		return nil
	}

	result := new(string)
	*result = fmt.Sprintf("Builder startup order sent. Startup is in '%s' state.", sar.Task.Progress)
	return result
}
