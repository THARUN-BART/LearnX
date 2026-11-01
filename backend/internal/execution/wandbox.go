package execution

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// WandboxClient communicates with the public, zero-auth Wandbox compiler API.
type WandboxClient struct {
	httpClient *http.Client
}

func NewWandboxClient(timeout time.Duration) *WandboxClient {
	if timeout <= 0 {
		timeout = 15 * time.Second
	}
	return &WandboxClient{
		httpClient: &http.Client{Timeout: timeout},
	}
}

func (w *WandboxClient) ProviderName() string {
	return "wandbox"
}

type wandboxPayload struct {
	Compiler string `json:"compiler"`
	Code     string `json:"code"`
	Stdin    string `json:"stdin,omitempty"`
}

type wandboxResponse struct {
	Status          string `json:"status"` // "0" for success
	Signal          string `json:"signal"`
	CompilerOutput  string `json:"compiler_output"`
	CompilerError   string `json:"compiler_error"`
	CompilerMessage string `json:"compiler_message"`
	ProgramOutput   string `json:"program_output"`
	ProgramError    string `json:"program_error"`
	ProgramMessage  string `json:"program_message"`
}

func (w *WandboxClient) Execute(ctx context.Context, req ExecutionRequest) (*ExecutionResponse, error) {
	compiler := w.resolveCompiler(req.Language)
	payload := wandboxPayload{
		Compiler: compiler,
		Code:     req.SourceCode,
		Stdin:    req.Stdin,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://wandbox.org/api/compile.json", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	startTime := time.Now()
	resp, err := w.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrExecutionServiceUnavailable, err)
	}
	defer resp.Body.Close()
	elapsedMs := int(time.Since(startTime).Milliseconds())

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var wResp wandboxResponse
	if err := json.Unmarshal(respBytes, &wResp); err != nil {
		return nil, fmt.Errorf("wandbox decode error: %w", err)
	}

	out := &ExecutionResponse{
		TimeMs:   elapsedMs,
		MemoryKB: 1024,
	}

	// 1. Check for compilation error
	if wResp.CompilerError != "" || (wResp.CompilerMessage != "" && wResp.ProgramOutput == "" && wResp.Status != "0") {
		compileMsg := wResp.CompilerError
		if compileMsg == "" {
			compileMsg = wResp.CompilerMessage
		}
		if compileMsg == "" {
			compileMsg = wResp.ProgramError
		}
		out.Status = StatusCompilationError
		out.CompileOutput = compileMsg
		out.Message = compileMsg
		return out, nil
	}

	// 2. Check for runtime signal / error
	if wResp.Signal != "" || (wResp.Status != "0" && wResp.Status != "") {
		if strings.Contains(wResp.Signal, "KILL") || strings.Contains(wResp.Signal, "XCPU") {
			out.Status = StatusTimeLimitExceeded
			return out, nil
		}
		out.Status = StatusRuntimeError
		out.Stderr = wResp.ProgramError
		out.Message = wResp.ProgramError
		return out, nil
	}

	out.Status = StatusAccepted
	out.Stdout = wResp.ProgramOutput
	return out, nil
}

func (w *WandboxClient) resolveCompiler(lang string) string {
	cleaned := strings.ToLower(strings.TrimSpace(lang))
	switch cleaned {
	case "cpp", "c++", "cplusplus":
		return "gcc-13.2.0"
	case "c":
		return "gcc-13.2.0-c"
	case "python", "py", "python3":
		return "cpython-3.12.7"
	case "go", "golang":
		return "go-1.23.2"
	case "java":
		return "openjdk-jdk-21+35"
	default:
		return "gcc-13.2.0"
	}
}
