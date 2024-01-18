package tag_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

const DEFAULT_PATH string = "http://localhost:8080/tag"
const USER_EMAIL string = "test@test.com"

func TestTagController(t *testing.T) {
	signInResp, err := http.Post("http://localhost:8080/user/sign-in", "application/json", bytes.NewBuffer([]byte(fmt.Sprintf(`{"email": "%s"}`, USER_EMAIL))))

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
	t.Run("add multiple user custom tag", func(t *testing.T) {

	})
	t.Run("attach multiple tags to an article", func(t *testing.T) {})
	t.Run("detach a tag from an article", func(t *testing.T) {})

}
