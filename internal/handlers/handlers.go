package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"waterlogger/internal/chemistry"
	"waterlogger/internal/config"
	"waterlogger/internal/middleware"
	"waterlogger/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		"BuildTime": c.MustGet("BuildTime"),
		"BuildDate": c.MustGet("BuildDate"),
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
		"BuildTime": c.MustGet("BuildTime"),
		"BuildDate": c.MustGet("BuildDate"),
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
		"BuildTime": c.MustGet("BuildTime"),
		"BuildDate": c.MustGet("BuildDate"),
	})
}

// Pools
func (h *Handlers) PoolsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "pools.html", gin.H{
		"title": "Pools - Waterlogger",
		"BuildTime": c.MustGet("BuildTime"),
		"BuildDate": c.MustGet("BuildDate"),
	})
}

func (h *Handlers) KitsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "kits.html", gin.H{
		"title": "Test Kits - Waterlogger",
		"BuildTime": c.MustGet("BuildTime"),
		"BuildDate": c.MustGet("BuildDate"),
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
		"BuildTime": c.MustGet("BuildTime"),
		"BuildDate": c.MustGet("BuildDate"),
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

	fmt.Printf("DEBUG: Parsed sample: %+v\n", sample)
	if sample.Measurements != nil {
		fmt.Printf("DEBUG: Sample measurements: %+v\n", sample.Measurements)
	}

	// Set user ID if not provided
	if sample.UserID == 0 {
		sample.UserID = getUserID(c)
	}

	// Create sample in database - GORM will automatically create associated measurements
	if err := h.db.Create(&sample).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create sample"})
		return
	}

	// Calculate water chemistry indices if we have measurements with the minimum required data
	if sample.Measurements != nil && sample.Measurements.PH != 0 {
		if indices, err := chemistry.CalculateIndices(sample.Measurements); err == nil {
			indices.SampleID = sample.ID
			if err := h.db.Create(indices).Error; err != nil {
				// Log error but don't fail the request
				fmt.Printf("Warning: Failed to create indices: %v\n", err)
			} else {
				sample.Indices = indices
			}
		}
	}

	// Load the complete sample with all relationships
	if err := h.db.Preload("Pool").Preload("Measurements").Preload("Indices").First(&sample, sample.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load complete sample"})
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
	if err := h.db.Preload("Measurements").Preload("Indices").First(&sample, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sample not found"})
		return
	}

	if err := c.ShouldBindJSON(&sample); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Store measurements separately to avoid GORM auto-save conflicts
	measurementsData := sample.Measurements
	sample.Measurements = nil
	sample.Indices = nil

	// Update sample
	if err := h.db.Save(&sample).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update sample"})
		return
	}

	// Restore measurements data for processing
	sample.Measurements = measurementsData

	// Update or create measurements if provided
	if sample.Measurements != nil {
		sample.Measurements.SampleID = sample.ID
		
		// Find existing measurements for this sample
		var existingMeasurements models.Measurements
		if err := h.db.Where("sample_id = ?", sample.ID).First(&existingMeasurements).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// No existing measurements, create new ones
				if err := h.db.Create(sample.Measurements).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create measurements"})
					return
				}
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query existing measurements"})
				return
			}
		} else {
			// Update existing measurements
			sample.Measurements.ID = existingMeasurements.ID
			if err := h.db.Save(sample.Measurements).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update measurements"})
				return
			}
		}

		// Always delete existing indices first
		h.db.Where("sample_id = ?", sample.ID).Delete(&models.Indices{})

		// Recalculate indices if we have the minimum required data
		if sample.Measurements.PH != 0 {
			if indices, err := chemistry.CalculateIndices(sample.Measurements); err == nil {
				indices.SampleID = sample.ID
				
				if err := h.db.Create(indices).Error; err != nil {
					fmt.Printf("Warning: Failed to update indices: %v\n", err)
				} else {
					sample.Indices = indices
				}
			} else {
				fmt.Printf("Warning: Failed to calculate indices: %v\n", err)
			}
		} else {
			// No pH data - indices were already deleted, so indices will be null
			sample.Indices = nil
		}
	}

	// Load the complete updated sample with all relationships
	if err := h.db.Preload("Pool").Preload("Measurements").Preload("Indices").First(&sample, sample.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load updated sample"})
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
		"BuildTime": c.MustGet("BuildTime"),
		"BuildDate": c.MustGet("BuildDate"),
	})
}


