package entities

import "time"

// Flow representa un flujo de conversación completo
type Flow struct {
	ID          string      `bson:"_id" json:"id"`
	Name        string      `bson:"name" json:"name"`
	Description string      `bson:"description" json:"description"`
	EntryNodeID string      `bson:"entry_node_id" json:"entry_node_id"`
	Nodes       []FlowNode  `bson:"nodes" json:"nodes"`
	Edges       []FlowEdge  `bson:"edges" json:"edges"`
	IsActive    bool        `bson:"is_active" json:"is_active"`
	IsDefault   bool        `bson:"is_default" json:"is_default"`
	TenantID    string      `bson:"tenant_id" json:"tenant_id"`
	InstanceID  string      `bson:"instance_id" json:"instance_id"`
	CreatedAt   time.Time   `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time   `bson:"updated_at" json:"updated_at"`
}

// FlowNode representa un nodo en el flujo
type FlowNode struct {
	ID     string                 `bson:"id" json:"id"`
	Type   string                 `bson:"type" json:"type"` // TEXT, BUTTONS, HTTP, CONDITION, RESPONSE, AUDIO
	Config map[string]interface{} `bson:"config" json:"config"`
}

// FlowEdge representa una conexión entre nodos
type FlowEdge struct {
	ID           string `bson:"id" json:"id"`
	From         string `bson:"from" json:"from"`
	To           string `bson:"to" json:"to"`
	SourceHandle string `bson:"source_handle" json:"source_handle"`
	TargetHandle string `bson:"target_handle" json:"target_handle"`
	Condition    string `bson:"condition,omitempty" json:"condition,omitempty"` // "yes", "no", "si", "default"
	DelayMs      int    `bson:"delay_ms" json:"delay_ms"`
}

// GetNodeByID busca un nodo por ID
func (f *Flow) GetNodeByID(nodeID string) *FlowNode {
	for i := range f.Nodes {
		if f.Nodes[i].ID == nodeID {
			return &f.Nodes[i]
		}
	}
	return nil
}

// GetOutgoingEdges obtiene los edges que salen de un nodo
func (f *Flow) GetOutgoingEdges(nodeID string) []FlowEdge {
	edges := []FlowEdge{}
	for _, edge := range f.Edges {
		if edge.From == nodeID {
			edges = append(edges, edge)
		}
	}
	return edges
}

// GetEdgeByCondition busca un edge por condición
func (f *Flow) GetEdgeByCondition(nodeID string, condition string) *FlowEdge {
	for i := range f.Edges {
		if f.Edges[i].From == nodeID && f.Edges[i].Condition == condition {
			return &f.Edges[i]
		}
	}
	return nil
}


