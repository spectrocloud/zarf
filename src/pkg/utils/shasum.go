// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors

// Package utils provides generic helper functions.
package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"

	"github.com/defenseunicorns/zarf/src/pkg/message"
)

// GetSha256Sum returns the computed SHA256 Sum of a given file.
func GetSha256Sum(path string) (string, error) {
	var data io.ReadCloser
	var err error

	if IsURL(path) {
		// Handle download from URL
		message.Warn("This is a remote source. If a published checksum is available you should use that rather than calculating it directly from the remote link.")
		data = Fetch(path)
	} else {
		// Handle local file
		data, err = os.Open(path)
		if err != nil {
			return "", err
		}
	}

	defer data.Close()

	hash := sha256.New()
	if _, err = io.Copy(hash, data); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
