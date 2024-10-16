package httpsrv

import "log"

type sessionStats struct {
	id   string
	sent int
}

func (w *sessionStats) print() {
	log.Printf("session %s has received %d messages\n", w.id, w.sent)
}

func (w *sessionStats) inc() {
	w.sent++
}

// func (s *Server) incStats(id string) {
// 	// Find and increment.
// 	for _, ws := range s.sessionStats {
// 		if ws.id == id {
// 			ws.inc()
// 			return
// 		}
// 	}
// 	// Not found, add new.
// 	s.sessionStats = append(s.sessionStats, sessionStats{id: id, sent: 1})
// }

// Problem #1
func (s *Server) incStats(id string) {
	s.statsLock.Lock()
	defer s.statsLock.Unlock()

	// Find the session by ID and increment.
	if stat, exists := s.sessionStats[id]; exists {
		stat.inc()
	} else {
		// Not found, create a new sessionStats entry.
		s.sessionStats[id] = &sessionStats{id: id, sent: 1}
	}
}
