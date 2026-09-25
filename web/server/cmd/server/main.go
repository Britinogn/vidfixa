package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/HalxDocs/bachs-go"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"gitlab.com/britinogn/vidfixa/config"
	"gitlab.com/britinogn/vidfixa/internal/db"
	"gitlab.com/britinogn/vidfixa/internal/downloader"
	"gitlab.com/britinogn/vidfixa/internal/handler"
	"gitlab.com/britinogn/vidfixa/internal/middleware"
	"gitlab.com/britinogn/vidfixa/internal/model"
	"gitlab.com/britinogn/vidfixa/internal/repository"
	"gitlab.com/britinogn/vidfixa/internal/routes"
	"gitlab.com/britinogn/vidfixa/internal/service"
	"gitlab.com/britinogn/vidfixa/internal/worker"
)

type handlers struct {
	auth         *handler.AuthHandler
	download     *handler.DownloadHandler
	usage        *handler.UsageHandler
	subscription *handler.SubscriptionHandler
	payment      *handler.PaymentHandler
	dashboard    *handler.DashboardHandler
	admin        *handler.AdminHandler
}

/*
initHandlers wires services, the worker pool, and handlers together.
The pool's ProcessFunc is built right here — this is the one place
that knows about both downloader (yt-dlp) and repository
(Postgres), keeping the worker package itself fully generic per the
CPC design discussion.
*/
func initHandlers(cfg *config.Config) (*handlers, *worker.Pool) {
	// --- Auth ---
	authService := service.NewAuthService()
	authHandler := handler.NewAuthHandler(authService)

	// --- Bachs client ---
	// bachs.WithProduction()
	bachsClient, err := bachs.NewClient(cfg.BachsAPIKey)
	if err != nil {
		log.Fatal("failed to create Bachs client:", err)
	}

	// --- Subscription ---
	subscriptionService := service.NewSubscriptionService(
		bachsClient,
		cfg.BachsPlusProductID,
		cfg.BachsProProductID,
		cfg.AppURL+"/subscription/success",
		cfg.AppURL+"/subscription/cancel",
	)

	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService)

	// --- User and admin dashboards ---
	usageService := service.NewUsageService()
	dashboardHandler := handler.NewDashboardHandler(usageService, subscriptionService)
	adminService := service.NewAdminService()
	adminHandler := handler.NewAdminHandler(adminService)

	// --- Worker pool ---
	processJob := func(job worker.Job) error {
		errMsg := ""

		if err := repository.UpdateDownloadStatus(job.Ctx, job.ID, model.StatusProcessing, nil, nil); err != nil {
			log.Printf("job %s: failed to mark processing: %v", job.ID, err)
		}

		filePath, err := downloader.Download(job.Ctx, cfg.YtdlpPath, job.URL, cfg.DownloadDir)
		if err != nil {
			errMsg = err.Error()
			_ = repository.UpdateDownloadStatus(context.Background(), job.ID, model.StatusFailed, nil, &errMsg)
			return err
		}

		if updateErr := repository.UpdateDownloadStatus(context.Background(), job.ID, model.StatusCompleted, &filePath, nil); updateErr != nil {
			log.Printf("job %s: failed to mark completed: %v", job.ID, updateErr)
		}

		return nil
	}

	pool := worker.NewPool(100, processJob) // 100 = buffered channel capacity

	// pool.Start(cfg.DownloadWorkers)
	workers, err := strconv.Atoi(cfg.DownloadWorkers)
	if err != nil || workers < 1 {
		workers = 4 // fallback default if unset/invalid
	}
	pool.Start(workers)

	// --- Download ---
	downloadService := service.NewDownloadService(pool, subscriptionService)
	downloadHandler := handler.NewDownloadHandler(downloadService)

	// --- Usage ---
	usageHandler := handler.NewUsageHandler(usageService, subscriptionService)

	// --- Payment / Webhook ---
	webhookService := service.NewWebhookService(
		bachsClient,
		cfg.BachsWebhookSecret,
		cfg.BachsPlusProductID,
		cfg.BachsProProductID,
	)

	paymentHandler := handler.NewPaymentHandler(webhookService)

	return &handlers{
		auth:         authHandler,
		download:     downloadHandler,
		usage:        usageHandler,
		subscription: subscriptionHandler,
		payment:      paymentHandler,
		dashboard:    dashboardHandler,
		admin:        adminHandler,
	}, pool
}

func main() {
	db.Init() 
	cfg := config.Load()

	_, err := db.ConnectPostgres(context.Background(), cfg)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()
	log.Println("✓ Database connected successfully")

	h, pool := initHandlers(cfg)

	router := chi.NewRouter()
	router.Use(chimw.RequestID)
	router.Use(chimw.Recoverer)
	router.Use(chimw.Logger)
	router.Use(middleware.CORS(cfg))

	// routes.SetupRoutes(router, cfg, h.auth, h.download, h.subscription, h.payment)
	routes.SetupRoutes(
		router,
		cfg,
		h.auth,
		h.download,
		h.usage,
		h.subscription,
		h.payment,
		h.dashboard,
		h.admin,
	)
	srv := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: router,
	}

	go func() {
		log.Println("Server running on port", cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Graceful shutdown: stop accepting new jobs, let in-flight
	// downloads finish (up to 30s) before the process exits — this
	// is what makes the worker pool's Shutdown() actually matter,
	// not just a method that never gets called.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Println("forced HTTP shutdown:", err)
	}
	pool.Shutdown()
	log.Println("shutdown complete")
}
