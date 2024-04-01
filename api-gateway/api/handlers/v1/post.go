package v1

import (
	"4microservice/api_gateway/api/handlers/models"
	pbp "4microservice/api_gateway/genproto/post_service"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"time"
)

// CreatePost godoc
// @Summary Create a new post
// @Description Creates a new post
// @Tags posts
// @Accept json
// @Produce json
// @Param post body models.PostSwag true "Post information"
// @Success 200 {object} models.Post "Created post"
// @Failure 400 {object} models.ResponseError "Invalid post data"
// @Failure 500 {object} models.ResponseError "Internal server error"
// @Router /v1/post/create [post]
func (h *HandlerV1) CreatePost(c *gin.Context) {
	var post pbp.Post

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Error: err.Error(),
			Code:  http.StatusBadRequest,
		})
		h.log.Error(err.Error())
		return
	}
	post.Id = uuid.NewString()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	createdPost, err := h.serviceManager.PostService().CreatePost(ctx, &post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})
		h.log.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, createdPost)
}

// UpdatePost godoc
// @Summary Update post information
// @Description Updates information of a specific post
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Param post body models.PostSwag true "Updated post information"
// @Success 200 {object} models.Post "Updated post information"
// @Failure 400 {object} models.ResponseError "Invalid post data"
// @Failure 404 {object} models.ResponseError "Post not found"
// @Failure 500 {object} models.ResponseError "Internal server error"
// @Router /v1/post/update/{id} [put]
func (h *HandlerV1) UpdatePost(c *gin.Context) {
	postID := c.Param("id")

	if postID == "" {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Error: "Invalid post ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	var updatedPost pbp.Post
	if err := c.ShouldBindJSON(&updatedPost); err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Error: err.Error(),
			Code:  http.StatusBadRequest,
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	updatedPost.Id = postID // Set post ID
	updatedPostResp, err := h.serviceManager.PostService().UpdatePost(ctx, &updatedPost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, updatedPostResp)
}

// GetPost godoc
// @Summary Get post information
// @Description Retrieves details of a specific post
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} models.Post "Retrieved post information"
// @Failure 400 {object} models.ResponseError "Invalid post ID"
// @Failure 404 {object} models.ResponseError "Post not found"
// @Failure 500 {object} models.ResponseError "Internal server error"
// @Router /v1/post/get/{id} [get]
func (h *HandlerV1) GetPost(c *gin.Context) {
	postID := c.Param("id")

	if postID == "" {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Error: "Invalid post ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	post, err := h.serviceManager.PostService().GetPost(ctx, &pbp.PostID{Id: postID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, post)
}

// GetAllPosts godoc
// @Summary Get all posts
// @Description Retrieves details of all posts
// @Tags posts
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} models.UserWithPosts "List of posts"
// @Failure 400 {object} models.ResponseError "Invalid request parameters"
// @Failure 500 {object} models.ResponseError "Internal server error"
// @Router /v1/posts/{page}/{limit} [get]
func (h *HandlerV1) GetAllPosts(c *gin.Context) {
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

	posts, err := h.serviceManager.PostService().GetAllPosts(ctx, &pbp.PostsRequest{
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

	c.JSON(http.StatusOK, posts)
}

// DeletePost godoc
// @Summary Delete a post
// @Description Deletes a post from the system
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Success 204 "Post deleted successfully"
// @Failure 400 {object} models.ResponseError "Invalid post ID"
// @Failure 404 {object} models.ResponseError "Post not found"
// @Failure 500 {object} models.ResponseError "Internal server error"
// @Router /v1/post/delete/{id} [delete]
func (h *HandlerV1) DeletePost(c *gin.Context) {
	postID := c.Param("id")

	if postID == "" {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Error: "Invalid post ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	_, err := h.serviceManager.PostService().DeletePost(ctx, &pbp.PostID{Id: postID})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			c.JSON(http.StatusInternalServerError, models.ResponseError{
				Error: "Request timed out",
				Code:  http.StatusInternalServerError,
			})
			return
		} else if err.Error() == "rpc error: code = NotFound desc = post not found" {
			c.JSON(http.StatusNotFound, models.ResponseError{
				Error: "Post not found",
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
