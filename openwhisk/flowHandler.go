/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package openwhisk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func flow(ap *ActionProxy, body []byte) {
	// execute the action
	DebugLimit("flow", body, 120)
	var objmap map[string]json.RawMessage
	err := json.Unmarshal(body, &objmap)
	if err != nil {
		Debug(fmt.Sprintf("Error: %s", err))
	}
	DebugLimit("value", objmap["value"], 120)
	DebugLimit("workflow", objmap["workflow"], 120)

	response, err := ap.theExecutor.Interact(body)
	// check for early termination
	if err != nil {
		Debug("WARNING! Command exited")
		ap.theExecutor = nil
		return
	}
	DebugLimit("received", response, 120)

	// check if the answer is an object map
	var answermap map[string]json.RawMessage
	err = json.Unmarshal(response, &answermap)
	if err != nil {
		Debug("WARNING! The action did not return a dictionary")
	}
}

func (ap *ActionProxy) flowHandler(w http.ResponseWriter, r *http.Request) {

	// parse the request
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		sendError(w, http.StatusBadRequest, fmt.Sprintf("Error reading request body: %v", err))
		return
	}
	Debug("done reading %d bytes", len(body))

	// check if you have an action
	if ap.theExecutor == nil {
		sendError(w, http.StatusInternalServerError, "no action defined yet")
		return
	}
	// check if the process exited
	if ap.theExecutor.Exited() {
		sendError(w, http.StatusInternalServerError, "command exited")
		return
	}

	// remove newlines
	body = bytes.Replace(body, []byte("\n"), []byte(""), -1)
	go flow(ap, body)
	sendOK(w)
}
