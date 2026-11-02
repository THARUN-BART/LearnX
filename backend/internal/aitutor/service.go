package aitutor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/THARUN-BART/LearnX/backend/internal/ai"
	"github.com/THARUN-BART/LearnX/backend/internal/mistakes"
	"github.com/THARUN-BART/LearnX/backend/internal/pedagogy"
)

var (
	ErrInvalidHintLevel   = errors.New("hint level must be between 1 and 5")
	ErrProblemNotFound    = errors.New("problem not found")
	ErrSubmissionNotFound = errors.New("submission not found")
	ErrUnauthorized       = errors.New("unauthorized conversation access")
)

type Service interface {
	GetProgressiveHint(ctx context.Context, studentID string, req HintRequest) (*HintResponse, error)
	GetUnlockedHints(ctx context.Context, studentID, problemID string) ([]AIHintHistory, error)
	ExplainMistake(ctx context.Context, studentID string, req ExplainMistakeRequest) (*ExplainMistakeResponse, error)
	WalkthroughCode(ctx context.Context, studentID string, req CodeWalkthroughRequest) (*CodeWalkthroughResponse, error)
	GenerateDryRun(ctx context.Context, studentID string, req DryRunRequest) (*DryRunResponse, error)
	AnalyzeComplexity(ctx context.Context, studentID string, req ComplexityRequest) (*ComplexityResponse, error)
	AskTutor(ctx context.Context, studentID string, req AskRequest) (*AskResponse, error)
	ListConversations(ctx context.Context, studentID string, limit, offset int) ([]AIConversation, error)
	GetConversationMessages(ctx context.Context, studentID, conversationID string) ([]AIMessage, error)
	SetPedagogyClient(client pedagogy.Client)
	SetMemoryService(memSvc interface {
		RecordHintEvent(ctx context.Context, studentID, problemID string, hintLevel int, hintContent string) error
		RecordWalkthroughEvent(ctx context.Context, studentID, problemID string) error
		BuildAIContext(ctx context.Context, studentID, topicID, problemID string) (string, error)
	})
}

type serviceImpl struct {
	llm             ai.LLMProvider
	repo            Repository
	mistakesService mistakes.Service
	ctxBuilder      *ContextBuilder
	logger          *slog.Logger
	pedagogyClient  pedagogy.Client
	memoryService   interface {
		RecordHintEvent(ctx context.Context, studentID, problemID string, hintLevel int, hintContent string) error
		RecordWalkthroughEvent(ctx context.Context, studentID, problemID string) error
		BuildAIContext(ctx context.Context, studentID, topicID, problemID string) (string, error)
	}
}

func NewService(
	llm ai.LLMProvider,
	repo Repository,
	mistakesService mistakes.Service,
	ctxBuilder *ContextBuilder,
	logger *slog.Logger,
) Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &serviceImpl{
		llm:             llm,
		repo:            repo,
		mistakesService: mistakesService,
		ctxBuilder:      ctxBuilder,
		logger:          logger,
	}
}

func (s *serviceImpl) SetPedagogyClient(client pedagogy.Client) {
	s.pedagogyClient = client
}

func (s *serviceImpl) SetMemoryService(memSvc interface {
	RecordHintEvent(ctx context.Context, studentID, problemID string, hintLevel int, hintContent string) error
	RecordWalkthroughEvent(ctx context.Context, studentID, problemID string) error
	BuildAIContext(ctx context.Context, studentID, topicID, problemID string) (string, error)
}) {
	s.memoryService = memSvc
}

