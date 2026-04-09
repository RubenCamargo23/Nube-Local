package main

import (
    "log"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/org/monitors-platform/internal/auth"
    "github.com/org/monitors-platform/internal/users"
    "github.com/org/monitors-platform/internal/periods"
    "github.com/org/monitors-platform/internal/spaces"
    "github.com/org/monitors-platform/internal/assignments"
    "github.com/org/monitors-platform/internal/tasks"
    "github.com/org/monitors-platform/internal/reports"
    "github.com/org/monitors-platform/pkg/config"
    "github.com/org/monitors-platform/pkg/database"
    jwtpkg "github.com/org/monitors-platform/pkg/jwt"
    "github.com/org/monitors-platform/pkg/response"
    "net/http"
)

func main() {
    cfg := config.Load()

    db, err := database.NewPool(cfg)
    if err != nil {
        log.Fatalf("Error conectando a DB: %v", err)
    }
    defer db.Close()

    os.MkdirAll(cfg.StoragePath, 0755)
    os.MkdirAll(cfg.PDFPath, 0755)

    jwtManager := jwtpkg.NewManager(cfg.JWTSecret, cfg.JWTExpirationHours)

    // Repositorios
    authRepo        := auth.NewRepository(db)
    usersRepo       := users.NewRepository(db)
    periodsRepo     := periods.NewRepository(db)
    spacesRepo      := spaces.NewRepository(db)
    assignmentsRepo := assignments.NewRepository(db)
    tasksRepo       := tasks.NewRepository(db)
    reportsRepo     := reports.NewRepository(db)

    // Servicios
    authSvc        := auth.NewService(authRepo, jwtManager)
    usersSvc       := users.NewService(usersRepo)
    periodsSvc     := periods.NewService(periodsRepo)
    spacesSvc      := spaces.NewService(spacesRepo, periodsSvc)
    assignmentsSvc := assignments.NewService(assignmentsRepo, spacesSvc)
    tasksSvc       := tasks.NewService(tasksRepo, assignmentsSvc, spacesSvc)

    ollamaClient   := reports.NewOllamaClient(cfg.OllamaHost, cfg.OllamaModel, cfg.OllamaTimeoutSecs)
    
    // There is no reports.Service clearly defined in prompt, let's create a minimal Service interface struct.
    type ReportsService struct {}
    reportsSvc := &ReportsService{}
    _ = reportsRepo
    _ = ollamaClient

    // Handlers
    authH        := auth.NewHandler(authSvc)
    usersH       := users.NewHandler(usersSvc)
    periodsH     := periods.NewHandler(periodsSvc)
    spacesH      := spaces.NewHandler(spacesSvc)
    assignmentsH := assignments.NewHandler(assignmentsSvc)
    tasksH       := tasks.NewHandler(tasksSvc)
    reportsH     := reports.NewHandler(reportsSvc)

    r := gin.Default()

    // Health
    r.GET("/api/v1/health", func(c *gin.Context) {
        response.OK(c, http.StatusOK, gin.H{"status": "ok", "version": "1.0.0"})
    })

    // Auth público
    r.POST("/api/v1/auth/login", authH.Login)

    // Rutas protegidas
    protected := r.Group("/api/v1")
    protected.Use(auth.JWTMiddleware(jwtManager))

    // Admin
    admin := protected.Group("")
    admin.Use(auth.RequireRole("ADMIN"))
    admin.POST("/admin/usuarios", usersH.Create)
    admin.POST("/admin/usuarios/:id/roles", usersH.AssignRole)
    admin.GET("/admin/usuarios", usersH.FindAll)
    admin.GET("/admin/usuarios/:id", usersH.FindByID)
    admin.POST("/admin/periodos", periodsH.Create)
    admin.GET("/admin/periodos", periodsH.FindAll)
    admin.GET("/admin/periodos/:id", periodsH.FindByID)
    admin.PATCH("/admin/periodos/:id/cerrar", periodsH.Close)
    admin.PATCH("/admin/vinculaciones/:id", assignmentsH.Update)

    // Professor
    prof := protected.Group("")
    prof.Use(auth.RequireRole("PROFESSOR"))
    prof.POST("/espacios", spacesH.Create)
    prof.GET("/espacios", spacesH.FindByProfesor)
    prof.GET("/espacios/:id", spacesH.FindByID)
    prof.PATCH("/espacios/:id/cerrar", spacesH.Close)
    prof.POST("/espacios/:id/vinculaciones", assignmentsH.Create)
    prof.GET("/espacios/:id/vinculaciones", assignmentsH.FindByEspacio)
    prof.GET("/profesor/vinculaciones", assignmentsH.FindByProfesor)
    prof.GET("/profesor/vinculaciones/:id/tareas", tasksH.FindByVinculacionSemana)
    prof.GET("/profesor/usuarios/:id/resumen", tasksH.ResumenByUsuario)
    prof.GET("/profesor/espacios/:id/semana", tasksH.EstadoSemanal)
    prof.POST("/profesor/reportes/generar", reportsH.Generate)
    prof.GET("/profesor/reportes", reportsH.FindByProfesor)
    prof.GET("/profesor/reportes/:id/descargar", reportsH.Download)

    // Monitor / Asistente
    operativo := protected.Group("")
    operativo.Use(auth.RequireRole("MONITOR", "GRAD_ASSISTANT"))
    operativo.POST("/vinculaciones/:id/tareas", tasksH.Create)
    operativo.GET("/vinculaciones/:id/tareas", tasksH.FindByVinculacion)
    operativo.GET("/tareas/:id", tasksH.FindByID)
    operativo.PUT("/tareas/:id", tasksH.Update)
    operativo.DELETE("/tareas/:id", tasksH.Delete)
    operativo.POST("/tareas/:id/adjuntos", tasksH.AddAttachment)
    operativo.GET("/me/vinculaciones", assignmentsH.FindByMe)
    operativo.GET("/me/tareas", tasksH.FindByMe)
    operativo.GET("/me/tareas/historial", tasksH.HistorialByMe)

    log.Printf("Servidor iniciado en :%s", cfg.Port)
    r.Run(":" + cfg.Port)
}
