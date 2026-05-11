package health

import (
	"context"
	"errors"
	"testing"
)

type mockPinger struct {
	pingErr error
}

func (m *mockPinger) PingContext(ctx context.Context) error {
	return m.pingErr
}

func TestAlive(t *testing.T) {
	s := &healthService{}
	err := s.Alive(context.Background())
	if err != nil {
		t.Fatalf("Alive error: %v", err)
	}
}

func TestReady_NilDB(t *testing.T) {
	s := &healthService{db: nil}
	err := s.Ready(context.Background())
	if err == nil {
		t.Fatal("Ready got nil")
	}
	if err.Error() != "DB not configured!" {
		t.Fatalf("wrong error: %s", err.Error())
	}
}

func TestReady_PingError(t *testing.T) {
	expectedErr := errors.New("ping failed")
	s := &healthService{db: &mockPinger{pingErr: expectedErr}}
	err := s.Ready(context.Background())
	if err == nil {
		t.Fatal("Ready got nil")
	}
	if err != expectedErr {
		t.Fatalf("expected error: %v, got error: %v", expectedErr, err)
	}
}

func TestReady_Success(t *testing.T) {
	s := &healthService{db: &mockPinger{pingErr: nil}}
	err := s.Ready(context.Background())
	if err != nil {
		t.Fatalf("Ready error: %v", err)
	}
}
