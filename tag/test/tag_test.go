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

const DEFAULT_PATH string = "http://localhost:8080"
const TAG string = DEFAULT_PATH + "/tag"
const ARTICLE string = DEFAULT_PATH + "/article"
const USER_EMAIL string = "tag_test@test.com"

func TestTagController(t *testing.T) {
	userTags := make([]models.Tag, 0)

	db := test_utils.GetRawDB()

	t.Cleanup(func() {
		db.Exec("DELETE FROM articles;")
		db.Exec("DELETE FROM tags;")

		db.Close()
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

	articles :=make([]models.Article, 0)
	articlesResp := tester.GET(fmt.Sprintf("%s/%s", ARTICLE, "all"))
	json.NewDecoder(articlesResp.Body).Decode(&articles)

	if len(articles) != 1 {
		panic("test setup failed")
	}

	targetArticleId := articles[0].ID

	t.Run("add user's custom tag", func(t *testing.T) {
		tester.POST(TAG, `{"title": "tag 1"}`)
		tester.POST(TAG, `{"title": "tag 2"}`)

		getAllResp := tester.GET(fmt.Sprintf("%s/%s", ARTICLE, "all"))

		defer getAllResp.Body.Close()

		tags := make([]models.Tag, 0)
		json.NewDecoder(getAllResp.Body).Decode(&tags)

		assert.Equal(t, 2, len(tags))

		userTags = tags
	})

	t.Run("attach tag to an article. an article can have multiple tags", func(t *testing.T) {
		tester.PATCH(fmt.Sprintf("%s/%s", TAG, "attach"), fmt.Sprintf(`{"article_id": %d, "tag_id": %d}`, targetArticleId, userTags[0].ID))
		tester.PATCH(fmt.Sprintf("%s/%s", TAG, "attach"), fmt.Sprintf(`{"article_id": %d, "tag_id": %d}`, targetArticleId, userTags[1].ID))

		getArticleResp := tester.GET(fmt.Sprintf("%s/%d", ARTICLE, targetArticleId))

		article := models.Article{}
		json.NewDecoder(getArticleResp.Body).Decode(&article)

		fmt.Println(article.Tags)
		assert.Equal(t, 2, len(article.Tags))
	})

	t.Run("detach a tag from an article", func(t *testing.T) {
		tester.PATCH(fmt.Sprintf("%s/%s", TAG, "detach"), fmt.Sprintf(`{"article_id": %d, "tag_id": %d}`, targetArticleId, userTags[1].ID))
		getArticleResp := tester.GET(fmt.Sprintf("%s/%d", ARTICLE, targetArticleId))

		article := models.Article{}
		json.NewDecoder(getArticleResp.Body).Decode(&article)

		fmt.Println(article.Tags)
		assert.Equal(t, 1, len(article.Tags))
		assert.Equal(t, article.Tags[0].ID, userTags[0].ID)
	})

	t.Run("delete a tag of a user. if tag already attached to an article, article automatically drops the tag", func(t *testing.T) {
		tester.DELETE(fmt.Sprintf("%s/%d", TAG, userTags[0].ID))

		// article에 속한 tag를 삭제했으므로 article을 불러왔을 때 태그가 없어야 한다
		// 물론 모든 태그를 조회했을 때도 존재하지 않아야 한다.
	})
}
