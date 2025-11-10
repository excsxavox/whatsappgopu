package entities

import "time"

// FlowData representa la estructura anidada de datos del flow
type FlowData struct {
	EntryNodeID string     `bson:"entryNodeId" json:"entryNodeId"`
	Nodes       []FlowNode `bson:"nodes" json:"nodes"`
	Edges       []FlowEdge `bson:"edges" json:"edges"`
	Version     string     `bson:"version,omitempty" json:"version,omitempty"`
}

// Flow representa un flujo de conversación completo
type Flow struct {
	ID          string    `bson:"_id" json:"id"`
	Name        string    `bson:"_name" json:"name"`
	Description string    `bson:"_description" json:"description"`
	FlowData    FlowData  `bson:"_flowData" json:"flowData"`
	IsActive    bool      `bson:"_isActive" json:"isActive"`
	IsDefault   bool      `bson:"is_default,omitempty" json:"isDefault,omitempty"`
	Status      string    `bson:"_status" json:"status"`
	TenantID    string    `bson:"tenant_id,omitempty" json:"tenantId,omitempty"`
	InstanceID  string    `bson:"instance_id,omitempty" json:"instanceId,omitempty"`
	CreatedAt   time.Time `bson:"_createdAt" json:"createdAt"`
	UpdatedAt   time.Time `bson:"_updatedAt" json:"updatedAt"`
}

// GetEntryNodeID obtiene el ID del nodo de entrada
func (f *Flow) GetEntryNodeID() string {
	return f.FlowData.EntryNodeID
}

// GetNodes obtiene los nodos del flow
func (f *Flow) GetNodes() []FlowNode {
	return f.FlowData.Nodes
}

// GetEdges obtiene los edges del flow
func (f *Flow) GetEdges() []FlowEdge {
	return f.FlowData.Edges
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
	nodes := f.GetNodes()
	for i := range nodes {
		if nodes[i].ID == nodeID {
			return &nodes[i]
		}
	}
	return nil
}

// GetOutgoingEdges obtiene los edges que salen de un nodo
func (f *Flow) GetOutgoingEdges(nodeID string) []FlowEdge {
	edges := []FlowEdge{}
	allEdges := f.GetEdges()
	for _, edge := range allEdges {
		if edge.From == nodeID {
			edges = append(edges, edge)
		}
	}
	return edges
}

// GetEdgeByCondition busca un edge por condición
func (f *Flow) GetEdgeByCondition(nodeID string, condition string) *FlowEdge {
	allEdges := f.GetEdges()
	for i := range allEdges {
		if allEdges[i].From == nodeID && allEdges[i].Condition == condition {
			return &allEdges[i]
		}
	}
	return nil
}


