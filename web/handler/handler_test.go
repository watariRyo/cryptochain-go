package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/watariRyo/cryptochain-go/web/block"
	mockBlock "github.com/watariRyo/cryptochain-go/web/block/mock"
	"github.com/watariRyo/cryptochain-go/web/redis"
	mockRedis "github.com/watariRyo/cryptochain-go/web/redis/mock"
	"go.uber.org/mock/gomock"
)

func Test_Mine(t *testing.T) {
	ctrl := gomock.NewController(t)
	mb := mockBlock.NewMockBlockChainInterface(ctrl)
	// TODO: Anyを使わない調査
	mb.EXPECT().AddBlock(gomock.Any()).MaxTimes(1)
	mb.EXPECT().GetBlock().Return([]*block.Block{}).MaxTimes(2)

	mr := mockRedis.NewMockRedisClientInterface(ctrl)
	// 非同期の都合、テスト終了時に呼ばれないパターンがあるためAnyTimesをつける
	mr.EXPECT().Publish(gomock.Any(), string(redis.BLOCKCHAIN), gomock.Any()).AnyTimes()

	testHandler := &Handler{
		BlockChain:  mb,
		RedisClient: mr,
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

// TODO: 通常Handlerと同じように(w, r)を受けないと、http.NewRequestをMockできなさそう
func Test_SyncChain(t *testing.T) {
	t.Skip("Not implemented")
}