func (s *serviceImpl) GetProgressiveHint(ctx context.Context, studentID string, req HintRequest) (*HintResponse, error) {
	if req.Level < 1 || req.Level > 5 {
		return nil, ErrInvalidHintLevel
	}

	lc, err := s.ctxBuilder.Build(ctx, studentID, req.ProblemID, req.SubmissionID)
	if err != nil {
		return nil, err
	}
	if lc.Problem == nil {
		return nil, ErrProblemNotFound
	}

	samplesStr := FormatSamples(lc.SampleTestCases)

	var parsed struct {
		Title           string `json:"title"`
		Content         string `json:"content"`
		GuidingQuestion string `json:"guiding_question"`
	}

	usedPedagogy := false
	isPython := strings.Contains(strings.ToLower(lc.Language), "python") || strings.Contains(strings.ToLower(lc.Language), "py")

	// Pedagogy Microservice Integration for Python DSA Problems
	if s.pedagogyClient != nil && isPython && req.Level >= 1 && req.Level <= 4 && strings.TrimSpace(lc.Code) != "" {
		hintLevelMap := map[int]string{1: "L1", 2: "L2", 3: "L3", 4: "L4"}
		pcRes, pcErr := s.pedagogyClient.Analyze(ctx, pedagogy.AnalyzeRequest{
			StudentCode: lc.Code,
			TestCode:    samplesStr,
			Skill:       "SKL-02",
			HintLevel:   hintLevelMap[req.Level],
			Timeout:     3.0,
		})
		if pcErr == nil && pcRes != nil && pcRes.Hint.Hint != "" {
			usedPedagogy = true
			parsed.Title = fmt.Sprintf("Level %d Hint (%s)", req.Level, pcRes.Diagnosis.Misconception)
			parsed.Content = pcRes.Hint.Hint
			parsed.GuidingQuestion = "Does this align with your approach?"
			s.logger.Info("progressive hint generated via Pedagogy microservice",
				slog.String("misconception", pcRes.Diagnosis.Misconception),
				slog.String("source", pcRes.Hint.HintSource),
			)
		} else if pcErr != nil {
			s.logger.Warn("pedagogy hint generation failed, falling back to LLM", slog.String("error", pcErr.Error()))
		}
	}

	if !usedPedagogy {
		messages := BuildHintMessages(
			req.Level,
			lc.Problem.Title,
			lc.Problem.Summary,
			"Standard DSA constraints",
			samplesStr,
			lc.Code,
			lc.Language,
		)

		llmResp, err := s.llm.Generate(ctx, ai.ChatCompletionRequest{
			Messages:       messages,
			ResponseFormat: "json_object",
			Temperature:    0.2,
			MaxTokens:      600,
		})
		if err != nil {
			s.logger.Error("failed to generate progressive hint from LLM", slog.String("error", err.Error()))
			return nil, err
		}

		if err := json.Unmarshal([]byte(llmResp.Content), &parsed); err != nil {
			// Fallback to plain content
			parsed.Title = fmt.Sprintf("Level %d Hint", req.Level)
			parsed.Content = llmResp.Content
		}
	}

	// Persist unlocked hint in history
	var subIDPtr *string
	if req.SubmissionID != "" {
		subIDPtr = &req.SubmissionID
	}
	_ = s.repo.SaveHintHistory(ctx, &AIHintHistory{
		StudentID:    studentID,
		ProblemID:    req.ProblemID,
		SubmissionID: subIDPtr,
		HintLevel:    req.Level,
		HintContent:  parsed.Content,
	})

	if s.memoryService != nil {
		_ = s.memoryService.RecordHintEvent(ctx, studentID, req.ProblemID, req.Level, parsed.Content)
	}

	return &HintResponse{
		ProblemID:          req.ProblemID,
		Level:              req.Level,
		Title:              parsed.Title,
		Content:            parsed.Content,
		GuidingQuestion:    parsed.GuidingQuestion,
		NextLevelAvailable: req.Level < 5,
		CreatedAt:          time.Now(),
	}, nil
}

func (s *serviceImpl) GetUnlockedHints(ctx context.Context, studentID, problemID string) ([]AIHintHistory, error) {
	return s.repo.ListHintHistory(ctx, studentID, problemID)
}

func (s *serviceImpl) ExplainMistake(ctx context.Context, studentID string, req ExplainMistakeRequest) (*ExplainMistakeResponse, error) {
	diag, err := s.mistakesService.GetDiagnosisBySubmission(ctx, studentID, req.SubmissionID)
	if err != nil {
		return nil, err
	}

	lc, err := s.ctxBuilder.Build(ctx, studentID, diag.ProblemID, req.SubmissionID)
	if err != nil {
		return nil, err
	}

	problemTitle := "DSA Problem"
	if lc.Problem != nil {
		problemTitle = lc.Problem.Title
	}

	errOutput := ""
	if lc.Submission != nil {
		errOutput = lc.Submission.ErrorMessage
	}

	messages := BuildExplainMistakeMessages(
		problemTitle,
		diag.Category,
		diag.Subtype,
		diag.Description,
		diag.SuggestedCorrection,
		errOutput,
		lc.Code,
		lc.Language,
	)

	llmResp, err := s.llm.Generate(ctx, ai.ChatCompletionRequest{
		Messages:       messages,
		ResponseFormat: "json_object",
		Temperature:    0.2,
		MaxTokens:      700,
	})
	if err != nil {
		return nil, err
	}

	var parsed struct {
		WhatWentWrong      string   `json:"what_went_wrong"`
		WhyItHappened      string   `json:"why_it_happened"`
		HowToFix           string   `json:"how_to_fix"`
		RelatedDSAConcepts []string `json:"related_dsa_concepts"`
	}
	if err := json.Unmarshal([]byte(llmResp.Content), &parsed); err != nil {
		parsed.WhatWentWrong = diag.Description
		parsed.HowToFix = diag.SuggestedCorrection
		parsed.RelatedDSAConcepts = []string{diag.RelatedConcept}
	}

	return &ExplainMistakeResponse{
		SubmissionID:       req.SubmissionID,
		Category:           diag.Category,
		Subtype:            diag.Subtype,
		WhatWentWrong:      parsed.WhatWentWrong,
		WhyItHappened:      parsed.WhyItHappened,
		HowToFix:           parsed.HowToFix,
		RelatedDSAConcepts: parsed.RelatedDSAConcepts,
	}, nil
}

