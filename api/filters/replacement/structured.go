// Copyright 2021 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package replacement

import (
	"encoding/json"
	"fmt"
	jsonpatch "github.com/evanphx/json-patch"
	"github.com/go-openapi/jsonpointer"
	yaml "gopkg.in/yaml.v3"
	"sigs.k8s.io/kustomize/api/filters/replacement/yamlpatch"
	"sigs.k8s.io/kustomize/api/filters/replacement/yptr"
	"sigs.k8s.io/kustomize/api/types"
)

func getJsonPathValue(options *types.FieldOptions, jsonValue string) (string, error) {
	p, err := jsonpointer.New(options.JSONPath)
	if err != nil {
		return "", err
	}

	var js interface{}

	if err := json.Unmarshal([]byte(jsonValue), &js); err != nil {
		return "", fmt.Errorf("json unmarshall error: %w", err)
	}

	v, _, err := p.Get(js)
	if err != nil {
		return "", fmt.Errorf("json pointer error: %w", err)
	}

	return fmt.Sprintf("%v", v), nil
}

func getJsonReplacementValue(options *types.FieldOptions, jsonValue string, replacementValue string) (string, error) {
	patchJSON := []byte(`[
		{"op": "replace", "path": "` + options.JSONPath + `", "value": "` + replacementValue + `"}
		]`)

	patch, err := jsonpatch.DecodePatch(patchJSON)
	if err != nil {
		return "", err
	}

	modified, err := patch.ApplyIndent([]byte(jsonValue), "  ")
	if err != nil {
		return "", err
	}

	return string(modified), nil
}

func getYAMLPathValue(options *types.FieldOptions, yamlValue string) (string, error) {
	var n yaml.Node
	if err := yaml.Unmarshal([]byte(yamlValue), &n); err != nil {
		return "", fmt.Errorf("yaml unmarshall error: %w", err)
	}

	v, err := yptr.Find(&n, options.YAMLPath)
	if err != nil {
		return "", fmt.Errorf("json pointer error: %w", err)
	}

	return fmt.Sprintf("%v", v), nil
}

func getYAMLReplacementValue(options *types.FieldOptions, yamlValue string, replacementValue string) (string, error) {
	patchYAML := []byte(`---
- op: replace
  path: ` + options.YAMLPath + `
  value: ` + replacementValue + `
`)

	patch, err := yamlpatch.DecodePatch(patchYAML)
	if err != nil {
		return "", err
	}

	modified, err := patch.Apply([]byte(yamlValue))
	if err != nil {
		return "", err
	}

	return string(modified), nil
}
