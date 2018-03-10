package models

type Registration struct {
	Identity Identity `json:"identity" binding:"required"`
}
