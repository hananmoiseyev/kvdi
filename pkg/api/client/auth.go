/*
Copyright 2020,2021 Avi Zimmerman

This file is part of kvdi.

kvdi is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

kvdi is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with kvdi.  If not, see <https://www.gnu.org/licenses/>.
*/

package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tinyzimmer/kvdi/pkg/types"
	"github.com/tinyzimmer/kvdi/pkg/util/errors"
)

// authenticate retrieves an access token for the API and starts a goroutine
// to refresh the token as needed.
func (c *Client) authenticate() error {
	loginRequest := &types.LoginRequest{
		Username: c.opts.Username,
		Password: c.opts.Password,
		State:    uuid.New().String(),
	}
	payload, err := json.Marshal(loginRequest)
	if err != nil {
		return err
	}
	res, err := c.httpClient.Post(c.getEndpoint("login"), "application/json", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return c.returnAPIError(body)
	}

	sessionResponse := &types.SessionResponse{}

	if err := json.Unmarshal(body, sessionResponse); err != nil {
		return err
	}

	if sessionResponse.State != loginRequest.State {
		return errors.New("State was malformed during authentication flow, your request might have been intercepted")
	}

	c.setAccessToken(sessionResponse.Token)

	if sessionResponse.Renewable {
		c.stopCh = make(chan struct{})
		go c.runTokenRefreshLoop(sessionResponse)
	}

	return nil
}

// runTokenRefreshLoop is used as a goroutine to request a new access token when the
// current one is about to expire.
func (c *Client) runTokenRefreshLoop(session *types.SessionResponse) {
	runIn := session.ExpiresAt - time.Now().Unix() - 10
	ticker := time.NewTicker(time.Duration(runIn) * time.Second)
	var err error

	for {
		select {
		case <-ticker.C:
			session, err = c.refreshToken()
			if err != nil {
				log.Println("Error refreshing client token, retrying in 2 seconds")
				ticker = time.NewTicker(time.Duration(2 * time.Second))
				continue
			}
			c.setAccessToken(session.Token)
			runIn = session.ExpiresAt - time.Now().Unix() - 10
			ticker = time.NewTicker(time.Duration(runIn) * time.Second)
		case <-c.stopCh:
			return
		}
	}
}

// refreshToken performs a refresh_token request and returns the response or any error.
func (c *Client) refreshToken() (*types.SessionResponse, error) {
	res, err := c.httpClient.Get(c.getEndpoint("refresh_token"))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, c.returnAPIError(body)
	}

	sessionResponse := &types.SessionResponse{}
	return sessionResponse, json.Unmarshal(body, sessionResponse)
}
