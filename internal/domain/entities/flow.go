package entities

import (
	"fmt"
	"time"
)

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
	ContextID   string    `bson:"idcontext,omitempty" json:"contextId,omitempty"` // ID del contexto asociado al flow
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
	ID           string      `bson:"id" json:"id"`
	From         string      `bson:"from" json:"from"`
	To           string      `bson:"to" json:"to"`
	SourceHandle string      `bson:"sourceHandle" json:"sourceHandle"`               // camelCase para coincidir con MongoDB
	TargetHandle string      `bson:"targetHandle" json:"targetHandle"`               // camelCase para coincidir con MongoDB
	Condition    interface{} `bson:"condition,omitempty" json:"condition,omitempty"` // "yes", "no", "si", "default" - puede ser string o documento
	DelayMs      int         `bson:"delayMs" json:"delayMs"`                         // camelCase para coincidir con MongoDB
}

// GetCondition retorna el condition como string, manejando tanto strings como documentos embebidos
func (e *FlowEdge) GetCondition() string {
	if e.Condition == nil {
		return ""
	}
	
	// Si ya es string, retornarlo
	if str, ok := e.Condition.(string); ok {
		return str
	}
	
	// Si es un documento embebido, intentar extraer un campo común o convertir a string
	if m, ok := e.Condition.(map[string]interface{}); ok {
		// Intentar extraer campos comunes
		if val, ok := m["value"].(string); ok {
			return val
		}
		if val, ok := m["condition"].(string); ok {
			return val
		}
		// Si no hay campo específico, retornar string vacío
		return ""
	}
	
	// Para cualquier otro tipo, convertir a string
	return fmt.Sprintf("%v", e.Condition)
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
		if allEdges[i].From == nodeID && allEdges[i].GetCondition() == condition {
			return &allEdges[i]
		}
	}
	return nil
}
