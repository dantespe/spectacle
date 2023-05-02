package handler_test

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "io"

    "github.com/go-test/deep"
    "github.com/dantespe/spectacle/handler"
    "github.com/gin-gonic/gin"
)

type JSONResult map[string]interface{}

func setupRouter() *gin.Engine {
    h, _ := handler.NewRestHandler()
    r := gin.Default()
    for k, v := range h.GetRoutes() {
        r.GET(k, v)
    }
    for k, v := range h.PostRoutes() {
        r.POST(k, v)
    }
    return r
}

func TestStatus(t *testing.T) {
    router := setupRouter()
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/status", nil)
    router.ServeHTTP(w, req)
    
    if w.Code != http.StatusOK  {
        t.Errorf("Got code: %d, want %d", w.Code, http.StatusOK)
    }

    bytes, err := io.ReadAll(w.Result().Body)
    if err != nil {
        t.Fatalf("failed to read Result().Body")
    }
    expectedResult := JSONResult {
        "num_records": float64(0),
        "status": "HEALTHY",
    }
    var i JSONResult
    json.Unmarshal(bytes, &i)
    if diff := deep.Equal(i, expectedResult); diff != nil {
        t.Errorf("Body failed with diff: %s", diff)
    }
}