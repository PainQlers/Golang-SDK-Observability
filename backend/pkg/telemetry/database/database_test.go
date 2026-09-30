package database

import (
	"context"
	"errors"
	"testing"
)

func TestObserveExecutesAndPropagatesError(t *testing.T) {
	called := false

	err := Observe(context.Background(), Select, "SELECT 1", func(ctx context.Context) error {
		called = true
		return nil
	})
	if err != nil {
		t.Fatalf("Observe returned unexpected error: %v", err)
	}
	if !called {
		t.Fatal("Observe did not execute the provided function")
	}
}

func TestObserveResultPropagatesError(t *testing.T) {
	expected := errors.New("boom")

	_, err := ObserveResult(context.Background(), Insert, "INSERT INTO users VALUES (1)", func(ctx context.Context) (string, error) {
		return "", expected
	})
	if !errors.Is(err, expected) {
		t.Fatalf("ObserveResult expected wrapped error %v, got %v", expected, err)
	}
}
