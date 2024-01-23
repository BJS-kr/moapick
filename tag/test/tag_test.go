package tag_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"moapick/db/models"
	"moapick/test_utils"
	"moapick/user"
	"net/http"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

const DEFAULT_PATH string = "http://localhost:8080"
const TAG string = DEFAULT_PATH + "/tag"
const ARTICLE string = DEFAULT_PATH + "/article"
const USER_EMAIL string = "tag_test@test.com"

func TestTagController(t *testing.T) {
	var targetTagId uint

	godotenv.Load("../../test.env")
	userTags := make([]models.Tag, 0)

	db := test_utils.GetRawDB()

	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM article_tags;")
		if err != nil {
			log.Println(err.Error())
		}
		_, err2 := db.Exec("DELETE FROM articles;")
		if err2 != nil {
			log.Println(err.Error())
		}
		_, err3 := db.Exec("DELETE FROM tags;")
		if err3 != nil {
			log.Println(err2.Error())
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

	accessToken := accessTokenBody.AccessToken
	tester := test_utils.Tester{AccessToken: accessToken}

	articleResp1 := tester.POST(ARTICLE, `{"link":"https://medium.com", "title":"개발 아티클이 많아요"}`)
	articleResp2 := tester.POST(ARTICLE, `{"link":"https://google.com", "title":"검색하기 좋아요"}`)

	if articleResp1.StatusCode != 201 || articleResp2.StatusCode != 201 {
		panic("test setup failed")
	}

	articles := make([]models.Article, 0)
	articlesResp := tester.GET(fmt.Sprintf("%s/%s", ARTICLE, "all"))
	json.NewDecoder(articlesResp.Body).Decode(&articles)

	if len(articles) != 2 {
		panic("test setup failed")
	}

	targetArticle1Id := articles[0].ID
	targetArticle2Id := articles[1].ID

	t.Run("add user's custom tag", func(t *testing.T) {
		tester.POST(TAG, `{"title": "tag 1"}`)
		tester.POST(TAG, `{"title": "tag 2"}`)

		getAllResp := tester.GET(fmt.Sprintf("%s/%s", TAG, "all"))

		defer getAllResp.Body.Close()

		tags := make([]models.Tag, 0)
		json.NewDecoder(getAllResp.Body).Decode(&tags)

		assert.Equal(t, 2, len(tags))

		userTags = tags
	})

	t.Run("attach tag to an article. an article can have multiple tags", func(t *testing.T) {
		tester.PATCH(fmt.Sprintf("%s/%s", TAG, "attach"), fmt.Sprintf(`{"article_id": %d, "tag_id": %d}`, targetArticle1Id, userTags[0].ID))
		tester.PATCH(fmt.Sprintf("%s/%s", TAG, "attach"), fmt.Sprintf(`{"article_id": %d, "tag_id": %d}`, targetArticle1Id, userTags[1].ID))

		getArticleResp := tester.GET(fmt.Sprintf("%s/%d", ARTICLE, targetArticle1Id))

		article := make(map[string]interface{})
		json.NewDecoder(getArticleResp.Body).Decode(&article)

		articleTags := article["tags"].([]interface{})

		assert.Equal(t, 2, len(articleTags))
	})

	t.Run("detach a tag from an article", func(t *testing.T) {
		tester.PATCH(fmt.Sprintf("%s/%s", TAG, "detach"), fmt.Sprintf(`{"article_id": %d, "tag_id": %d}`, targetArticle1Id, userTags[1].ID))
		getArticleResp := tester.GET(fmt.Sprintf("%s/%d", ARTICLE, targetArticle1Id))

		article := models.Article{}
		json.NewDecoder(getArticleResp.Body).Decode(&article)

		assert.Equal(t, 1, len(article.Tags))
		assert.Equal(t, article.Tags[0].ID, userTags[0].ID)
	})

	t.Run("delete a tag of a user. if tag already attached to an article, article automatically drops the tag", func(t *testing.T) {
		tester.DELETE(fmt.Sprintf("%s/%d", TAG, userTags[0].ID))

		// article에 속한 tag를 삭제했으므로 article을 불러왔을 때 태그가 없어야 한다
		// 물론 모든 태그를 조회했을 때도 존재하지 않아야 한다.
		getArticleResp := tester.GET(fmt.Sprintf("%s/%d", ARTICLE, targetArticle1Id))
		getAllTagsResp := tester.GET(fmt.Sprintf("%s/%s", TAG, "all"))

		article := models.Article{}
		tags := make([]models.Tag, 0)

		json.NewDecoder(getArticleResp.Body).Decode(&article)
		json.NewDecoder(getAllTagsResp.Body).Decode(&tags)

		// 총 두 개의 user custom tag 중, 아티클에 붙어있던 tag를 삭제했으니
		// article에는 태그가 없어지고
		// tag의 총 갯수는 하나가 된다.
		assert.Equal(t, 0, len(article.Tags))
		assert.Equal(t, 1, len(tags))

		targetTagId = tags[0].ID
	})

	t.Run("get users's articles by tag", func(t *testing.T) {
		tester.PATCH(fmt.Sprintf("%s/%s", TAG, "attach"), fmt.Sprintf(`{"article_id": %d, "tag_id": %d}`, targetArticle1Id, targetTagId))
		tester.PATCH(fmt.Sprintf("%s/%s", TAG, "attach"), fmt.Sprintf(`{"article_id": %d, "tag_id": %d}`, targetArticle2Id, targetTagId))

		getArticlesByTagResp := tester.GET(fmt.Sprintf("%s/%s/%d", TAG, "articles", targetTagId))

		articles := make([]models.Article, 0)
		json.NewDecoder(getArticlesByTagResp.Body).Decode(&articles)

		assert.Equal(t, 2, len(articles))
	})
}
