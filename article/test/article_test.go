package article_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveArticle(t *testing.T) {
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

	saveReq, _ :=http.NewRequest("POST", "http://localhost:8080/article", bytes.NewBuffer([]byte(`{"link":"https://stackoverflow.com", "title":"개발 핵꿀팁 저장소"}`)))
	saveReq.Header.Set("Authorization", "Bearer " + accessToken)
	saveReq.Header.Set("Content-Type", "application/json")
	saveResp, err := http.DefaultClient.Do(saveReq)

	if err != nil {
		t.Error(err.Error())
	}

	defer saveResp.Body.Close()
	body, _ := io.ReadAll(saveResp.Body)
	strBody := string(body)
	fmt.Println(strBody)

	assert.Equal(t, 201, saveResp.StatusCode, "response status code must be 201")
}