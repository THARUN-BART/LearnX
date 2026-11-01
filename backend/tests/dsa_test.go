package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/THARUN-BART/LearnX/backend/internal/analytics"
	"github.com/THARUN-BART/LearnX/backend/internal/auth"
	"github.com/THARUN-BART/LearnX/backend/internal/config"
	"github.com/THARUN-BART/LearnX/backend/internal/crypto"
	"github.com/THARUN-BART/LearnX/backend/internal/dsa"
	"github.com/THARUN-BART/LearnX/backend/internal/dsa/model"
	"github.com/THARUN-BART/LearnX/backend/internal/importer"
	"github.com/THARUN-BART/LearnX/backend/internal/middleware"
	"github.com/THARUN-BART/LearnX/backend/internal/otp"
	"github.com/THARUN-BART/LearnX/backend/internal/problems"
	"github.com/THARUN-BART/LearnX/backend/internal/profile"
	"github.com/THARUN-BART/LearnX/backend/internal/progress"
	"github.com/THARUN-BART/LearnX/backend/internal/roadmap"
	"github.com/THARUN-BART/LearnX/backend/internal/server"
	"github.com/THARUN-BART/LearnX/backend/internal/topics"
)

func sampleSeedPath(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve test file path")
	}
	return filepath.Join(filepath.Dir(file), "..", "data", "strivers-a2z-sheet.sample.json")
}

func loadSampleCatalog(t *testing.T) *model.Memory {
	t.Helper()
	seed, err := importer.LoadSeedFile(sampleSeedPath(t))
	if err != nil {
		t.Fatalf("load seed: %v", err)
	}
	mem := model.NewMemory()
	if _, err := importer.ImportMemory(mem, seed); err != nil {
		t.Fatalf("import seed: %v", err)
	}
	return mem
}

func TestRoadmapAndTopicRetrieval(t *testing.T) {
	mem := loadSampleCatalog(t)
	roadmapSvc := roadmap.NewService(roadmap.NewMemoryRepository(mem), topics.NewMemoryRepository(mem))
	topicsSvc := topics.NewService(topics.NewMemoryRepository(mem))

	roadmaps, err := roadmapSvc.List(context.Background())
	if err != nil {
		t.Fatalf("list roadmaps: %v", err)
	}
	if len(roadmaps) != 1 {
		t.Fatalf("expected 1 roadmap, got %d", len(roadmaps))
	}
	if roadmaps[0].Slug != "strivers-a2z-dsa-sheet" {
		t.Fatalf("unexpected roadmap slug %s", roadmaps[0].Slug)
	}

	item, err := roadmapSvc.Get(context.Background(), "strivers-a2z-dsa-sheet")
	if err != nil {
		t.Fatalf("get roadmap: %v", err)
	}
	topicsList, err := roadmapSvc.ListTopics(context.Background(), item.ID)
	if err != nil {
		t.Fatalf("list topics: %v", err)
	}
	if len(topicsList) == 0 {
		t.Fatal("expected official A2Z topics")
	}
	if topicsList[0].TopicOrder != 1 || topicsList[0].Name != "Beginner Problems" {
		t.Fatalf("unexpected first topic: %+v", topicsList[0])
	}

	topic, err := topicsSvc.Get(context.Background(), topicsList[0].ID)
	if err != nil {
		t.Fatalf("get topic: %v", err)
	}
	subs, err := topicsSvc.ListSubtopics(context.Background(), topic.ID)
	if err != nil {
		t.Fatalf("list subtopics: %v", err)
	}
	if len(subs) == 0 || subs[0].SubtopicOrder != 1 {
		t.Fatalf("expected ordered subtopics, got %+v", subs)
	}
}