func (h *Handlers) ExportPage(c *gin.Context) {
	c.HTML(http.StatusOK, "export.html", gin.H{
		"title": "Export - Waterlogger",
		"BuildTime": c.MustGet("BuildTime"),
		"BuildDate": c.MustGet("BuildDate"),
	})
}

func (h *Handlers) SettingsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "settings.html", gin.H{
		"title": "Settings - Waterlogger",
		"BuildTime": c.MustGet("BuildTime"),
		"BuildDate": c.MustGet("BuildDate"),
	})
}

func (h *Handlers) GetUsers(c *gin.Context) {
	var users []models.User
	if err := h.db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	
	// Don't return passwords in the response
	for i := range users {
		users[i].Password = ""
	}
	
	c.JSON(http.StatusOK, users)
}

func (h *Handlers) CreateUser(c *gin.Context) {
	var createData struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	
	if err := c.ShouldBindJSON(&createData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate password requirements
	if errors := middleware.ValidatePassword(createData.Password); len(errors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password validation failed", "details": errors})
		return
	}

	// Hash password
	hashedPassword, err := middleware.HashPassword(createData.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user model
	user := models.User{
		Username: createData.Username,
		Email:    createData.Email,
		Password: hashedPassword,
	}

	// Create user
	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Don't return the password in the response
	user.Password = ""
	c.JSON(http.StatusCreated, user)
}

func (h *Handlers) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var updateData struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	if updateData.Username != "" {
		user.Username = updateData.Username
	}
	if updateData.Email != "" {
		user.Email = updateData.Email
	}
	
	// Update password if provided
	if updateData.Password != "" {
		// Validate password requirements
		if errors := middleware.ValidatePassword(updateData.Password); len(errors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password validation failed", "details": errors})
			return
		}

		// Hash password
		hashedPassword, err := middleware.HashPassword(updateData.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		user.Password = hashedPassword
	}

	// Save user
	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Don't return the password in the response
	user.Password = ""
	c.JSON(http.StatusOK, user)
}

func (h *Handlers) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Prevent deletion of the last user
	var userCount int64
	if err := h.db.Model(&models.User{}).Count(&userCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user count"})
		return
	}
	
	if userCount <= 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete the last user"})
		return
	}

	// Delete user
	if err := h.db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *Handlers) GetKits(c *gin.Context) {
	var kits []models.Kit
	if err := h.db.Find(&kits).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch kits"})
		return
	}
	c.JSON(http.StatusOK, kits)
}

// getUserID extracts the user ID from the request context
func getUserID(c *gin.Context) uint {
	// For now, return a default user ID of 1 (admin user)
	// In a real implementation, this would get the user ID from the session/token
	return 1
}

func (h *Handlers) CreateKit(c *gin.Context) {
	var kit models.Kit
	if err := c.ShouldBindJSON(&kit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	// Validate required fields
	if kit.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kit name is required"})
		return
	}

	// Set audit context
	ctx := context.WithValue(c.Request.Context(), "user_id", getUserID(c))
	if err := h.db.WithContext(ctx).Create(&kit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create kit"})
		return
	}

	c.JSON(http.StatusCreated, kit)
}

