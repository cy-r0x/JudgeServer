package cluster

type Node struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels,omitempty"`
}

type Handler struct {
	nodes          []Node
	availableNodes map[string]string
}

func NewHandler() *Handler {
	return &Handler{
		nodes:          []Node{},
		availableNodes: make(map[string]string),
	}
}
