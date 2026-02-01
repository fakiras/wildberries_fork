package service

import (
	"context"
)

type AIService struct {
	// Внешний AI-клиент (OpenAI/другой) — заглушка
}

func NewAIService() *AIService {
	return &AIService{}
}

// GenerateThemes возвращает список тем (value, label).
func (s *AIService) GenerateThemes(ctx context.Context) ([]struct{ Value, Label string }, error) {
	_ = ctx
	// TODO: вызов внешнего API
	return nil, nil
}

// GenerateSegments по теме возвращает сегменты (name, category_name).
func (s *AIService) GenerateSegments(ctx context.Context, theme string, limit int) ([]struct{ Name, CategoryName string }, error) {
	_, _, _ = ctx, theme, limit
	return nil, nil
}

// GenerateQuestions по теме возвращает вопросы с опциями.
func (s *AIService) GenerateQuestions(ctx context.Context, theme string) ([]struct {
	Text    string
	Options []struct{ Text, Value string }
}, error) {
	_, _ = ctx, theme
	return nil, nil
}

// GenerateAnswerTree по теме возвращает узлы дерева (node_id, parent_node_id, label, value).
func (s *AIService) GenerateAnswerTree(ctx context.Context, theme string) ([]struct {
	NodeID, ParentNodeID, Label, Value string
}, error) {
	_, _ = ctx, theme
	return nil, nil
}

// GenerateText генерирует текст; при указании segmentID обновляет segment.text.
func (s *AIService) GenerateText(ctx context.Context, params map[string]string, segmentID *int64) (text string, err error) {
	_, _, _ = params, segmentID, text
	return "", nil
}
