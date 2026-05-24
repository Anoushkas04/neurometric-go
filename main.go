package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// --- MODELS ---

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"uniqueIndex" json:"username"`
	PasswordHash string    `json:"-"`
	FullName     string    `json:"full_name"`
	Age          int       `json:"age"`
	Gender       string    `json:"gender"`
	CreatedAt    time.Time `json:"created_at"`
}

type Assessment struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `json:"user_id"`
	OverallScore  int       `json:"overall_score"`
	Accuracy      int       `json:"accuracy"`
	TotalErrors   int       `json:"total_errors"`
	TimeSec       float64   `json:"time_sec"`
	AvgRT         float64   `json:"avg_rt_ms"`
	DetailedStats string    `json:"detailed_stats"` // JSON string for breakdown
	CreatedAt     time.Time `json:"created_at"`
}

var db *gorm.DB
var jwtSecret = []byte("clinical-secret-key-2024") // In production, use env var

func initDB() {
	// Ensure data directory exists (crucial for Docker/Persistent volumes)
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		os.Mkdir("data", 0755)
	}

	var err error
	db, err = gorm.Open(sqlite.Open("data/neurometric.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto Migrate the schemas
	db.AutoMigrate(&User{}, &Assessment{})
}

// --- HELPERS ---

func generateToken(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	return token.SignedString(jwtSecret)
}

// --- MIDDLEWARE ---

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		c.Set("user_id", uint(claims["user_id"].(float64)))
		c.Next()
	}
}

// --- HANDLERS ---

func register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		FullName string `json:"full_name"`
		Age      int    `json:"age"`
		Gender   string `json:"gender"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	user := User{
		Username:     input.Username,
		PasswordHash: string(hashedPassword),
		FullName:     input.FullName,
		Age:          input.Age,
		Gender:       input.Gender,
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

func login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user User
	if err := db.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user": gin.H{
			"username":  user.Username,
			"full_name": user.FullName,
		},
	})
}

func submitAssessment(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var input struct {
		OverallScore  int     `json:"Overall_Score"`
		Accuracy      int     `json:"Accuracy_Percent"`
		TotalErrors   int     `json:"Total_Errors"`
		TimeSec       float64 `json:"Completion_Time_Sec"`
		AvgRT         float64 `json:"Avg_Response_Time_Ms"`
		DetailedStats string  `json:"detailed_stats"` // We'll receive this as a string or raw JSON
	}

	// For simplicity, we'll bind the whole payload
	// In a real app, you'd want to validate each field
	if err := c.ShouldBindJSON(&input); err != nil {
		// If detailedStats is not a string but an object, we need to handle it
		// Let's try to capture the whole raw body for detailedStats if needed
	}

	assessment := Assessment{
		UserID:        userID,
		OverallScore:  input.OverallScore,
		Accuracy:      input.Accuracy,
		TotalErrors:   input.TotalErrors,
		TimeSec:       input.TimeSec,
		AvgRT:         input.AvgRT,
		DetailedStats: input.DetailedStats,
	}

	if err := db.Create(&assessment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save assessment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Assessment saved successfully", "id": assessment.ID})
}

func main() {
	initDB()

	r := gin.Default()

	// API Routes (Register these FIRST)
	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", register)
			auth.POST("/login", login)
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(AuthMiddleware())
		{
			protected.GET("/user/profile", func(c *gin.Context) {
				userID := c.MustGet("user_id").(uint)
				var user User
				db.First(&user, userID)
				c.JSON(http.StatusOK, user)
			})

			protected.POST("/assessments", submitAssessment)
			protected.GET("/assessments", func(c *gin.Context) {
				userID := c.MustGet("user_id").(uint)
				var assessments []Assessment
				db.Where("user_id = ?", userID).Order("created_at desc").Find(&assessments)
				c.JSON(http.StatusOK, assessments)
			})
		}
	}

	// Static Files (Serve these AFTER API routes)
	r.StaticFile("/", "./public/index.html")
	r.StaticFile("/index.html", "./public/index.html")
	r.Static("/js", "./public/js")
	
	// Fallback for other files in public/
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		c.File("./public" + path)
	})

	log.Println("Server starting on http://localhost:8080")
	r.Run(":8080")
}
