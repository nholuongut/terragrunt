package test_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nholuongut/terragrunt/options"
	"github.com/nholuongut/terragrunt/shell"
	"github.com/nholuongut/terragrunt/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	terragruntDebugFile = "terragrunt-debug.tfvars.json"

	fixtureMultiIncludeDependency = "fixtures/multiinclude-dependency"
	fixtureRenderJSON             = "fixtures/render-json"
	fixtureRenderJSONRegression   = "fixtures/render-json-regression"
)

var (
	fixtureRenderJSONMainModulePath = filepath.Join(fixtureRenderJSON, "main")
	fixtureRenderJSONDepModulePath  = filepath.Join(fixtureRenderJSON, "dep")
)

func TestDebugGeneratedInputs(t *testing.T) {
	t.Parallel()

	cleanupTerraformFolder(t, testFixtureInputs)
	tmpEnvPath := copyEnvironment(t, testFixtureInputs)
	rootPath := util.JoinPath(tmpEnvPath, testFixtureInputs)

	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}

	require.NoError(
		t,
		runTerragruntCommand(t, "terragrunt plan --terragrunt-non-interactive --terragrunt-log-level debug --terragrunt-debug --terragrunt-working-dir "+rootPath, &stdout, &stderr),
	)

	debugFile := util.JoinPath(rootPath, terragruntDebugFile)
	assert.True(t, util.FileExists(debugFile))

	assert.Contains(t, stderr.String(), fmt.Sprintf("-chdir=\"%s\"", getPathRelativeTo(t, rootPath, rootPath)))

	// If the debug file is generated correctly, we should be able to run terraform apply using the generated var file
	// without going through terragrunt.
	mockOptions, err := options.NewTerragruntOptionsForTest("integration_test")
	require.NoError(t, err)
	mockOptions.WorkingDir = rootPath
	require.NoError(
		t,
		shell.RunTerraformCommand(context.Background(), mockOptions, "apply", "-auto-approve", "-var-file", debugFile),
	)

	stdout = bytes.Buffer{}
	stderr = bytes.Buffer{}
	require.NoError(
		t,
		runTerragruntCommand(t, "terragrunt output -no-color -json --terragrunt-non-interactive --terragrunt-working-dir "+rootPath, &stdout, &stderr),
	)

	outputs := map[string]TerraformOutput{}
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &outputs))
	validateInputs(t, outputs)

	// Also make sure the undefined variable is not included in the json file
	debugJSONContents, err := os.ReadFile(debugFile)
	require.NoError(t, err)
	var data map[string]interface{}
	require.NoError(t, json.Unmarshal(debugJSONContents, &data))
	_, isDefined := data["undefined_var"]
	assert.False(t, isDefined)
}

func TestTerragruntInputsWithDashes(t *testing.T) {
	t.Parallel()

	cleanupTerraformFolder(t, testFixtureInputs)
	tmpEnvPath := copyEnvironment(t, testFixtureInputs)
	rootPath := util.JoinPath(tmpEnvPath, testFixtureInputs)

	runTerragrunt(t, fmt.Sprintf("terragrunt init --terragrunt-working-dir=%s --terragrunt-log-level=debug", rootPath))
}

func TestTerragruntValidateInputs(t *testing.T) {
	t.Parallel()

	moduleDirs, err := filepath.Glob(filepath.Join("fixtures/validate-inputs", "*"))
	require.NoError(t, err)

	for _, module := range moduleDirs {
		// capture range var within range scope so it doesn't change as the tests are spun to the background in the
		// t.Parallel call.
		module := module

		name := filepath.Base(module)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			nameDashSplit := strings.Split(name, "-")
			runTerragruntValidateInputs(t, module, []string{"--terragrunt-strict-validate"}, nameDashSplit[0] == "success")
		})
	}
}

func TestTerragruntValidateInputsWithCLIVars(t *testing.T) {
	t.Parallel()

	moduleDir := filepath.Join("fixtures/validate-inputs", "fail-no-inputs")
	args := []string{"-var=input=from_env"}
	runTerragruntValidateInputs(t, moduleDir, args, true)
}

func TestTerragruntValidateInputsWithCLIVarFile(t *testing.T) {
	t.Parallel()

	curdir, err := os.Getwd()
	require.NoError(t, err)

	moduleDir := filepath.Join("fixtures/validate-inputs", "fail-no-inputs")
	args := []string{fmt.Sprintf("-var-file=%s/fixtures/validate-inputs/success-var-file/varfiles/main.tfvars", curdir)}
	runTerragruntValidateInputs(t, moduleDir, args, true)
}