func TestProblemFilteringSearchAndOrder(t *testing.T) {
	mem := loadSampleCatalog(t)
	svc := problems.NewService(problems.NewMemoryRepository(mem), topics.NewMemoryRepository(mem))

	all, total, err := svc.List(context.Background(), model.ProblemFilter{Page: 1, Limit: 20, Sort: "problem_order"})
	if err != nil {
		t.Fatalf("list problems: %v", err)
	}
	if total == 0 || len(all) == 0 {
		t.Fatal("expected imported problems")
	}
	if len(all) > 1 && all[0].ProblemOrder > all[1].ProblemOrder {
		t.Fatalf("problems are not sorted by problem_order")
	}

	easy, easyTotal, err := svc.List(context.Background(), model.ProblemFilter{Difficulty: "Easy", Page: 1, Limit: 20})
	if err != nil {
		t.Fatalf("filter difficulty: %v", err)
	}
	if easyTotal == 0 {
		t.Fatal("expected Easy problems in the official sample")
	}
	for _, p := range easy {
		if p.Difficulty != "Easy" {
			t.Fatalf("difficulty filter leaked %s", p.Difficulty)
		}
	}

	topicID := all[0].TopicID
	byTopic, topicTotal, err := svc.List(context.Background(), model.ProblemFilter{TopicID: topicID, Page: 1, Limit: 50})
	if err != nil {
		t.Fatalf("filter topic: %v", err)
	}
	if topicTotal == 0 {
		t.Fatal("expected topic filter matches")
	}
	for _, p := range byTopic {
		if p.TopicID != topicID {
			t.Fatalf("topic filter leaked %s", p.TopicID)
		}
	}

	subID := all[0].SubtopicID
	bySub, subTotal, err := svc.List(context.Background(), model.ProblemFilter{SubtopicID: subID, Page: 1, Limit: 50})
	if err != nil {
		t.Fatalf("filter subtopic: %v", err)
	}
	if subTotal == 0 {
		t.Fatal("expected subtopic filter matches")
	}
	for _, p := range bySub {
		if p.SubtopicID != subID {
			t.Fatalf("subtopic filter leaked %s", p.SubtopicID)
		}
	}

	found, foundTotal, err := svc.List(context.Background(), model.ProblemFilter{Query: "two-sum", Page: 1, Limit: 20})
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if foundTotal != 0 && len(found) == 0 {
		t.Fatal("search total/result mismatch")
	}

	pattern, patternTotal, err := svc.List(context.Background(), model.ProblemFilter{Query: "pattern", Page: 1, Limit: 20})
	if err != nil {
		t.Fatalf("search pattern: %v", err)
	}
	if patternTotal == 0 {
		t.Fatal("expected title search for pattern")
	}
	_ = pattern

	page1, _, err := svc.List(context.Background(), model.ProblemFilter{Page: 1, Limit: 2, Sort: "problem_order"})
	if err != nil {
		t.Fatal(err)
	}
	if len(page1) != 2 {
		t.Fatalf("expected pagination limit 2, got %d", len(page1))
	}
}

func TestStudentProgressIsolationAndReset(t *testing.T) {
	mem := loadSampleCatalog(t)
	problemRepo := problems.NewMemoryRepository(mem)
	progressRepo := progress.NewMemoryRepository(mem)
	svc := progress.NewService(progressRepo, problemRepo)

	items, _, err := problemRepo.List(context.Background(), model.ProblemFilter{Page: 1, Limit: 5})
	if err != nil || len(items) == 0 {
		t.Fatalf("need problems: %v", err)
	}
	problemID := items[0].ID
	alice := "11111111-1111-4111-8111-111111111111"
	bob := "22222222-2222-4222-8222-222222222222"

	if _, err := svc.MarkAttempted(context.Background(), alice, problemID); err != nil {
		t.Fatalf("alice attempt: %v", err)
	}
	solved, err := svc.MarkSolved(context.Background(), alice, problemID)
	if err != nil {
		t.Fatalf("alice solve: %v", err)
	}
	if solved.Status != model.StatusSolved || solved.SuccessfulAttempts < 1 {
		t.Fatalf("unexpected alice progress %+v", solved)
	}

	bobProgress, err := svc.Get(context.Background(), bob, problemID)
	if err != nil {
		t.Fatalf("bob get: %v", err)
	}
	if bobProgress.Status != model.StatusNotStarted {
		t.Fatalf("bob should be isolated, got %+v", bobProgress)
	}

	aliceList, err := svc.List(context.Background(), alice)
	if err != nil {
		t.Fatal(err)
	}
	bobList, err := svc.List(context.Background(), bob)
	if err != nil {
		t.Fatal(err)
	}
	if len(aliceList) == 0 {
		t.Fatal("alice should have progress rows")
	}
	if len(bobList) != 0 {
		t.Fatalf("bob should have no progress rows, got %d", len(bobList))
	}

	reset, err := svc.Update(context.Background(), alice, problemID, progress.UpdateRequest{Reset: true})
	if err != nil {
		t.Fatalf("reset: %v", err)
	}
	if reset.Status != model.StatusNotStarted || reset.TotalAttempts != 0 {
		t.Fatalf("reset failed: %+v", reset)
	}
}

