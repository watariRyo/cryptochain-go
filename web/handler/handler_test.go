package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mockUseCase "github.com/watariRyo/cryptochain-go/web/usecase/mock"
	"go.uber.org/mock/gomock"
)

func Test_Mine(t *testing.T) {
	ctrl := gomock.NewController(t)
	mu := mockUseCase.NewMockUseCaseInterface(ctrl)
	mu.EXPECT().Mine("Test").Times(1)
	mu.EXPECT().GetBlock().Times(1)

	testHandler := &Handler{
		usecase: mu,
	}

	postBody := map[string]interface{}{
		"data": "Test",
	}

	body, _ := json.Marshal(postBody)

	req, _ := http.NewRequest("POST", "/api/mine", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(testHandler.Mine)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected http.StatusOK but got %d", rr.Code)
	}
}