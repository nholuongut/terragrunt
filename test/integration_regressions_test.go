package test_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/nholuongut/terragrunt/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testFixtureRegressions = "fixtures/regressions"
)

func TestNoAutoInit(t *testing.T) {
	t.Parallel()

	cleanupTerraformFolder(t, testFixtureRegressions)
	tmpEnvPath := copyEnvironment(t, testFixtureRegressions)
	rootPath := util.JoinPath(tmpEnvPath, testFixtureRegressions, "skip-init")

	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}
	err := runTerragruntCommand(t, "terragrunt apply --terragrunt-no-auto-init --terragrunt-log-level debug --terragrunt-non-interactive --terragrunt-working-dir "+rootPath, &stdout, &stderr)
	logBufferContentsLineByLine(t, stdout, "no force apply stdout")
	logBufferContentsLineByLine(t, stderr, "no force apply stderr")
	require.Error(t, err)
	assert.Contains(t, stderr.String(), "This module is not yet installed.")
}

// Test case for yamldecode bug: https://github.com/nholuongut/terragrunt/issues/834
func TestYamlDecodeRegressions(t *testing.T) {
	t.Parallel()

	cleanupTerraformFolder(t, testFixtureRegressions)
	tmpEnvPath := copyEnvironment(t, testFixtureRegressions)
	rootPath := util.JoinPath(tmpEnvPath, testFixtureRegressions, "yamldecode")

	runTerragrunt(t, "terragrunt apply -auto-approve --terragrunt-non-interactive --terragrunt-working-dir "+rootPath)

	// Check the output of yamldecode and make sure it doesn't parse the string incorrectly
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}

	require.NoError(
		t,
		runTerragruntCommand(t, "terragrunt output -no-color -json --terragrunt-non-interactive --terragrunt-working-dir "+rootPath, &stdout, &stderr),
	)

	outputs := map[string]TerraformOutput{}
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &outputs))
	assert.Equal(t, "003", outputs["test1"].Value)
	assert.Equal(t, "1.00", outputs["test2"].Value)
	assert.Equal(t, "0ba", outputs["test3"].Value)
}

func TestMockOutputsMergeWithState(t *testing.T) {
	t.Parallel()

	cleanupTerraformFolder(t, testFixtureRegressions)
	tmpEnvPath := copyEnvironment(t, testFixtureRegressions)
	rootPath := util.JoinPath(tmpEnvPath, testFixtureRegressions, "mocks-merge-with-state")

	modulePath := util.JoinPath(rootPath, "module")
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}
	err := runTerragruntCommand(t, "terragrunt apply --terragrunt-log-level debug --terragrunt-non-interactive -auto-approve --terragrunt-working-dir "+modulePath, &stdout, &stderr)
	logBufferContentsLineByLine(t, stdout, "module-executed")
	require.NoError(t, err)

	deepMapPath := util.JoinPath(rootPath, "deep-map")
	stdout = bytes.Buffer{}
	stderr = bytes.Buffer{}
	err = runTerragruntCommand(t, "terragrunt apply --terragrunt-log-level debug --terragrunt-non-interactive -auto-approve --terragrunt-working-dir "+deepMapPath, &stdout, &stderr)
	logBufferContentsLineByLine(t, stdout, "deep-map-executed")
	require.NoError(t, err)

	shallowPath := util.JoinPath(rootPath, "shallow")
	stdout = bytes.Buffer{}
	stderr = bytes.Buffer{}
	err = runTerragruntCommand(t, "terragrunt apply --terragrunt-log-level debug --terragrunt-non-interactive -auto-approve --terragrunt-working-dir "+shallowPath, &stdout, &stderr)
	logBufferContentsLineByLine(t, stdout, "shallow-map-executed")
	require.NoError(t, err)
}
