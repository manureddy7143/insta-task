package controller

import (
	"strings"
	"time"

	"net/http"

	"github.com/dgrijalva/jwt-go"
	models "github.com/manureddy7143/GolangStarter/source/model"
	"github.com/manureddy7143/GolangStarter/source/repository"
	dto "github.com/manureddy7143/GolangStarter/source/service"
	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

//UserManagement controller
type UserManagement struct{}

// LoginPath - URL Path for Login
const LoginPath = "/login"

// RegisterPath - URL Path for Register
const RegisterPath = "/register"

// ProfilePath - URL Path for GetProfile
const ProfilePath = "/profile"

const secretKey = "secret"

var authRepository repository.AuthRepository = repository.AuthRepository{}

// Register function
// @Summary Register In API
// @Description Register  API for Usres and to get Access Token
// @Request Body containing username, password, email, firstname, lastname"
// @Produce json
// @Router /auth/users/register [post]
func (userManagement UserManagement) Register(c *gin.Context) {
	//Defining User
	var person dto.RequestRegister
	err := c.ShouldBindJSON(&person)
	if err != nil {
		log.Panic().Msgf("Binding error")
		errorHandling(c, "err-5001", err)
	}
	//Generate Password using bcrypt
	password, _ := bcrypt.GenerateFromPassword([]byte(person.Password), 14)
	if person.Username == "" || person.Firstname == "" || person.Lastname == "" || person.Email == "" {
		log.Panic().Msgf("Bad Request")
		errorHandling(c, "err-4000", err)
		return
	}
	// Verifying the username and email
	aa, err := authRepository.FindUsers(map[string]interface{}{"username": person.Username, "email": person.Email})
	if err != nil {
		log.Panic().Msgf("Database Error")
		errorHandling(c, "err-5000", err)
		return
	}
	if len(aa) > 0 {
		log.Panic().Msgf("User Already Exists")
		errorHandling(c, "err-4001", err)
		return
	}
	//Defining user to create db entry
	user := models.Users{
		Username:  person.Username,
		Firstname: person.Firstname,
		Lastname:  person.Lastname,
		Email:     person.Email,
		Password:  password,
	}
	//Calling Db method to create User
	UserId, err := authRepository.CreateUsers(user)
	if err != nil {
		log.Panic().Msgf("Database Error")
		errorHandling(c, "err-5000", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Registed", "UserId": UserId,
	})
}

// Login function
// @Summary Sign In API
// @Description Sign In API for authentication and to get Access token
// @ Login body dto.RequestLogin true "Request body containing user credentials"
// @Produce json
// @Success 200 {object} dto.RespLogin "Access,Refresh Tokens"
// @Failure 401 {object} dto.ErrorDTO "Invalid Credentials"
// @Failure 500 {object} dto.ErrorDTO "Internal Server Error"
// @Router /auth/users/login [post]
func (a UserManagement) Login(c *gin.Context) {
	//Defining User
	var person dto.RequestLogin
	err := c.ShouldBind(&person)
	if err != nil {
		log.Panic().Msgf("Binding error")
		errorHandling(c, "err-5001", err)
	}
	// Verifying the username and email
	if person.Email == "" || person.Password == "" {
		log.Panic().Msgf("Bad Request")
		errorHandling(c, "err-4000", err)
		return
	}
	//Calling db method to find Users
	Users, err := authRepository.FindUsers(map[string]interface{}{"email": person.Email})
	if err != nil {
		log.Panic().Msgf("Database Error")
		errorHandling(c, "err-4001", err)
		return
	}
	if len(Users) == 0 {
		log.Panic().Msgf("User Dont Exists")
		errorHandling(c, "err-4002", err)
		return
	}
	//Comparing Passwords
	if err := bcrypt.CompareHashAndPassword(Users[0].Password, []byte(person.Password)); err != nil {
		c.Status(http.StatusNotFound)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Incorrect Password"})
		return
	}
	//Adding issuer claims with email
	clamins := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    person.Email,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})
	//Generating Token
	token, err := clamins.SignedString([]byte(secretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not login",
		})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"message": "Logged In successfully", "accessToken": token,
	})
}

// Profile  function
// @Summary Profile  In API
// @Description Sign In API for authentication and to get Access, Refresh token
// @Request Body containers access token
// @Produce json
// @Success 200 {object} dto.ResponseLogin "Access,Refresh Tokens"
// @Failure 401 {object} dto.ErrorDTO "Invalid Credentials"
// @Failure 500 {object} dto.ErrorDTO "Internal Server Error"
// @Router /auth/users/profile[post]
func (a UserManagement) Profile(c *gin.Context) {
	// Retriving The Token
	reqToken := c.Request.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	accessToken := splitToken[1]
	ac := &jwt.StandardClaims{}
	//Parsing claims
	token, err := jwt.ParseWithClaims(accessToken, ac, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Go To Login Page",
		})
		return
	}
	claims := token.Claims.(*jwt.StandardClaims)
	// Finding Users
	Users, err := authRepository.FindUsers(map[string]interface{}{"email": claims.Issuer})
	if err != nil {
		log.Panic().Msgf("Database Error")
		errorHandling(c, "err-4001", err)
		return
	}
	ResponseUserProfile := dto.RespProfile{
		Username:  Users[0].Firstname,
		Email:     Users[0].Email,
		Firstname: Users[0].Firstname,
		Lastname:  Users[0].Lastname,
	}
	c.JSON(http.StatusAccepted, ResponseUserProfile)
}
