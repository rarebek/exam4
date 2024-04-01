package v1

import (
	models "4microservice/api_gateway/api/handlers/models"
	pbu "4microservice/api_gateway/genproto/user_service"
	"4microservice/api_gateway/pkg/codegen"
	"4microservice/api_gateway/pkg/logger"
	"4microservice/api_gateway/pkg/password"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/k0kubun/pp"
	"google.golang.org/protobuf/encoding/protojson"
)

type PageData struct {
	VerificationLink string
}

// RegisterUser godoc
// @Summary Register a User
// @Description Registers a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User information"
// @Success 200 {object} models.RegisterUserResponse
// @Failure 400 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
// @Router /v1/user/register [post]
func (h *HandlerV1) RegisterUser(c *gin.Context) {
	var (
		body        models.UserModelForRegister
		code        string
		jsbpMarshal protojson.MarshalOptions
	)
	jsbpMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
	}

	body.Email = strings.TrimSpace(body.Email)
	body.Email = strings.ToLower(body.Email)

	//err = body.Validate()
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, models.ResponseError{
	//		Code:  http.StatusBadRequest,
	//		Error: "wrong format, correct your email or password format amd try again",
	//	})
	//	return
	//}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	exists, err := h.serviceManager.UserService().CheckField(ctx, &pbu.CheckFieldRequest{
		Field: "email",
		Value: body.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		h.log.Error("failed to check uniqueness: ", logger.Error(err))
		return
	}

	if exists.Unique {
		c.JSON(http.StatusConflict, models.ResponseError{
			Code:  http.StatusConflict,
			Error: "this email is already in use",
		})
		h.log.Error("email is already exist in database")
		return
	}
	code = codegen.GenerateCode()
	type PageData struct {
		OTP string
	}
	tpl := template.Must(template.ParseFiles("verify_email_template.html"))
	data := PageData{
		OTP: code,
	}
	var buf bytes.Buffer
	tpl.Execute(&buf, data)
	htmlContent := buf.Bytes()

	auth := smtp.PlainAuth("", "nodirbekgolang@gmail.com", "ecncwhvfdyvjghux", "smtp.gmail.com")
	err = smtp.SendMail("smtp.gmail.com:587", auth, "nodirbekgolang@gmail.com", []string{body.Email}, []byte("To: "+body.Email+"\r\nSubject: Email verification\r\nMIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"+string(htmlContent)))
	if err != nil {
		log.Fatalf("Error sending otp to email: %v", err)
	}
	log.Println("Email sent successfully")
	body.OTP = code

	byteUser, err := json.Marshal(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		h.log.Error("failed while marshalling user data")
		return
	}

	if err := h.reds.SetWithTTL(body.Email, string(byteUser), int(time.Second)*300); err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		h.log.Error("cannot set redis")
		return
	}

	c.JSON(http.StatusOK, models.RegisterUserResponse{
		Message: "One time verification password sent to your email. Please verify",
	})
}

// Verify User
// @Summary verify user
// @Tags users
// @Description Verify a user with code sent to their email
// @Accept json
// @Product json
// @Param email path string true "email"
// @Param code path string true "code"
// @Success 200 {object} models.User
// @Failure 400 string error models.ResponseError
// @Failure 400 string error models.ResponseError
// @Router /v1/user/verify/{email}/{code} [post]
func (h *HandlerV1) Verify(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	userEmail := c.Param("email")
	code := c.Param("code")

	user, err := redis.Bytes(h.reds.Get(userEmail))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Code:  http.StatusUnauthorized,
			Error: "code is expired, try again.",
		})
		h.log.Error("Code is expired, TTL is over.")
		return
	}

	var respUser models.UserModelForRegister
	if err := json.Unmarshal(user, &respUser); err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		h.log.Error("cannot unmarshal user from redis", logger.Error(err))
		fmt.Println(respUser)
		return
	}

	if respUser.OTP != code {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Code:  http.StatusBadRequest,
			Error: "code is incorrect, try again.",
		})
		h.log.Error("verification failed", logger.Error(err))
		return
	}
	pp.Println(respUser)
	respUser.Password, err = password.HashPassword(respUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		h.log.Error("cannot hash the password", logger.Error(err))
		return
	}

	access, refresh, err := h.jwthandler.GenerateAuthJWT("user", respUser.Username, "byRegister")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error(err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	id := uuid.New().String()

	createdUser, err := h.serviceManager.UserService().CreateUser(ctx, &pbu.User{
		Id:           id,
		Username:     respUser.Username,
		Email:        respUser.Email,
		Password:     respUser.Password,
		FirstName:    respUser.FirstName,
		LastName:     respUser.LastName,
		Bio:          respUser.Bio,
		Website:      respUser.Website,
		RefreshToken: refresh,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("cannot create user", logger.Error(err))
		return
	}

	c.Header("Authorization", "Bearer "+access)
	c.JSON(http.StatusOK, createdUser)
}

// Login LoginUser godoc
// @Summary Log In User
// @Description Api for Logging in
// @Tags users
// @Accept json
// @Produce json
// @Param email path string true "Email"
// @Param password path string true "Password"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} models.ResponseError
// @Failure 500 {object} models.ResponseError
// @Router /v1/user/login/{email}/{password} [post]
func (h *HandlerV1) Login(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	email := c.Param("email")
	userpass := c.Param("password")

	jspbMarshal.UseProtoNames = true

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	resp, err := h.serviceManager.UserService().GetUserByEmail(ctx, &pbu.Email{
		Email: email,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get user by email", logger.Error(err))
		return
	}

	if !password.CompareHashPassword(resp.Password, userpass) {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Code:  http.StatusBadRequest,
			Error: "wrong password",
		})
		return
	}

	access, _, err := h.jwthandler.GenerateAuthJWT("user", resp.Username, "byRegister")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		h.log.Error("cannot create access token", logger.Error(err))
		return
	}

	res := models.LoginResponse{
		Message:     "Successfully logged in",
		AccessToken: access,
	}

	c.Header("Authorization", "Bearer "+access)
	c.JSON(http.StatusOK, res)
}

