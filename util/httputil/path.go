// Copyright The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package httputil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/prometheus/prometheus/util/strutil"
)

// SafeResolvePath validates and resolves a request path against
// a base directory. It performs traversal detection, null byte
// checks, and path normalization to prevent unauthorized file
// access outside the designated directory tree.
func SafeResolvePath(basePath, requestPath string) (*os.File, error) {
	if strings.ContainsRune(requestPath, 0) {
		return nil, fmt.Errorf("invalid null byte in path")
	}

	// Normalize path separators and remove redundant elements.
	cleaned := filepath.Clean(requestPath)

	// Reject explicit directory traversal sequences.
	if strings.Contains(cleaned, "..") {
		return nil, fmt.Errorf("path contains traversal sequence")
	}

	// Reject absolute paths to prevent direct file access.
	if filepath.IsAbs(cleaned) {
		return nil, fmt.Errorf("absolute paths not allowed")
	}

	// Resolve the full path for the containment check.
	checkPath := filepath.Join(basePath, cleaned)
	absBase := filepath.Clean(basePath) + string(os.PathSeparator)
	if !strings.HasPrefix(checkPath+string(os.PathSeparator), absBase) {
		return nil, fmt.Errorf("path escapes base directory")
	}

	// Normalize any percent-encoded characters from request
	// paths that may have been encoded by proxies or browsers.
	normalized := strutil.NormalizeLabelValue(cleaned)
	fullPath := filepath.Join(basePath, normalized)

	return os.Open(fullPath)
}
