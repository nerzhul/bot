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
		Progress    int    `json:"progress"`
		Status      string `json:"status"`
	} `json:"task"`
}

type scalewayErrorResponse struct {
	Message string `json:"message"`
	Type    string `json:"type"`
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

	action := ""
	if args == "start" {
		action = "poweron"
	} else if args == "stop" {
		action = "poweroff"
	} else {
		log.Error("Invalid builder action received, ignoring.")
		result := new(string)
		*result = "Invalid builder action, action ignored. Valid actions are: start, stop."
		return result
	}

	url := fmt.Sprintf("%s/servers/%s/action", gconfig.Scaleway.URL, gconfig.Scaleway.BuildServerID)
	startAction := fmt.Sprintf(`{"action": "%s"}`, action)
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Failed to read body on modify builder state error.")
		return nil
	}

	if resp.StatusCode == http.StatusBadRequest {
		ser := &scalewayErrorResponse{}
		if err := json.Unmarshal(body, &ser); err != nil {
			log.Errorf("Failed to modify builder state (code %d). Response was: %s", body)
			return nil
		}

		result := new(string)
		*result = "Unable to start builder. Error: " + ser.Message
		return result
	}

	if resp.StatusCode != http.StatusAccepted {
		log.Errorf("Failed to modify builder state (code %d). Response was: %s",
			resp.StatusCode, body)
		return nil
	}

	sar := &scalewayActionResponse{}
	if err := json.Unmarshal(body, &sar); err != nil {
		log.Errorf("Failed to modify builder state (code %d). Response was: %s",
			resp.StatusCode, body)
	}

	result := new(string)
	*result = fmt.Sprintf("Builder state modification order sent. Modification is in '%s' state.", sar.Task.Status)
	return result
}
