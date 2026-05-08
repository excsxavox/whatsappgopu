package entities

import (
	"fmt"
	"time"
)

// Contexto representa el contexto de una empresa con información para la IA
type Contexto struct {
	ID           string                 `json:"_id" bson:"_id"`
	Name         string                 `json:"_name" bson:"_name"`
	Description  string                 `json:"_description" bson:"_description"`
	Prompt       string                 `json:"_prompt" bson:"_prompt"`
	BusinessInfo map[string]interface{} `json:"_businessInfo" bson:"_businessInfo"` // Información de la empresa, productos, direcciones, etc.
	Instructions string                 `json:"_instructions" bson:"_instructions"`
	Tone         string                 `json:"_tone" bson:"_tone"`
	Language     string                 `json:"_language" bson:"_language"`
	ChannelID    *string                `json:"_channelId" bson:"_channelId"`
	CompanyID    string                 `json:"_companyId" bson:"_companyId"`
	IsActive     bool                   `json:"_isActive" bson:"_isActive"`
	CreatedAt    time.Time              `json:"_createdAt" bson:"_createdAt"`
	UpdatedAt    time.Time              `json:"_updatedAt" bson:"_updatedAt"`
	Version      int                    `json:"__v" bson:"__v"`
}

// GetFullPrompt retorna el prompt completo con toda la información del contexto
func (c *Contexto) GetFullPrompt() string {
	prompt := c.Prompt

	if c.Instructions != "" {
		prompt += "\n\nInstrucciones: " + c.Instructions
	}

	if c.BusinessInfo != nil && len(c.BusinessInfo) > 0 {
		// Convertir BusinessInfo a JSON string para incluirlo en el prompt
		businessInfoJSON := ""
		for key, value := range c.BusinessInfo {
			businessInfoJSON += fmt.Sprintf("%s: %v\n", key, value)
		}
		if businessInfoJSON != "" {
			prompt += "\n\nInformación de la empresa:\n" + businessInfoJSON
		}
	}

	if c.Tone != "" {
		prompt += "\n\nTono: " + c.Tone
	}

	return prompt
}
