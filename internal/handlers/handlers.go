package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"waterlogger/internal/config"
	"waterlogger/internal/middleware"
	"waterlogger/internal/models"
)

type Handlers struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewHandlers(db *gorm.DB, cfg *config.Config) *Handlers {
	return &Handlers{
		db:  db,
		cfg: cfg,
	}
}

// Setup Wizard
func (h *Handlers) SetupWizardPage(c *gin.Context) {
	c.HTML(http.StatusOK, "setup.html", gin.H{
		"title": "Setup Waterlogger",
	})
}

func (h *Handlers) SetupWizardAPI(c *gin.Context) {
	log.Printf("Setup wizard API called from %s", c.ClientIP())
	
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
		
		DatabaseType string `json:"database_type" binding:"required"`
		DBHost       string `json:"db_host"`
		DBPort       int    `json:"db_port"`
		DBUsername   string `json:"db_username"`
		DBPassword   string `json:"db_password"`
		DBName       string `json:"db_name"`
		
		ServerPort int `json:"server_port"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Setup wizard JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}
	
	log.Printf("Setup wizard request: username=%s, email=%s, db_type=%s", req.Username, req.Email, req.DatabaseType)

	// Validate password
	if errors := middleware.ValidatePassword(req.Password); len(errors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password validation failed", "details": errors})
		return
	}

	// Hash password
	hashedPassword, err := middleware.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create admin user
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := h.db.Create(&user).Error; err != nil {
		log.Printf("Failed to create user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": err.Error()})
		return
	}

	// Create user preferences
	preferences := models.UserPreferences{
		UserID:     user.ID,
		UnitSystem: "imperial",
	}

	if err := h.db.Create(&preferences).Error; err != nil {
		log.Printf("Failed to create user preferences: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user preferences", "details": err.Error()})
		return
	}

	// Update configuration
	if req.DatabaseType != "" {
		h.cfg.Database.Type = req.DatabaseType
		
		// Update MariaDB configuration if provided
		if req.DatabaseType == "mariadb" {
			if req.DBHost != "" {
				h.cfg.Database.MariaDB.Host = req.DBHost
			}
			if req.DBPort > 0 {
				h.cfg.Database.MariaDB.Port = req.DBPort
			}
			if req.DBUsername != "" {
				h.cfg.Database.MariaDB.Username = req.DBUsername
			}
			if req.DBPassword != "" {
				h.cfg.Database.MariaDB.Password = req.DBPassword
			}
			if req.DBName != "" {
				h.cfg.Database.MariaDB.Database = req.DBName
			}
		}
	}
	if req.ServerPort > 0 {
		h.cfg.Server.Port = req.ServerPort
	}

	// Save configuration
	if err := h.cfg.Save("config.yaml"); err != nil {
		log.Printf("Failed to save configuration: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save configuration", "details": err.Error()})
		return
	}

	log.Printf("Setup completed successfully for user: %s", req.Username)
	c.JSON(http.StatusOK, gin.H{"message": "Setup completed successfully"})
}

// Authentication
func (h *Handlers) LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login - Waterlogger",
	})
}

func (h *Handlers) LoginAPI(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !middleware.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Create session
	sessionID := middleware.CreateSession(user.ID)
	c.SetCookie("session", sessionID, 86400, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func (h *Handlers) LogoutAPI(c *gin.Context) {
	c.SetCookie("session", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// Dashboard
func (h *Handlers) Dashboard(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"title": "Dashboard - Waterlogger",
	})
}

// Pools
func (h *Handlers) PoolsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "pools.html", gin.H{
		"title": "Pools - Waterlogger",
	})
}

func (h *Handlers) GetPools(c *gin.Context) {
	var pools []models.Pool
	if err := h.db.Find(&pools).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pools"})
		return
	}
	c.JSON(http.StatusOK, pools)
}

func (h *Handlers) CreatePool(c *gin.Context) {
	var pool models.Pool
	if err := c.ShouldBindJSON(&pool); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&pool).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create pool"})
		return
	}

	c.JSON(http.StatusCreated, pool)
}

func (h *Handlers) UpdatePool(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pool ID"})
		return
	}

	var pool models.Pool
	if err := h.db.First(&pool, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pool not found"})
		return
	}

	if err := c.ShouldBindJSON(&pool); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Save(&pool).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update pool"})
		return
	}

	c.JSON(http.StatusOK, pool)
}

func (h *Handlers) DeletePool(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pool ID"})
		return
	}

	if err := h.db.Delete(&models.Pool{}, uint(id)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete pool"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pool deleted successfully"})
}

// Samples
func (h *Handlers) SamplesPage(c *gin.Context) {
	c.HTML(http.StatusOK, "samples.html", gin.H{
		"title": "Samples - Waterlogger",
	})
}

func (h *Handlers) GetSamples(c *gin.Context) {
	var samples []models.Sample
	if err := h.db.Preload("Pool").Preload("User").Preload("Kit").
		Preload("Measurements").Preload("Indices").Find(&samples).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch samples"})
		return
	}
	c.JSON(http.StatusOK, samples)
}

func (h *Handlers) CreateSample(c *gin.Context) {
	var sample models.Sample
	if err := c.ShouldBindJSON(&sample); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&sample).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create sample"})
		return
	}

	c.JSON(http.StatusCreated, sample)
}

func (h *Handlers) UpdateSample(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sample ID"})
		return
	}

	var sample models.Sample
	if err := h.db.First(&sample, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sample not found"})
		return
	}

	if err := c.ShouldBindJSON(&sample); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Save(&sample).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update sample"})
		return
	}

	c.JSON(http.StatusOK, sample)
}

func (h *Handlers) DeleteSample(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sample ID"})
		return
	}

	if err := h.db.Delete(&models.Sample{}, uint(id)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete sample"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sample deleted successfully"})
}

// Placeholder handlers for remaining endpoints
func (h *Handlers) ChartsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "charts.html", gin.H{
		"title": "Charts - Waterlogger",
	})
}

func (h *Handlers) ExportPage(c *gin.Context) {
	c.HTML(http.StatusOK, "export.html", gin.H{
		"title": "Export - Waterlogger",
	})
}

func (h *Handlers) SettingsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "settings.html", gin.H{
		"title": "Settings - Waterlogger",
	})
}

func (h *Handlers) GetUsers(c *gin.Context) {
	var users []models.User
	if err := h.db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *Handlers) CreateUser(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

func (h *Handlers) UpdateUser(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

func (h *Handlers) DeleteUser(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

func (h *Handlers) GetKits(c *gin.Context) {
	var kits []models.Kit
	if err := h.db.Find(&kits).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch kits"})
		return
	}
	c.JSON(http.StatusOK, kits)
}

func (h *Handlers) CreateKit(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

func (h *Handlers) UpdateKit(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

func (h *Handlers) DeleteKit(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

func (h *Handlers) GetChartData(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

func (h *Handlers) ExportExcel(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

func (h *Handlers) ExportMarkdown(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

func (h *Handlers) GetSettings(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

func (h *Handlers) UpdateSettings(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}