func (h *Handlers) UpdateKit(c *gin.Context) {
	kitID := c.Param("id")
	if kitID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kit ID is required"})
		return
	}

	// Check if kit exists
	var existingKit models.Kit
	if err := h.db.First(&existingKit, kitID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Kit not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch kit"})
		}
		return
	}

	// Bind updated data
	var updates models.Kit
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	// Validate required fields
	if updates.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kit name is required"})
		return
	}

	// Set audit context and update
	ctx := context.WithValue(c.Request.Context(), "user_id", getUserID(c))
	if err := h.db.WithContext(ctx).Model(&existingKit).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update kit"})
		return
	}

	// Return updated kit
	if err := h.db.First(&existingKit, kitID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated kit"})
		return
	}

	c.JSON(http.StatusOK, existingKit)
}

func (h *Handlers) DeleteKit(c *gin.Context) {
	kitID := c.Param("id")
	if kitID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kit ID is required"})
		return
	}

	// Check if kit exists
	var kit models.Kit
	if err := h.db.First(&kit, kitID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Kit not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch kit"})
		}
		return
	}

	// Check if kit is being used by any samples
	var sampleCount int64
	if err := h.db.Model(&models.Sample{}).Where("kit_id = ?", kitID).Count(&sampleCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check kit usage"})
		return
	}

	if sampleCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete kit that is being used by samples"})
		return
	}

	// Delete the kit
	if err := h.db.Delete(&kit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete kit"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kit deleted successfully"})
}

func (h *Handlers) GetChartData(c *gin.Context) {
	poolID := c.Query("pool_id")
	parameter := c.Query("parameter")
	days := c.Query("days")
	
	fmt.Printf("DEBUG: Chart request - pool_id: %s, parameter: %s, days: %s\n", poolID, parameter, days)
	
	// If no parameter is specified, return all samples (for charts page)
	if parameter == "" {
		// Get query parameters
		daysInt := 30
		if days != "" {
			if d, err := strconv.Atoi(days); err == nil {
				daysInt = d
			}
		}
		
		// Calculate date range
		startDate := time.Now().AddDate(0, 0, -daysInt)
		
		// Build query
		query := h.db.Preload("Pool").Preload("Measurements").Preload("Indices")
		
		// Filter by pool if specified
		if poolID != "" {
			query = query.Where("pool_id = ?", poolID)
		}
		
		// Filter by date range using proper column name
		query = query.Where("sample_datetime >= ?", startDate)
		
		// Order by date
		query = query.Order("sample_datetime ASC")
		
		var samples []models.Sample
		if err := query.Find(&samples).Error; err != nil {
			log.Printf("ChartData DB error: %v", err) // Add this line
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chart data"})
			return
		}
		
		c.JSON(http.StatusOK, samples)
		return
	}
	
	// Original parameter-specific logic for backward compatibility
	if poolID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pool_id is required"})
		return
	}
	
	// Get samples with measurements for the specified pool
	var samples []models.Sample
	query := h.db.Preload("Measurements").Preload("Indices").Where("pool_id = ?", poolID)
	
	// Add date range filter if provided
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("sample_datetime >= ?", startDate)
	}
	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("sample_datetime <= ?", endDate)
	}
	
	if err := query.Order("sample_datetime ASC").Find(&samples).Error; err != nil {
		fmt.Printf("DEBUG: Chart query error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chart data"})
		return
	}
	
	fmt.Printf("DEBUG: Found %d samples\n", len(samples))
	
	// Extract data points for the specified parameter
	var dataPoints []map[string]interface{}
	
	for _, sample := range samples {
		point := map[string]interface{}{
			"date": sample.SampleDateTime,
		}
		
		// Get the value based on the parameter
		if sample.Measurements != nil {
			switch parameter {
			case "ph":
				point["value"] = sample.Measurements.PH
			case "fc":
				point["value"] = sample.Measurements.FC
			case "tc":
				point["value"] = sample.Measurements.TC
			case "ta":
				point["value"] = sample.Measurements.TA
			case "ch":
				point["value"] = sample.Measurements.CH
			case "cya":
				if sample.Measurements.CYA != nil {
					point["value"] = *sample.Measurements.CYA
				}
			case "temperature":
				point["value"] = sample.Measurements.Temperature
			case "salinity":
				if sample.Measurements.Salinity != nil {
					point["value"] = *sample.Measurements.Salinity
				}
			}
		}
		
		// Add indices data
		if sample.Indices != nil {
			if parameter == "lsi" && sample.Indices.LSI != nil {
				point["value"] = *sample.Indices.LSI
			} else if parameter == "rsi" && sample.Indices.RSI != nil {
				point["value"] = *sample.Indices.RSI
			}
		}
		
		// Only add point if it has a value
		if point["value"] != nil {
			dataPoints = append(dataPoints, point)
		}
	}
	
	fmt.Printf("DEBUG: Generated %d data points for parameter %s\n", len(dataPoints), parameter)
	
	c.JSON(http.StatusOK, gin.H{
		"parameter": parameter,
		"pool_id":   poolID,
		"data":      dataPoints,
	})
}

