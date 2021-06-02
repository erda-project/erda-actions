// Copyright (c) 2021 Terminus, Inc.
//
// This program is free software: you can use, redistribute, and/or modify
// it under the terms of the GNU Affero General Public License, version 3
// or later ("AGPL"), as published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/erda-project/erda/apistructs"
	"github.com/pkg/errors"
)

const (
	loginAPI       = "/login"
	extenstionsAPI = "/api/extensions"
)

func New(host, username, password string) (*Client, error) {
	if !strings.HasPrefix(host, "http://") || !strings.HasPrefix(host, "https://") {
		host = "https://" + host
	}

	return &Client{
		Host:     host,
		Username: username,
		Password: password,
		client: func() *http.Client {
			return &http.Client{Timeout: time.Second * time.Duration(180)}
		},
		status: nil,
	}, nil
}

type Client struct {
	Host     string
	Username string
	Password string

	status *StatusInfo
	client func() *http.Client
}

func (c *Client) Logged() bool {
	return c.status != nil
}

func (c *Client) Login() error {
	parameters := make(url.Values)
	parameters.Set("username", c.Username)
	parameters.Set("password", c.Password)

	body := strings.NewReader(parameters.Encode())
	request, err := http.NewRequest(http.MethodPost, c.Host+loginAPI, body)
	if err != nil {
		return errors.Wrap(err, "failed to NewRequest")
	}
	request.Header.Set("content-type", "application/x-www-form-urlencoded")

	response, err := c.client().Do(request)
	if err != nil {
		return errors.Wrap(err, "failed to Do request")
	}
	defer response.Body.Close()

	var status StatusInfo
	if err = json.NewDecoder(response.Body).Decode(&status); err != nil {
		return errors.Wrapf(err, "failed to parse login response body, request: %v+", request)
	}
	expiredAt := time.Now().Add(time.Hour * 12)
	status.ExpiredAt = &expiredAt

	c.status = &status

	return nil
}

func (c *Client) Push(payload *apistructs.ExtensionVersionCreateRequest) error {
	if payload == nil {
		return errors.New("payload is nil")
	}
	if !c.Logged() {
		if err := c.Login(); err != nil {
			return err
		}
	}

	var body = bytes.NewBuffer(nil)
	if err := json.NewEncoder(body).Encode(payload); err != nil {
		return errors.Wrap(err, "failed to Encode payload")
	}

	uri := c.Host + filepath.Join(extenstionsAPI, payload.Name)
	request, err := http.NewRequest(http.MethodPost, uri, body)
	if err != nil {
		return errors.Wrapf(err, "failed to NewRequest, url: %s", uri)
	}
	request.Header.Set("use-token", "true")
	request.AddCookie(&http.Cookie{Name: "OPENAPISESSION", Value: c.status.SessionID})

	response, err := c.client().Do(request)
	if err != nil {
		return errors.Wrapf(err, "failed to Do request, request: %+v", request)
	}
	defer response.Body.Close()

	var (
		resp apistructs.ExtensionVersionCreateResponse
	)
	if err = json.NewDecoder(response.Body).Decode(&resp); err != nil {
		return errors.Wrap(err, "failed to Decode")
	}

	return nil
}
