package v1

import (
	"4microservice/api_gateway/api/handlers/models"
	pbc "4microservice/api_gateway/genproto/comment_service"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

// CreateComment godoc
// @Summary Create a new comment
// @Description Creates a new comment
// @Tags comments
// @Accept json
// @Produce json
// @Param comment body models.Comment true "Comment information"
// @Success 200 {object} models.Comment "Created comment"
// @Failure 400 {object} models.ResponseError "Invalid comment data"
// @Failure 500 {object} models.ResponseError "Internal server error"
// @Router /v1/comment/create [post]
func (h *HandlerV1) CreateComment(c *gin.Context) {
	var comment pbc.Comment

	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Error: err.Error(),
			Code:  http.StatusBadRequest,
		})
		h.log.Error(err.Error())
		return
	}
	comment.Id = uuid.NewString()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	createdComment, err := h.serviceManager.CommentService().CreateComment(ctx, &comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})
		h.log.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, createdComment)
}

// UpdateComment godoc
// @Summary Update comment information
// @Description Updates information of a specific comment
// @Tags comments
// @Accept json
// @Produce json
// @Param id path string true "Comment ID"
// @Param comment body models.Comment true "Updated comment information"
// @Success 200 {object} models.Comment "Updated comment information"
// @Failure 400 {object} models.ResponseError "Invalid comment data"
// @Failure 404 {object} models.ResponseError "Comment not found"
// @Failure 500 {object} models.ResponseError "Internal server error"
// @Router /v1/comment/update/{id} [put]
func (h *HandlerV1) UpdateComment(c *gin.Context) {
	commentID := c.Param("id")

	if commentID == "" {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Error: "Invalid comment ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	var updatedComment pbc.Comment
	if err := c.ShouldBindJSON(&updatedComment); err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Error: err.Error(),
			Code:  http.StatusBadRequest,
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	updatedComment.Id = commentID // Set comment ID
	updatedCommentResp, err := h.serviceManager.CommentService().UpdateComment(ctx, &updatedComment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, updatedCommentResp)
}

// GetComment godoc
// @Summary Get comment information
// @Description Retrieves details of a specific comment
// @Tags comments
// @Accept json
// @Produce json
// @Param id path string true "Comment ID"
// @Success 200 {object} models.Comment "Retrieved comment information"
// @Failure 400 {object} models.ResponseError "Invalid comment ID"
// @Failure 404 {object} models.ResponseError "Comment not found"
// @Failure 500 {object} models.ResponseError "Internal server error"
// @Router /v1/comment/get/{id} [get]
func (h *HandlerV1) GetComment(c *gin.Context) {
	commentID := c.Param("id")

	if commentID == "" {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Error: "Invalid comment ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	comment, err := h.serviceManager.CommentService().GetComment(ctx, &pbc.CommentId{Id: commentID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// GetAllComments godoc
// @Summary Get all comments
// @Description Retrieves details of all comments
// @Tags comments
// @Accept json
// @Produce json
// @Param post_id query string false "Post ID"
// @Success 200 {object} models.CommentsResponse "List of comments"
// @Failure 400 {object} models.ResponseError "Invalid request parameters"
// @Failure 500 {object} models.ResponseError "Internal server error"
// @Router /v1/comments/{post_id} [get]
func (h *HandlerV1) GetAllComments(c *gin.Context) {
	postID := c.Query("post_id")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	comments, err := h.serviceManager.CommentService().GetAllComments(ctx, &pbc.CommentsRequest{PostId: postID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, comments)
}

// DeleteComment godoc
// @Summary Delete a comment
// @Description Deletes a comment from the system
// @Tags comments
// @Accept json
// @Produce json
// @Param id path string true "Comment ID"
// @Success 204 "Comment deleted successfully"
// @Failure 400 {object} models.ResponseError "Invalid comment ID"
// @Failure 404 {object} models.ResponseError "Comment not found"
// @Failure 500 {object} models.ResponseError "Internal server error"
// @Router /v1/comment/delete/{id} [delete]
func (h *HandlerV1) DeleteComment(c *gin.Context) {
	commentID := c.Param("id")

	if commentID == "" {
		c.JSON(http.StatusBadRequest, models.ResponseError{
			Error: "Invalid comment ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	_, err := h.serviceManager.CommentService().DeleteComment(ctx, &pbc.CommentId{Id: commentID})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			c.JSON(http.StatusInternalServerError, models.ResponseError{
				Error: "Request timed out",
				Code:  http.StatusInternalServerError,
			})
			return
		} else if err.Error() == "rpc error: code = NotFound desc = comment not found" {
			c.JSON(http.StatusNotFound, models.ResponseError{
				Error: "Comment not found",
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
