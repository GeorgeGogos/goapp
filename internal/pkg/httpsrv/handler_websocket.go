package httpsrv

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/GeorgeGogos/goaap/internal/pkg/watcher"
	"github.com/GeorgeGogos/goaap/pkg/util"
	"github.com/gorilla/websocket"
)

func (s *Server) handlerWebSocket(w http.ResponseWriter, r *http.Request) {
	// Create and start a watcher.
	var watch = watcher.New()
	if err := watch.Start(); err != nil {
		s.error(w, http.StatusInternalServerError, fmt.Errorf("failed to start watcher: %w", err))
		return
	}
	defer watch.Stop()

	s.addWatcher(watch)
	defer s.removeWatcher(watch)

	// Start WS.
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.error(w, http.StatusInternalServerError, fmt.Errorf("failed to upgrade connection: %w", err))
		return
	}
	defer func() { _ = c.Close() }()

	log.Printf("websocket started for watcher %s\n", watch.GetWatcherId())
	defer func() {
		log.Printf("websocket stopped for watcher %s\n", watch.GetWatcherId())

		//Problem #1: Print the total number of messages received for the session
		s.statsLock.RLock()
		if stat, exists := s.sessionStats[watch.GetWatcherId()]; exists {
			stat.print()
		}
		s.statsLock.RUnlock()
	}()

	// Read done.
	readDoneCh := make(chan struct{})

	// All done.
	doneCh := make(chan struct{})
	defer close(doneCh)

	go func() {
		defer close(readDoneCh)
		for {
			select {
			default:
				_, p, err := c.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived) {
						log.Printf("failed to read message: %v\n", err)
					}
					return
				}
				var m watcher.CounterReset
				if err := json.Unmarshal(p, &m); err != nil {
					log.Printf("failed to unmarshal message: %v\n", err)
					continue
				}
				watch.ResetCounter()
			case <-doneCh:
				return
			case <-s.quitChannel:
				return
			}
		}
	}()

	for {
		select {
		case cv := <-watch.Recv():
			//data, _ := json.Marshal(cv)

			//New Feature #2: Generate a new random HEX string value and print it
			randHexValue := util.RandHexString(16)
			response := fmt.Sprintf("{\"iteration\": %d, \"value\": \"%s\"}", cv.Iteration, randHexValue)
			err = c.WriteMessage(websocket.TextMessage, []byte(response))
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("failed to write message: %v\n", err)
				}
				return
			}
		case <-readDoneCh:
			return
		case <-s.quitChannel:
			return
		}
	}
}
