package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/THARUN-BART/LearnX/backend/internal/ai"
	"github.com/THARUN-BART/LearnX/backend/internal/aitutor"
	"github.com/THARUN-BART/LearnX/backend/internal/auth"
	"github.com/THARUN-BART/LearnX/backend/internal/dsa/model"
	"github.com/THARUN-BART/LearnX/backend/internal/execution"
	"github.com/THARUN-BART/LearnX/backend/internal/middleware"
	"github.com/THARUN-BART/LearnX/backend/internal/mistakes"
	"github.com/THARUN-BART/LearnX/backend/internal/otp"
	"github.com/THARUN-BART/LearnX/backend/internal/problems"
	"github.com/THARUN-BART/LearnX/backend/internal/profile"
	"github.com/THARUN-BART/LearnX/backend/internal/server"
	"github.com/THARUN-BART/LearnX/backend/internal/submissions"
	"github.com/THARUN-BART/LearnX/backend/internal/supabase"
	"github.com/THARUN-BART/LearnX/backend/internal/testcases"
)

func setupAITutorTestEnv(t *testing.T) (http.Handler, string, string) {
	studentID := "stu-ai-42"
	token := "mock-jwt-token-for-" + studentID
	probID := "prob-two-sum"

	// Mock LLM & repos
	mockLLM := ai.NewMockProvider()
	aiRepo := aitutor.NewMemoryRepository()
	mistakesRepo := mistakes.NewMemoryRepository()
	mistakesSvc := mistakes.NewService(mistakesRepo, mistakes.NewDetector(), nil)

	// DSA problems repo with a sample problem
	memoryModel := model.NewMemory()
	memoryModel.Problems[probID] = &model.Problem{
		ID:         probID,
		Title:      "Two Sum",
		Slug:       "two-sum",
		Summary:    "Find two indices that sum up to target.",
		Difficulty: model.DifficultyEasy,
	}
	problemsRepo := problems.NewMemoryRepository(memoryModel)

	// TestCases repo with 1 public sample and 1 hidden case
	testCasesRepo := testcases.NewMemoryRepository()
	_ = testCasesRepo.Create(context.Background(), &testcases.TestCase{
		ID:             "tc-sample-1",
		ProblemID:      probID,
		InputData:      "[2,7,11,15], target = 9",
		ExpectedOutput: "[0,1]",
		IsHidden:       false,
		TestOrder:      1,
	})
	_ = testCasesRepo.Create(context.Background(), &testcases.TestCase{
		ID:             "tc-hidden-secret",
		ProblemID:      probID,
		InputData:      "[3,2,4], target = 6 (SECRET_HIDDEN_CASE)",
		ExpectedOutput: "[1,2]",
		IsHidden:       true,
		TestOrder:      2,
	})

	submissionsRepo := submissions.NewMemoryRepository()
	_ = submissionsRepo.Create(context.Background(), &submissions.Submission{
		ID:           "sub-fail-1",
		StudentID:    studentID,
		ProblemID:    probID,
		Language:     "python",
		Status:       execution.StatusWrongAnswer,
		ErrorMessage: "AssertionError: expected [0, 1] but got None",
		SourceCode:   "def twoSum(nums, target): return None",
	})

	ctxBuilder := aitutor.NewContextBuilder(problemsRepo, testCasesRepo, submissionsRepo)
	aiSvc := aitutor.NewService(mockLLM, aiRepo, mistakesSvc, ctxBuilder, nil)
	aiHandler := aitutor.NewHandler(aiSvc)

	sbClient := newRouterMockSupabase()
	sbClient.users["student@learnx.com"] = &supabase.SupabaseUser{
		ID:    studentID,
		Email: "student@learnx.com",
	}

	authRepo := auth.NewMemoryRepository()
	authSvc := auth.NewService(authRepo, otp.NewService(&routerMockOTPClient{}, time.Minute, 5), sbClient, []byte("32-byte-key-for-aes-256-encryption"), 10*time.Minute, 15*time.Minute)
	profileSvc := profile.NewService(profile.NewMemoryRepository())
	authMW := middleware.NewAuthMiddleware(sbClient, "")

	router := server.BuildRouter(
		nil,
		auth.NewHandler(authSvc),
		profile.NewHandler(profileSvc),
		authMW,
		nil,
		server.WithIntelligence(&server.IntelligenceHandlers{
			AITutor: aiHandler,
		}),
	)

	return router, token, probID
}

