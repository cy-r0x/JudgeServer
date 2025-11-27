package standings

import "time"

func (h *Handler) MemoryEviction() {
	h.mu.Lock()
	defer h.mu.Unlock()

	for contestId, entry := range h.Last_standings {
		if entry.timestamp != nil {
			if entry.timestamp.Add(1 * time.Hour).Before(time.Now()) {
				delete(h.Last_standings, contestId)
			}
		}
	}
}