func (h *Handlers) ExportExcel(c *gin.Context) {
	// Get samples with related data
	var samples []models.Sample
	if err := h.db.Preload("Pool").Preload("Measurements").Preload("Indices").Find(&samples).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch samples"})
		return
	}
	
	// Generate CSV content (simplified Excel export)
	csvContent := "Sample Date,Pool Name,pH,Free Chlorine (ppm),Total Chlorine (ppm),Total Alkalinity (ppm),Calcium Hardness (ppm),Cyanuric Acid (ppm),Temperature (°F),Salinity (ppm),LSI,RSI,Notes\n"
	
	for _, sample := range samples {
		poolName := ""
		if sample.Pool != nil {
			poolName = sample.Pool.Name
		}
		
		// Format date
		date := sample.SampleDateTime.Format("2006-01-02 15:04:05")
		
		// Get measurement values
		ph := ""
		fc := ""
		tc := ""
		ta := ""
		ch := ""
		cya := ""
		temp := ""
		salinity := ""
		lsi := ""
		rsi := ""
		
		if sample.Measurements != nil {
			if sample.Measurements.PH != 0 {
				ph = fmt.Sprintf("%.2f", sample.Measurements.PH)
			}
			if sample.Measurements.FC != 0 {
				fc = fmt.Sprintf("%.2f", sample.Measurements.FC)
			}
			if sample.Measurements.TC != 0 {
				tc = fmt.Sprintf("%.2f", sample.Measurements.TC)
			}
			if sample.Measurements.TA != 0 {
				ta = fmt.Sprintf("%.2f", sample.Measurements.TA)
			}
			if sample.Measurements.CH != 0 {
				ch = fmt.Sprintf("%.2f", sample.Measurements.CH)
			}
			if sample.Measurements.CYA != nil {
				cya = fmt.Sprintf("%.2f", *sample.Measurements.CYA)
			}
			if sample.Measurements.Temperature != 0 {
				temp = fmt.Sprintf("%.1f", sample.Measurements.Temperature)
			}
			if sample.Measurements.Salinity != nil {
				salinity = fmt.Sprintf("%.2f", *sample.Measurements.Salinity)
			}
		}
		
		if sample.Indices != nil {
			if sample.Indices.LSI != nil {
				lsi = fmt.Sprintf("%.2f", *sample.Indices.LSI)
			}
			if sample.Indices.RSI != nil {
				rsi = fmt.Sprintf("%.2f", *sample.Indices.RSI)
			}
		}
		
		// Create CSV row
		csvContent += fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,\"%s\"\n",
			date, poolName, ph, fc, tc, ta, ch, cya, temp, salinity, lsi, rsi, sample.Notes)
	}
	
	// Set headers for file download
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=waterlogger_export.csv")
	c.String(http.StatusOK, csvContent)
}

