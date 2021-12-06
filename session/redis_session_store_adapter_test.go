package session_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/sangianpatrick/devoria-article-service/session"

	redismock "github.com/go-redis/redismock/v8"
)

func TestRedisSessionStoreAdapter_Get_Success(t *testing.T) {
	rdb, mock := redismock.NewClientMock()
	mock.ExpectGet("test").SetVal("test data")

	sess := session.NewRedisSessionStoreAdapter(rdb, time.Second*5)
	data, err := sess.Get(context.TODO(), "test")

	assert.NoError(t, err)
	assert.Equal(t, "test data", string(data))

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestRedisSessionStoreAdapter_Get_ErrorSessionNotFound(t *testing.T) {
	rdb, mock := redismock.NewClientMock()
	mock.ExpectGet("test").RedisNil()

	sess := session.NewRedisSessionStoreAdapter(rdb, time.Second*5)
	_, err := sess.Get(context.TODO(), "test")

	assert.Error(t, err)
	assert.Equal(t, session.ErrSessionNotFound, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestRedisSessionStoreAdapter_Get_ErrorUnexpected(t *testing.T) {
	rdb, mock := redismock.NewClientMock()
	mock.ExpectGet("test").SetErr(fmt.Errorf("unexpected"))

	sess := session.NewRedisSessionStoreAdapter(rdb, time.Second*5)
	_, err := sess.Get(context.TODO(), "test")

	assert.Error(t, err)
	assert.Equal(t, session.ErrUnexpected, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestRedisSessionStoreAdapter_Set_Error(t *testing.T) {
	rdb, mock := redismock.NewClientMock()
	mock.ExpectSetEX("test", []byte("test data"), time.Second*1).SetErr(fmt.Errorf("unexpected"))

	sess := session.NewRedisSessionStoreAdapter(rdb, time.Second*1)
	err := sess.Set(context.TODO(), "test", []byte("test data"))

	assert.Error(t, err)
	assert.Equal(t, session.ErrUnexpected, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestRedisSessionStoreAdapter_Set_Success(t *testing.T) {
	rdb, mock := redismock.NewClientMock()
	mock.ExpectSetEX("test", []byte("test data"), time.Second*1).SetVal("1")

	sess := session.NewRedisSessionStoreAdapter(rdb, time.Second*1)
	err := sess.Set(context.TODO(), "test", []byte("test data"))

	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestRedisSessionAdapter_Update_Success(t *testing.T) {
	watchTxIDMock := "watch:transaction:session:testtx"
	value := []byte("testvalue")
	rdb, mock := redismock.NewClientMock()
	mock.ExpectWatch(watchTxIDMock).SetErr(nil)
	mock.ExpectTTL("testtx").SetVal(time.Second * 3600)
	mock.ExpectTxPipeline()

	setEXResultMock := mock.ExpectSetEX("testtx", value, time.Second*3600)
	setEXResultMock.SetErr(nil)
	setEXResultMock.SetVal("OK")

	mock.ExpectTxPipelineExec().SetErr(nil)

	sess := session.NewRedisSessionStoreAdapter(rdb, time.Second*1)
	err := sess.Update(context.TODO(), "testtx", value)

	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
