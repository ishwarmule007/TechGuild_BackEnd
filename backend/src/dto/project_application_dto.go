package dto

// Apply project request and response

type ApplyProjectRequest struct {
	CoverLetter      string `json:"cover_letter" binding:"required"`
	ProposedBudget   float64 `json:"proposed_budget" binding:"required"`
	Currency         string `json:"currency"`
	EstimatedDuration string `json:"estimated_duration"`
	TeamID           string  `json:"team_id,omitempty"`
}

type ApplyProjectResponse struct {
	Message       string `json:"message"`
	ApplicationID string `json:"application_id"`
}

// Withdraw application response

type WithdrawApplicationResponse struct {
	Message string `json:"message"`
}

// Accept application response

type AcceptApplicationResponse struct {
	Message string `json:"message"`
}

// Reject application response

type RejectApplicationResponse struct {
	Message string `json:"message"`
}

// Shortlist application response

type ShortlistApplicationResponse struct {
	Message string `json:"message"`
}

// Project application response

type ProjectApplicationResponse struct {
	ID string `json:"id"`

	ProjectID string `json:"project_id"`

	ApplicantID string `json:"applicant_id,omitempty"`

	TeamID string `json:"team_id,omitempty"`

	CoverLetter string `json:"cover_letter"`

	ProposedBudget float64 `json:"proposed_budget"`

	Currency string `json:"currency"`

	EstimatedDuration string `json:"estimated_duration"`

	Status string `json:"status"`

	ClientMessage string `json:"client_message"`

	AppliedAt string `json:"applied_at"`

	ReviewedAt string `json:"reviewed_at,omitempty"`

	CreatedAt string `json:"created_at"`

	UpdatedAt string `json:"updated_at"`
}

// List of project applications response

type ProjectApplicationListResponse struct {
	Applications []ProjectApplicationResponse `json:"applications"`

	Total int `json:"total"`
}