func (h *Handlers) ExportBackup(c *gin.Context) {
	// Get all data for backup
	var users []models.User
	var pools []models.Pool
	var kits []models.Kit
	var samples []models.Sample
	
	if err := h.db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	
	if err := h.db.Find(&pools).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pools"})
		return
	}
	
	if err := h.db.Find(&kits).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch kits"})
		return
	}
	
	if err := h.db.Preload("Pool").Preload("Measurements").Preload("Indices").Find(&samples).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch samples"})
		return
	}
	
	// Create backup data structure
	backupData := map[string]interface{}{
		"users": users,
		"pools": pools,
		"kits": kits,
		"samples": samples,
		"exported_at": time.Now().Format("2006-01-02 15:04:05"),
		"version": "1.0.0",
	}
	
	// Generate JSON backup
	jsonData, err := json.MarshalIndent(backupData, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate backup"})
		return
	}
	
	// Generate filename
	filename := fmt.Sprintf("WL_backup_%s.json", time.Now().Format("20060102_150405"))
	
	// Set headers for file download
	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Length", fmt.Sprintf("%d", len(jsonData)))
	
	// Send the backup file
	c.Data(http.StatusOK, "application/json", jsonData)
}

func (h *Handlers) ExportMarkdown(c *gin.Context) {
	// Get samples with related data
	var samples []models.Sample
	if err := h.db.Preload("Pool").Preload("Measurements").Preload("Indices").Find(&samples).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch samples"})
		return
	}
	
	// Generate Markdown content
	mdContent := "# Waterlogger Export\n\n"
	mdContent += fmt.Sprintf("Generated on: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))
	
	if len(samples) == 0 {
		mdContent += "No samples found.\n"
	} else {
		mdContent += fmt.Sprintf("## Water Test Results (%d samples)\n\n", len(samples))
		
		for _, sample := range samples {
			poolName := "Unknown Pool"
			if sample.Pool != nil {
				poolName = sample.Pool.Name
			}
			
			mdContent += fmt.Sprintf("### %s - %s\n\n", poolName, sample.SampleDateTime.Format("2006-01-02 15:04:05"))
			
			if sample.Measurements != nil {
				mdContent += "**Chemical Measurements:**\n"
				if sample.Measurements.PH != 0 {
					mdContent += fmt.Sprintf("- pH: %.2f\n", sample.Measurements.PH)
				}
				if sample.Measurements.FC != 0 {
					mdContent += fmt.Sprintf("- Free Chlorine: %.2f ppm\n", sample.Measurements.FC)
				}
				if sample.Measurements.TC != 0 {
					mdContent += fmt.Sprintf("- Total Chlorine: %.2f ppm\n", sample.Measurements.TC)
				}
				if sample.Measurements.TA != 0 {
					mdContent += fmt.Sprintf("- Total Alkalinity: %.2f ppm\n", sample.Measurements.TA)
				}
				if sample.Measurements.CH != 0 {
					mdContent += fmt.Sprintf("- Calcium Hardness: %.2f ppm\n", sample.Measurements.CH)
				}
				if sample.Measurements.CYA != nil {
					mdContent += fmt.Sprintf("- Cyanuric Acid: %.2f ppm\n", *sample.Measurements.CYA)
				}
				if sample.Measurements.Temperature != 0 {
					mdContent += fmt.Sprintf("- Temperature: %.1f°F\n", sample.Measurements.Temperature)
				}
				if sample.Measurements.Salinity != nil {
					mdContent += fmt.Sprintf("- Salinity: %.2f ppm\n", *sample.Measurements.Salinity)
				}
				mdContent += "\n"
			}
			
			if sample.Indices != nil {
				mdContent += "**Water Balance Indices:**\n"
				if sample.Indices.LSI != nil {
					mdContent += fmt.Sprintf("- LSI (Langelier Saturation Index): %.2f\n", *sample.Indices.LSI)
				}
				if sample.Indices.RSI != nil {
					mdContent += fmt.Sprintf("- RSI (Ryznar Stability Index): %.2f\n", *sample.Indices.RSI)
				}
				mdContent += "\n"
			}
			
			if sample.Notes != "" {
				mdContent += fmt.Sprintf("**Notes:** %s\n\n", sample.Notes)
			}
			
			mdContent += "---\n\n"
		}
	}
	
	// Add appendices at the bottom
	mdContent += "\n\n---\n\n# Appendices\n\n"
	
	// Read all markdown files from the appendices directory
	appendicesDir := "appendices"
	if files, err := ioutil.ReadDir(appendicesDir); err == nil {
		// Sort files by name for consistent ordering
		sort.Slice(files, func(i, j int) bool {
			return files[i].Name() < files[j].Name()
		})
		
		for _, file := range files {
			if filepath.Ext(file.Name()) == ".md" {
				filePath := filepath.Join(appendicesDir, file.Name())
				if content, err := ioutil.ReadFile(filePath); err == nil {
					mdContent += string(content) + "\n\n"
				} else {
					log.Printf("Warning: Failed to read appendix file %s: %v", filePath, err)
				}
			}
		}
	} else {
		log.Printf("Warning: Failed to read appendices directory: %v", err)
	}
	
	// Set headers for file download
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=\"waterlogger_markdown.md\"")
	c.String(http.StatusOK, mdContent)
}

