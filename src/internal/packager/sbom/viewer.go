// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors

// Package sbom contains tools for generating SBOMs.
package sbom

import (
	"encoding/json"
	"fmt"
	"html/template"

	"github.com/defenseunicorns/zarf/src/types"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
)

func (b *Builder) createSBOMViewerAsset(identifier string, jsonData []byte) error {

	// Create the sbom viewer file for the image
	sbomViewerFile, err := b.createSBOMFile("sbom-viewer-%s.html", identifier)
	if err != nil {
		return err
	}

	defer sbomViewerFile.Close()

	// Create the sbomviewer template data
	tplData := struct {
		ThemeCSS  template.CSS
		ViewerCSS template.CSS
		List      template.JS
		Data      template.JS
		LibraryJS template.JS
		ViewerJS  template.JS
	}{
		ThemeCSS:  b.loadFileCSS("theme.css"),
		ViewerCSS: b.loadFileCSS("styles.css"),
		List:      template.JS(b.jsonList),
		Data:      template.JS(jsonData),
		LibraryJS: b.loadFileJS("library.js"),
		ViewerJS:  b.loadFileJS("viewer.js"),
	}

	// Render the sbomviewer template
	tpl, err := template.ParseFS(viewerAssets, "viewer/template.gohtml")
	if err != nil {
		return err
	}

	// Write the sbomviewer template to disk
	return tpl.Execute(sbomViewerFile, tplData)
}

func (b *Builder) loadFileCSS(name string) template.CSS {
	data, _ := viewerAssets.ReadFile("viewer/" + name)
	return template.CSS(data)
}

func (b *Builder) loadFileJS(name string) template.JS {
	data, _ := viewerAssets.ReadFile("viewer/" + name)
	return template.JS(data)
}

// This could be optimized, but loop over all the images and components to create a list of json files.
func (b *Builder) generateJSONList(componentToFiles map[string]*types.ComponentSBOM, tagToImage map[name.Tag]v1.Image) ([]byte, error) {
	var jsonList []string

	for tag := range tagToImage {
		normalized := b.getNormalizedFileName(tag.String())
		jsonList = append(jsonList, normalized)
	}

	for component := range componentToFiles {
		normalized := b.getNormalizedFileName(fmt.Sprintf("%s%s", componentPrefix, component))
		jsonList = append(jsonList, normalized)
	}

	return json.Marshal(jsonList)
}