// CreateUser godoc
// @Summary Create a new user
// @Description Creates a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User information"
// @Success 200 {object} models.User "Created user"
// @Failure 400 {object} models.ResponseError "Invalid user data"
// @Failure 401 {object} models.ResponseError "Unauthorized"
// @Failure 500 {object} models.ResponseError "Internal server error"
// @Router /v1/admin/user/create [post]
func (h *HandlerV1) CreateUser(c *gin.Context) {
	var (
		body        models.User
		jsbpMarshal protojson.MarshalOptions
	)

	jsbpMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Error: err.Error(),
			Code:  http.StatusBadRequest,
		})
		h.log.Error(err.Error())
		return
	}
	var user models.User
	access, refresh, err := h.jwthandler.GenerateAuthJWT("user", body.Username, "byAdmin")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error(err.Error())
		return
	}

	user.RefreshToken = refresh

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	createdUser, err := h.serviceManager.UserService().CreateUser(ctx, &pbu.User{
		Id:           body.Id,
		Username:     body.Username,
		Email:        body.Email,
		Password:     body.Password,
		FirstName:    body.FirstName,
		LastName:     body.LastName,
		Bio:          body.Bio,
		Website:      body.Website,
		RefreshToken: body.RefreshToken,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error(err.Error())
		return
	}

	userJSON, err := json.Marshal(user)

	h.natsConn.Publish("user.created", userJSON)

	c.Header("Authorization", "Bearer "+access)
	c.JSON(http.StatusOK, createdUser)
}

// GetUser godoc
// @Summary Get user information
// @Description Retrieves details of a specific user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.GetUserResponse "Retrieved user information"
// @Failure 400 {object} models.ResponseError "Invalid user ID"
// @Failure 404 {object} models.ResponseError "User not found"
// @Failure 500 {object} models.ResponseError "Internal server error"
// @Router /v1/user/get/{id} [get]
func (h *HandlerV1) GetUser(c *gin.Context) {
	userId := c.Param("id")

	if userId == "" {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Error: "Invalid user ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	user, err := h.serviceManager.UserService().GetUser(ctx, &pbu.UserIdd{Id: userId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Deletes a user from the system
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 204 "User deleted successfully"
// @Failure 400 {object} models.ResponseError "Invalid user ID"
// @Failure 404 {object} models.ResponseError "User not found"
// @Failure 500 {object} models.ResponseError "Internal server error"
// @Router /v1/user/delete/{id} [delete]
func (h *HandlerV1) DeleteUser(c *gin.Context) {
	userId := c.Param("id")

	if userId == "" {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Error: "Invalid user ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	_, err := h.serviceManager.UserService().DeleteUser(ctx, &pbu.UserIdd{Id: userId})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			c.JSON(http.StatusInternalServerError, models.ResponseError{
				Error: "Request timed out",
				Code:  http.StatusInternalServerError,
			})
			return
		} else if err.Error() == "rpc error: code = NotFound desc = user not found" {
			c.JSON(http.StatusNotFound, models.ResponseError{
				Error: "User not found",
				Code:  http.StatusNotFound,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})
		return
	}
	c.Status(http.StatusNoContent)
}

// UpdateUser godoc
// @Summary Update user information
// @Description Updates information of a specific user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body models.User true "Updated user information"
// @Success 200 {object} models.User "Updated user information"
// @Failure 400 {object} models.ResponseError "Invalid user data"
// @Failure 404 {object} models.ResponseError "User not found"
// @Failure 500 {object} models.ResponseError "Internal server error"
// @Router /v1/user/update/{id} [put]
func (h *HandlerV1) UpdateUser(c *gin.Context) {
	userID := c.Param("id")

	if userID == "" {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Error: "Invalid user ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	var updatedUser models.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Error: err.Error(),
			Code:  http.StatusBadRequest,
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	updatedUser.Id = userID // Set user ID
	updatedUserresp, err := h.serviceManager.UserService().UpdateUser(ctx, &pbu.User{
		Id:           updatedUser.Id,
		Username:     updatedUser.Username,
		Email:        updatedUser.Email,
		Password:     updatedUser.Password,
		FirstName:    updatedUser.FirstName,
		LastName:     updatedUser.LastName,
		Bio:          updatedUser.Bio,
		Website:      updatedUser.Website,
		RefreshToken: updatedUser.RefreshToken,
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			c.JSON(http.StatusInternalServerError, models.ResponseError{
				Error: "Request timed out",
				Code:  http.StatusInternalServerError,
			})
			return
		} else if err.Error() == "rpc error: code = NotFound desc = user not found" {
			c.JSON(http.StatusNotFound, models.ResponseError{
				Error: "User not found",
				Code:  http.StatusNotFound,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, updatedUserresp)
}

// GetAllUser godoc
// @Summary Get all users
// @Description Retrieves details of all users
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} models.GetAllResponse "List of users"
// @Failure 400 {object} models.ResponseError "Invalid request parameters"
// @Failure 500 {object} models.ResponseError "Internal server error"
// @Router /v1/users/{page}/{limit} [get]
func (h *HandlerV1) GetAllUser(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Error: "Invalid page number",
			Code:  http.StatusBadRequest,
		})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Error: "Invalid limit value",
			Code:  http.StatusBadRequest,
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	users, err := h.serviceManager.UserService().GetAllUser(ctx, &pbu.GetAllRequest{
		Page:  int64(page),
		Limit: int64(limit),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, users)
}
