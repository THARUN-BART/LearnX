package learningpath

import "context"

// Repository defines data access operations for student adaptive learning paths.
type Repository interface {
	// Path & Nodes
	GetAdaptivePath(ctx context.Context, studentID string) (*AdaptivePath, error)
	UpsertAdaptivePath(ctx context.Context, path *AdaptivePath) error
	GetPathNodes(ctx context.Context, studentID string) ([]PathNode, error)
	UpsertPathNode(ctx context.Context, node *PathNode) error
	BatchUpsertPathNodes(ctx context.Context, studentID string, nodes []PathNode) error

	// Recommendations
	GetActiveRecommendation(ctx context.Context, studentID string) (*LearningRecommendation, error)
	CreateRecommendation(ctx context.Context, rec *LearningRecommendation) error
	CompleteRecommendation(ctx context.Context, studentID, recID string) error
	ListRecommendationHistory(ctx context.Context, studentID string, limit, offset int) ([]LearningRecommendation, int, error)
}
