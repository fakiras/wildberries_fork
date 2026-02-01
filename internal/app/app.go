package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"wildberries/internal/config"
	"wildberries/internal/handler"
	"wildberries/internal/repository"
	"wildberries/internal/service"
)

type App struct {
	cfg    *config.Config
	pool   *pgxpool.Pool
	router *chi.Mux
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	pool, err := pgxpool.New(ctx, cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.New: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("db ping: %w", err)
	}

	// Repositories
	promoRepo := repository.NewPromotionPostgres(pool)
	segRepo := repository.NewSegmentPostgres(pool)
	slotRepo := repository.NewSlotPostgres(pool)
	productRepo := repository.NewProductPostgres(pool)
	modRepo := repository.NewModerationPostgres(pool)
	pollRepo := repository.NewPollPostgres(pool)
	auctionRepo := repository.NewAuctionPostgres(pool)
	betRepo := repository.NewBetPostgres(pool)

	// Services
	promoSvc := service.NewPromotionService(promoRepo, segRepo)
	segSvc := service.NewSegmentService(segRepo)
	productSvc := service.NewProductService(productRepo)
	slotSvc := service.NewSlotService(slotRepo, modRepo, auctionRepo, betRepo)
	modSvc := service.NewModerationService(modRepo, slotRepo)
	identSvc := service.NewIdentificationService(promoRepo, pollRepo, segRepo)
	aiSvc := service.NewAIService()

	// Handlers
	buyerH := handler.NewBuyerHandler(promoSvc, segSvc, slotSvc, productSvc, identSvc)
	adminH := handler.NewAdminHandler(promoSvc, segSvc, slotSvc, modSvc, identSvc, aiSvc)
	sellerH := handler.NewSellerHandler(productSvc, promoSvc, slotSvc)
	aiH := handler.NewAIHandler(aiSvc, segSvc)

	r := chi.NewRouter()
	r.Route("/promotions", func(r chi.Router) {
		r.Get("/current", buyerH.GetCurrentPromotion)
		r.Get("/{promotionId}/segments/{segmentId}/products", buyerH.GetSegmentProducts)
	})
	r.Route("/identification", func(r chi.Router) {
		r.Post("/start", buyerH.StartIdentification)
		r.Post("/answer", buyerH.Answer)
	})

	r.Route("/admin", func(r chi.Router) {
		r.Post("/promotions", adminH.CreatePromotion)
		r.Get("/promotions/{id}", adminH.GetPromotion)
		r.Patch("/promotions/{id}", adminH.UpdatePromotion)
		r.Delete("/promotions/{id}", adminH.DeletePromotion)
		r.Put("/promotions/{id}/fixed-prices", adminH.SetFixedPrices)
		r.Put("/promotions/{id}/status", adminH.ChangeStatus)
		r.Get("/promotions/{id}/moderation/applications", adminH.GetModerationApplications)
		r.Post("/promotions/{id}/segments/generate", adminH.GenerateSegments)
		r.Post("/promotions/{id}/segments", adminH.CreateSegment)
		r.Patch("/promotions/{id}/segments/{segmentId}", adminH.UpdateSegment)
		r.Delete("/promotions/{id}/segments/{segmentId}", adminH.DeleteSegment)
		r.Post("/promotions/{id}/segments/shuffle-categories", adminH.ShuffleSegmentCategories)
		r.Post("/promotions/{id}/poll/generate", adminH.GeneratePoll)
		r.Post("/promotions/{id}/poll/questions", adminH.SetPollQuestions)
		r.Post("/promotions/{id}/poll/answer-tree", adminH.SetAnswerTree)
		r.Post("/moderation/{applicationId}/approve", adminH.ApproveModeration)
		r.Post("/moderation/{applicationId}/reject", adminH.RejectModeration)
	})
	r.Post("/horoscope/products", adminH.SetSlotProduct)

	r.Route("/seller", func(r chi.Router) {
		r.Get("/actions", sellerH.GetSellerActions)
		r.Get("/bets/list", sellerH.GetSellerBetsList)
		r.Post("/bets/make", sellerH.MakeBet)
		r.Post("/bets/remove", sellerH.RemoveBet)
	})
	r.Get("/products/list-by", sellerH.ListProductsBy)

	r.Route("/ai", func(r chi.Router) {
		r.Post("/themes", aiH.GenerateThemes)
		r.Post("/segments", aiH.GenerateSegments)
		r.Post("/questions", aiH.GenerateQuestions)
		r.Post("/answer-tree", aiH.GenerateAnswerTree)
		r.Post("/get-text", aiH.GetText)
	})

	return &App{
		cfg:    cfg,
		pool:   pool,
		router: r,
	}, nil
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

func (a *App) Shutdown(ctx context.Context) {
	a.pool.Close()
}