func (s *serviceImpl) WalkthroughCode(ctx context.Context, studentID string, req CodeWalkthroughRequest) (*CodeWalkthroughResponse, error) {
	lc, err := s.ctxBuilder.Build(ctx, studentID, req.ProblemID, "")
	if err != nil {
		return nil, err
	}
	problemTitle := "DSA Problem"
	if lc.Problem != nil {
		problemTitle = lc.Problem.Title
	}

	codeToUse := req.SourceCode
	if codeToUse == "" && lc.Code != "" {
		codeToUse = lc.Code
	}
	langToUse := req.Language
	if langToUse == "" && lc.Language != "" {
		langToUse = lc.Language
	}

	messages := BuildWalkthroughMessages(problemTitle, codeToUse, langToUse)
	llmResp, err := s.llm.Generate(ctx, ai.ChatCompletionRequest{
		Messages:       messages,
		ResponseFormat: "json_object",
		Temperature:    0.2,
		MaxTokens:      900,
	})
	if err != nil {
		return nil, err
	}

	var res CodeWalkthroughResponse
	if err := json.Unmarshal([]byte(llmResp.Content), &res); err != nil {
		res.OverallSummary = llmResp.Content
	}

	if s.memoryService != nil && req.ProblemID != "" {
		_ = s.memoryService.RecordWalkthroughEvent(ctx, studentID, req.ProblemID)
	}

	return &res, nil
}

func (s *serviceImpl) GenerateDryRun(ctx context.Context, studentID string, req DryRunRequest) (*DryRunResponse, error) {
	lc, err := s.ctxBuilder.Build(ctx, studentID, req.ProblemID, "")
	if err != nil {
		return nil, err
	}
	problemTitle := "DSA Problem"
	if lc.Problem != nil {
		problemTitle = lc.Problem.Title
	}

	inputTestCase := req.CustomInput
	if inputTestCase == "" && len(lc.SampleTestCases) > 0 {
		inputTestCase = lc.SampleTestCases[0].InputData
	}
	if inputTestCase == "" {
		inputTestCase = "standard problem sample input"
	}

	codeToUse := req.SourceCode
	if codeToUse == "" && lc.Code != "" {
		codeToUse = lc.Code
	}
	langToUse := req.Language
	if langToUse == "" && lc.Language != "" {
		langToUse = lc.Language
	}

	messages := BuildDryRunMessages(problemTitle, inputTestCase, codeToUse, langToUse)
	llmResp, err := s.llm.Generate(ctx, ai.ChatCompletionRequest{
		Messages:       messages,
		ResponseFormat: "json_object",
		Temperature:    0.1,
		MaxTokens:      1000,
	})
	if err != nil {
		return nil, err
	}

	var res DryRunResponse
	if err := json.Unmarshal([]byte(llmResp.Content), &res); err != nil {
		res.TestCaseInput = inputTestCase
		res.FinalOutput = "Executed trace"
	}
	return &res, nil
}

func (s *serviceImpl) AnalyzeComplexity(ctx context.Context, studentID string, req ComplexityRequest) (*ComplexityResponse, error) {
	lc, err := s.ctxBuilder.Build(ctx, studentID, req.ProblemID, "")
	if err != nil {
		return nil, err
	}
	problemTitle := "DSA Problem"
	if lc.Problem != nil {
		problemTitle = lc.Problem.Title
	}

	codeToUse := req.SourceCode
	if codeToUse == "" && lc.Code != "" {
		codeToUse = lc.Code
	}
	langToUse := req.Language
	if langToUse == "" && lc.Language != "" {
		langToUse = lc.Language
	}

	messages := BuildComplexityMessages(problemTitle, "Standard problem constraints", codeToUse, langToUse)
	llmResp, err := s.llm.Generate(ctx, ai.ChatCompletionRequest{
		Messages:       messages,
		ResponseFormat: "json_object",
		Temperature:    0.1,
		MaxTokens:      600,
	})
	if err != nil {
		return nil, err
	}

	var res ComplexityResponse
	if err := json.Unmarshal([]byte(llmResp.Content), &res); err != nil {
		res.TimeComplexity = "O(N)"
		res.TimeExplanation = llmResp.Content
		res.SpaceComplexity = "O(1)"
		res.IsOptimal = true
	}
	return &res, nil
}

