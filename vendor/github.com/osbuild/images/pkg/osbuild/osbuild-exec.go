package osbuild

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/osbuild/images/data/dependencies"
	"github.com/osbuild/images/pkg/datasizes"

	"github.com/hashicorp/go-version"
)

// Run an instance of osbuild, returning a parsed osbuild.Result.
//
// Note that osbuild returns non-zero when the pipeline fails. This function
// does not return an error in this case. Instead, the failure is communicated
// with its corresponding logs through osbuild.Result.
func RunOSBuild(manifest []byte, store, outputDirectory string, exports, checkpoints, extraEnv []string, result bool, errorWriter io.Writer) (*Result, error) {
	if err := CheckMinimumOSBuildVersion(); err != nil {
		return nil, err
	}

	var stdoutBuffer bytes.Buffer
	var res Result

	cmd := exec.Command(
		"osbuild",
		"--store", store,
		"--output-directory", outputDirectory,
		"-",
	)

	for _, export := range exports {
		cmd.Args = append(cmd.Args, "--export", export)
	}

	for _, checkpoint := range checkpoints {
		cmd.Args = append(cmd.Args, "--checkpoint", checkpoint)
	}

	if len(checkpoints) > 0 {
		// set the cache-max-size to a reasonable size that the checkpoints actually get stored
		cmd.Args = append(cmd.Args, "--cache-max-size", fmt.Sprint(20*datasizes.GiB))
	}

	if result {
		cmd.Args = append(cmd.Args, "--json")
		cmd.Stdout = &stdoutBuffer
	} else {
		cmd.Stdout = os.Stdout
	}

	if len(extraEnv) > 0 {
		cmd.Env = append(os.Environ(), extraEnv...)
	}

	cmd.Stderr = errorWriter
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("error setting up stdin for osbuild: %v", err)
	}

	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("error starting osbuild: %v", err)
	}

	_, err = stdin.Write(manifest)
	if err != nil {
		return nil, fmt.Errorf("error writing osbuild manifest: %v", err)
	}

	err = stdin.Close()
	if err != nil {
		return nil, fmt.Errorf("error closing osbuild's stdin: %v", err)
	}

	err = cmd.Wait()

	if result {
		// try to decode the output even though the job could have failed
		if stdoutBuffer.Len() == 0 {
			return nil, fmt.Errorf("osbuild did not return any output")
		}
		decodeErr := json.Unmarshal(stdoutBuffer.Bytes(), &res)
		if decodeErr != nil {
			return nil, fmt.Errorf("error decoding osbuild output: %v\nthe raw output:\n%s", decodeErr, stdoutBuffer.String())
		}
	}

	if err != nil {
		// ignore ExitError if output could be decoded correctly (only if running with --json)
		if _, isExitError := err.(*exec.ExitError); !isExitError || !result {
			return nil, fmt.Errorf("running osbuild failed: %v", err)
		}
	}

	return &res, nil
}

func CheckMinimumOSBuildVersion() error {
	osbuildVersion, err := OSBuildVersion()
	if err != nil {
		return fmt.Errorf("error getting osbuild version: %v", err)
	}

	minVersion, err := version.NewVersion(dependencies.MinimumOSBuildVersion())
	if err != nil {
		return fmt.Errorf("error parsing minimum osbuild version: %v", err)
	}

	currentVersion, err := version.NewVersion(osbuildVersion)
	if err != nil {
		return fmt.Errorf("error parsing current osbuild version: %v", err)
	}

	if currentVersion.LessThan(minVersion) {
		return fmt.Errorf("osbuild version %q is lower than the minimum required version %q",
			osbuildVersion, dependencies.MinimumOSBuildVersion())
	}

	return nil
}

// OSBuildVersion returns the version of osbuild.
func OSBuildVersion() (string, error) {
	var stdoutBuffer bytes.Buffer
	cmd := exec.Command("osbuild", "--version")
	cmd.Stdout = &stdoutBuffer

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("running osbuild failed: %v", err)
	}

	// osbuild --version prints the version in the form of "osbuild VERSION". Extract the version.
	version := strings.TrimPrefix(stdoutBuffer.String(), "osbuild ")
	// Remove the trailing newline.
	version = strings.TrimSpace(version)
	return version, nil
}
