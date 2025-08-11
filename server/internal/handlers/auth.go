package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/Motionists/mt-order/internal/models"
)

type registerReq struct {
	Mobile   string `json:"mobile" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}
type loginReq struct {
	Mobile   string `json:"mobile" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pw, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := models.User{Mobile: req.Mobile, Password: string(pw)}
	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "mobile exists"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": user.ID})
}

func (h *Handler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	if err := h.db.Where("mobile = ?", req.Mobile).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	claims := jwt.MapClaims{
		"uid": user.ID,
		"exp": time.Now().Add(time.Duration(h.cfg.JWT.TTLMin) * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, _ := token.SignedString([]byte(h.cfg.JWT.Secret))
	c.JSON(http.StatusOK, gin.H{"token": ss})
}

func (h *Handler) Me(c *gin.Context) {
	var user models.User
	if err := h.db.First(&user, uid(c)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": user.ID, "mobile": user.Mobile, "nickname": user.Nickname, "avatar": user.AvatarURL})
}