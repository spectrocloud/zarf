// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors

// Package test provides e2e tests for Zarf.
package test

import (
	"context"
	"testing"
	"time"

	"github.com/defenseunicorns/zarf/src/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestZarfInit(t *testing.T) {
	t.Log("E2E: Zarf init (limit to 10 minutes)")
	e2e.setupWithCluster(t)
	defer e2e.teardown(t)

	initComponents := "logging,git-server"
	// Add k3s component in appliance mode
	if e2e.applianceMode {
		initComponents = "k3s,logging,git-server"
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Minute)
	defer cancel()

	// run `zarf init`
	_, _, err := utils.ExecCommandWithContext(ctx, true, e2e.zarfBinPath, "init", "--components="+initComponents, "--confirm")
	require.NoError(t, err)

	// Check that gitea is actually running and healthy
	stdOut, _, err := utils.ExecCommandWithContext(ctx, true, "kubectl", "get", "pods", "-l", "app in (gitea)", "-n", "zarf", "-o", "jsonpath={.items[*].status.phase}")
	require.NoError(t, err)
	require.Contains(t, stdOut, "Running")

	// Check that the logging stack is actually running and healthy
	stdOut, _, err = utils.ExecCommandWithContext(ctx, true, "kubectl", "get", "pods", "-l", "app in (loki)", "-n", "zarf", "-o", "jsonpath={.items[*].status.phase}")
	require.NoError(t, err)
	require.Contains(t, stdOut, "Running")
	stdOut, _, err = utils.ExecCommandWithContext(ctx, true, "kubectl", "get", "pods", "-l", "app.kubernetes.io/name in (grafana)", "-n", "zarf", "-o", "jsonpath={.items[*].status.phase}")
	require.NoError(t, err)
	require.Contains(t, stdOut, "Running")
	stdOut, _, err = utils.ExecCommandWithContext(ctx, true, "kubectl", "get", "pods", "-l", "app.kubernetes.io/name in (promtail)", "-n", "zarf", "-o", "jsonpath={.items[*].status.phase}")
	require.NoError(t, err)
	require.Contains(t, stdOut, "Running")

	// Special sizing-hacking for reducing resources where Kind + CI eats a lot of free cycles (ignore errors)
	_, _, _ = utils.ExecCommandWithContext(ctx, true, "kubectl", "scale", "deploy", "-n", "kube-system", "coredns", "--replicas=1")
	_, _, _ = utils.ExecCommandWithContext(ctx, true, "kubectl", "scale", "deploy", "-n", "zarf", "agent-hook", "--replicas=1")
}
