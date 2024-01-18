package article_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"moapick/db/models"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArticleController(t *testing.T) {
	signInResp, err := http.Post("http://localhost:8080/user/sign-in", "application/json", bytes.NewBuffer([]byte(`{"email": "test@test.com"}`)))

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

	t.Run("save article", func(t *testing.T) {
		saveReq, _ := http.NewRequest("POST", "http://localhost:8080/article", bytes.NewBuffer([]byte(`{"link":"https://stackoverflow.com", "title":"개발 핵꿀팁 저장소"}`)))
		SetHeaders(saveReq, accessToken)

		saveResp, err := http.DefaultClient.Do(saveReq)

		if err != nil {
			t.Error(err.Error())
		}

		defer saveResp.Body.Close()

		assert.Equal(t, 201, saveResp.StatusCode, "response status code must be 201")
	})
	t.Run("get all articles of a user", func(t *testing.T) {
		// 1. save more articles
		saveReq_1, _ := http.NewRequest("POST", "http://localhost:8080/article", bytes.NewBuffer([]byte(`{"link":"https://google.com", "title":"검색은 역시 구글"}`)))
		saveReq_2, _ := http.NewRequest("POST", "http://localhost:8080/article", bytes.NewBuffer([]byte(`{"link":"https://naver.com", "title":"한국인이라면 제발 네이버 씁시다"}`)))

		SetHeaders(saveReq_1, accessToken)
		SetHeaders(saveReq_2, accessToken)

		http.DefaultClient.Do(saveReq_1)
		http.DefaultClient.Do(saveReq_2)

		// 첫 번째 테스트에서 1회, 이번 테스트에서 2회 추가로 저장했으니 총 3개의 아티클이 반환되어야 한다.
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

		assert.Equal(t, 3, len(getAllRespBody), "response body length must be 3")
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

		article := models.Article{}

		json.NewDecoder(getResp.Body).Decode(&article)

		assert.Equal(t, targetArticleId, article.ID)
		assert.Equal(t, 200, getResp.StatusCode, "response status code must be 200")
	})

	t.Run("delete article by id", func(t *testing.T) {
		delReq, _ := http.NewRequest("DELETE", fmt.Sprintf("http://localhost:8080/article/%d", targetArticleId), nil)
		SetHeaders(delReq, accessToken)
		delResp, _ := http.DefaultClient.Do(delReq)

		getAllReq, _ := http.NewRequest("GET", "http://localhost:8080/article/all", nil)
		SetHeaders(getAllReq, accessToken)
		getAllResp, _ := http.DefaultClient.Do(getAllReq)

		defer getAllResp.Body.Close()

		articles := make([]models.Article, 0)
		json.NewDecoder(getAllResp.Body).Decode(&articles)

		// targetId에 해당하는 아티클이 삭제되었으니 새로운 아이디로 교체
		targetArticleId = articles[0].ID

		assert.Equal(t, 2, len(articles), "all article length must be 2")
		assert.Equal(t, 200, delResp.StatusCode, "delete response code must be 200")
	})

	t.Run("update title of saved article", func(t *testing.T) {
		updateReq, _ := http.NewRequest("PATCH", fmt.Sprintf("http://localhost:8080/article/title/%d", targetArticleId), bytes.NewBuffer([]byte(`{"title": "내 맘대로 정해보는 타이틀 후후"}`)))
		SetHeaders(updateReq, accessToken)
		http.DefaultClient.Do(updateReq)

		getReq, _ := http.NewRequest("GET", fmt.Sprintf("http://localhost:8080/article/%d", targetArticleId), nil)
		SetHeaders(getReq, accessToken)
		getResp, _ := http.DefaultClient.Do(getReq)

		defer getResp.Body.Close()

		article := models.Article{}
		json.NewDecoder(getResp.Body).Decode(&article)

		assert.Equal(t, article.Title, "내 맘대로 정해보는 타이틀 후후", "title must be updated as expected")
	})

	t.Run("delete all article", func(t *testing.T) {
		deleteAllReq, _ := http.NewRequest("DELETE", "http://localhost:8080/article/all", nil)
		SetHeaders(deleteAllReq, accessToken)
		http.DefaultClient.Do(deleteAllReq)

		getAllReq, _ := http.NewRequest("GET", "http://localhost:8080/article/all", nil)
		SetHeaders(getAllReq, accessToken)
		getAllResp, _ := http.DefaultClient.Do(getAllReq)

		defer getAllResp.Body.Close()

		articles := make([]models.Article, 0)

		json.NewDecoder(getAllResp.Body).Decode(&articles)

		assert.Equal(t, 0, len(articles), "articles length must be zero because all articles have been deleted")
	})
}
