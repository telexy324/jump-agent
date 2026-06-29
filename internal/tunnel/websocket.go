package tunnel

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"jump-agent/internal/model"

	"nhooyr.io/websocket"
)

const idleTimeout = 10 * time.Minute

type Session struct {
	listener   net.Listener
	done       chan struct{}
	closeOnce  sync.Once
	active     atomic.Int64
	lastActive atomic.Int64
}

func Start(payload *model.SessionPayload) (*Session, *model.SessionPayload, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, nil, err
	}

	t := &Session{
		listener: listener,
		done:     make(chan struct{}),
	}
	t.touch()

	local := *payload
	local.BastionHost = "127.0.0.1"
	local.BastionPort = listener.Addr().(*net.TCPAddr).Port

	go t.acceptLoop(payload)

	return t, &local, nil
}

func (t *Session) Wait() {
	<-t.done
}

func (t *Session) Close() {
	t.closeOnce.Do(func() {
		_ = t.listener.Close()
		close(t.done)
	})
}

func (t *Session) acceptLoop(payload *model.SessionPayload) {
	defer t.Close()

	for {
		if tcpListener, ok := t.listener.(*net.TCPListener); ok {
			_ = tcpListener.SetDeadline(time.Now().Add(time.Second))
		}

		conn, err := t.listener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				if t.active.Load() == 0 && time.Since(time.Unix(t.lastActive.Load(), 0)) > idleTimeout {
					return
				}
				continue
			}
			return
		}

		t.touch()
		t.active.Add(1)
		go func() {
			defer t.active.Add(-1)
			t.proxy(conn, payload)
		}()
	}
}

func (t *Session) proxy(localConn net.Conn, payload *model.SessionPayload) {
	defer localConn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wsConn, err := dialWebSocket(ctx, payload)
	if err != nil {
		log.Printf("websocket dial failed: %v", err)
		return
	}
	defer wsConn.Close(websocket.StatusNormalClosure, "")

	remoteConn := websocket.NetConn(ctx, wsConn, websocket.MessageBinary)
	defer remoteConn.Close()

	var once sync.Once
	closeBoth := func() {
		once.Do(func() {
			_ = localConn.Close()
			_ = remoteConn.Close()
			cancel()
		})
	}

	go func() {
		_, _ = io.Copy(remoteConn, localConn)
		closeBoth()
	}()

	_, _ = io.Copy(localConn, remoteConn)
	closeBoth()
	t.touch()
}

func (t *Session) touch() {
	t.lastActive.Store(time.Now().Unix())
}

func dialWebSocket(ctx context.Context, payload *model.SessionPayload) (*websocket.Conn, error) {
	var lastErr error
	for _, endpoint := range wsEndpoints(payload) {
		wsConn, _, err := websocket.Dial(ctx, endpoint, nil)
		if err == nil {
			log.Printf("websocket tunnel connected: %s", endpoint)
			return wsConn, nil
		}
		lastErr = err
	}
	return nil, lastErr
}

func wsEndpoints(payload *model.SessionPayload) []string {
	if payload.WebSocketURL != "" {
		return []string{payload.WebSocketURL}
	}

	base := strings.TrimSpace(payload.BastionHost)
	if base == "" {
		base = "127.0.0.1"
	}

	if strings.HasPrefix(base, "ws://") || strings.HasPrefix(base, "wss://") {
		return []string{withDefaultPath(base, "/api/jumpServer/ws")}
	}

	scheme := "ws"
	host := base
	if strings.HasPrefix(base, "http://") || strings.HasPrefix(base, "https://") {
		u, err := url.Parse(base)
		if err == nil {
			host = u.Host
			if u.Scheme == "https" {
				scheme = "wss"
			}
		}
	}

	if payload.BastionPort > 0 && !strings.Contains(host, ":") {
		host = net.JoinHostPort(host, strconv.Itoa(payload.BastionPort))
	}
	if payload.BastionPort == 443 {
		scheme = "wss"
	}

	if payload.WebSocketPath != "" {
		return []string{fmt.Sprintf("%s://%s%s", scheme, host, ensureLeadingSlash(payload.WebSocketPath))}
	}

	return []string{
		fmt.Sprintf("%s://%s/api/jumpServer/ws", scheme, host),
		fmt.Sprintf("%s://%s/jumpServer/ws", scheme, host),
	}
}

func withDefaultPath(rawURL, defaultPath string) string {
	u, err := url.Parse(rawURL)
	if err != nil || u.Path != "" {
		return rawURL
	}
	u.Path = defaultPath
	return u.String()
}

func ensureLeadingSlash(path string) string {
	if strings.HasPrefix(path, "/") {
		return path
	}
	return "/" + path
}
