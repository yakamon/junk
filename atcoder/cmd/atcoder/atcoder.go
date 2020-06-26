package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	// parse args
	cmdName := os.Args[0]
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s <filename>\n", cmdName)
		os.Exit(1)
	}
	src := os.Args[1]

	inputSrcList := findInputSrc(src)
	buildAndRun(src, inputSrcList)

	if err := copyToClipboard(src); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func findInputSrc(src string) (inputSrcList []string) {
	dir := filepath.Dir(src)
	inputDir := filepath.Join(dir, "input")
	err := filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		inputSrcList = append(inputSrcList, path)
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() failed: %v", err)
		os.Exit(1)
	}
	return inputSrcList
}

func buildAndRun(src string, inputSrcList []string) {
	if len(inputSrcList) == 0 {
		return
	}

	fmt.Println(inputSrcList)

	switch filepath.Ext(src) {
	case ".cpp":
		if err := buildAndRunCpp(src, inputSrcList); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case ".go":
		if err := buildAndRunGo(src, inputSrcList); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case ".py":
		if err := buildAndRunPy(src, inputSrcList); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func buildAndRunCpp(src string, inputSrcList []string) error {
	out := strings.TrimSuffix(src, filepath.Ext(src))

	buildCmd := stdoutCommand("g++", "-o", out, src)
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("failed to build C++ file: %v", err)
	}
	defer removeFile(out)

	for _, inputSrc := range inputSrcList {
		cmdStr := fmt.Sprintf("%s < %s", out, inputSrc)
		runCmd := stdoutCommand("sh", "-c", cmdStr)
		if err := runCommand(runCmd, inputSrc); err != nil {
			return fmt.Errorf("some error occurred while running C++ program: %v", err)
		}
	}
	return nil
}

func buildAndRunGo(src string, inputSrcList []string) error {
	out := strings.TrimSuffix(src, filepath.Ext(src))

	buildCmd := stdoutCommand("go", "build", "-o", out, src)
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("failed to build Go file: %v", err)
	}
	defer removeFile(out)

	for _, inputSrc := range inputSrcList {
		cmdStr := fmt.Sprintf("%s < %s", out, inputSrc)
		runCmd := stdoutCommand("sh", "-c", cmdStr)
		if err := runCommand(runCmd, inputSrc); err != nil {
			return fmt.Errorf("some error occurred while running Go program: %v", err)
		}
	}
	return nil
}

func buildAndRunPy(src string, inputSrcList []string) error {
	for _, inputSrc := range inputSrcList {
		cmdStr := fmt.Sprintf("python %s < %s", src, inputSrc)
		runCmd := stdoutCommand("sh", "-c", cmdStr)
		if err := runCommand(runCmd, inputSrc); err != nil {
			return fmt.Errorf("some error occurred while running Python program: %v", err)
		}
	}
	return nil
}

func stdoutCommand(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	return cmd
}

func runCommand(cmd *exec.Cmd, inputSrc string) error {
	input, err := ioutil.ReadFile(inputSrc)
	if err != nil {
		return fmt.Errorf("failed to read input file: %v", err)
	}
	printBar("=", inputSrc)
	printBar("-", "input")
	fmt.Println(string(input))
	printBar("-", "output")
	result, err := cmd.Output()
	if err != nil {
		return err
	}
	fmt.Println(string(result))
	printBar("=", "")
	return nil
}

func printBar(rep, message string) {
	bar := bytes.Repeat([]byte(rep), 100)
	barLen := len(string(bar))
	messageLen := len(message)
	if messageLen == 0 {
		fmt.Println(string(bar))
		return
	}
	start := (barLen - messageLen - 1) / 2
	copy(bar[start:start+messageLen+2], " "+message+" ")
	fmt.Println(string(bar))
}

func removeFile(src string) {
	if _, err := os.Open(src); err != nil {
		fmt.Printf("os.Open(%s) failed: %v", src, err)
	} else {
		if err := os.Remove(src); err != nil {
			fmt.Printf("os.Remove(%s) failed: %v", src, err)
		}
	}
}

func copyToClipboard(src string) error {
	cmdStr := fmt.Sprintf("cat %s | pbcopy", src)
	copyCmd := stdoutCommand("sh", "-c", cmdStr)
	if err := copyCmd.Run(); err != nil {
		return fmt.Errorf("src copy failed: %v", err)
	}
	return nil
}