func TestTerragruntValidateInputsWithStrictMode(t *testing.T) {
	t.Parallel()

	moduleDir := filepath.Join("fixtures/validate-inputs", "success-inputs-only")
	args := []string{"--terragrunt-strict-validate"}
	runTerragruntValidateInputs(t, moduleDir, args, true)
}

func TestTerragruntValidateInputsWithStrictModeDisabledAndUnusedVar(t *testing.T) {
	t.Parallel()

	moduleDir := filepath.Join("fixtures/validate-inputs", "success-inputs-only")
	args := []string{"-var=testvariable=testvalue"}
	runTerragruntValidateInputs(t, moduleDir, args, true)
}

func TestTerragruntValidateInputsWithStrictModeEnabledAndUnusedVar(t *testing.T) {
	t.Parallel()

	moduleDir := filepath.Join("fixtures/validate-inputs", "success-inputs-only")
	args := []string{"-var=testvariable=testvalue", "--terragrunt-strict-validate"}
	runTerragruntValidateInputs(t, moduleDir, args, false)
}

func TestTerragruntValidateInputsWithStrictModeEnabledAndUnusedInputs(t *testing.T) {
	t.Parallel()

	moduleDir := filepath.Join("fixtures/validate-inputs", "fail-unused-inputs")
	cleanupTerraformFolder(t, moduleDir)
	tmpEnvPath, _ := filepath.EvalSymlinks(copyEnvironment(t, moduleDir))
	rootPath := util.JoinPath(tmpEnvPath, moduleDir)

	args := []string{"--terragrunt-strict-validate"}
	runTerragruntValidateInputs(t, rootPath, args, false)
}

func TestTerragruntValidateInputsWithStrictModeDisabledAndUnusedInputs(t *testing.T) {
	t.Parallel()

	moduleDir := filepath.Join("fixtures/validate-inputs", "fail-unused-inputs")
	cleanupTerraformFolder(t, moduleDir)
	tmpEnvPath, _ := filepath.EvalSymlinks(copyEnvironment(t, moduleDir))
	rootPath := util.JoinPath(tmpEnvPath, moduleDir)

	args := []string{}
	runTerragruntValidateInputs(t, rootPath, args, true)
}

func TestRenderJSONConfig(t *testing.T) {
	t.Parallel()

	tmpDir, err := os.MkdirTemp("", "terragrunt-render-json-*")
	require.NoError(t, err)
	jsonOut := filepath.Join(tmpDir, "terragrunt_rendered.json")
	defer os.RemoveAll(tmpDir)

	cleanupTerraformFolder(t, fixtureRenderJSONMainModulePath)
	cleanupTerraformFolder(t, fixtureRenderJSONDepModulePath)

	runTerragrunt(t, "terragrunt run-all apply -auto-approve --terragrunt-non-interactive --terragrunt-log-level debug --terragrunt-working-dir "+fixtureRenderJSON)
	runTerragrunt(t, fmt.Sprintf("terragrunt render-json --terragrunt-non-interactive --terragrunt-log-level debug --terragrunt-working-dir %s --terragrunt-json-out %s", fixtureRenderJSONMainModulePath, jsonOut))

	jsonBytes, err := os.ReadFile(jsonOut)
	require.NoError(t, err)

	var rendered map[string]interface{}
	require.NoError(t, json.Unmarshal(jsonBytes, &rendered))

	// Make sure all terraform block is visible
	terraformBlock, hasTerraform := rendered["terraform"]
	if assert.True(t, hasTerraform) {
		source, hasSource := terraformBlock.(map[string]interface{})["source"]
		assert.True(t, hasSource)
		assert.Equal(t, "./module", source)
	}

	// Make sure included remoteState is rendered out
	remoteState, hasRemoteState := rendered["remote_state"]
	if assert.True(t, hasRemoteState) {
		assert.Equal(
			t,
			map[string]interface{}{
				"backend": "local",
				"generate": map[string]interface{}{
					"path":      "backend.tf",
					"if_exists": "overwrite_terragrunt",
				},
				"config": map[string]interface{}{
					"path": "foo.tfstate",
				},
				"disable_init":                    false,
				"disable_dependency_optimization": false,
			},
			remoteState.(map[string]interface{}),
		)
	}

	// Make sure dependency blocks are rendered out
	dependencyBlocks, hasDependency := rendered["dependency"]
	if assert.True(t, hasDependency) {
		assert.Equal(
			t,
			map[string]interface{}{
				"dep": map[string]interface{}{
					"name":         "dep",
					"config_path":  "../dep",
					"outputs":      nil,
					"inputs":       nil,
					"mock_outputs": nil,
					"enabled":      nil,
					"mock_outputs_allowed_terraform_commands": nil,
					"mock_outputs_merge_strategy_with_state":  nil,
					"mock_outputs_merge_with_state":           nil,
					"skip":                                    nil,
				},
			},
			dependencyBlocks.(map[string]interface{}),
		)
	}

	// Make sure included generate block is rendered out
	generateBlocks, hasGenerate := rendered["generate"]
	if assert.True(t, hasGenerate) {
		assert.Equal(
			t,
			map[string]interface{}{
				"provider": map[string]interface{}{
					"path":              "provider.tf",
					"comment_prefix":    "# ",
					"disable_signature": false,
					"disable":           false,
					"if_exists":         "overwrite_terragrunt",
					"if_disabled":       "skip",
					"contents": `provider "aws" {
  region = "us-east-1"
}
`,
				},
			},
			generateBlocks.(map[string]interface{}),
		)
	}

	// Make sure all inputs are merged together
	inputsBlock, hasInputs := rendered["inputs"]
	if assert.True(t, hasInputs) {
		assert.Equal(
			t,
			map[string]interface{}{
				"env":        "qa",
				"name":       "dep",
				"type":       "main",
				"aws_region": "us-east-1",
			},
			inputsBlock.(map[string]interface{}),
		)
	}
}

