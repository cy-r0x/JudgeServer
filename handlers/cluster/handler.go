package cluster

type Node struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels,omitempty"`
}

type Handler struct {
	nodes []Node
}

func NewHandler() *Handler {
	return &Handler{
		nodes: []Node{},
	}
}
