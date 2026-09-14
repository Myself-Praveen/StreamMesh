package logger

import (
	"context"
	"testing"
)

func TestInitLogger(t *testing.T) {
	InitLogger("development")
	if Log == nil {
		t.Errorf("expected global logger to be initialized")
	}

	InitLogger("production")
	if Log == nil {
		t.Errorf("expected global logger to be initialized for production")
	}
}

func TestWithContext(t *testing.T) {
	InitLogger("development")
	ctx := context.Background()
	l := WithContext(ctx)
	if l == nil {
		t.Errorf("expected logger from context")
	}
}