func TestRenderJSONConfigWithIncludesDependenciesAndLocals(t *testing.T) {
	t.Parallel()

	tmpDir, err := os.MkdirTemp("", "terragrunt-render-json-*")
	require.NoError(t, err)
	jsonOut := filepath.Join(tmpDir, "terragrunt_rendered.json")
	defer os.RemoveAll(tmpDir)

	tmpEnvPath := copyEnvironment(t, fixtureRenderJSONRegression)
	workDir := filepath.Join(tmpEnvPath, fixtureRenderJSONRegression)

	runTerragrunt(t, "terragrunt run-all apply -auto-approve --terragrunt-non-interactive --terragrunt-log-level debug --terragrunt-working-dir "+workDir)

	runTerragrunt(t, fmt.Sprintf("terragrunt render-json --terragrunt-non-interactive --terragrunt-log-level debug --terragrunt-working-dir %s --terragrunt-json-out ", workDir)+jsonOut)

	jsonBytes, err := os.ReadFile(jsonOut)
	require.NoError(t, err)

	var rendered map[string]interface{}
	require.NoError(t, json.Unmarshal(jsonBytes, &rendered))

	// Make sure all terraform block is visible
	terraformBlock, hasTerraform := rendered["terraform"]
	if assert.True(t, hasTerraform) {
		source, hasSource := terraformBlock.(map[string]interface{})["source"]
		assert.True(t, hasSource)
		assert.Equal(t, "./foo", source)
	}

	// Make sure top level locals are rendered out
	locals, hasLocals := rendered["locals"]
	if assert.True(t, hasLocals) {
		assert.Equal(
			t,
			map[string]interface{}{
				"foo": "bar",
			},
			locals.(map[string]interface{}),
		)
	}

	// Make sure included dependency block is rendered out, and with the outputs rendered
	dependencyBlocks, hasDependency := rendered["dependency"]
	if assert.True(t, hasDependency) {
		assert.Equal(
			t,
			map[string]interface{}{
				"baz": map[string]interface{}{
					"name":         "baz",
					"config_path":  "./baz",
					"outputs":      nil,
					"inputs":       nil,
					"mock_outputs": nil,
					"enabled":      nil,
					"mock_outputs_allowed_terraform_commands": nil,
					"mock_outputs_merge_strategy_with_state":  nil,
					"mock_outputs_merge_with_state":           nil,
					"skip":                                    nil,
				},
			},
			dependencyBlocks.(map[string]interface{}),
		)
	}

	// Make sure generate block is rendered out
	generateBlocks, hasGenerate := rendered["generate"]
	if assert.True(t, hasGenerate) {
		assert.Equal(
			t,
			map[string]interface{}{
				"provider": map[string]interface{}{
					"path":              "provider.tf",
					"comment_prefix":    "# ",
					"disable_signature": false,
					"disable":           false,
					"if_exists":         "overwrite",
					"if_disabled":       "skip",
					"contents":          "# This is just a test",
				},
			},
			generateBlocks.(map[string]interface{}),
		)
	}

	// Make sure all inputs are merged together
	inputsBlock, hasInputs := rendered["inputs"]
	if assert.True(t, hasInputs) {
		assert.Equal(
			t,
			map[string]interface{}{
				"foo":       "bar",
				"baz":       "blah",
				"another":   "baz",
				"from_root": "Hi",
			},
			inputsBlock.(map[string]interface{}),
		)
	}
}