func (s *serviceImpl) AskTutor(ctx context.Context, studentID string, req AskRequest) (*AskResponse, error) {
	var conv *AIConversation
	var err error

	if req.ConversationID != "" {
		conv, err = s.repo.GetConversation(ctx, req.ConversationID)
		if err != nil {
			return nil, err
		}
		if conv.StudentID != studentID {
			return nil, ErrUnauthorized
		}
	} else {
		var probIDPtr *string
		if req.ProblemID != "" {
			probIDPtr = &req.ProblemID
		}
		title := req.Question
		if len(title) > 40 {
			title = title[:40] + "..."
		}
		conv = &AIConversation{
			StudentID:        studentID,
			ProblemID:        probIDPtr,
			ConversationType: "concept_tutoring",
			Title:            title,
		}
		if err := s.repo.CreateConversation(ctx, conv); err != nil {
			s.logger.WarnContext(ctx, "failed to persist conversation, using ephemeral session", "error", err)
			conv.ID = "ephemeral-session-" + time.Now().Format("20060102150405")
		}
	}

	// Fetch previous messages in conversation
	pastMsgs, _ := s.repo.ListMessagesByConversation(ctx, conv.ID)
	history := make([]ai.ChatMessage, 0, len(pastMsgs))
	for _, m := range pastMsgs {
		history = append(history, ai.ChatMessage{
			Role:    m.Role,
			Content: m.Message,
		})
	}

	problemContext := ""
	if req.ProblemID != "" {
		lc, _ := s.ctxBuilder.Build(ctx, studentID, req.ProblemID, "")
		if lc != nil && lc.Problem != nil {
			problemContext = fmt.Sprintf("Problem: %s\nSummary: %s", lc.Problem.Title, lc.Problem.Summary)
		}
	}

	// Pull learning memory context for targeted pedagogical grounding
	if s.memoryService != nil && studentID != "" {
		memContext, err := s.memoryService.BuildAIContext(ctx, studentID, "", req.ProblemID)
		if err == nil && memContext != "" {
			if problemContext != "" {
				problemContext += "\n\n" + memContext
			} else {
				problemContext = memContext
			}
		}
	}

	promptMessages := BuildAskMessages(history, req.Question, problemContext)

	llmResp, err := s.llm.Generate(ctx, ai.ChatCompletionRequest{
		Messages:    promptMessages,
		Temperature: 0.5,
		MaxTokens:   800,
	})
	if err != nil {
		return nil, err
	}

	// Record student message
	_ = s.repo.SaveMessage(ctx, &AIMessage{
		ConversationID: conv.ID,
		Role:           "user",
		Message:        req.Question,
		ResponseType:   "text",
	})

	// Record assistant reply
	_ = s.repo.SaveMessage(ctx, &AIMessage{
		ConversationID: conv.ID,
		Role:           "assistant",
		Message:        llmResp.Content,
		ResponseType:   "text",
		TokenUsage: map[string]interface{}{
			"prompt_tokens":     llmResp.Usage.PromptTokens,
			"completion_tokens": llmResp.Usage.CompletionTokens,
			"total_tokens":      llmResp.Usage.TotalTokens,
		},
	})

	return &AskResponse{
		ConversationID: conv.ID,
		Answer:         llmResp.Content,
		CreatedAt:      time.Now(),
	}, nil
}

func (s *serviceImpl) ListConversations(ctx context.Context, studentID string, limit, offset int) ([]AIConversation, error) {
	return s.repo.ListConversationsByStudent(ctx, studentID, limit, offset)
}

func (s *serviceImpl) GetConversationMessages(ctx context.Context, studentID, conversationID string) ([]AIMessage, error) {
	conv, err := s.repo.GetConversation(ctx, conversationID)
	if err != nil {
		return nil, err
	}
	if conv.StudentID != studentID {
		return nil, ErrUnauthorized
	}
	return s.repo.ListMessagesByConversation(ctx, conversationID)
}
