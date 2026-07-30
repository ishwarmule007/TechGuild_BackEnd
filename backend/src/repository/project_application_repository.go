package repository

import (
	"techguild-backend/src/database/postgres"
	"techguild-backend/src/models"

	"github.com/google/uuid"
)

type ProjectApplicationRepository struct {
}

func NewProjectApplicationRepository() *ProjectApplicationRepository {
	return &ProjectApplicationRepository{}
}

// Create application
func (r *ProjectApplicationRepository) Create(
	application *models.ProjectApplication,
) error {
	return postgres.DB.Create(application).Error
}

// Update application
func (r *ProjectApplicationRepository) Update(
	application *models.ProjectApplication,
) error {
	return postgres.DB.Save(application).Error
}

// Delete application
func (r *ProjectApplicationRepository) Delete(
	application *models.ProjectApplication,
) error {
	return postgres.DB.Delete(application).Error
}

// Find by UUID
func (r *ProjectApplicationRepository) FindByUUID(
	applicationID uuid.UUID,
) (*models.ProjectApplication, error) {

	var application models.ProjectApplication

	err := postgres.DB.
		Preload("Project").
		Preload("Applicant").
		Preload("Team").
		First(&application, "id = ?", applicationID).Error

	if err != nil {
		return nil, err
	}

	return &application, nil
}

// Find by application ID string
func (r *ProjectApplicationRepository) FindByID(
	applicationID string,
) (*models.ProjectApplication, error) {

	id, err := uuid.Parse(applicationID)
	if err != nil {
		return nil, err
	}

	return r.FindByUUID(id)
}

// Find project applications by project ID
func (r *ProjectApplicationRepository) FindByProject(
	projectID uuid.UUID,
) ([]models.ProjectApplication, error) {

	var applications []models.ProjectApplication

	err := postgres.DB.
		Where("project_id = ?", projectID).
		Preload("Applicant").
		Preload("Team").
		Order("created_at DESC").
		Find(&applications).Error

	return applications, err
}

// Find applications by applicant
func (r *ProjectApplicationRepository) FindByApplicant(
	applicantID uuid.UUID,
) ([]models.ProjectApplication, error) {

	var applications []models.ProjectApplication

	err := postgres.DB.
		Where("applicant_id = ?", applicantID).
		Preload("Project").
		Preload("Team").
		Order("created_at DESC").
		Find(&applications).Error

	return applications, err
}

// Find applications by team
func (r *ProjectApplicationRepository) FindByTeam(
	teamID uuid.UUID,
) ([]models.ProjectApplication, error) {

	var applications []models.ProjectApplication

	err := postgres.DB.
		Where("team_id = ?", teamID).
		Preload("Project").
		Preload("Applicant").
		Order("created_at DESC").
		Find(&applications).Error

	return applications, err
}

// Find existing application by project and applicant
func (r *ProjectApplicationRepository) FindExisting(
	projectID uuid.UUID,
	applicantID uuid.UUID,
) (*models.ProjectApplication, error) {

	var application models.ProjectApplication

	err := postgres.DB.
		Where(
			"project_id = ? AND applicant_id = ?",
			projectID,
			applicantID,
		).
		First(&application).Error

	if err != nil {
		return nil, err
	}

	return &application, nil
}

// Find existing application by project and team
func (r *ProjectApplicationRepository) FindExistingByTeam(
	projectID uuid.UUID,
	teamID uuid.UUID,
) (*models.ProjectApplication, error) {

	var application models.ProjectApplication

	err := postgres.DB.
		Where(
			"project_id = ? AND team_id = ?",
			projectID,
			teamID,
		).
		First(&application).Error

	if err != nil {
		return nil, err
	}

	return &application, nil
}

// Update application status
func (r *ProjectApplicationRepository) UpdateStatus(
	applicationID uuid.UUID,
	status models.ApplicationStatus,
) error {

	return postgres.DB.
		Model(&models.ProjectApplication{}).
		Where("id = ?", applicationID).
		Update("status", status).Error
}