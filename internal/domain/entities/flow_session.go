package entities

import "time"

// FlowSession representa el estado actual de una conversación dentro de un flujo
type FlowSession struct {
	ID                  string                 `bson:"_id" json:"id"`
	ConversationID      string                 `bson:"conversation_id" json:"conversation_id"`
	FlowID              string                 `bson:"flow_id" json:"flow_id"`
	CurrentNodeID       string                 `bson:"current_node_id" json:"current_node_id"`
	Variables           map[string]interface{} `bson:"variables" json:"variables"`
	WaitingForResponse  bool                   `bson:"waiting_for_response" json:"waiting_for_response"`
	WaitingForVariable  string                 `bson:"waiting_for_variable" json:"waiting_for_variable"`
	ExecutedNodes       []string               `bson:"executed_nodes" json:"executed_nodes"`
	Status              string                 `bson:"status" json:"status"` // active, completed, abandoned, error
	CreatedAt           time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt           time.Time              `bson:"updated_at" json:"updated_at"`
	LastActivityAt      time.Time              `bson:"last_activity_at" json:"last_activity_at"`
	TenantID            string                 `bson:"tenant_id" json:"tenant_id"`
	InstanceID          string                 `bson:"instance_id" json:"instance_id"`
}

// NewFlowSession crea una nueva sesión de flujo
func NewFlowSession(conversationID, flowID, entryNodeID, tenantID, instanceID string) *FlowSession {
	now := time.Now()
	return &FlowSession{
		ConversationID:     conversationID,
		FlowID:             flowID,
		CurrentNodeID:      entryNodeID,
		Variables:          make(map[string]interface{}),
		WaitingForResponse: false,
		WaitingForVariable: "",
		ExecutedNodes:      []string{},
		Status:             "active",
		CreatedAt:          now,
		UpdatedAt:          now,
		LastActivityAt:     now,
		TenantID:           tenantID,
		InstanceID:         instanceID,
	}
}

// SetVariable guarda una variable en la sesión
func (fs *FlowSession) SetVariable(name string, value interface{}) {
	fs.Variables[name] = value
	fs.UpdatedAt = time.Now()
}

// GetVariable obtiene una variable de la sesión
func (fs *FlowSession) GetVariable(name string) (interface{}, bool) {
	value, exists := fs.Variables[name]
	return value, exists
}

// MoveToNode actualiza el nodo actual
func (fs *FlowSession) MoveToNode(nodeID string) {
	fs.CurrentNodeID = nodeID
	fs.ExecutedNodes = append(fs.ExecutedNodes, nodeID)
	fs.UpdatedAt = time.Now()
	fs.LastActivityAt = time.Now()
}

// SetWaitingForResponse marca la sesión como esperando respuesta
func (fs *FlowSession) SetWaitingForResponse(variableName string) {
	fs.WaitingForResponse = true
	fs.WaitingForVariable = variableName
	fs.UpdatedAt = time.Now()
}

// ClearWaitingForResponse limpia el estado de espera
func (fs *FlowSession) ClearWaitingForResponse() {
	fs.WaitingForResponse = false
	fs.WaitingForVariable = ""
	fs.UpdatedAt = time.Now()
}

// Complete marca la sesión como completada
func (fs *FlowSession) Complete() {
	fs.Status = "completed"
	fs.UpdatedAt = time.Now()
}

// Abandon marca la sesión como abandonada
func (fs *FlowSession) Abandon() {
	fs.Status = "abandoned"
	fs.UpdatedAt = time.Now()
}

// MarkError marca la sesión con error
func (fs *FlowSession) MarkError() {
	fs.Status = "error"
	fs.UpdatedAt = time.Now()
}

// UpdateActivity actualiza la última actividad
func (fs *FlowSession) UpdateActivity() {
	fs.LastActivityAt = time.Now()
	fs.UpdatedAt = time.Now()
}


