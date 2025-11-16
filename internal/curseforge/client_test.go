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
	"io"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_client_GameVersionTypes(t *testing.T) {
	s := NewFakeServer(8080)
	go func() {
		_ = s.ListenAndServe()
	}()
	t.Cleanup(func() {
		_ = s.Close()
	})

	u, err := url.Parse("http://localhost:8080")
	assert.NoError(t, err)

	validAPIToken := "beafbeaf-beaf-4000-8000-beafbeafbeaf"
	invalidAPIToken := "deaddead-dead-4000-8000-deaddeaddead"
	s.RegisterAPIToken(validAPIToken)

	tests := []struct {
		name         string
		apiToken     string
		isServerDown bool
		input        *GameVersionTypesInput
		want         *GameVersionTypesOutput
		wantErr      bool
	}{
		{
			name:         "happy path",
			apiToken:     validAPIToken,
			isServerDown: false,
			input:        &GameVersionTypesInput{},
			want: &GameVersionTypesOutput{
				{ID: 2, Name: "Java", Slug: "java"},
				{ID: 615, Name: "Addons", Slug: "addons"},
				{ID: 68441, Name: "Modloader", Slug: "modloader"},
				{ID: 75125, Name: "Minecraft 1.20", Slug: "minecraft-1-20"},
				{ID: 75208, Name: "Environment", Slug: "environment"},
				{ID: 77784, Name: "Minecraft 1.21", Slug: "minecraft-1-21"},
			},
			wantErr: false,
		},
		{
			name:         "edge path: missing API token",
			apiToken:     "",
			isServerDown: false,
			input:        &GameVersionTypesInput{},
			want:         nil,
			wantErr:      true,
		},
		{
			name:         "edge path: malformed API token",
			apiToken:     "malformed-api-token",
			isServerDown: false,
			input:        &GameVersionTypesInput{},
			want:         nil,
			wantErr:      true,
		},
		{
			name:         "edge path: invalid API token",
			apiToken:     invalidAPIToken,
			isServerDown: false,
			input:        &GameVersionTypesInput{},
			want:         nil,
			wantErr:      true,
		},
		{
			name:         "edge path: server down",
			apiToken:     validAPIToken,
			isServerDown: true,
			input:        &GameVersionTypesInput{},
			want:         nil,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer s.Up()
			if tt.isServerDown {
				s.Down()
			}

			c, err := NewClient(tt.apiToken, WithBaseURI(u))
			assert.NoError(t, err)

			got, err := c.GameVersionTypes(context.Background(), tt.input)
			if tt.wantErr {
				e := new(Error)
				assert.ErrorAs(t, err, &e)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_client_GameVersions(t *testing.T) {
	s := NewFakeServer(8080)
	go func() {
		_ = s.ListenAndServe()
	}()
	t.Cleanup(func() {
		_ = s.Close()
	})

	u, err := url.Parse("http://localhost:8080")
	assert.NoError(t, err)

	validAPIToken := "beafbeaf-beaf-4000-8000-beafbeafbeaf"
	invalidAPIToken := "deaddead-dead-4000-8000-deaddeaddead"
	s.RegisterAPIToken(validAPIToken)

	empty := ""

	tests := []struct {
		name         string
		apiToken     string
		isServerDown bool
		input        *GameVersionsInput
		want         *GameVersionsOutput
		wantErr      bool
	}{
		{
			name:         "happy path",
			apiToken:     validAPIToken,
			isServerDown: false,
			input:        &GameVersionsInput{},
			want: &GameVersionsOutput{
				{ID: 7498, GameVersionTypeID: 68441, Name: "Forge", Slug: "forge", APIVersion: &empty},
				{ID: 7499, GameVersionTypeID: 68441, Name: "Fabric", Slug: "fabric", APIVersion: &empty},
				{ID: 9638, GameVersionTypeID: 75208, Name: "Client", Slug: "client", APIVersion: &empty},
				{ID: 9639, GameVersionTypeID: 75208, Name: "Server", Slug: "server", APIVersion: &empty},
				{ID: 11135, GameVersionTypeID: 2, Name: "Java 21", Slug: "java-21", APIVersion: &empty},
				{ID: 9970, GameVersionTypeID: 615, Name: "1.20", Slug: "1-20", APIVersion: &empty},
				{ID: 9971, GameVersionTypeID: 75125, Name: "1.20", Slug: "1-20", APIVersion: &empty},
				{ID: 11457, GameVersionTypeID: 77784, Name: "1.21", Slug: "1-21", APIVersion: &empty},
				{ID: 11458, GameVersionTypeID: 615, Name: "1.21", Slug: "1-21", APIVersion: &empty},
				{ID: 11515, GameVersionTypeID: 1, Name: "1.21", Slug: "1-21", APIVersion: &empty},
				{ID: 11779, GameVersionTypeID: 77784, Name: "1.21.1", Slug: "1-21-1", APIVersion: nil},
				{ID: 12079, GameVersionTypeID: 77784, Name: "1.21.2", Slug: "1-21-2", APIVersion: &empty},
				{ID: 12084, GameVersionTypeID: 77784, Name: "1.21.3", Slug: "1-21-3", APIVersion: nil},
				{ID: 12281, GameVersionTypeID: 77784, Name: "1.21.4", Slug: "1-21-4", APIVersion: nil},
				{ID: 12735, GameVersionTypeID: 1, Name: "1.21.1", Slug: "1-21-1", APIVersion: &empty},
				{ID: 12736, GameVersionTypeID: 1, Name: "1.21.2", Slug: "1-21-2", APIVersion: &empty},
				{ID: 12737, GameVersionTypeID: 1, Name: "1.21.3", Slug: "1-21-3", APIVersion: &empty},
				{ID: 12738, GameVersionTypeID: 1, Name: "1.21.4", Slug: "1-21-4", APIVersion: &empty},
				{ID: 12934, GameVersionTypeID: 77784, Name: "1.21.5", Slug: "1-21-5", APIVersion: nil},
				{ID: 12988, GameVersionTypeID: 1, Name: "1.21.5", Slug: "1-21-5", APIVersion: &empty},
				{ID: 13422, GameVersionTypeID: 77784, Name: "1.21.6", Slug: "1-21-6", APIVersion: nil},
				{ID: 13473, GameVersionTypeID: 1, Name: "1.21.6", Slug: "1-21-6", APIVersion: &empty},
				{ID: 13506, GameVersionTypeID: 77784, Name: "1.21.7", Slug: "1-21-7", APIVersion: nil},
				{ID: 13574, GameVersionTypeID: 1, Name: "1.21.7", Slug: "1-21-7", APIVersion: &empty},
				{ID: 13620, GameVersionTypeID: 77784, Name: "1.21.8", Slug: "1-21-8", APIVersion: nil},
				{ID: 13683, GameVersionTypeID: 1, Name: "1.21.8", Slug: "1-21-8", APIVersion: &empty},
				{ID: 13927, GameVersionTypeID: 77784, Name: "1.21.9", Slug: "1-21-9", APIVersion: nil},
				{ID: 13933, GameVersionTypeID: 1, Name: "1.21.9", Slug: "1-21-9", APIVersion: &empty},
				{ID: 13964, GameVersionTypeID: 77784, Name: "1.21.10", Slug: "1-21-10", APIVersion: nil},
				{ID: 13966, GameVersionTypeID: 1, Name: "1.21.10", Slug: "1-21-10", APIVersion: &empty},
			},
			wantErr: false,
		},
		{
			name:         "edge path: missing API token",
			apiToken:     "",
			isServerDown: false,
			input:        &GameVersionsInput{},
			want:         nil,
			wantErr:      true,
		},
		{
			name:         "edge path: malformed API token",
			apiToken:     "malformed-api-token",
			isServerDown: false,
			input:        &GameVersionsInput{},
			want:         nil,
			wantErr:      true,
		},
		{
			name:         "edge path: invalid API token",
			apiToken:     invalidAPIToken,
			isServerDown: false,
			input:        &GameVersionsInput{},
			want:         nil,
			wantErr:      true,
		},
		{
			name:         "edge path: server down",
			apiToken:     validAPIToken,
			isServerDown: true,
			input:        &GameVersionsInput{},
			want:         nil,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer s.Up()
			if tt.isServerDown {
				s.Down()
			}

			c, err := NewClient(tt.apiToken, WithBaseURI(u))
			assert.NoError(t, err)

			got, err := c.GameVersions(context.Background(), tt.input)
			if tt.wantErr {
				e := new(Error)
				assert.ErrorAs(t, err, &e)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_client_ProjectUploadFile(t *testing.T) {
	s := NewFakeServer(8080)
	go func() {
		_ = s.ListenAndServe()
	}()
	t.Cleanup(func() {
		_ = s.Close()
	})

	u, err := url.Parse("http://localhost:8080")
	assert.NoError(t, err)

	validAPIToken := "beafbeaf-beaf-4000-8000-beafbeafbeaf"
	invalidAPIToken := "deaddead-dead-4000-8000-deaddeaddead"
	s.RegisterAPIToken(validAPIToken)

	validProjectID := 123456
	invalidProjectID := 999999
	s.RegisterProjectID(validProjectID)

	tests := []struct {
		name         string
		apiToken     string
		isServerDown bool
		input        *ProjectUploadFileInput
		want         *ProjectUploadFileOutput
		wantErr      bool
	}{
		{
			name:         "happy path",
			apiToken:     validAPIToken,
			isServerDown: false,
			input: &ProjectUploadFileInput{
				ProjectID: validProjectID,
				FileName:  "example-mod-1.2.3.jar",
				File:      bytes.NewBufferString("mod.jar"),
				Metadata: &ProjectUploadFileMetadata{
					DisplayName: "Example Mod v1.2.3",
					GameVersions: []int{
						7499,
						9638,
						9639,
						11135,
						13927,
						13964,
					},
					ReleaseType:   ProjectUploadFileMetadataReleaseTypeRelease,
					ChangelogType: ProjectUploadFileMetadataChangelogTypeMarkdown,
					Changelog:     "[v1.2.3](https://github.com/FabricMC/fabric-example-mod/releases/tag/v1.2.3)",
					Relations: &ProjectUploadFileMetadataRelations{
						Projects: []ProjectUploadFileMetadataRelationsProject{
							{Slug: "cloth-config", Type: ProjectUploadFileMetadataRelationsProjectTypeEmbeddedLibrary},
							{Slug: "fabric-api", Type: ProjectUploadFileMetadataRelationsProjectTypeRequiredDependency},
							{Slug: "modmenu", Type: ProjectUploadFileMetadataRelationsProjectTypeOptionalDependency},
						},
					},
				},
			},
			want: &ProjectUploadFileOutput{
				ID: 1234567,
			},
			wantErr: false,
		},
		{
			name:         "edge path: missing API token",
			apiToken:     "",
			isServerDown: false,
			input: &ProjectUploadFileInput{
				ProjectID: validProjectID,
				FileName:  "example-mod-1.2.3.jar",
				File:      bytes.NewBufferString("mod.jar"),
				Metadata: &ProjectUploadFileMetadata{
					DisplayName: "Example Mod v1.2.3",
					GameVersions: []int{
						7499,
						9638,
						9639,
						11135,
						13927,
						13964,
					},
					ReleaseType:   ProjectUploadFileMetadataReleaseTypeRelease,
					ChangelogType: ProjectUploadFileMetadataChangelogTypeMarkdown,
					Changelog:     "[v1.2.3](https://github.com/FabricMC/fabric-example-mod/releases/tag/v1.2.3)",
					Relations: &ProjectUploadFileMetadataRelations{
						Projects: []ProjectUploadFileMetadataRelationsProject{
							{Slug: "cloth-config", Type: ProjectUploadFileMetadataRelationsProjectTypeEmbeddedLibrary},
							{Slug: "fabric-api", Type: ProjectUploadFileMetadataRelationsProjectTypeRequiredDependency},
							{Slug: "modmenu", Type: ProjectUploadFileMetadataRelationsProjectTypeOptionalDependency},
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:         "edge path: malformed API token",
			apiToken:     "malformed-api-token",
			isServerDown: false,
			input: &ProjectUploadFileInput{
				ProjectID: validProjectID,
				FileName:  "example-mod-1.2.3.jar",
				File:      bytes.NewBufferString("mod.jar"),
				Metadata: &ProjectUploadFileMetadata{
					DisplayName: "Example Mod v1.2.3",
					GameVersions: []int{
						7499,
						9638,
						9639,
						11135,
						13927,
						13964,
					},
					ReleaseType:   ProjectUploadFileMetadataReleaseTypeRelease,
					ChangelogType: ProjectUploadFileMetadataChangelogTypeMarkdown,
					Changelog:     "[v1.2.3](https://github.com/FabricMC/fabric-example-mod/releases/tag/v1.2.3)",
					Relations: &ProjectUploadFileMetadataRelations{
						Projects: []ProjectUploadFileMetadataRelationsProject{
							{Slug: "cloth-config", Type: ProjectUploadFileMetadataRelationsProjectTypeEmbeddedLibrary},
							{Slug: "fabric-api", Type: ProjectUploadFileMetadataRelationsProjectTypeRequiredDependency},
							{Slug: "modmenu", Type: ProjectUploadFileMetadataRelationsProjectTypeOptionalDependency},
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:         "edge path: invalid API token",
			apiToken:     invalidAPIToken,
			isServerDown: false,
			input: &ProjectUploadFileInput{
				ProjectID: validProjectID,
				FileName:  "example-mod-1.2.3.jar",
				File:      bytes.NewBufferString("mod.jar"),
				Metadata: &ProjectUploadFileMetadata{
					DisplayName: "Example Mod v1.2.3",
					GameVersions: []int{
						7499,
						9638,
						9639,
						11135,
						13927,
						13964,
					},
					ReleaseType:   ProjectUploadFileMetadataReleaseTypeRelease,
					ChangelogType: ProjectUploadFileMetadataChangelogTypeMarkdown,
					Changelog:     "[v1.2.3](https://github.com/FabricMC/fabric-example-mod/releases/tag/v1.2.3)",
					Relations: &ProjectUploadFileMetadataRelations{
						Projects: []ProjectUploadFileMetadataRelationsProject{
							{Slug: "cloth-config", Type: ProjectUploadFileMetadataRelationsProjectTypeEmbeddedLibrary},
							{Slug: "fabric-api", Type: ProjectUploadFileMetadataRelationsProjectTypeRequiredDependency},
							{Slug: "modmenu", Type: ProjectUploadFileMetadataRelationsProjectTypeOptionalDependency},
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:         "edge path: missing required properties",
			apiToken:     validAPIToken,
			isServerDown: false,
			input: &ProjectUploadFileInput{
				ProjectID: validProjectID,
				FileName:  "example-mod-1.2.3.jar",
				File:      bytes.NewBufferString("mod.jar"),
				Metadata:  &ProjectUploadFileMetadata{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:         "edge path: invalid project ID",
			apiToken:     validAPIToken,
			isServerDown: false,
			input: &ProjectUploadFileInput{
				ProjectID: invalidProjectID,
				FileName:  "example-mod-1.2.3.jar",
				File:      bytes.NewBufferString("mod.jar"),
				Metadata: &ProjectUploadFileMetadata{
					DisplayName: "Example Mod v1.2.3",
					GameVersions: []int{
						7499,
						9638,
						9639,
						11135,
						13927,
						13964,
					},
					ReleaseType:   ProjectUploadFileMetadataReleaseTypeRelease,
					ChangelogType: ProjectUploadFileMetadataChangelogTypeMarkdown,
					Changelog:     "[v1.2.3](https://github.com/FabricMC/fabric-example-mod/releases/tag/v1.2.3)",
					Relations: &ProjectUploadFileMetadataRelations{
						Projects: []ProjectUploadFileMetadataRelationsProject{
							{Slug: "cloth-config", Type: ProjectUploadFileMetadataRelationsProjectTypeEmbeddedLibrary},
							{Slug: "fabric-api", Type: ProjectUploadFileMetadataRelationsProjectTypeRequiredDependency},
							{Slug: "modmenu", Type: ProjectUploadFileMetadataRelationsProjectTypeOptionalDependency},
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:         "edge path: server down",
			apiToken:     validAPIToken,
			isServerDown: true,
			input: &ProjectUploadFileInput{
				ProjectID: validProjectID,
				FileName:  "example-mod-1.2.3.jar",
				File:      bytes.NewBufferString("mod.jar"),
				Metadata: &ProjectUploadFileMetadata{
					DisplayName: "Example Mod v1.2.3",
					GameVersions: []int{
						7499,
						9638,
						9639,
						11135,
						13927,
						13964,
					},
					ReleaseType:   ProjectUploadFileMetadataReleaseTypeRelease,
					ChangelogType: ProjectUploadFileMetadataChangelogTypeMarkdown,
					Changelog:     "[v1.2.3](https://github.com/FabricMC/fabric-example-mod/releases/tag/v1.2.3)",
					Relations: &ProjectUploadFileMetadataRelations{
						Projects: []ProjectUploadFileMetadataRelationsProject{
							{Slug: "cloth-config", Type: ProjectUploadFileMetadataRelationsProjectTypeEmbeddedLibrary},
							{Slug: "fabric-api", Type: ProjectUploadFileMetadataRelationsProjectTypeRequiredDependency},
							{Slug: "modmenu", Type: ProjectUploadFileMetadataRelationsProjectTypeOptionalDependency},
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer s.Up()
			if tt.isServerDown {
				s.Down()
			}

			c, err := NewClient(tt.apiToken, WithBaseURI(u))
			assert.NoError(t, err)

			got, err := c.ProjectUploadFile(context.Background(), tt.input)
			if tt.wantErr {
				e := new(Error)
				assert.ErrorAs(t, err, &e)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_client_createProjectUploadFileContent(t *testing.T) {
	tests := []struct {
		name     string
		input    *ProjectUploadFileInput
		boundary string
		want     string
		want2    string
		wantErr  bool
	}{
		{
			name: "happy path",
			input: &ProjectUploadFileInput{
				ProjectID: 123456,
				FileName:  "example-mod-1.2.3.jar",
				File:      bytes.NewBufferString("mod.jar"),
				Metadata: &ProjectUploadFileMetadata{
					DisplayName: "Example Mod v1.2.3",
					GameVersions: []int{
						7499,
						9638,
						9639,
						11135,
						13927,
						13964,
					},
					ReleaseType:   ProjectUploadFileMetadataReleaseTypeRelease,
					ChangelogType: ProjectUploadFileMetadataChangelogTypeMarkdown,
					Changelog:     "[v1.2.3](https://github.com/FabricMC/fabric-example-mod/releases/tag/v1.2.3)",
					Relations: &ProjectUploadFileMetadataRelations{
						Projects: []ProjectUploadFileMetadataRelationsProject{
							{Slug: "cloth-config", Type: ProjectUploadFileMetadataRelationsProjectTypeEmbeddedLibrary},
							{Slug: "fabric-api", Type: ProjectUploadFileMetadataRelationsProjectTypeRequiredDependency},
							{Slug: "modmenu", Type: ProjectUploadFileMetadataRelationsProjectTypeOptionalDependency},
						},
					},
				},
			},
			boundary: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789ab",
			want: strings.Join([]string{
				`--0123456789abcdef0123456789abcdef0123456789abcdef0123456789ab`,
				`Content-Disposition: form-data; name="file"; filename="example-mod-1.2.3.jar"`,
				`Content-Type: application/octet-stream`,
				``,
				`mod.jar`,
				`--0123456789abcdef0123456789abcdef0123456789abcdef0123456789ab`,
				`Content-Disposition: form-data; name="metadata"`,
				``,
				`{"displayName":"Example Mod v1.2.3","gameVersions":[7499,9638,9639,11135,13927,13964],"releaseType":"release","changelogType":"markdown","changelog":"[v1.2.3](https://github.com/FabricMC/fabric-example-mod/releases/tag/v1.2.3)","relations":{"projects":[{"slug":"cloth-config","type":"embeddedLibrary"},{"slug":"fabric-api","type":"requiredDependency"},{"slug":"modmenu","type":"optionalDependency"}]}}`,
				`--0123456789abcdef0123456789abcdef0123456789abcdef0123456789ab--`,
				``,
			}, "\r\n"),
			want2:   "multipart/form-data; boundary=0123456789abcdef0123456789abcdef0123456789abcdef0123456789ab",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c client
			got, got2, err := c.createProjectUploadFileContent(context.Background(), tt.input, tt.boundary)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.Empty(t, got2)
			} else {
				assert.NoError(t, err)
				body, err := io.ReadAll(got)
				assert.NoError(t, err)
				assert.Equal(t, tt.want, string(body))
				assert.Equal(t, tt.want2, got2)
			}
		})
	}
}
