/*
Copyright 2021 TriggerMesh Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"log"
	"time"

	pdk "github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
)

const Version = "0.0.1"
const Priority = 1

var eventType = "com.konghq.ce-plugin"

type Config struct {
	EventType        string `json:"eventType"`
	DiscardCEContext bool   `json:"discardCEContext"`
}

func main() {
	server.StartServer(New, Version, Priority)
}

func New() interface{} {
	return &Config{}
}

func (conf *Config) Access(kong *pdk.PDK) {
	source, err := kong.Request.GetHeader("User-Agent")
	if err != nil {
		log.Printf("cannot read user agent header: %v", err)
	}
	if source == "" {
		source = "kong"
	}
	if conf.EventType != "" {
		eventType = conf.EventType
	}
	kong.ServiceRequest.SetHeader("ce-id", uuid.New().String())
	kong.ServiceRequest.SetHeader("ce-type", eventType)
	kong.ServiceRequest.SetHeader("ce-time", time.Now().Format(time.RFC3339))
	kong.ServiceRequest.SetHeader("ce-source", source)
	kong.ServiceRequest.SetHeader("ce-specversion", cloudevents.VersionV1)
}

func (conf *Config) Response(kong *pdk.PDK) {
	if !conf.DiscardCEContext {
		return
	}

	headers, err := kong.ServiceResponse.GetHeaders(-1)
	if err != nil {
	}

	body, err := kong.ServiceResponse.GetRawBody()
	if err != nil {
	}

	switch {
	case isBinaryCE(headers):
		removeCEHeaders()
	case isStructuredCE(body):
		extractData()
	}
}

func isBinaryCE(map[string][]string) bool {
	return false
}

func isStructuredCE(body string) bool {
	return false
}

func removeCEHeaders() {}

func extractData() {}
