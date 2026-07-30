package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ApplicationStatus string

const (
	ApplicationPending     ApplicationStatus = "pending"
	ApplicationShortlisted ApplicationStatus = "shortlisted"
	ApplicationAccepted    ApplicationStatus = "accepted"
	ApplicationRejected    ApplicationStatus = "rejected"
	ApplicationWithdrawn   ApplicationStatus = "withdrawn"
)

type ProjectApplication struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	// Project relationship
	ProjectID uuid.UUID `gorm:"type:uuid;not null;index"`
	Project   Project `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`

	// Individual / Agency applicant
	ApplicantID *uuid.UUID `gorm:"type:uuid;index"`
	Applicant   *User      `gorm:"foreignKey:ApplicantID;constraint:OnDelete:CASCADE"`

	// Team applicant
	TeamID *uuid.UUID `gorm:"type:uuid;index"`
	Team   *Team      `gorm:"foreignKey:TeamID;constraint:OnDelete:CASCADE"`

	// Proposal details
	CoverLetter      string  `gorm:"type:text;not null"`
	ProposedBudget   float64 `gorm:"not null"`
	Currency         string  `gorm:"size:10;default:'INR'"`
	EstimatedDuration string `gorm:"size:100"`

	// Contract relationship
	Contract *ProjectContract `gorm:"foreignKey:ApplicationID"`

	// Application status
	Status ApplicationStatus `gorm:"type:varchar(30);default:'pending'"`

	ClientMessage string `gorm:"type:text"`

	// Dates
	AppliedAt  time.Time
	ReviewedAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *ProjectApplication) BeforeCreate(tx *gorm.DB) error {
	p.ID = uuid.New()

	if p.AppliedAt.IsZero() {
		p.AppliedAt = time.Now()
	}

	return nil
}