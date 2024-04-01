package handler

import (
	"4microservice/api_gateway/api-test/storage"
	"4microservice/api_gateway/api-test/storage/kv"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

// User

func RegisterUser(c *gin.Context) {
	var newUser storage.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userJson, err := json.Marshal(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	id := uuid.NewString()

	if err := kv.Set(id, string(userJson)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"content": "OTP sent",
	})
}

func Verify(c *gin.Context) {
	userCode := c.Param("code")

	if userCode != "12345" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Incorrect code",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
	})
}

func GetUser(c *gin.Context) {
	id := c.Query("id")

	userGet, err := kv.Get(id)
	fmt.Println(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var res storage.User
	if err := json.Unmarshal([]byte(userGet), &res); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)

}

func CreateUser(c *gin.Context) {
	var newUser storage.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	newUser.Id = uuid.NewString()

	userJson, err := json.Marshal(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := kv.Set(newUser.Id, string(userJson)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, newUser)

}

func DeleteUser(c *gin.Context) {
	id := c.Query("id")

	if err := kv.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user deleted",
	})
}

// Post
func CreatePost(c *gin.Context) {
	var body storage.Post

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	body.Id = uuid.NewString()

	prodJ, err := json.Marshal(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := kv.Set(body.Id, string(prodJ)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, body)
}

func UpdatePost(c *gin.Context) {
	id := c.Query("id")

	var updatedPost storage.Post
	if err := c.ShouldBindJSON(&updatedPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	postJson, err := json.Marshal(updatedPost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := kv.Set(id, string(postJson)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, updatedPost)
}

func DeletePost(c *gin.Context) {
	id := c.Query("id")

	if err := kv.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "post deleted",
	})
}

// Comment
func CreateComment(c *gin.Context) {
	var newComment storage.Comment
	if err := c.ShouldBindJSON(&newComment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	newComment.Id = uuid.NewString()

	commentJson, err := json.Marshal(newComment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := kv.Set(newComment.Id, string(commentJson)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, newComment)
}

func UpdateComment(c *gin.Context) {
	id := c.Query("id")

	var updatedComment storage.Comment
	if err := c.ShouldBindJSON(&updatedComment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	commentJson, err := json.Marshal(updatedComment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := kv.Set(id, string(commentJson)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, updatedComment)
}

func DeleteComment(c *gin.Context) {
	id := c.Query("id")

	if err := kv.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "comment deleted",
	})
}