func (h *Handlers) GetSettings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get user preferences
	var preferences models.UserPreferences
	if err := h.db.Where("user_id = ?", userID).First(&preferences).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create default preferences
			preferences = models.UserPreferences{
				UserID:     userID.(uint),
				UnitSystem: "imperial",
			}
			preferences.CreatedBy = userID.(uint)
			preferences.UpdatedBy = userID.(uint)
			
			if err := h.db.Create(&preferences).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create default preferences"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load preferences"})
			return
		}
	}

	// Get system information
	systemInfo := gin.H{
		"database_type": h.cfg.Database.Type,
		"server_port":   h.cfg.Server.Port,
		"app_version":   h.cfg.App.Version,
	}

	c.JSON(http.StatusOK, gin.H{
		"preferences": preferences,
		"system":      systemInfo,
	})
}

func (h *Handlers) UpdateSettings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		UnitSystem string `json:"unit_system" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Validate unit system
	if req.UnitSystem != "imperial" && req.UnitSystem != "metric" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid unit system. Must be 'imperial' or 'metric'"})
		return
	}

	// Update or create user preferences
	var preferences models.UserPreferences
	result := h.db.Where("user_id = ?", userID).First(&preferences)
	
	if result.Error == gorm.ErrRecordNotFound {
		// Create new preferences
		preferences = models.UserPreferences{
			UserID:     userID.(uint),
			UnitSystem: req.UnitSystem,
		}
		preferences.CreatedBy = userID.(uint)
		preferences.UpdatedBy = userID.(uint)
		
		if err := h.db.Create(&preferences).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create preferences"})
			return
		}
	} else if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load preferences"})
		return
	} else {
		// Update existing preferences
		preferences.UnitSystem = req.UnitSystem
		preferences.UpdatedBy = userID.(uint)
		
		if err := h.db.Save(&preferences).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update preferences"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Settings updated successfully",
		"preferences": preferences,
	})
}

// Unit conversion endpoints
func (h *Handlers) ConvertUnits(c *gin.Context) {
	var req struct {
		Value     float64 `json:"value" binding:"required"`
		Parameter string  `json:"parameter" binding:"required"`
		FromSystem string `json:"from_system" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var fromSystem chemistry.UnitSystem
	if req.FromSystem == "metric" {
		fromSystem = chemistry.Metric
	} else {
		fromSystem = chemistry.Imperial
	}

	converted := chemistry.ConvertMeasurement(req.Value, req.Parameter, fromSystem)
	c.JSON(http.StatusOK, converted)
}