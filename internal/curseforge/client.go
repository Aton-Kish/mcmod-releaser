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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"
)

type Client interface {
	GameVersionTypes(ctx context.Context, input *GameVersionTypesInput) (output *GameVersionTypesOutput, err error)
	GameVersions(ctx context.Context, input *GameVersionsInput) (output *GameVersionsOutput, err error)
	ProjectUploadFile(ctx context.Context, input *ProjectUploadFileInput) (output *ProjectUploadFileOutput, err error)
}

type client struct {
	client   *http.Client
	baseURI  *url.URL
	apiToken string
}

type ClientOptionFunc func(*client)

func NewClient(token string, optFns ...ClientOptionFunc) (Client, error) {
	u, err := url.Parse("https://minecraft.curseforge.com")
	if err != nil {
		return nil, err
	}

	c := &client{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURI:  u,
		apiToken: token,
	}

	for _, fn := range optFns {
		fn(c)
	}

	return c, nil
}

func WithHTTPClient(httpClient *http.Client) ClientOptionFunc {
	return func(c *client) {
		c.client = httpClient
	}
}

func WithBaseURI(u *url.URL) ClientOptionFunc {
	return func(c *client) {
		c.baseURI = u
	}
}

func (c *client) do(ctx context.Context, method string, url string, header http.Header, body io.Reader) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header = header
	req.Header.Set("User-Agent", "atonkish/0.1.0")
	req.Header.Set("X-Api-Token", c.apiToken)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("%s %s: %s: %s", method, url, resp.Status, data)
	}

	return data, nil
}

type GameVersionType struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Slug string `json:"slug,omitempty"`
}

type GameVersionTypesInput struct {
}

type GameVersionTypesOutput []GameVersionType

func (c *client) GameVersionTypes(ctx context.Context, input *GameVersionTypesInput) (output *GameVersionTypesOutput, err error) {
	defer handleError(&err)

	data, err := c.do(ctx, http.MethodGet, c.baseURI.JoinPath("/api/game/version-types").String(), http.Header{}, nil)
	if err != nil {
		return nil, err
	}

	var resp GameVersionTypesOutput
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

type GameVersion struct {
	ID                int     `json:"id,omitempty"`
	GameVersionTypeID int     `json:"gameVersionTypeID,omitempty"`
	Name              string  `json:"name,omitempty"`
	Slug              string  `json:"slug,omitempty"`
	APIVersion        *string `json:"apiVersion,omitempty"`
}

type GameVersionsInput struct {
}

type GameVersionsOutput []GameVersion

func (c *client) GameVersions(ctx context.Context, input *GameVersionsInput) (output *GameVersionsOutput, err error) {
	defer handleError(&err)

	data, err := c.do(ctx, http.MethodGet, c.baseURI.JoinPath("/api/game/versions").String(), http.Header{}, nil)
	if err != nil {
		return nil, err
	}

	var resp GameVersionsOutput
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

type ProjectUploadFileMetadata struct {
	DisplayName              string                                 `json:"displayName,omitempty"`
	IsMarkedForManualRelease bool                                   `json:"isMarkedForManualRelease,omitempty"`
	ParentFileID             int                                    `json:"parentFileID,omitempty"`
	GameVersions             []int                                  `json:"gameVersions,omitempty"`
	ReleaseType              ProjectUploadFileMetadataReleaseType   `json:"releaseType,omitempty"`
	ChangelogType            ProjectUploadFileMetadataChangelogType `json:"changelogType,omitempty"`
	Changelog                string                                 `json:"changelog,omitempty"`
	Relations                *ProjectUploadFileMetadataRelations    `json:"relations,omitempty"`
}

type ProjectUploadFileMetadataChangelogType string

const (
	ProjectUploadFileMetadataChangelogTypeText     ProjectUploadFileMetadataChangelogType = "text"
	ProjectUploadFileMetadataChangelogTypeHTML     ProjectUploadFileMetadataChangelogType = "html"
	ProjectUploadFileMetadataChangelogTypeMarkdown ProjectUploadFileMetadataChangelogType = "markdown"
)

type ProjectUploadFileMetadataReleaseType string

const (
	ProjectUploadFileMetadataReleaseTypeRelease ProjectUploadFileMetadataReleaseType = "release"
	ProjectUploadFileMetadataReleaseTypeBeta    ProjectUploadFileMetadataReleaseType = "beta"
	ProjectUploadFileMetadataReleaseTypeAlpha   ProjectUploadFileMetadataReleaseType = "alpha"
)

type ProjectUploadFileMetadataRelations struct {
	Projects []ProjectUploadFileMetadataRelationsProject `json:"projects,omitempty"`
}

type ProjectUploadFileMetadataRelationsProject struct {
	ProjectID int                                           `json:"projectID,omitempty"`
	Slug      string                                        `json:"slug,omitempty"`
	Type      ProjectUploadFileMetadataRelationsProjectType `json:"type,omitempty"`
}

type ProjectUploadFileMetadataRelationsProjectType string

const (
	ProjectUploadFileMetadataRelationsProjectTypeRequiredDependency ProjectUploadFileMetadataRelationsProjectType = "requiredDependency"
	ProjectUploadFileMetadataRelationsProjectTypeOptionalDependency ProjectUploadFileMetadataRelationsProjectType = "optionalDependency"
	ProjectUploadFileMetadataRelationsProjectTypeIncompatible       ProjectUploadFileMetadataRelationsProjectType = "incompatible"
	ProjectUploadFileMetadataRelationsProjectTypeEmbeddedLibrary    ProjectUploadFileMetadataRelationsProjectType = "embeddedLibrary"
	ProjectUploadFileMetadataRelationsProjectTypeTool               ProjectUploadFileMetadataRelationsProjectType = "tool"
)

type ProjectUploadFileInput struct {
	ProjectID int
	File      io.Reader
	Metadata  *ProjectUploadFileMetadata
}

type ProjectUploadFileOutput struct {
	ID int `json:"id,omitempty"`
}

func (c *client) ProjectUploadFile(ctx context.Context, input *ProjectUploadFileInput) (output *ProjectUploadFileOutput, err error) {
	defer handleError(&err)

	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)

	filePart, err := w.CreateFormFile("file", "mod.jar")
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(filePart, input.File); err != nil {
		return nil, err
	}

	metadata, err := json.Marshal(input.Metadata)
	if err != nil {
		return nil, err
	}
	if err := w.WriteField("metadata", string(metadata)); err != nil {
		return nil, err
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	data, err := c.do(ctx, http.MethodPost, c.baseURI.JoinPath(fmt.Sprintf("/api/projects/%d/upload-file", input.ProjectID)).String(), http.Header{"Content-Type": {w.FormDataContentType()}}, body)
	if err != nil {
		return nil, err
	}

	var resp ProjectUploadFileOutput
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
