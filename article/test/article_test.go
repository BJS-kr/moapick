package article_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"moapick/db/models"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArticleController(t *testing.T) {
	signInResp, err:= http.Post("http://localhost:8080/user/sign-in", "application/json", bytes.NewBuffer([]byte(`{"email": "test@test.com"}`)))
	
	if err != nil {
		t.Error(err.Error())
	}

	defer signInResp.Body.Close()
	
	accessTokenBody := make(map[string]string)

	json.NewDecoder(signInResp.Body).Decode(&accessTokenBody)
	accessToken, ok := accessTokenBody["access_token"]

	if !ok {
		t.Error("failed to get access token")
	}

	var targetArticleId uint

	t.Run("save article", func (t *testing.T) {
		saveReq, _ := http.NewRequest("POST", "http://localhost:8080/article", bytes.NewBuffer([]byte(`{"link":"https://stackoverflow.com", "title":"개발 핵꿀팁 저장소"}`)))
		SetHeaders(saveReq, accessToken)
		
		saveResp, err := http.DefaultClient.Do(saveReq)

		if err != nil {
			t.Error(err.Error())
		}

		defer saveResp.Body.Close()

		PrintReturnBody(saveResp)

		assert.Equal(t, 201, saveResp.StatusCode, "response status code must be 201")
	})
	t.Run("get all articles of a user", func(t *testing.T) {
		getAllReq, _ := http.NewRequest("GET", "http://localhost:8080/article/all", nil)
		SetHeaders(getAllReq, accessToken)

		getAllResp, err := http.DefaultClient.Do(getAllReq)

		if err != nil {
			t.Error(err.Error())
		}

		defer getAllResp.Body.Close()

		getAllRespBody := make([]models.Article, 0)

		json.NewDecoder(getAllResp.Body).Decode(&getAllRespBody)
		targetArticleId = getAllRespBody[0].ID

		PrintReturnBody(getAllResp)

		assert.Equal(t, 200, getAllResp.StatusCode, "response status code must be 200")
	})

	t.Run("get an article by id", func(t *testing.T) {
		getReq, _ := http.NewRequest("GET", fmt.Sprintf("http://localhost:8080/article/%d", targetArticleId), nil)
		SetHeaders(getReq, accessToken)

		getResp, err := http.DefaultClient.Do(getReq)

		if err != nil {
			t.Error(err.Error())
		}

		defer getResp.Body.Close()

		PrintReturnBody(getResp)

		assert.Equal(t, 200, getResp.StatusCode, "response status code must be 200")
	})
}

func SetHeaders(req *http.Request, accessToken string) {
	req.Header.Set("Authorization", "Bearer " + accessToken)
	req.Header.Set("Content-Type", "application/json")
}

func PrintReturnBody(resp *http.Response) {
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}

