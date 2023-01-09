// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors

// Package git contains functions for interacting with git repositories.
package git

import (
	"fmt"
	"hash/crc32"
	"regexp"

	"github.com/defenseunicorns/zarf/src/pkg/message"
)

// For further explanation: https://regex101.com/r/zq64q4/1.
var gitURLRegex = regexp.MustCompile(`^(?P<proto>[a-z]+:\/\/)(?P<hostPath>.+?)\/(?P<repo>[\w\-\.]+?)(?P<git>\.git)?(?P<atRef>@(?P<ref>[\w\-\.]+))?$`)

// MutateGitURLsInText changes the gitURL hostname to use the repository Zarf is configured to use.
func (g *Git) MutateGitURLsInText(text string) string {
	extractPathRegex := regexp.MustCompilePOSIX(`https?://[^/]+/(.*\.git)`)
	output := extractPathRegex.ReplaceAllStringFunc(text, func(match string) string {
		output, err := g.transformURL(match)
		if err != nil {
			message.Warnf("Unable to transform the git url, using the original url we have: %s", match)
			output = match
		}
		return output
	})
	return output
}

// TransformURLtoRepoName takes a git url and returns a Zarf-compatible repo name.
func (g *Git) TransformURLtoRepoName(url string) (string, error) {
	matches := gitURLRegex.FindStringSubmatch(url)
	idx := gitURLRegex.SubexpIndex

	if len(matches) == 0 {
		// Unable to find a substring match for the regex
		return "", fmt.Errorf("unable to get extract the repoName from the url %s", url)
	}

	repoName := matches[idx("repo")]
	// NOTE: We remove the .git and protocol so that https://zarf.dev/repo.git and http://zarf.dev/repo
	// resolve to the same repp (as they would in real life)
	sanitizedURL := fmt.Sprintf("%s/%s%s", matches[idx("hostPath")], repoName, matches[idx("atRef")])

	// Add crc32 hash of the repoName to the end of the repo
	table := crc32.MakeTable(crc32.IEEE)
	checksum := crc32.Checksum([]byte(sanitizedURL), table)
	newRepoName := fmt.Sprintf("%s-%d", repoName, checksum)

	return newRepoName, nil
}

func (g *Git) transformURL(url string) (string, error) {
	repoName, err := g.TransformURLtoRepoName(url)
	if err != nil {
		return "", err
	}
	output := fmt.Sprintf("%s/%s/%s", g.Server.Address, g.Server.PushUsername, repoName)
	message.Debugf("Rewrite git URL: %s -> %s", url, output)
	return output, nil
}