func TestAnalyticsCalculationsFromRecords(t *testing.T) {
	mem := loadSampleCatalog(t)
	repos := dsa.NewMemoryRepositories(mem)
	svc := analytics.NewService(repos.Roadmap, repos.Topics, repos.Problems, repos.Progress)

	items, _, err := repos.Problems.List(context.Background(), model.ProblemFilter{Page: 1, Limit: 100})
	if err != nil || len(items) < 3 {
		t.Fatalf("need sample problems: %v len=%d", err, len(items))
	}
	studentID := "33333333-3333-4333-8333-333333333333"
	if _, err := repos.Progress.MarkSolved(context.Background(), studentID, items[0].ID); err != nil {
		t.Fatal(err)
	}
	if _, err := repos.Progress.MarkSolved(context.Background(), studentID, items[1].ID); err != nil {
		t.Fatal(err)
	}
	if _, err := repos.Progress.MarkAttempted(context.Background(), studentID, items[2].ID); err != nil {
		t.Fatal(err)
	}

	summary, err := svc.Summary(context.Background(), studentID, "")
	if err != nil {
		t.Fatalf("summary: %v", err)
	}
	if summary.Overview.Total != len(items) {
		t.Fatalf("expected total %d, got %d", len(items), summary.Overview.Total)
	}
	if summary.Overview.Solved != 2 {
		t.Fatalf("expected 2 solved, got %d", summary.Overview.Solved)
	}
	if summary.Overview.Attempted < 1 {
		t.Fatalf("expected attempted >= 1, got %d", summary.Overview.Attempted)
	}
	if summary.Overview.Pending != summary.Overview.Total-summary.Overview.Solved-summary.Overview.Attempted {
		t.Fatalf("pending mismatch: %+v", summary.Overview)
	}
	expected := model.CompletionPercent(2, summary.Overview.Total)
	if summary.Overview.Completion != expected {
		t.Fatalf("completion got %v want %v", summary.Overview.Completion, expected)
	}
	if len(summary.Topics) == 0 {
		t.Fatal("expected topic-wise progress")
	}
	if len(summary.Difficulty) != 4 {
		t.Fatalf("expected 4 difficulty buckets, got %d", len(summary.Difficulty))
	}

	activity, err := svc.Activity(context.Background(), studentID, 7)
	if err != nil {
		t.Fatal(err)
	}
	if activity.Days != 7 || len(activity.Activity) != 7 {
		t.Fatalf("unexpected activity window %+v", activity)
	}
}

func TestDuplicateImportDoesNotDeleteProgress(t *testing.T) {
	seed, err := importer.LoadSeedFile(sampleSeedPath(t))
	if err != nil {
		t.Fatal(err)
	}
	mem := model.NewMemory()
	first, err := importer.ImportMemory(mem, seed)
	if err != nil {
		t.Fatal(err)
	}
	if first.ProblemsUpserted == 0 {
		t.Fatal("expected problems on first import")
	}

	var problemID string
	for _, p := range mem.Problems {
		problemID = p.ID
		break
	}
	studentID := "44444444-4444-4444-8444-444444444444"
	if _, err := mem.MarkSolved(studentID, problemID); err != nil {
		t.Fatal(err)
	}

	second, err := importer.ImportMemory(mem, seed)
	if err != nil {
		t.Fatal(err)
	}
	if second.ProblemsUpserted != first.ProblemsUpserted {
		t.Fatalf("repeatable import changed count %d vs %d", second.ProblemsUpserted, first.ProblemsUpserted)
	}
	if second.DuplicatesFound == 0 {
		t.Fatal("expected duplicate detection on second import")
	}
	if len(mem.Problems) != first.ProblemsUpserted {
		t.Fatalf("problem catalogue duplicated: %d", len(mem.Problems))
	}

	rec, ok := mem.GetProgress(studentID, problemID)
	if !ok || rec.Status != model.StatusSolved {
		t.Fatalf("student progress was disturbed: ok=%v rec=%+v", ok, rec)
	}
}

