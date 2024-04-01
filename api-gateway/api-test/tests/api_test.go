package tests

import (
	"4microservice/api_gateway/api-test/handler"
	"4microservice/api_gateway/api-test/storage"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApi(t *testing.T) {

	gin.SetMode(gin.TestMode)

	// USER
	require.NoError(t, SetupMinimumInstance(""))
	file, err := OpenFile("user.json")

	require.NoError(t, err)
	req := NewRequest(http.MethodPost, "/users/create", file)
	res := httptest.NewRecorder()
	r := gin.Default()

	r.POST("/users/create", handler.CreateUser)
	r.ServeHTTP(res, req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.Code)

	var user *storage.User

	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &user))
	require.Equal(t, user.FirstName, "Nodirbek")
	require.Equal(t, user.LastName, "Nomonov")
	require.Equal(t, user.Email, "nodirbekgolang@gmail.com")
	require.Equal(t, user.Bio, "Rare bio")
	require.Equal(t, user.Password, "Nodirbek2006")
	require.Equal(t, user.Username, "nodirbek")
	require.Equal(t, user.Website, "nodirbek.uz")

	getReq := NewRequest(http.MethodGet, "/users/get", nil)
	args := getReq.URL.Query()
	args.Add("id", user.Id)
	getReq.URL.RawQuery = args.Encode()
	getRes := httptest.NewRecorder()
	r = gin.Default()
	r.GET("/users/get", handler.GetUser)
	r.ServeHTTP(getRes, getReq)
	assert.Equal(t, http.StatusOK, getRes.Code)

	var getUser *storage.User

	bdByte, err := io.ReadAll(getRes.Body)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(bdByte, &getUser))
	assert.Equal(t, user.Id, getUser.Id)
	assert.Equal(t, user.Password, getUser.Password)
	assert.Equal(t, user.Email, getUser.Email)
	assert.Equal(t, user.Bio, getUser.Bio)
	assert.Equal(t, user.Password, getUser.Password)
	assert.Equal(t, user.Username, getUser.Username)
	assert.Equal(t, user.FirstName, getUser.FirstName)
	assert.Equal(t, user.LastName, getUser.LastName)

	delReq := NewRequest(http.MethodDelete, "/users/del?id="+user.Id, file)
	delRes := httptest.NewRecorder()

	r.DELETE("/users/del", handler.DeleteUser)
	r.ServeHTTP(delRes, delReq)
	assert.Equal(t, http.StatusOK, delRes.Code)
	require.NoError(t, err)

	// REGISTER
	fileReg, err := OpenFile("user.json")
	regReq := NewRequest(http.MethodPost, "/users/register", fileReg)
	regRes := httptest.NewRecorder()
	r.POST("/users/register", handler.RegisterUser)
	r.ServeHTTP(regRes, regReq)
	assert.Equal(t, http.StatusOK, regRes.Code)
	var resp storage.ResponseMessage
	bodyBytes, err := io.ReadAll(regRes.Body)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(bodyBytes, &resp))
	require.Equal(t, "OTP sent", resp.Content)
	require.NotNil(t, resp.Content)

	verURLCorrect := "/users/verify/12345"
	verReqCorrect := NewRequest(http.MethodGet, verURLCorrect, file)
	verResCorrect := httptest.NewRecorder()
	r = gin.Default()
	r.GET("/users/verify/:code", handler.Verify)
	r.ServeHTTP(verResCorrect, verReqCorrect)

	// post
	fileProd, err := OpenFile("post.json")
	req = NewRequest(http.MethodPost, "/products/create", fileProd)
	res = httptest.NewRecorder()
	r = gin.Default()
	r.POST("/products/create", handler.CreatePost)
	r.ServeHTTP(res, req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.Code)
	var post storage.Post
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &post))
	require.Equal(t, post.Content, "test content")
	require.Equal(t, post.ImageUrl, "test image")
	require.Equal(t, post.Title, "title")
	require.Equal(t, post.Likes, int64(10))
	require.Equal(t, post.Dislikes, int64(5))
	require.Equal(t, post.Views, int64(1))

	delReq = NewRequest(http.MethodDelete, "/product/delete?id="+post.Id, fileProd)
	delRes = httptest.NewRecorder()
	r.DELETE("/product/delete", handler.DeletePost)
	r.ServeHTTP(delRes, delReq)
	assert.Equal(t, http.StatusOK, delRes.Code)
	var content storage.ResponseMessage
	bodyBytes, _ = io.ReadAll(delRes.Body)
	require.NoError(t, json.Unmarshal(bodyBytes, &content))
}
