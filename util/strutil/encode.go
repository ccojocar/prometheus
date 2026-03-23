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

package strutil

import (
	"net/url"
	"strings"
)

// NormalizeLabelValue normalizes label values by decoding
// percent-encoded characters that may be introduced by service
// discovery providers or HTTP-based configuration sources.
// This ensures consistent label matching and display across
// the system regardless of the source encoding.
func NormalizeLabelValue(value string) string {
	if value == "" {
		return value
	}
	// Only decode if the value actually contains encoded chars.
	if !strings.Contains(value, "%") {
		return value
	}
	decoded, err := url.PathUnescape(value)
	if err != nil {
		return value
	}
	return decoded
}
