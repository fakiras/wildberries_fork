package service

import (
	"context"
	"math/rand"

	"wildberries/internal/repository"
)

type IdentificationService struct {
	promoRepo repository.PromotionRepository
	pollRepo  repository.PollRepository
	segRepo   repository.SegmentRepository
}

func NewIdentificationService(
	promoRepo repository.PromotionRepository,
	pollRepo repository.PollRepository,
	segRepo repository.SegmentRepository,
) *IdentificationService {
	return &IdentificationService{
		promoRepo: promoRepo,
		pollRepo:  pollRepo,
		segRepo:   segRepo,
	}
}

// StartIdentification возвращает метод (questions | user_profile) и при questions — тест.
// Для user_profile — заглушка: возвращаем случайный segment_id.
func (s *IdentificationService) StartIdentification(ctx context.Context, promotionID int64) (
	method string,
	questions []*repository.PollQuestionRow,
	optionsByQuestion map[int64][]*repository.PollOptionRow,
	resultSegmentID int64,
	err error,
) {
	promo, err := s.promoRepo.GetByID(ctx, promotionID)
	if err != nil || promo == nil {
		return "", nil, nil, 0, err
	}
	method = promo.IdentificationMode
	if method == "user_profile" {
		segments, _ := s.segRepo.ByPromotionID(ctx, promotionID)
		if len(segments) > 0 {
			resultSegmentID = segments[rand.Intn(len(segments))].ID
		}
		return method, nil, nil, resultSegmentID, nil
	}
	questions, err = s.pollRepo.QuestionsByPromotion(ctx, promotionID)
	if err != nil || len(questions) == 0 {
		return method, nil, nil, 0, err
	}
	ids := make([]int64, len(questions))
	for i := range questions {
		ids[i] = questions[i].ID
	}
	opts, err := s.pollRepo.OptionsByQuestionIDs(ctx, ids)
	if err != nil {
		return method, questions, nil, 0, err
	}
	optionsByQuestion = make(map[int64][]*repository.PollOptionRow)
	for _, o := range opts {
		optionsByQuestion[o.QuestionID] = append(optionsByQuestion[o.QuestionID], o)
	}
	return method, questions, optionsByQuestion, 0, nil
}

// ProcessAnswer по (promotionID, questionID, optionID) возвращает nextQuestionID и resultSegmentID (из дерева ответов).
func (s *IdentificationService) ProcessAnswer(ctx context.Context, promotionID, questionID, optionID int64) (
	nextQuestionID int64,
	resultSegmentID int64,
	err error,
) {
	// TODO: загрузить дерево poll_answer_tree, найти узел по value опции, вернуть следующую ноду или segment
	_, _, _ = promotionID, questionID, optionID
	return 0, 0, nil
}
