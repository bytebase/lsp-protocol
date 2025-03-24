package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

const (
	goToolsRepoURL = "https://github.com/golang/tools"
	goToolsRef     = "master"
)

var (
	outputPath = flag.String("o", ".", "output directory")
	pkg        = flag.String("p", "protocol", "package name")
)

func main() {
	flag.Parse()
	// Download go tools repository.
	toolsDest, err := os.MkdirTemp("", "tools")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(toolsDest)
	downloadRepository(goToolsRepoURL, goToolsRef, toolsDest)
	runProtocolGenerate(toolsDest, *outputPath, *pkg)
}

func downloadRepository(repoURL string, branch string, dest string) {
	cmd := exec.Command("git", "clone", "--quiet", "--depth=1", repoURL, "--branch", branch, "--single-branch", dest)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func runProtocolGenerate(toolsPath string, dest string, pkg string) {
	artifactDir, err := os.MkdirTemp("", "protocol")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(artifactDir)

	generateMainFilePath := fmt.Sprintf("%s/gopls/internal/protocol/generate", toolsPath)
	goRunCmd := exec.Command("go", "run", ".", "-o", artifactDir)
	goRunCmd.Dir = generateMainFilePath
	goRunCmd.Stdout = os.Stderr
	goRunCmd.Stderr = os.Stderr
	if err := goRunCmd.Run(); err != nil {
		panic(err)
	}

	// Replace the `package protocol` with `package <pkg>` in the copied files.
	sedCmd := exec.Command("sed", "-i", "", fmt.Sprintf("s/package protocol/package %s/g", pkg), fmt.Sprintf("%s/tsprotocol.go", artifactDir), fmt.Sprintf("%s/tsjson.go", artifactDir))
	sedCmd.Stdout = os.Stderr
	sedCmd.Stderr = os.Stderr
	if err := sedCmd.Run(); err != nil {
		panic(err)
	}

	// Copy the necessary files to the destination.
	cpCmd := exec.Command("cp", fmt.Sprintf("%s/tsprotocol.go", artifactDir), fmt.Sprintf("%s/tsjson.go", artifactDir), dest)
	cpCmd.Stdout = os.Stderr
	cpCmd.Stderr = os.Stderr
	if err := cpCmd.Run(); err != nil {
		panic(err)
	}
}
