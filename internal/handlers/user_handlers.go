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

// RegisterUser gère l'inscription d'un utilisateur
func RegisterUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Requête invalide"},
		)
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Echec du hachage du mot de passe."},
		)
		return
	}

	user.Password = string(hashedPass)

	// Enregistre l'utilisateur dans la BDD
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Echec de l'inscription"},
		)
		log.Println("Erreur: echec de la création de l'utilisateur", user, err)
		return
	}
	c.JSON(
		http.StatusCreated,
		gin.H{"message": "Utilisateur enregistré avec succès !"},
	)
}

// gèrela connexion de l'utilisateur
func LoginHandler(c *gin.Context) {
	var req models.User
	var user models.User

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Requête invalide"},
		)
		return
	}

	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "Utilisateur non trouve."},
		)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(req.Password), []byte(user.Password)); err != nil {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": " Mot de passe invalide."},
		)
		return
	}

	token, err := utils.GenerateJWT(strconv.Itoa(int(user.ID)), user.Email)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Echec de la création du token jwt."},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{"token": token},
	)
}
