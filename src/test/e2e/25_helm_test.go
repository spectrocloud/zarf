// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors

// Package test provides e2e tests for Zarf.
package test

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHelm(t *testing.T) {
	t.Log("E2E: Helm chart")
	e2e.setupWithCluster(t)
	defer e2e.teardown(t)

	path := fmt.Sprintf("build/zarf-package-test-helm-releasename-%s.tar.zst", e2e.arch)

	// Deploy the charts
	stdOut, stdErr, err := e2e.execZarfCommand("package", "deploy", path, "--confirm")
	require.NoError(t, err, stdOut, stdErr)

	// Verify multiple helm installs of different release names were deployed
	kubectlOut, _ := exec.Command("kubectl", "get", "pods", "-n=helm-releasename", "--no-headers").Output()
	assert.Contains(t, string(kubectlOut), "cool-name-podinfo")

	stdOut, stdErr, err = e2e.execZarfCommand("package", "remove", "test-helm-releasename", "--confirm")
	require.NoError(t, err, stdOut, stdErr)
}
