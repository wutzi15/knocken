package utils_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wutzi15/knocken/utils"
)

func TestGetHTML(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
	}))
	defer svr.Close()
	content, err := utils.GetHTML(svr.URL)
	if err != nil {
		t.Error(err)
	}
	if string(content) != "test" {
		t.Errorf("Expected 'test' but got %s", content)
	}
}
