package diffcheck_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	diffCheck "github.com/wutzi15/knocken/diffCheck"
)

func TestGetHTML(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
	}))
	defer svr.Close()
	content, err := diffCheck.GetHTML(svr.URL)
	if err != nil {
		t.Error(err)
	}
	if string(content) != "test" {
		t.Errorf("Expected 'test' but got %s", content)
	}
}

func TestWriteFile(t *testing.T) {
	err := diffCheck.WriteFile("test", []byte("test"))
	defer os.Remove("./html/test")
	if err != nil {
		t.Error(err)
	}
}
