// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors

// Package test provides e2e tests for Zarf.
package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestMismatchedArchitectures ensures that zarf produces an error
// when the package architecture doesn't match the target cluster architecture.
func TestMismatchedArchitectures(t *testing.T) {
	t.Log("E2E: Mismatched architectures")
	e2e.SetupWithCluster(t)
	defer e2e.Teardown(t)

	var (
		mismatchedArch         = e2e.GetMismatchedArch()
		mismatchedGamesPackage = fmt.Sprintf("zarf-package-dos-games-%s.tar.zst", mismatchedArch)
		initPackageVersion     = "UnknownVersion"
		mismatchedInitPackage  = fmt.Sprintf("zarf-init-%s-%s.tar.zst", mismatchedArch, initPackageVersion)
		expectedErrorMessage   = fmt.Sprintf("this package architecture is %s", mismatchedArch)
	)

	// Build init package with different arch than the cluster arch.
	stdOut, stdErr, err := e2e.ExecZarfCommand("package", "create", ".", "--architecture", mismatchedArch, "--confirm")
	require.NoError(t, err, stdOut, stdErr)
	defer e2e.CleanFiles(mismatchedInitPackage)

	// Build dos-games package with different arch than the cluster arch.
	stdOut, stdErr, err = e2e.ExecZarfCommand("package", "create", "examples/dos-games/", "--architecture", mismatchedArch, "--confirm")
	require.NoError(t, err, stdOut, stdErr)
	defer e2e.CleanFiles(mismatchedGamesPackage)

	// Ensure zarf init returns an error because of the mismatched architectures.
	// We need to use the --architecture flag here to force zarf to find the package.
	_, stdErr, err = e2e.ExecZarfCommand("init", "--architecture", mismatchedArch, "--confirm")
	require.Error(t, err, stdErr)
	require.Contains(t, stdErr, expectedErrorMessage)

	// Ensure zarf package deploy returns an error because of the mismatched architectures.
	_, stdErr, err = e2e.ExecZarfCommand("package", "deploy", mismatchedGamesPackage, "--confirm")
	require.Error(t, err, stdErr)
	require.Contains(t, stdErr, expectedErrorMessage)
}
