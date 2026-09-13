// Package wsverify tests the transparent WS routing of go-binance/v2 against
// a local fake Binance API serving both plain HTTP and the WebSocket API.
package wsverify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/gorilla/websocket"
)

type hub struct {
	mu            sync.Mutex
	lastTransport string
	lastMethod    string
	wsCount       int
}

func (h *hub) record(transport, method string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.lastTransport = transport
	h.lastMethod = method
}

func (h *hub) stats() (string, string, int) {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.lastTransport, h.lastMethod, h.wsCount
}

func (h *hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.EqualFold(r.Header.Get("Upgrade"), "websocket") {
		h.serveWS(w, r)
		return
	}

	h.record("http", r.URL.Path)
	switch r.URL.Path {
	case "/api/v3/time":
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"serverTime":%d}`, time.Now().UnixMilli())
	default:
		http.NotFound(w, r)
	}
}

func (h *hub) serveWS(w http.ResponseWriter, r *http.Request) {
	up := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	conn, err := up.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	h.mu.Lock()
	h.wsCount++
	h.mu.Unlock()

	_ = conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		var req struct {
			Id     string         `json:"id"`
			Method string         `json:"method"`
			Params map[string]any `json:"params"`
		}
		if err := json.Unmarshal(msg, &req); err != nil {
			continue
		}

		h.record("ws", req.Method)
		reply := map[string]any{
			"id":     req.Id,
			"status": 200,
		}
		switch req.Method {
		case "ping":
			reply["result"] = map[string]any{}
		case "time":
			reply["result"] = map[string]any{"serverTime": time.Now().UnixMilli()}
		default:
			reply["result"] = map[string]any{}
		}
		if err := conn.WriteJSON(reply); err != nil {
			return
		}
	}
}

func newFixture(t *testing.T) (*hub, *httptest.Server, *binance.Client) {
	t.Helper()
	h := &hub{}
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws-api/v3"
	binance.WsAPIMainURL = wsURL

	client := binance.NewClient("fake-key", "fake-secret")
	client.BaseURL = srv.URL

	// allow the initial client connection to be established
	deadline := time.Now().Add(3 * time.Second)
	for !client.WsConnected() && time.Now().Before(deadline) {
		time.Sleep(20 * time.Millisecond)
	}
	if !client.WsConnected() {
		t.Fatal("ws connection was not established")
	}
	t.Cleanup(client.Close)

	return h, srv, client
}

func TestTransparentWsRouting(t *testing.T) {
	h, _, client := newFixture(t)
	ctx := context.Background()

	t.Run("enabled routes through ws", func(t *testing.T) {
		st, err := client.NewServerTimeService().Do(ctx)
		if err != nil {
			t.Fatal(err)
		}
		tr, m, _ := h.stats()
		if st == 0 || tr != "ws" || m != "time" {
			t.Fatalf("expected ws transport for method time, got transport=%s method=%s st=%d", tr, m, st)
		}
	})

	t.Run("disabled falls back to http", func(t *testing.T) {
		client.DisableWsAPI()
		st, err := client.NewServerTimeService().Do(ctx)
		if err != nil {
			t.Fatal(err)
		}
		tr, m, _ := h.stats()
		if st == 0 || tr != "http" || m != "/api/v3/time" {
			t.Fatalf("expected http fallback, got transport=%s method=%s st=%d", tr, m, st)
		}
	})

	t.Run("re-enabled routes through ws again", func(t *testing.T) {
		client.EnableWsAPI()
		st, err := client.NewServerTimeService().Do(ctx)
		if err != nil {
			t.Fatal(err)
		}
		tr, m, _ := h.stats()
		if st == 0 || tr != "ws" || m != "time" {
			t.Fatalf("expected ws transport again, got transport=%s method=%s st=%d", tr, m, st)
		}
	})
}