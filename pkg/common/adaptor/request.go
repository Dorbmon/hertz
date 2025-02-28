/*
 * Copyright 2022 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package adaptor

import (
	"bytes"
	"net/http"

	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/henrylee2cn/goutil"
)

// GetCompatRequest only support basic function of Request, not for all.
func GetCompatRequest(req *protocol.Request) (*http.Request, error) {
	r, err := http.NewRequest(goutil.BytesToString(req.Method()), goutil.BytesToString(req.URI().FullURI()), bytes.NewReader(req.Body()))
	if err != nil {
		return r, err
	}

	h := make(map[string][]string)
	req.Header.VisitAll(func(k, v []byte) {
		h[goutil.BytesToString(k)] = append(h[goutil.BytesToString(k)], goutil.BytesToString(v))
	})

	r.Header = h
	return r, nil
}

// CopyToHertzRequest copy uri, host, method, protocol, header, but share body reader from http.Request to protocol.Request.
func CopyToHertzRequest(req *http.Request, hreq *protocol.Request) error {
	hreq.Header.SetRequestURI(req.RequestURI)
	hreq.Header.SetHost(req.Host)
	hreq.Header.SetMethod(req.Method)
	hreq.Header.SetProtocol(req.Proto)
	for k, v := range req.Header {
		for _, vv := range v {
			hreq.Header.Add(k, vv)
		}
	}
	if req.Body != nil {
		hreq.SetBodyStream(req.Body, hreq.Header.ContentLength())
	}
	return nil
}
