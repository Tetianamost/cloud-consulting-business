package domain

import (
	"time"
)

// Inquiry represents a service inquiry from a client
type Inquiry struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Company   string    `json:"company"`
	Phone     string    `json:"phone"`
	Services  []string  `json:"services"`
	Message   string    `json:"message"`
	Status    string    `json:"status"`
	Priority  string    `json:"priority"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateInquiryRequest represents the request to create a new inquiry
type CreateInquiryRequest struct {
	Name     string   `json:"name" binding:"required"`
	Email    string   `json:"email" binding:"required,email"`
	Company  string   `json:"company"`
	Phone    string   `json:"phone"`
	Services []string `json:"services" binding:"required"`
	Message  string   `json:"message" binding:"required"`
	Source   string   `json:"source"`
}