func TestDSAHTTPRoutes(t *testing.T) {
	mem := loadSampleCatalog(t)
	cfg := &config.Config{
		OTPAppName:               "LearnX",
		PendingSignupExpiry:      10 * time.Minute,
		PasswordResetTokenExpiry: 15 * time.Minute,
		EncryptionKey:            crypto.DeriveKey("test-aes-key-for-dsa-router!!"),
		ClientURL:                "*",
	}
	authRepo := auth.NewMemoryRepository()
	profileRepo := profile.NewMemoryRepository()
	sbClient := newRouterMockSupabase()
	otpService := otp.NewService(&routerMockOTPClient{}, 1*time.Second, 5)
	authService := auth.NewService(authRepo, otpService, sbClient, cfg.EncryptionKey, cfg.PendingSignupExpiry, cfg.PasswordResetTokenExpiry)
	dsaHandlers := dsa.NewHandlers(dsa.NewMemoryRepositories(mem))
	router := server.BuildRouter(
		cfg,
		auth.NewHandler(authService),
		profile.NewHandler(profile.NewService(profileRepo)),
		middleware.NewAuthMiddleware(sbClient, ""),
		middleware.NewRateLimiter(100, time.Minute),
		server.WithDSA(dsaHandlers),
	)

	unauth := httptest.NewRequest(http.MethodGet, "/api/v1/roadmaps", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, unauth)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for unauthenticated roadmap access, got %d", rec.Code)
	}

	signupBody, _ := json.Marshal(auth.SignupRequest{
		FullName:        "DSA Student",
		Email:           "dsa@learnx.com",
		Password:        "Password123!",
		ConfirmPassword: "Password123!",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signup", bytes.NewBuffer(signupBody))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("signup failed: %s", rec.Body.String())
	}
	verifyBody, _ := json.Marshal(auth.VerifySignupOTPRequest{Email: "dsa@learnx.com", Code: "123456"})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/verify-signup-otp", bytes.NewBuffer(verifyBody))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("verify failed: %s", rec.Body.String())
	}
	var authResp auth.AuthResponse
	_ = json.NewDecoder(rec.Body).Decode(&authResp)
	token := authResp.Session.AccessToken

	authGet := func(path string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		return rec
	}

	rec = authGet("/api/v1/roadmaps")
	if rec.Code != http.StatusOK {
		t.Fatalf("roadmaps: %d %s", rec.Code, rec.Body.String())
	}
	var roadmapResp struct {
		Roadmaps []model.Roadmap `json:"roadmaps"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &roadmapResp)
	if len(roadmapResp.Roadmaps) != 1 {
		t.Fatalf("expected seeded roadmap, got %+v", roadmapResp)
	}
	roadmapID := roadmapResp.Roadmaps[0].ID

	rec = authGet("/api/v1/roadmaps/" + roadmapID + "/topics")
	if rec.Code != http.StatusOK {
		t.Fatalf("topics: %s", rec.Body.String())
	}

	rec = authGet("/api/v1/problems?difficulty=Easy&limit=5&sort=problem_order")
	if rec.Code != http.StatusOK {
		t.Fatalf("problems filter: %s", rec.Body.String())
	}
	rec = authGet("/api/v1/problems/search?q=pattern")
	if rec.Code != http.StatusOK {
		t.Fatalf("search: %s", rec.Body.String())
	}

	var problemResp struct {
		Problems []model.Problem `json:"problems"`
	}
	listRec := authGet("/api/v1/problems?limit=1")
	_ = json.Unmarshal(listRec.Body.Bytes(), &problemResp)
	if len(problemResp.Problems) == 0 {
		t.Fatal("expected a problem for progress routes")
	}
	problemID := problemResp.Problems[0].ID

	attemptReq := httptest.NewRequest(http.MethodPost, "/api/v1/student/progress/problems/"+problemID+"/attempt", nil)
	attemptReq.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, attemptReq)
	if rec.Code != http.StatusOK {
		t.Fatalf("attempt: %s", rec.Body.String())
	}

	patchBody, _ := json.Marshal(progress.UpdateRequest{Status: model.StatusSolved})
	patchReq := httptest.NewRequest(http.MethodPatch, "/api/v1/student/progress/problems/"+problemID, bytes.NewBuffer(patchBody))
	patchReq.Header.Set("Authorization", "Bearer "+token)
	patchReq.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, patchReq)
	if rec.Code != http.StatusOK {
		t.Fatalf("solve: %s", rec.Body.String())
	}

	for _, path := range []string{
		"/api/v1/student/progress",
		"/api/v1/student/progress/summary",
		"/api/v1/student/analytics/overview",
		"/api/v1/student/analytics/topics",
		"/api/v1/student/analytics/difficulty",
		"/api/v1/student/analytics/activity",
	} {
		rec = authGet(path)
		if rec.Code != http.StatusOK {
			t.Fatalf("%s failed: %d %s", path, rec.Code, rec.Body.String())
		}
	}
}
