package models

import "time"

// ModelStats holds various stats for each model in the database and can be displayed in the models dashboard
type ModelStats struct {
	ID         int       `schema:"id"`
	CreatedAt  time.Time `schema:"created"`
	UpdatedAt  time.Time `schema:"updated"`
	Model      string    `schema:"model"`      // Name of model
	Total      uint      `schema:"total"`      // Total is the total number of models in existence. Avaiable in template as {{.models_total}}
	References uint      `schema:"referenced"` // References is how many times per day the model is referenced.Avaiable in template as {{.model_references}}
	PctActive  uint      `schema:"pct_active"` // Is the rounded, total percent of models that are active. Avaiable in template as {{.model_pct_active}}
	Unused     uint      `schema:"unused"`     // Unused is the total number of models that are not being referenced. Avaiable in template as {{.model_unused}}
	Active     uint      `schema:"active"`     // Active is the total number of models that are being referenced. Avaiable in template as {{.model_active}}
	Archived   uint      `schema:"archived"`   // Archived is the total number of models that are currently archived. Avaiable in template as {{.model_archived}}
}

// CRUD is called by each models CRUD handler immediately after a sucessful SQL transaction
// TODO: Find a better way of doing this that is more efficient and does not involve programmer action.
func (m *ModelStats) NewModelStats() {

}

func (m *ModelStats) CRUD() {

}
