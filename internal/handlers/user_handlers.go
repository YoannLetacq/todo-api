package handlers

import (
	"YoannLetacq/todo-api.git/config"
	"YoannLetacq/todo-api.git/internal/models"
	"YoannLetacq/todo-api.git/internal/utils"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Requête invalide"})
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Échec du hachage du mot de passe."})
		return
	}

	user.Password = string(hashedPass)

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Échec de l'inscription"})
		log.Println("Erreur : échec de la création de l'utilisateur", user, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Utilisateur enregistré avec succès !"})
}

func LoginHandler(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var user models.User

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Requête invalide"})
		return
	}

	// Vérifier si l'utilisateur existe
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	// Vérifier le mot de passe
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Mot de passe invalide"})
		return
	}

	// Générer le token JWT en convertissant correctement l'ID en string
	token, err := utils.GenerateJWT(strconv.Itoa(int(user.ID)), user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Échec de la génération du token JWT"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
