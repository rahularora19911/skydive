/*
 * Copyright (C) 2016 Red Hat, Inc.
 *
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 *
 */

package tests

import (
	"testing"

	"github.com/skydive-project/skydive/api"
	shttp "github.com/skydive-project/skydive/http"
)

func TestAlertAPI(t *testing.T) {
	client, err := api.NewCrudClientFromConfig(&shttp.AuthenticationOpts{})
	if err != nil {
		t.Fatal(err.Error())
	}

	alert := api.NewAlert()
	alert.Expression = "G.V().Has('MTU', GT(1500))"
	if err := client.Create("alert", alert); err != nil {
		t.Errorf("Failed to create alert: %s", err.Error())
	}

	alert2 := api.NewAlert()
	alert2.Expression = "G.V().Has('MTU', Gt(1500))"
	if err := client.Get("alert", alert.UUID, &alert2); err != nil {
		t.Error(err)
	}

	if *alert != *alert2 {
		t.Errorf("Alert corrupted: %+v != %+v", alert, alert2)
	}

	var alerts map[string]api.Alert
	if err := client.List("alert", &alerts); err != nil {
		t.Error(err)
	} else {
		if len(alerts) != 1 {
			t.Errorf("Wrong number of alerts: got %d, expected 1", len(alerts))
		}
	}

	if alerts[alert.UUID] != *alert {
		t.Errorf("Alert corrupted: %+v != %+v", alerts[alert.UUID], alert)
	}

	if err := client.Delete("alert", alert.UUID); err != nil {
		t.Errorf("Failed to delete alert: %s", err.Error())
	}

	var alerts2 map[string]api.Alert
	if err := client.List("alert", &alerts2); err != nil {
		t.Errorf("Failed to list alerts: %s", err.Error())
	} else {
		if len(alerts2) != 0 {
			t.Errorf("Wrong number of alerts: got %d, expected 0 (%+v)", len(alerts2), alerts2)
		}
	}
}

func TestCaptureAPI(t *testing.T) {
	client, err := api.NewCrudClientFromConfig(&shttp.AuthenticationOpts{})
	if err != nil {
		t.Fatal(err.Error())
	}

	capture := api.NewCapture("G.V().Has('Name', 'br-int')", "port 80")
	if err := client.Create("capture", capture); err != nil {
		t.Fatalf("Failed to create alert: %s", err.Error())
	}

	capture2 := &api.Capture{}
	if err := client.Get("capture", capture.ID(), &capture2); err != nil {
		t.Error(err)
	}

	if *capture != *capture2 {
		t.Errorf("Capture corrupted: %+v != %+v", capture, capture2)
	}

	var captures map[string]api.Capture
	if err := client.List("capture", &captures); err != nil {
		t.Error(err)
	} else {
		if len(captures) != 1 {
			t.Errorf("Wrong number of captures: got %d, expected 1", len(captures))
		}
	}

	if captures[capture.ID()] != *capture {
		t.Errorf("Capture corrupted: %+v != %+v", captures[capture.ID()], capture)
	}

	if err := client.Delete("capture", capture.ID()); err != nil {
		t.Errorf("Failed to delete capture: %s", err.Error())
	}

	var captures2 map[string]api.Capture
	if err := client.List("capture", &captures2); err != nil {
		t.Errorf("Failed to list captures: %s", err.Error())
	} else {
		if len(captures2) != 0 {
			t.Errorf("Wrong number of captures: got %d, expected 0 (%+v)", len(captures2), captures2)
		}
	}

	if err := client.Get("capture", capture.ID(), &capture2); err == nil {
		t.Errorf("Found delete capture: %s", capture.ID())
	}
}
