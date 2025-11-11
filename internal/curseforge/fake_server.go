// Copyright (c) 2025 Aton-Kish
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package curseforge

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type FakeServer interface {
	ListenAndServe() error
	Close() error

	Up()
	Down()

	RegisterAPIToken(token string)
	RegisterProjectID(projectID int)
}

type fakeServer struct {
	*http.Server

	apiToken  string
	projectID int

	isServerDown bool
}

func NewFakeServer(port int) FakeServer {
	mux := http.NewServeMux()

	s := &fakeServer{
		Server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
	}

	s.registerGameVersionTypesHandler(mux)
	s.registerGameVersionsHandler(mux)
	s.registerProjectUploadFileHandler(mux)

	return s
}

func (s *fakeServer) Up() {
	s.isServerDown = false
}

func (s *fakeServer) Down() {
	s.isServerDown = true
}

func (s *fakeServer) RegisterAPIToken(token string) {
	s.apiToken = token
}

func (s *fakeServer) RegisterProjectID(projectID int) {
	s.projectID = projectID
}

func (s *fakeServer) registerHandler(mux *http.ServeMux, method string, path string, handler func(w http.ResponseWriter, r *http.Request)) {
	mux.HandleFunc(fmt.Sprintf("%s %s", method, path), func(w http.ResponseWriter, r *http.Request) {
		if s.isServerDown {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`Oops! Something went wrong.`))
			return
		}

		token := r.Header.Get("X-Api-Token")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"errorCode":401,"errorMessage":"You must provide an API token using the ` + "`" + `X-Api-Token` + "`" + ` header, the ` + "`" + `token` + "`" + ` query string parameter, your email address and an API token using HTTP basic authentication."}`))
			return
		}

		if _, err := uuid.Parse(token); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"errorCode":3,"errorMessage":"API token is malformed. Token provided: ` + token + `"}`))
			return
		}

		if token != s.apiToken {
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"errorCode":403,"errorMessage":"Invalid API token provided."}`))
			return
		}

		handler(w, r)
	})
}

func (s *fakeServer) registerGameVersionTypesHandler(mux *http.ServeMux) {
	s.registerHandler(mux, http.MethodGet, "/api/game/version-types", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"id":2,"name":"Java","slug":"java"},{"id":615,"name":"Addons","slug":"addons"},{"id":68441,"name":"Modloader","slug":"modloader"},{"id":75125,"name":"Minecraft 1.20","slug":"minecraft-1-20"},{"id":75208,"name":"Environment","slug":"environment"},{"id":77784,"name":"Minecraft 1.21","slug":"minecraft-1-21"}]`))
	})
}

func (s *fakeServer) registerGameVersionsHandler(mux *http.ServeMux) {
	s.registerHandler(mux, http.MethodGet, "/api/game/versions", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"id":7498,"gameVersionTypeID":68441,"name":"Forge","slug":"forge","apiVersion":""},{"id":7499,"gameVersionTypeID":68441,"name":"Fabric","slug":"fabric","apiVersion":""},{"id":9638,"gameVersionTypeID":75208,"name":"Client","slug":"client","apiVersion":""},{"id":9639,"gameVersionTypeID":75208,"name":"Server","slug":"server","apiVersion":""},{"id":11135,"gameVersionTypeID":2,"name":"Java 21","slug":"java-21","apiVersion":""},{"id":9970,"gameVersionTypeID":615,"name":"1.20","slug":"1-20","apiVersion":""},{"id":9971,"gameVersionTypeID":75125,"name":"1.20","slug":"1-20","apiVersion":""},{"id":11457,"gameVersionTypeID":77784,"name":"1.21","slug":"1-21","apiVersion":""},{"id":11458,"gameVersionTypeID":615,"name":"1.21","slug":"1-21","apiVersion":""},{"id":11515,"gameVersionTypeID":1,"name":"1.21","slug":"1-21","apiVersion":""},{"id":11779,"gameVersionTypeID":77784,"name":"1.21.1","slug":"1-21-1","apiVersion":null},{"id":12079,"gameVersionTypeID":77784,"name":"1.21.2","slug":"1-21-2","apiVersion":""},{"id":12084,"gameVersionTypeID":77784,"name":"1.21.3","slug":"1-21-3","apiVersion":null},{"id":12281,"gameVersionTypeID":77784,"name":"1.21.4","slug":"1-21-4","apiVersion":null},{"id":12735,"gameVersionTypeID":1,"name":"1.21.1","slug":"1-21-1","apiVersion":""},{"id":12736,"gameVersionTypeID":1,"name":"1.21.2","slug":"1-21-2","apiVersion":""},{"id":12737,"gameVersionTypeID":1,"name":"1.21.3","slug":"1-21-3","apiVersion":""},{"id":12738,"gameVersionTypeID":1,"name":"1.21.4","slug":"1-21-4","apiVersion":""},{"id":12934,"gameVersionTypeID":77784,"name":"1.21.5","slug":"1-21-5","apiVersion":null},{"id":12988,"gameVersionTypeID":1,"name":"1.21.5","slug":"1-21-5","apiVersion":""},{"id":13422,"gameVersionTypeID":77784,"name":"1.21.6","slug":"1-21-6","apiVersion":null},{"id":13473,"gameVersionTypeID":1,"name":"1.21.6","slug":"1-21-6","apiVersion":""},{"id":13506,"gameVersionTypeID":77784,"name":"1.21.7","slug":"1-21-7","apiVersion":null},{"id":13574,"gameVersionTypeID":1,"name":"1.21.7","slug":"1-21-7","apiVersion":""},{"id":13620,"gameVersionTypeID":77784,"name":"1.21.8","slug":"1-21-8","apiVersion":null},{"id":13683,"gameVersionTypeID":1,"name":"1.21.8","slug":"1-21-8","apiVersion":""},{"id":13927,"gameVersionTypeID":77784,"name":"1.21.9","slug":"1-21-9","apiVersion":null},{"id":13933,"gameVersionTypeID":1,"name":"1.21.9","slug":"1-21-9","apiVersion":""},{"id":13964,"gameVersionTypeID":77784,"name":"1.21.10","slug":"1-21-10","apiVersion":null},{"id":13966,"gameVersionTypeID":1,"name":"1.21.10","slug":"1-21-10","apiVersion":""}]`))
	})
}

func (s *fakeServer) registerProjectUploadFileHandler(mux *http.ServeMux) {
	s.registerHandler(mux, http.MethodPost, "/api/projects/{projectID}/upload-file", func(w http.ResponseWriter, r *http.Request) {
		metadata := make(map[string]any)
		if err := json.Unmarshal([]byte(r.FormValue("metadata")), &metadata); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"errorCode":1002,"errorMessage":"Error in field ` + "`" + `metadata` + "`" + `:\nMust be a JSON object."}`))
			return
		}

		_, ok1 := metadata["releaseType"]
		_, ok2 := metadata["changelog"]
		if !ok1 || !ok2 {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"errorCode":1002,"errorMessage":"Error in field ` + "`" + `metadata` + "`" + `:\n* Required properties are missing from object."}`))
			return
		}

		projectID, _ := strconv.Atoi(r.PathValue("projectID"))
		if projectID != s.projectID {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"errorCode":1005,"errorMessage":"Invalid ` + "`" + `projectID` + "`" + `."}`))
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":1234567}`))
	})
}
