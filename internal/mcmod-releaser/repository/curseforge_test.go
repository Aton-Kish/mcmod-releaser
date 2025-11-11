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

package command

import (
	"bytes"
	"context"
	"io"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Aton-Kish/mcmod-releaser/internal/curseforge"
	"github.com/Aton-Kish/mcmod-releaser/internal/mcmod-releaser/model"
)

func Test_curseForgeRepository_CreateMod(t *testing.T) {
	s := curseforge.NewFakeServer(8080)
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
		file         io.Reader
		mod          *model.Mod
		want         *model.Mod
		wantErr      bool
	}{
		{
			name:         "happy path",
			apiToken:     validAPIToken,
			isServerDown: false,
			file:         bytes.NewBufferString("mod.jar"),
			mod: &model.Mod{
				ProjectID:    strconv.Itoa(validProjectID),
				ReleaseType:  model.ModReleaseTypeRelease,
				Version:      "1.2.3",
				Name:         "Example Mod v1.2.3",
				Changelog:    "[v1.2.3](https://github.com/FabricMC/fabric-example-mod/releases/tag/v1.2.3)",
				Environments: []string{"Client", "Server"},
				Loaders:      []string{"Fabric"},
				JavaVersions: []string{"Java 21"},
				GameVersions: []string{"1.21.9", "1.21.10"},
				Dependencies: map[string]model.ModDependencyType{
					"fabric-api":   model.ModDependencyTypeRequired,
					"cloth-config": model.ModDependencyTypeEmbedded,
					"modmenu":      model.ModDependencyTypeOptional,
				},
			},
			want: &model.Mod{
				ID:           "1234567",
				ProjectID:    strconv.Itoa(validProjectID),
				ReleaseType:  model.ModReleaseTypeRelease,
				Version:      "1.2.3",
				Name:         "Example Mod v1.2.3",
				Changelog:    "[v1.2.3](https://github.com/FabricMC/fabric-example-mod/releases/tag/v1.2.3)",
				Environments: []string{"Client", "Server"},
				Loaders:      []string{"Fabric"},
				JavaVersions: []string{"Java 21"},
				GameVersions: []string{"1.21.9", "1.21.10"},
				Dependencies: map[string]model.ModDependencyType{
					"fabric-api":   model.ModDependencyTypeRequired,
					"cloth-config": model.ModDependencyTypeEmbedded,
					"modmenu":      model.ModDependencyTypeOptional,
				},
			},
			wantErr: false,
		},
		{
			name:         "edge path: missing API token",
			apiToken:     "",
			isServerDown: false,
			file:         bytes.NewBufferString("mod.jar"),
			mod: &model.Mod{
				ProjectID:    strconv.Itoa(validProjectID),
				ReleaseType:  model.ModReleaseTypeRelease,
				Version:      "1.2.3",
				Name:         "Example Mod v1.2.3",
				Changelog:    "[v1.2.3](https://github.com/FabricMC/fabric-example-mod/releases/tag/v1.2.3)",
				Environments: []string{"Client", "Server"},
				Loaders:      []string{"Fabric"},
				JavaVersions: []string{"Java 21"},
				GameVersions: []string{"1.21.9", "1.21.10"},
				Dependencies: map[string]model.ModDependencyType{
					"fabric-api":   model.ModDependencyTypeRequired,
					"cloth-config": model.ModDependencyTypeEmbedded,
					"modmenu":      model.ModDependencyTypeOptional,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:         "edge path: malformed API token",
			apiToken:     "malformed-api-token",
			isServerDown: false,
			file:         bytes.NewBufferString("mod.jar"),
			mod: &model.Mod{
				ProjectID:    strconv.Itoa(validProjectID),
				ReleaseType:  model.ModReleaseTypeRelease,
				Version:      "1.2.3",
				Name:         "Example Mod v1.2.3",
				Changelog:    "[v1.2.3](https://github.com/FabricMC/fabric-example-mod/releases/tag/v1.2.3)",
				Environments: []string{"Client", "Server"},
				Loaders:      []string{"Fabric"},
				JavaVersions: []string{"Java 21"},
				GameVersions: []string{"1.21.9", "1.21.10"},
				Dependencies: map[string]model.ModDependencyType{
					"fabric-api":   model.ModDependencyTypeRequired,
					"cloth-config": model.ModDependencyTypeEmbedded,
					"modmenu":      model.ModDependencyTypeOptional,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:         "edge path: invalid API token",
			apiToken:     invalidAPIToken,
			isServerDown: false,
			file:         bytes.NewBufferString("mod.jar"),
			mod: &model.Mod{
				ProjectID:    strconv.Itoa(validProjectID),
				ReleaseType:  model.ModReleaseTypeRelease,
				Version:      "1.2.3",
				Name:         "Example Mod v1.2.3",
				Changelog:    "[v1.2.3](https://github.com/FabricMC/fabric-example-mod/releases/tag/v1.2.3)",
				Environments: []string{"Client", "Server"},
				Loaders:      []string{"Fabric"},
				JavaVersions: []string{"Java 21"},
				GameVersions: []string{"1.21.9", "1.21.10"},
				Dependencies: map[string]model.ModDependencyType{
					"fabric-api":   model.ModDependencyTypeRequired,
					"cloth-config": model.ModDependencyTypeEmbedded,
					"modmenu":      model.ModDependencyTypeOptional,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:         "edge path: missing required properties",
			apiToken:     validAPIToken,
			isServerDown: false,
			file:         bytes.NewBufferString("mod.jar"),
			mod:          &model.Mod{},
			want:         nil,
			wantErr:      true,
		},
		{
			name:         "edge path: invalid project ID",
			apiToken:     validAPIToken,
			isServerDown: false,
			file:         bytes.NewBufferString("mod.jar"),
			mod: &model.Mod{
				ProjectID:    strconv.Itoa(invalidProjectID),
				ReleaseType:  model.ModReleaseTypeRelease,
				Version:      "1.2.3",
				Name:         "Example Mod v1.2.3",
				Changelog:    "[v1.2.3](https://github.com/FabricMC/fabric-example-mod/releases/tag/v1.2.3)",
				Environments: []string{"Client", "Server"},
				Loaders:      []string{"Fabric"},
				JavaVersions: []string{"Java 21"},
				GameVersions: []string{"1.21.9", "1.21.10"},
				Dependencies: map[string]model.ModDependencyType{
					"fabric-api":   model.ModDependencyTypeRequired,
					"cloth-config": model.ModDependencyTypeEmbedded,
					"modmenu":      model.ModDependencyTypeOptional,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:         "edge path: server down",
			apiToken:     validAPIToken,
			isServerDown: true,
			file:         bytes.NewBufferString("mod.jar"),
			mod: &model.Mod{
				ProjectID:    strconv.Itoa(validProjectID),
				ReleaseType:  model.ModReleaseTypeRelease,
				Version:      "1.2.3",
				Name:         "Example Mod v1.2.3",
				Changelog:    "[v1.2.3](https://github.com/FabricMC/fabric-example-mod/releases/tag/v1.2.3)",
				Environments: []string{"Client", "Server"},
				Loaders:      []string{"Fabric"},
				JavaVersions: []string{"Java 21"},
				GameVersions: []string{"1.21.9", "1.21.10"},
				Dependencies: map[string]model.ModDependencyType{
					"fabric-api":   model.ModDependencyTypeRequired,
					"cloth-config": model.ModDependencyTypeEmbedded,
					"modmenu":      model.ModDependencyTypeOptional,
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

			c, err := curseforge.NewClient(tt.apiToken, curseforge.WithBaseURI(u))
			assert.NoError(t, err)

			r := &curseForgeRepository{
				client: c,
			}

			got, err := r.CreateMod(context.Background(), tt.file, tt.mod)
			if tt.wantErr {
				t.Log(err)
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
