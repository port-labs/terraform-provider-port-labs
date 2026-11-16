package cli

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteBlueprintWaitsForRelationDelete(t *testing.T) {
	const blueprint = `{
		"ok": true,
		"blueprint": {
			"identifier": "burp_site",
			"title": "Burp Site",
			"schema": {"properties": {}},
			"relations": {
				"parent": {"target": "other", "many": false, "required": false}
			}
		}
	}`

	getStarted := make(chan struct{})
	releaseGet := make(chan struct{})
	deleteStarted := make(chan struct{})
	var once sync.Once
	var putFinished atomic.Bool
	var deleteBeforePut atomic.Bool

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			once.Do(func() { close(getStarted) })
			select {
			case <-releaseGet:
			case <-r.Context().Done():
				return
			}
			_, _ = w.Write([]byte(blueprint))
		case http.MethodPut:
			putFinished.Store(true)
			_, _ = w.Write([]byte(blueprint))
		case http.MethodDelete:
			if !putFinished.Load() {
				deleteBeforePut.Store(true)
			}
			select {
			case <-deleteStarted:
			default:
				close(deleteStarted)
			}
			_, _ = w.Write([]byte(`{"ok":true}`))
		default:
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
	}))
	t.Cleanup(server.Close)

	client := &PortClient{Client: resty.New().SetBaseURL(server.URL)}
	ctx := context.Background()
	errCh := make(chan error, 2)

	go func() {
		_, err := client.DeleteBlueprintRelation(ctx, "burp_site", "parent")
		errCh <- err
	}()

	<-getStarted

	deleteEntered := make(chan struct{})
	go func() {
		close(deleteEntered)
		errCh <- client.DeleteBlueprint(ctx, "burp_site")
	}()
	<-deleteEntered

	select {
	case <-deleteStarted:
		t.Fatal("DELETE entered while relation delete held the blueprint lock")
	case <-time.After(50 * time.Millisecond):
	}

	close(releaseGet)
	require.NoError(t, <-errCh)
	require.NoError(t, <-errCh)
	assert.False(t, deleteBeforePut.Load(), "DELETE ran before the relation PUT finished")
}
