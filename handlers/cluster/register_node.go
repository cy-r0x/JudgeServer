package cluster

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) registerNode(w http.ResponseWriter, r *http.Request) {
	// Implementation for registering a node in the cluster
	decoder := json.NewDecoder(r.Body)
	var node Node
	if err := decoder.Decode(&node); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	node.Labels = make(map[string]string)
	node.Labels["node"] = "node_" + strconv.Itoa(len(h.nodes)+1)
	h.nodes = append(h.nodes, node)
	utils.SendResponse(w, 200, node)
}
