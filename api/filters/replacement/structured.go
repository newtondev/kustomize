// Copyright 2021 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package replacement

import (
	"encoding/json"
	"fmt"
	jsonpatch "github.com/evanphx/json-patch"
	"github.com/go-openapi/jsonpointer"
	"sigs.k8s.io/kustomize/api/types"
)

func getJsonPathValue(options *types.FieldOptions, jsonValue string) (string, error) {
	p, err := jsonpointer.New(options.FormatPath)
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
		{"op": "replace", "path": "` + options.FormatPath + `", "value": "` + replacementValue + `"}
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
