package httpsrv

import (
	"github.com/GeorgeGogos/goaap/internal/pkg/watcher"
)

func (s *Server) addWatcher(w *watcher.Watcher) {
	s.watchersLock.Lock()
	defer s.watchersLock.Unlock()
	s.watchers[w.GetWatcherId()] = w
}

func (s *Server) removeWatcher(w *watcher.Watcher) {
	s.watchersLock.Lock()
	defer s.watchersLock.Unlock()
	// Print satistics before removing watcher.
	for i := range s.sessionStats {
		if s.sessionStats[i].id == w.GetWatcherId() {
			s.sessionStats[i].print()
		}
	}
	// Remove watcher.
	delete(s.watchers, w.GetWatcherId())

	s.statsLock.Lock()
	delete(s.sessionStats, w.GetWatcherId()) //Problem #2: Free up memory by deleting stats after session ends
	s.statsLock.Unlock()
}

func (s *Server) notifyWatchers(str string) {
	s.watchersLock.RLock()
	defer s.watchersLock.RUnlock()

	// Send message to all watchers and increment stats.
	for id := range s.watchers {
		s.watchers[id].Send(str)
		s.incStats(id)
	}
}
