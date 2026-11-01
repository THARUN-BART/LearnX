package execution

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// LocalRunner executes code using local compilers (gcc, g++, python3, go, java) with timeout contexts.
type LocalRunner struct{}

func NewLocalRunner() *LocalRunner {
	return &LocalRunner{}
}

func (r *LocalRunner) ProviderName() string {
	return "local"
}

func (r *LocalRunner) Execute(ctx context.Context, req ExecutionRequest) (*ExecutionResponse, error) {
	tmpDir, err := os.MkdirTemp("", "learnx-exec-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp execution sandbox: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	timeout := req.CPUTimeLimit
	if timeout <= 0 {
		timeout = 5 * time.Second
	}

	execCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	lang := strings.ToLower(strings.TrimSpace(req.Language))
	switch lang {
	case "python", "py", "python3":
		return r.runPython(execCtx, tmpDir, req.SourceCode, req.Stdin)
	case "cpp", "c++":
		return r.runCpp(execCtx, tmpDir, req.SourceCode, req.Stdin)
	case "c":
		return r.runC(execCtx, tmpDir, req.SourceCode, req.Stdin)
	case "go", "golang":
		return r.runGo(execCtx, tmpDir, req.SourceCode, req.Stdin)
	case "java":
		return r.runJava(execCtx, tmpDir, req.SourceCode, req.Stdin)
	default:
		return nil, ErrUnsupportedLanguage
	}
}

func (r *LocalRunner) runPython(ctx context.Context, tmpDir, sourceCode, stdinData string) (*ExecutionResponse, error) {
	srcPath := filepath.Join(tmpDir, "solution.py")
	if err := os.WriteFile(srcPath, []byte(sourceCode), 0600); err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(ctx, "python3", srcPath)
	return r.executeCommand(ctx, cmd, stdinData)
}

func (r *LocalRunner) runCpp(ctx context.Context, tmpDir, sourceCode, stdinData string) (*ExecutionResponse, error) {
	srcPath := filepath.Join(tmpDir, "solution.cpp")
	binPath := filepath.Join(tmpDir, "solution.out")
	if err := os.WriteFile(srcPath, []byte(sourceCode), 0600); err != nil {
		return nil, err
	}

	// 1. Compile with g++
	compileCmd := exec.CommandContext(ctx, "g++", "-O2", "-o", binPath, srcPath)
	var compileErr bytes.Buffer
	compileCmd.Stderr = &compileErr
	if err := compileCmd.Run(); err != nil {
		return &ExecutionResponse{
			Status:        StatusCompilationError,
			CompileOutput: strings.TrimSpace(compileErr.String()),
			Message:       "Compilation Error",
		}, nil
	}

	// 2. Execute compiled binary
	cmd := exec.CommandContext(ctx, binPath)
	return r.executeCommand(ctx, cmd, stdinData)
}

func (r *LocalRunner) runC(ctx context.Context, tmpDir, sourceCode, stdinData string) (*ExecutionResponse, error) {
	srcPath := filepath.Join(tmpDir, "solution.c")
	binPath := filepath.Join(tmpDir, "solution.out")
	if err := os.WriteFile(srcPath, []byte(sourceCode), 0600); err != nil {
		return nil, err
	}

	// 1. Compile with gcc
	compileCmd := exec.CommandContext(ctx, "gcc", "-O2", "-o", binPath, srcPath, "-lm")
	var compileErr bytes.Buffer
	compileCmd.Stderr = &compileErr
	if err := compileCmd.Run(); err != nil {
		return &ExecutionResponse{
			Status:        StatusCompilationError,
			CompileOutput: strings.TrimSpace(compileErr.String()),
			Message:       "Compilation Error",
		}, nil
	}

	// 2. Execute compiled binary
	cmd := exec.CommandContext(ctx, binPath)
	return r.executeCommand(ctx, cmd, stdinData)
}

func (r *LocalRunner) runGo(ctx context.Context, tmpDir, sourceCode, stdinData string) (*ExecutionResponse, error) {
	srcPath := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(srcPath, []byte(sourceCode), 0600); err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(ctx, "go", "run", srcPath)
	return r.executeCommand(ctx, cmd, stdinData)
}

func (r *LocalRunner) runJava(ctx context.Context, tmpDir, sourceCode, stdinData string) (*ExecutionResponse, error) {
	className := "Main"
	if strings.Contains(sourceCode, "public class Solution") {
		className = "Solution"
	}
	srcPath := filepath.Join(tmpDir, className+".java")
	if err := os.WriteFile(srcPath, []byte(sourceCode), 0600); err != nil {
		return nil, err
	}

	// 1. Compile with javac
	compileCmd := exec.CommandContext(ctx, "javac", srcPath)
	var compileErr bytes.Buffer
	compileCmd.Stderr = &compileErr
	if err := compileCmd.Run(); err != nil {
		return &ExecutionResponse{
			Status:        StatusCompilationError,
			CompileOutput: strings.TrimSpace(compileErr.String()),
			Message:       "Compilation Error",
		}, nil
	}

	// 2. Execute with java
	cmd := exec.CommandContext(ctx, "java", "-cp", tmpDir, className)
	return r.executeCommand(ctx, cmd, stdinData)
}

func (r *LocalRunner) executeCommand(ctx context.Context, cmd *exec.Cmd, stdinData string) (*ExecutionResponse, error) {
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if stdinData != "" {
		cmd.Stdin = strings.NewReader(stdinData)
	}

	startTime := time.Now()
	err := cmd.Run()
	elapsed := time.Since(startTime)
	elapsedMs := int(elapsed.Milliseconds())

	out := &ExecutionResponse{
		TimeMs:   elapsedMs,
		MemoryKB: 1024,
	}

	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) || strings.Contains(err.Error(), "signal: killed") {
			out.Status = StatusTimeLimitExceeded
			out.Message = "Time Limit Exceeded"
			return out, nil
		}

		stderrStr := strings.TrimSpace(stderr.String())
		if strings.Contains(stderrStr, "syntax") || strings.Contains(stderrStr, "expected") {
			out.Status = StatusCompilationError
			out.CompileOutput = stderrStr
			out.Message = stderrStr
			return out, nil
		}

		out.Status = StatusRuntimeError
		out.Stderr = stderrStr
		out.Message = stderrStr
		return out, nil
	}

	out.Status = StatusAccepted
	out.Stdout = stdout.String()
	return out, nil
}
