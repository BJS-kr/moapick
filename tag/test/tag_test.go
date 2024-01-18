package tag_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"moapick/db/models"
	"moapick/test_utils"
	"moapick/user"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const DEFAULT_PATH string = "http://localhost:8080/tag"
const USER_EMAIL string = "tag_test@test.com"

func TestTagController(t *testing.T) {
	var targetTagId int

	db := test_utils.GetRawDB()

	defer db.Close()

	t.Cleanup(func() {
		db.Exec()
	})

	signInResp, err := http.Post("http://localhost:8080/user/sign-in", "application/json", bytes.NewBuffer([]byte(fmt.Sprintf(`{"email": "%s"}`, USER_EMAIL))))

	if err != nil {
		t.Error(err.Error())
	}

	defer signInResp.Body.Close()

	accessTokenBody := user.JwtAccessToken{}

	json.NewDecoder(signInResp.Body).Decode(&accessTokenBody)

	accessToken := accessTokenBody.AccessToken
	tester := test_utils.MakeHTTPTester(accessToken)

	articleResp := tester.POST("http://localhost:8080/article", `{"link":"https://medium.com", "title":"개발 아티클이 많아요"}`)

	if articleResp.StatusCode != 201 {
		panic("test setup failed")
	}

	articles :=make( []models.Article, 0)
	articlesResp := tester.GET("http://localhost:8080/article/all")
	json.NewDecoder(articlesResp.Body).Decode(&articles)

	if len(articles) != 1 {
		panic("test setup failed")
	}

	targetArticle := articles[0]

	t.Run("add user's custom tag", func(t *testing.T) {
		tester.POST(DEFAULT_PATH, `{"title": "tag 1"}`)
		tester.POST(DEFAULT_PATH, `{"title": "tag 2"}`)

		getAllResp := tester.GET(fmt.Sprintf("%s/%s", DEFAULT_PATH, "all"))

		defer getAllResp.Body.Close()

		tags := make([]models.Tag, 0)
		json.NewDecoder(getAllResp.Body).Decode(&tags)

		assert.Equal(t, 2, len(tags))
	})

	t.Skip("attach multiple tags to an article", func(t *testing.T) {
		tester.POST(fmt.Sprintf("%s/%s", DEFAULT_PATH, "attach"), fmt.Sprintf(`{"article_id": %d, "tag_id": %d}`))
	})
	t.Skip("detach a tag from an article", func(t *testing.T) {})

}
