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

package repository

import (
	"context"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strconv"

	"github.com/Aton-Kish/mcmod-releaser/internal/curseforge"
	"github.com/Aton-Kish/mcmod-releaser/internal/mcmod-releaser/model"
)

type CurseForgeRepository interface {
	CreateMod(ctx context.Context, path string, mod *model.Mod) (*model.Mod, error)
}

type curseForgeRepository struct {
	client curseforge.Client
}

func NewCurseForgeRepository(client curseforge.Client) CurseForgeRepository {
	return &curseForgeRepository{
		client: client,
	}
}

func (r *curseForgeRepository) CreateMod(ctx context.Context, path string, mod *model.Mod) (*model.Mod, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	projectID, err := strconv.Atoi(mod.ProjectID)
	if err != nil {
		return nil, err
	}

	releaseType, err := r.convertReleaseType(ctx, mod.ReleaseType)
	if err != nil {
		return nil, err
	}

	relations, err := r.convertDependencies(ctx, mod.Dependencies)
	if err != nil {
		return nil, err
	}

	gameVersions, err := r.convertGameVersions(ctx, slices.Concat(mod.Environments, mod.Loaders, mod.JavaVersions, mod.GameVersions))
	if err != nil {
		return nil, err
	}

	uploaded, err := r.client.ProjectUploadFile(ctx, &curseforge.ProjectUploadFileInput{
		ProjectID: projectID,
		FileName:  filepath.Base(file.Name()),
		File:      file,
		Metadata: &curseforge.ProjectUploadFileMetadata{
			DisplayName:   mod.Name,
			GameVersions:  gameVersions,
			ReleaseType:   releaseType,
			ChangelogType: curseforge.ProjectUploadFileMetadataChangelogTypeMarkdown,
			Changelog:     mod.Changelog,
			Relations:     relations,
		},
	})
	if err != nil {
		return nil, err
	}

	return &model.Mod{
		ID:           strconv.Itoa(uploaded.ID),
		ProjectID:    mod.ProjectID,
		ReleaseType:  mod.ReleaseType,
		Version:      mod.Version,
		Name:         mod.Name,
		Changelog:    mod.Changelog,
		Environments: mod.Environments,
		Loaders:      mod.Loaders,
		JavaVersions: mod.JavaVersions,
		GameVersions: mod.GameVersions,
		Dependencies: maps.Clone(mod.Dependencies),
	}, nil
}

func (r *curseForgeRepository) convertReleaseType(ctx context.Context, v model.ModReleaseType) (curseforge.ProjectUploadFileMetadataReleaseType, error) {
	switch v {
	case model.ModReleaseTypeRelease:
		return curseforge.ProjectUploadFileMetadataReleaseTypeRelease, nil
	case model.ModReleaseTypeBeta:
		return curseforge.ProjectUploadFileMetadataReleaseTypeBeta, nil
	case model.ModReleaseTypeAlpha:
		return curseforge.ProjectUploadFileMetadataReleaseTypeAlpha, nil
	default:
		return "", fmt.Errorf("invalid release type: %s", v)
	}
}

func (r *curseForgeRepository) convertDependency(ctx context.Context, v model.ModDependencyType) (curseforge.ProjectUploadFileMetadataRelationsProjectType, error) {
	switch v {
	case model.ModDependencyTypeRequired:
		return curseforge.ProjectUploadFileMetadataRelationsProjectTypeRequiredDependency, nil
	case model.ModDependencyTypeOptional:
		return curseforge.ProjectUploadFileMetadataRelationsProjectTypeOptionalDependency, nil
	case model.ModDependencyTypeIncompatible:
		return curseforge.ProjectUploadFileMetadataRelationsProjectTypeIncompatible, nil
	case model.ModDependencyTypeEmbedded:
		return curseforge.ProjectUploadFileMetadataRelationsProjectTypeEmbeddedLibrary, nil
	case model.ModDependencyTypeTool:
		return curseforge.ProjectUploadFileMetadataRelationsProjectTypeTool, nil
	default:
		return "", fmt.Errorf("invalid dependency type: %s", v)
	}
}

func (r *curseForgeRepository) convertDependencies(ctx context.Context, v map[string]model.ModDependencyType) (*curseforge.ProjectUploadFileMetadataRelations, error) {
	ps := make([]curseforge.ProjectUploadFileMetadataRelationsProject, 0, len(v))
	for k, v := range v {
		d, err := r.convertDependency(ctx, v)
		if err != nil {
			return nil, err
		}

		ps = append(ps, curseforge.ProjectUploadFileMetadataRelationsProject{
			Slug: k,
			Type: d,
		})
	}

	return &curseforge.ProjectUploadFileMetadataRelations{Projects: ps}, nil
}

func (r *curseForgeRepository) convertGameVersions(ctx context.Context, v []string) ([]int, error) {
	gameVersionTypes, err := r.client.GameVersionTypes(ctx, &curseforge.GameVersionTypesInput{})
	if err != nil {
		return nil, err
	}

	gameVersionTypeIDs := make([]int, 0, len(*gameVersionTypes))
	for _, gameVersionType := range *gameVersionTypes {
		gameVersionTypeIDs = append(gameVersionTypeIDs, gameVersionType.ID)
	}

	gameVersions, err := r.client.GameVersions(ctx, &curseforge.GameVersionsInput{})
	if err != nil {
		return nil, err
	}

	gameVersionIDs := make([]int, 0, len(v))
	for _, gameVersion := range *gameVersions {
		if slices.Contains(gameVersionTypeIDs, gameVersion.GameVersionTypeID) && slices.Contains(v, gameVersion.Name) {
			gameVersionIDs = append(gameVersionIDs, gameVersion.ID)
		}
	}

	return gameVersionIDs, nil
}
