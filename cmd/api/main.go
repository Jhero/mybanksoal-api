package main

import (
	"context"
	"log"
	"strings"

	"github.com/casbin/casbin/v3" // Use v3
	"github.com/casbin/casbin/v3/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	oapiMiddleware "github.com/oapi-codegen/echo-middleware"
	"github.com/jovan/mybanksoal-api/config"
	// "github.com/jovan/mybanksoal-api/internal/entity"
	"github.com/jovan/mybanksoal-api/internal/handler"
	"github.com/jovan/mybanksoal-api/internal/middleware"
	"github.com/jovan/mybanksoal-api/internal/repository"
	"github.com/jovan/mybanksoal-api/internal/usecase"
	"github.com/jovan/mybanksoal-api/pkg/database"
	"github.com/jovan/mybanksoal-api/pkg/utils"
)

func main() {
	// 1. Load Config
	cfg := config.LoadConfig()

	// 2. Setup Database
	db := database.NewDatabase(cfg)

	// 3. Auto Migrate - Disabled in favor of versioned migrations
	// err := db.AutoMigrate(&entity.User{}, &entity.Question{})
	// if err != nil {
	// 	log.Fatalf("Failed to migrate database: %v", err)
	// }
	log.Println("Note: AutoMigrate is disabled. Use 'go run cmd/migrate/main.go up' to migrate database.")

	// 4. Setup Casbin RBAC
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		log.Fatalf("Failed to initialize casbin adapter: %v", err)
	}

	// Use in-memory model to avoid file reading issues
	m, err := model.NewModelFromString(`
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && keyMatch2(r.obj, p.obj) && r.act == p.act
`)
	if err != nil {
		log.Fatalf("Failed to create casbin model: %v", err)
	}

	enforcer, err := casbin.NewEnforcer()
	if err != nil {
		log.Fatalf("Failed to create enforcer: %v", err)
	}
	
	enforcer.SetModel(m)
	enforcer.SetAdapter(adapter)
	
	if err := enforcer.LoadPolicy(); err != nil {
		log.Fatalf("Failed to load policy: %v", err)
	}

	// Seed RBAC policies if empty
	if hasPolicy, _ := enforcer.HasPolicy("admin", "/*", "*"); !hasPolicy {
		enforcer.AddPolicy("admin", "/*", "*")            // Admin can access everything
		enforcer.AddPolicy("user", "/questions", "GET")   // User can view questions
		enforcer.AddPolicy("user", "/questions/*", "GET") // User can view single question
		enforcer.AddPolicy("user", "/crosswords/levels", "GET")
		enforcer.AddPolicy("user", "/crosswords/levels/*", "GET")
		enforcer.AddPolicy("user", "/crosswords/levels/*/submit", "POST")
		enforcer.AddPolicy("user", "/crosswords/leaderboard", "GET")
		enforcer.AddPolicy("user", "/crosswords/questions", "GET")
		enforcer.AddPolicy("user", "/crosswords/questions/*", "GET")
		enforcer.AddPolicy("editor", "/questions", "POST")
		enforcer.AddPolicy("editor", "/questions/*", "PUT")
		enforcer.SavePolicy()
	}

	// Seed mobile_reader policies
	mobilePolicies := [][]interface{}{
		{"mobile_reader", "/questions", "GET"},
		{"mobile_reader", "/questions/*", "GET"},
		{"mobile_reader", "/crosswords/levels", "GET"},
		{"mobile_reader", "/crosswords/levels/*", "GET"},
		{"mobile_reader", "/crosswords/levels/*/submit", "POST"},
		{"mobile_reader", "/crosswords/leaderboard", "GET"},
	}

	for _, policy := range mobilePolicies {
		if added, _ := enforcer.AddPolicy(policy...); added {
			log.Printf("Added policy: %v", policy)
		}
	}

	if err := enforcer.SavePolicy(); err != nil {
		log.Printf("Failed to save policy: %v", err)
	}

	// 5. Setup Layers
	userRepo := repository.NewUserRepository(db)
	questionRepo := repository.NewQuestionRepository(db)
	levelRepo := repository.NewLevelRepository(db)
	crosswordQuestionRepo := repository.NewCrosswordQuestionRepository(db)

	userUseCase := usecase.NewUserUseCase(userRepo, cfg)
	questionUseCase := usecase.NewQuestionUseCase(questionRepo)
	levelUseCase := usecase.NewLevelUseCase(levelRepo)
	crosswordQuestionUseCase := usecase.NewCrosswordQuestionUseCase(crosswordQuestionRepo)

	userHandler := handler.NewUserHandler(userUseCase)
	questionHandler := handler.NewQuestionHandler(questionUseCase)
	levelHandler := handler.NewLevelHandler(levelUseCase)
	crosswordQuestionHandler := handler.NewCrosswordQuestionHandler(crosswordQuestionUseCase)

	// 6. Setup Echo
	e := echo.New()
	e.Validator = &utils.CustomValidator{Validator: validator.New()}
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORS())

	// Swagger Spec
	if cfg.App.Env != "production" {
		e.Static("/docs", "docs")
		e.File("/api.yml", "api.yml")
		log.Println("Swagger UI enabled at /docs")
	}

	// Load OpenAPI spec for validation
	swagger, err := openapi3.NewLoader().LoadFromFile("api.yml")
	if err != nil {
		log.Printf("Warning: Failed to load api.yml for validation: %v", err)
	} else {
		// Enable OpenAPI Validation Middleware
		// Skip validation for docs and swagger-ui
		options := &oapiMiddleware.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: func(c context.Context, input *openapi3filter.AuthenticationInput) error {
					return nil // Skip auth validation here, handled by custom middleware
				},
			},
			Skipper: func(c echo.Context) bool {
				path := c.Request().URL.Path
				// In production, docs are disabled, so we don't need to skip validation for them,
				// but skipping doesn't hurt.
				return strings.HasPrefix(path, "/docs") || path == "/api.yml"
			},
		}
		e.Use(oapiMiddleware.OapiRequestValidatorWithOptions(swagger, options))
		log.Println("OpenAPI Validation Middleware Enabled")
	}

	// Routes
	auth := e.Group("/auth")
	auth.POST("/register", userHandler.Register)
	auth.POST("/login", userHandler.Login)

	// Protected Routes
	api := e.Group("")
	api.Use(middleware.AuthMiddleware(cfg, userRepo))
	api.Use(middleware.CasbinMiddleware(enforcer))

	// Question Routes
	// Note: Casbin matches path/method.
	// We need to be careful with path params in Casbin matching.
	// For simplicity, we can use wildcards in policy or handle regex.
	// The current casbin middleware passes raw path.
	// e.g., /questions/1.
	// Policy: /questions/*
	// Matcher: keyMatch2(r.obj, p.obj) in model handles :id or *

	api.POST("/questions", questionHandler.Create)
	api.GET("/questions", questionHandler.GetAll)
	api.GET("/questions/:id", questionHandler.GetByID)
	api.PUT("/questions/:id", questionHandler.Update)
	api.PATCH("/questions/:id/status", questionHandler.UpdateStatus)
	api.DELETE("/questions/:id", questionHandler.Delete)

	// Crossword Routes
	api.POST("/crosswords/import", levelHandler.ImportLevels)
	api.GET("/crosswords/levels", levelHandler.GetLevels)
	api.GET("/crosswords/levels/:id", levelHandler.GetLevel)
	api.POST("/crosswords/levels", levelHandler.CreateLevel)
	api.PUT("/crosswords/levels/:id", levelHandler.UpdateLevel)
	api.DELETE("/crosswords/levels/:id", levelHandler.DeleteLevel)
	
	api.POST("/crosswords/levels/:id/submit", levelHandler.SubmitLevel)
	api.GET("/crosswords/leaderboard", levelHandler.GetLeaderboard)

	// Crossword Question Routes
	api.GET("/crosswords/questions", crosswordQuestionHandler.GetAll)
	api.POST("/crosswords/questions", crosswordQuestionHandler.Create)
	api.GET("/crosswords/questions/:id", crosswordQuestionHandler.GetByID)
	api.PUT("/crosswords/questions/:id", crosswordQuestionHandler.Update)
	api.DELETE("/crosswords/questions/:id", crosswordQuestionHandler.Delete)

	// Start Server
	e.Logger.Fatal(e.Start(":" + cfg.App.Port))
}
