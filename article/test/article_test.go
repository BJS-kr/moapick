package article_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"moapick/db/models"
	"moapick/test_utils"
	"moapick/user"
	"net/http"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/stretchr/testify/assert"
)

const DEFAULT_PATH string = "http://localhost:8080/article"
const USER_EMAIL string = "article_test@test.com"
func TestArticleController(t *testing.T) {
	var targetArticleId uint

	godotenv.Load("../../test.env")

	db := test_utils.GetRawDB()

	t.Cleanup(func ()  {
		_, err := db.Exec("DELETE FROM articles;")
	
		if err!=nil {
			panic(err)
		}

		db.Close()
	})

	signInResp, err := http.Post("http://localhost:8080/user/sign-in", "application/json", bytes.NewBuffer([]byte(fmt.Sprintf(`{"email": "%s"}`, USER_EMAIL))))

	if err != nil {
		t.Error(err.Error())
	}

	defer signInResp.Body.Close()

	accessTokenBody := user.JwtAccessToken{}

	json.NewDecoder(signInResp.Body).Decode(&accessTokenBody)
	
	tester := test_utils.MakeHTTPTester(accessTokenBody.AccessToken)
	t.Run("save article", func(t *testing.T) {
		saveResp := tester.POST(DEFAULT_PATH, `{"link":"https://stackoverflow.com", "title":"개발 핵꿀팁 저장소"}`)

		defer saveResp.Body.Close()

		assert.Equal(t, 201, saveResp.StatusCode, "response status code must be 201")
	})
	t.Run("get all articles of a user", func(t *testing.T) {
		// 1. save more articles
		body1 := `{"link":"https://google.com", "title":"검색은 역시 구글"}`
		body2 := `{"link":"https://naver.com", "title":"한국인이라면 제발 네이버 씁시다"}`
		tester.POST(DEFAULT_PATH, body1)
		tester.POST(DEFAULT_PATH, body2)

		// 첫 번째 테스트에서 1회, 이번 테스트에서 2회 추가로 저장했으니 총 3개의 아티클이 반환되어야 한다.
		getAllResp := tester.GET(DEFAULT_PATH + "/all")

		defer getAllResp.Body.Close()

		getAllRespBody := make([]models.Article, 0)

		json.NewDecoder(getAllResp.Body).Decode(&getAllRespBody)
		targetArticleId = getAllRespBody[0].ID

		assert.Equal(t, 3, len(getAllRespBody), "response body length must be 3")
		assert.Equal(t, 200, getAllResp.StatusCode, "response status code must be 200")
	})

	t.Run("get an article by id", func(t *testing.T) {
		getResp := tester.GET(fmt.Sprintf("%s/%d", DEFAULT_PATH, targetArticleId))

		defer getResp.Body.Close()

		article := models.Article{}

		json.NewDecoder(getResp.Body).Decode(&article)

		assert.Equal(t, targetArticleId, article.ID)
		assert.Equal(t, 200, getResp.StatusCode, "response status code must be 200")
	})

	t.Run("delete article by id", func(t *testing.T) {
		delResp := tester.DELETE(fmt.Sprintf("%s/%d", DEFAULT_PATH, targetArticleId))
		getAllResp := tester.GET(fmt.Sprintf("%s/%s", DEFAULT_PATH, "all"))

		defer getAllResp.Body.Close()

		articles := make([]models.Article, 0)
		json.NewDecoder(getAllResp.Body).Decode(&articles)

		// targetId에 해당하는 아티클이 삭제되었으니 새로운 아이디로 교체
		targetArticleId = articles[0].ID

		assert.Equal(t, 2, len(articles), "all article length must be 2")
		assert.Equal(t, 200, delResp.StatusCode, "delete response code must be 200")
	})

	t.Run("update title of saved article", func(t *testing.T) {
		tester.PATCH(fmt.Sprintf("%s/title/%d", DEFAULT_PATH, targetArticleId), `{"title": "내 맘대로 정해보는 타이틀 후후"}`)
		getResp := tester.GET(fmt.Sprintf("%s/%d", DEFAULT_PATH, targetArticleId))

		defer getResp.Body.Close()

		article := models.Article{}
		json.NewDecoder(getResp.Body).Decode(&article)

		assert.Equal(t, "내 맘대로 정해보는 타이틀 후후", article.Title, "title must be updated as expected")
	})

	t.Run("delete all articles of a user", func(t *testing.T) {
		tester.DELETE(fmt.Sprintf("%s/%s", DEFAULT_PATH, "all"))
		getAllResp := tester.GET(fmt.Sprintf("%s/%s", DEFAULT_PATH, "all"))

		defer getAllResp.Body.Close()

		articles := make([]models.Article, 0)

		json.NewDecoder(getAllResp.Body).Decode(&articles)

		assert.Equal(t, 0, len(articles), "articles length must be zero because all articles of a user have been deleted")
	})
}