func TestAITutor_ProgressiveHints(t *testing.T) {
	router, token, probID := setupAITutorTestEnv(t)

	// Request Level 1 hint
	body, _ := json.Marshal(map[string]interface{}{
		"problem_id": probID,
		"level":      1,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/ai/hint", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 for level 1 hint, got %d: %s", rec.Code, rec.Body.String())
	}

	var res struct {
		Success bool                 `json:"success"`
		Data    aitutor.HintResponse `json:"data"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &res)
	if res.Data.Level != 1 {
		t.Errorf("expected level 1, got %d", res.Data.Level)
	}
	if !res.Data.NextLevelAvailable {
		t.Errorf("expected next level to be available for level 1")
	}

	// Request Invalid Level (0 and 6)
	for _, badLevel := range []int{0, 6} {
		badBody, _ := json.Marshal(map[string]interface{}{
			"problem_id": probID,
			"level":      badLevel,
		})
		badReq := httptest.NewRequest(http.MethodPost, "/api/v1/ai/hint", bytes.NewReader(badBody))
		badReq.Header.Set("Authorization", "Bearer "+token)
		badRec := httptest.NewRecorder()
		router.ServeHTTP(badRec, badReq)
		if badRec.Code != http.StatusBadRequest {
			t.Errorf("expected 400 for bad hint level %d, got %d", badLevel, badRec.Code)
		}
	}

	// Verify unlocked hints endpoint
	reqList := httptest.NewRequest(http.MethodGet, "/api/v1/problems/"+probID+"/hints", nil)
	reqList.Header.Set("Authorization", "Bearer "+token)
	recList := httptest.NewRecorder()
	router.ServeHTTP(recList, reqList)

	if recList.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recList.Code)
	}
	var listRes struct {
		Data []aitutor.AIHintHistory `json:"data"`
	}
	_ = json.Unmarshal(recList.Body.Bytes(), &listRes)
	if len(listRes.Data) != 1 || listRes.Data[0].HintLevel != 1 {
		t.Errorf("expected 1 unlocked hint at level 1, got %+v", listRes.Data)
	}
}

func TestAITutor_WalkthroughDryRunAndComplexity(t *testing.T) {
	router, token, probID := setupAITutorTestEnv(t)

	// 1. Walkthrough
	wtBody, _ := json.Marshal(map[string]interface{}{
		"problem_id":  probID,
		"source_code": "def twoSum(nums, target): return [0, 1]",
		"language":    "python",
	})
	reqWT := httptest.NewRequest(http.MethodPost, "/api/v1/ai/walkthrough", bytes.NewReader(wtBody))
	reqWT.Header.Set("Authorization", "Bearer "+token)
	recWT := httptest.NewRecorder()
	router.ServeHTTP(recWT, reqWT)

	if recWT.Code != http.StatusOK {
		t.Fatalf("walkthrough failed: %d", recWT.Code)
	}

	// 2. Dry Run
	dryBody, _ := json.Marshal(map[string]interface{}{
		"problem_id":  probID,
		"source_code": "def twoSum(nums, target): return [0, 1]",
		"language":    "python",
	})
	reqDry := httptest.NewRequest(http.MethodPost, "/api/v1/ai/dry-run", bytes.NewReader(dryBody))
	reqDry.Header.Set("Authorization", "Bearer "+token)
	recDry := httptest.NewRecorder()
	router.ServeHTTP(recDry, reqDry)

	if recDry.Code != http.StatusOK {
		t.Fatalf("dry run failed: %d", recDry.Code)
	}
	var dryRes struct {
		Data aitutor.DryRunResponse `json:"data"`
	}
	_ = json.Unmarshal(recDry.Body.Bytes(), &dryRes)
	if len(dryRes.Data.Steps) == 0 {
		t.Errorf("expected simulated dry-run steps, got 0")
	}

	// 3. Complexity
	cBody, _ := json.Marshal(map[string]interface{}{
		"problem_id":  probID,
		"source_code": "def twoSum(nums, target): return [0, 1]",
		"language":    "python",
	})
	reqC := httptest.NewRequest(http.MethodPost, "/api/v1/ai/complexity", bytes.NewReader(cBody))
	reqC.Header.Set("Authorization", "Bearer "+token)
	recC := httptest.NewRecorder()
	router.ServeHTTP(recC, reqC)

	if recC.Code != http.StatusOK {
		t.Fatalf("complexity failed: %d", recC.Code)
	}
	var cRes struct {
		Data aitutor.ComplexityResponse `json:"data"`
	}
	_ = json.Unmarshal(recC.Body.Bytes(), &cRes)
	if cRes.Data.TimeComplexity == "" {
		t.Errorf("expected non-empty time complexity")
	}
}

func TestAITutor_ConversationsAndAsk(t *testing.T) {
	router, token, probID := setupAITutorTestEnv(t)

	// Ask tutor initial question
	askBody, _ := json.Marshal(map[string]interface{}{
		"problem_id": probID,
		"question":   "Why do we use a hash map instead of a nested loop?",
	})
	reqAsk := httptest.NewRequest(http.MethodPost, "/api/v1/ai/ask", bytes.NewReader(askBody))
	reqAsk.Header.Set("Authorization", "Bearer "+token)
	recAsk := httptest.NewRecorder()
	router.ServeHTTP(recAsk, reqAsk)

	if recAsk.Code != http.StatusOK {
		t.Fatalf("ask tutor failed: %d: %s", recAsk.Code, recAsk.Body.String())
	}
	var askRes struct {
		Data aitutor.AskResponse `json:"data"`
	}
	_ = json.Unmarshal(recAsk.Body.Bytes(), &askRes)
	convID := askRes.Data.ConversationID
	if convID == "" {
		t.Fatal("expected conversation ID to be returned")
	}

	// Follow-up question in the same conversation
	askBody2, _ := json.Marshal(map[string]interface{}{
		"conversation_id": convID,
		"question":        "What is the space complexity overhead of that map?",
	})
	reqAsk2 := httptest.NewRequest(http.MethodPost, "/api/v1/ai/ask", bytes.NewReader(askBody2))
	reqAsk2.Header.Set("Authorization", "Bearer "+token)
	recAsk2 := httptest.NewRecorder()
	router.ServeHTTP(recAsk2, reqAsk2)

	if recAsk2.Code != http.StatusOK {
		t.Fatalf("follow-up ask failed: %d", recAsk2.Code)
	}

	// List messages in conversation
	reqMsgs := httptest.NewRequest(http.MethodGet, "/api/v1/ai/conversations/"+convID, nil)
	reqMsgs.Header.Set("Authorization", "Bearer "+token)
	recMsgs := httptest.NewRecorder()
	router.ServeHTTP(recMsgs, reqMsgs)

	if recMsgs.Code != http.StatusOK {
		t.Fatalf("list messages failed: %d", recMsgs.Code)
	}
	var msgsRes struct {
		Data []aitutor.AIMessage `json:"data"`
	}
	_ = json.Unmarshal(recMsgs.Body.Bytes(), &msgsRes)
	// Should contain 2 user messages + 2 assistant messages = 4 messages
	if len(msgsRes.Data) != 4 {
		t.Errorf("expected 4 messages in conversation, got %d", len(msgsRes.Data))
	}
}

func TestAITutor_ContextBuilderHiddenTestCaseIsolation(t *testing.T) {
	probID := "prob-security-test"
	testCasesRepo := testcases.NewMemoryRepository()

	// 1 public sample
	_ = testCasesRepo.Create(context.Background(), &testcases.TestCase{
		ID:             "tc-pub",
		ProblemID:      probID,
		InputData:      "PUBLIC_INPUT_SAMPLE",
		ExpectedOutput: "PUBLIC_OUTPUT",
		IsHidden:       false,
	})
	// 1 secret hidden testcase
	_ = testCasesRepo.Create(context.Background(), &testcases.TestCase{
		ID:             "tc-secret",
		ProblemID:      probID,
		InputData:      "SECRET_TOP_SECRET_INPUT_CASE",
		ExpectedOutput: "SECRET_TOP_SECRET_OUTPUT",
		IsHidden:       true,
	})

	ctxBuilder := aitutor.NewContextBuilder(nil, testCasesRepo, nil)
	lc, err := ctxBuilder.Build(context.Background(), "stu-1", probID, "")
	if err != nil {
		t.Fatalf("ContextBuilder error: %v", err)
	}

	if len(lc.SampleTestCases) != 1 {
		t.Fatalf("expected exactly 1 sample testcase, got %d", len(lc.SampleTestCases))
	}
	if lc.SampleTestCases[0].InputData != "PUBLIC_INPUT_SAMPLE" {
		t.Errorf("expected public sample, got %s", lc.SampleTestCases[0].InputData)
	}

	samplesFormatted := aitutor.FormatSamples(lc.SampleTestCases)
	if bytes.Contains([]byte(samplesFormatted), []byte("SECRET_TOP_SECRET")) {
		t.Errorf("CRITICAL SECURITY VIOLATION: Hidden test case data was leaked into prompt context!")
	}
}

func TestGeminiProvider_SuccessAndErrorHandling(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("key") != "test-gemini-key" {
			http.Error(w, `{"error":{"message":"API key not valid"}}`, http.StatusUnauthorized)
			return
		}

		resp := map[string]interface{}{
			"candidates": []map[string]interface{}{
				{
					"content": map[string]interface{}{
						"parts": []map[string]interface{}{
							{"text": "Consider sorting the array first and using two pointers."},
						},
						"role": "model",
					},
					"finishReason": "STOP",
				},
			},
			"usageMetadata": map[string]interface{}{
				"promptTokenCount":     42,
				"candidatesTokenCount": 18,
				"totalTokenCount":      60,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	provider := ai.NewGeminiProvider("test-gemini-key", ts.URL, "gemini-1.5-flash", 5*time.Second)
	if provider.Name() != "gemini" {
		t.Errorf("expected provider name 'gemini', got %s", provider.Name())
	}

	req := ai.ChatCompletionRequest{
		Messages: []ai.ChatMessage{
			{Role: "system", Content: "You are an expert DSA tutor."},
			{Role: "user", Content: "Give me an approach for 3Sum."},
		},
		ResponseFormat: "json_object",
		Temperature:    0.2,
		MaxTokens:      500,
	}

	resp, err := provider.Generate(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error from gemini provider: %v", err)
	}

	if !bytes.Contains([]byte(resp.Content), []byte("two pointers")) {
		t.Errorf("expected content to contain 'two pointers', got %s", resp.Content)
	}
	if resp.Usage.TotalTokens != 60 {
		t.Errorf("expected 60 total tokens, got %d", resp.Usage.TotalTokens)
	}

	// Test missing API key
	emptyProvider := ai.NewGeminiProvider("", ts.URL, "gemini-1.5-flash", time.Second)
	_, err = emptyProvider.Generate(context.Background(), req)
	if err == nil {
		t.Errorf("expected error when api key is missing")
	}
}