func TestRenderJSONConfigRunAll(t *testing.T) {
	t.Parallel()

	tmpEnvPath := copyEnvironment(t, fixtureRenderJSONRegression)
	workDir := filepath.Join(tmpEnvPath, fixtureRenderJSONRegression)

	// NOTE: bar is not rendered out because it is considered a parent terragrunt.hcl config.

	bazJSONOut := filepath.Join(workDir, "baz", "terragrunt_rendered.json")
	rootChildJSONOut := filepath.Join(workDir, "terragrunt_rendered.json")

	defer os.Remove(bazJSONOut)
	defer os.Remove(rootChildJSONOut)

	runTerragrunt(t, "terragrunt run-all apply -auto-approve --terragrunt-non-interactive --terragrunt-log-level debug --terragrunt-working-dir "+workDir)

	runTerragrunt(t, "terragrunt run-all render-json --terragrunt-non-interactive --terragrunt-log-level debug --terragrunt-working-dir "+workDir)

	bazJSONBytes, err := os.ReadFile(bazJSONOut)
	require.NoError(t, err)

	var bazRendered map[string]interface{}
	require.NoError(t, json.Unmarshal(bazJSONBytes, &bazRendered))

	// Make sure top level locals are rendered out
	bazLocals, bazHasLocals := bazRendered["locals"]
	if assert.True(t, bazHasLocals) {
		assert.Equal(
			t,
			map[string]interface{}{
				"self": "baz",
			},
			bazLocals.(map[string]interface{}),
		)
	}

	rootChildJSONBytes, err := os.ReadFile(rootChildJSONOut)
	require.NoError(t, err)

	var rootChildRendered map[string]interface{}
	require.NoError(t, json.Unmarshal(rootChildJSONBytes, &rootChildRendered))

	// Make sure top level locals are rendered out
	rootChildLocals, rootChildHasLocals := rootChildRendered["locals"]
	if assert.True(t, rootChildHasLocals) {
		assert.Equal(
			t,
			map[string]interface{}{
				"foo": "bar",
			},
			rootChildLocals.(map[string]interface{}),
		)
	}
}

func TestDependencyGraphWithMultiInclude(t *testing.T) {
	t.Parallel()

	cleanupTerraformFolder(t, fixtureMultiIncludeDependency)
	tmpEnvPath := copyEnvironment(t, fixtureMultiIncludeDependency)
	rootPath := util.JoinPath(tmpEnvPath, fixtureMultiIncludeDependency)

	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}
	require.NoError(
		t,
		runTerragruntCommand(t, "terragrunt graph-dependencies --terragrunt-non-interactive --terragrunt-working-dir "+rootPath, &stdout, &stderr),
	)
	stdoutStr := stdout.String()

	assert.Contains(t, stdoutStr, `"main" -> "depa";`)
	assert.Contains(t, stdoutStr, `"main" -> "depb";`)
	assert.Contains(t, stdoutStr, `"main" -> "depc";`)
}

func runTerragruntValidateInputs(t *testing.T, moduleDir string, extraArgs []string, isSuccessTest bool) {
	t.Helper()

	maybeNested := filepath.Join(moduleDir, "module")
	if util.FileExists(maybeNested) {
		// Nested module test case with included file, so run terragrunt from the nested module.
		moduleDir = maybeNested
	}

	cmd := fmt.Sprintf("terragrunt validate-inputs %s --terragrunt-log-level debug --terragrunt-non-interactive --terragrunt-working-dir %s", strings.Join(extraArgs, " "), moduleDir)
	t.Logf("Command: %s", cmd)
	_, _, err := runTerragruntCommandWithOutput(t, cmd)
	if isSuccessTest {
		require.NoError(t, err)
	} else {
		require.Error(t, err)
	}
}