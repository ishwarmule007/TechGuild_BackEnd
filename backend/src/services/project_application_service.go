package services

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"techguild-backend/src/dto"
	"techguild-backend/src/models"
	"techguild-backend/src/repository"
)

type ProjectApplicationService struct {
	applicationRepo *repository.ProjectApplicationRepository
	projectRepo     *repository.ProjectRepository
	teamRepo        *repository.TeamRepository
}

func NewProjectApplicationService() *ProjectApplicationService {
	return &ProjectApplicationService{
		applicationRepo: repository.NewProjectApplicationRepository(),
		projectRepo:     repository.NewProjectRepository(),
		teamRepo:        repository.NewTeamRepository(),
	}
}

// Apply for a project
func (s *ProjectApplicationService) ApplyProject(
	applicantID string,
	projectID string,
	req dto.ApplyProjectRequest,
) (*dto.ApplyProjectResponse, error) {

	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		return nil, errors.New("invalid project id")
	}

	project, err := s.projectRepo.FindByUUID(projectUUID)
	if err != nil {
		return nil, errors.New("project not found")
	}

	if project.Status != models.ProjectPublished {
		return nil, errors.New("project is not accepting applications")
	}

	application := models.ProjectApplication{
		ProjectID:       projectUUID,
		CoverLetter:     req.CoverLetter,
		ProposedBudget:  req.ProposedBudget,
		Currency:        req.Currency,
		EstimatedDuration: req.EstimatedDuration,
		Status:          models.ApplicationPending,
		AppliedAt:       time.Now(),
	}

	// Team application
	if req.TeamID != "" {

		teamUUID, err := uuid.Parse(req.TeamID)
		if err != nil {
			return nil, errors.New("invalid team id")
		}

		team, err := s.teamRepo.GetByID(teamUUID)
		if err != nil {
			return nil, errors.New("team not found")
		}

		if team.Status != models.TeamActive {
			return nil, errors.New("team is not active")
		}

		userUUID, err := uuid.Parse(applicantID)
		if err != nil {
			return nil, errors.New("invalid applicant id")
		}

		// User must be a member of the team
		isMember := false

		for _, member := range team.Members {
			if member.UserID == userUUID &&
				member.Status == models.MemberActive {
				isMember = true
				break
			}
		}

		if !isMember {
			return nil, errors.New("you are not an active member of this team")
		}

		_, err = s.applicationRepo.FindExistingByTeam(
			projectUUID,
			teamUUID,
		)

		if err == nil {
			return nil, errors.New("this team has already applied for this project")
		}

		application.TeamID = &teamUUID

	} else {

		// Freelancer / Agency application
		applicantUUID, err := uuid.Parse(applicantID)
		if err != nil {
			return nil, errors.New("invalid applicant id")
		}

		_, err = s.applicationRepo.FindExisting(
			projectUUID,
			applicantUUID,
		)

		if err == nil {
			return nil, errors.New("you have already applied for this project")
		}

		application.ApplicantID = &applicantUUID
	}

	if err := s.applicationRepo.Create(&application); err != nil {
		return nil, err
	}

	return &dto.ApplyProjectResponse{
		Message:       "Application submitted successfully",
		ApplicationID: application.ID.String(),
	}, nil
}

// Withdraw application
func (s *ProjectApplicationService) WithdrawApplication(
	applicantID string,
	applicationID string,
) error {

	application, err := s.applicationRepo.FindByID(applicationID)
	if err != nil {
		return errors.New("application not found")
	}

	userUUID, err := uuid.Parse(applicantID)
	if err != nil {
		return errors.New("invalid applicant id")
	}

	authorized := false

	if application.ApplicantID != nil &&
		*application.ApplicantID == userUUID {
		authorized = true
	}

	if !authorized && application.TeamID != nil {

		team, err := s.teamRepo.GetByID(*application.TeamID)
		if err != nil {
			return errors.New("team not found")
		}

		if team.LeaderID == userUUID {
			authorized = true
		}
	}

	if !authorized {
		return errors.New("unauthorized")
	}

	if application.Status != models.ApplicationPending {
		return errors.New("only pending applications can be withdrawn")
	}

	application.Status = models.ApplicationWithdrawn

	return s.applicationRepo.Update(application)
}

// Accept application
func (s *ProjectApplicationService) AcceptApplication(
	applicationID string,
) error {

	application, err := s.applicationRepo.FindByID(applicationID)
	if err != nil {
		return errors.New("application not found")
	}

	if application.Status != models.ApplicationPending &&
		application.Status != models.ApplicationShortlisted {
		return errors.New("application cannot be accepted")
	}

	now := time.Now()

	application.Status = models.ApplicationAccepted
	application.ReviewedAt = &now

	return s.applicationRepo.Update(application)
}

// Reject application
func (s *ProjectApplicationService) RejectApplication(
	applicationID string,
) error {

	application, err := s.applicationRepo.FindByID(applicationID)
	if err != nil {
		return errors.New("application not found")
	}

	if application.Status == models.ApplicationAccepted {
		return errors.New("accepted application cannot be rejected")
	}

	now := time.Now()

	application.Status = models.ApplicationRejected
	application.ReviewedAt = &now

	return s.applicationRepo.Update(application)
}

// Shortlist application
func (s *ProjectApplicationService) ShortlistApplication(
	applicationID string,
) error {

	application, err := s.applicationRepo.FindByID(applicationID)
	if err != nil {
		return errors.New("application not found")
	}

	if application.Status != models.ApplicationPending {
		return errors.New("only pending applications can be shortlisted")
	}

	now := time.Now()

	application.Status = models.ApplicationShortlisted
	application.ReviewedAt = &now

	return s.applicationRepo.Update(application)
}

// Get my applications
func (s *ProjectApplicationService) GetMyApplications(
	applicantID string,
) (*dto.ProjectApplicationListResponse, error) {

	applicantUUID, err := uuid.Parse(applicantID)
	if err != nil {
		return nil, errors.New("invalid applicant id")
	}

	applications, err := s.applicationRepo.FindByApplicant(applicantUUID)
	if err != nil {
		return nil, err
	}

	response := dto.ProjectApplicationListResponse{}

	for _, application := range applications {
		response.Applications = append(
			response.Applications,
			s.convertToApplicationResponse(&application),
		)
	}

	response.Total = len(response.Applications)

	return &response, nil
}

// Get project applications
func (s *ProjectApplicationService) GetProjectApplications(
	projectID string,
) (*dto.ProjectApplicationListResponse, error) {

	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		return nil, errors.New("invalid project id")
	}

	applications, err := s.applicationRepo.FindByProject(projectUUID)
	if err != nil {
		return nil, err
	}

	response := dto.ProjectApplicationListResponse{}

	for _, application := range applications {
		response.Applications = append(
			response.Applications,
			s.convertToApplicationResponse(&application),
		)
	}

	response.Total = len(response.Applications)

	return &response, nil
}

// Convert application model to response DTO
func (s *ProjectApplicationService) convertToApplicationResponse(
	application *models.ProjectApplication,
) dto.ProjectApplicationResponse {

	response := dto.ProjectApplicationResponse{
		ID:                application.ID.String(),
		ProjectID:         application.ProjectID.String(),
		CoverLetter:       application.CoverLetter,
		ProposedBudget:    application.ProposedBudget,
		Currency:          application.Currency,
		EstimatedDuration: application.EstimatedDuration,
		Status:            string(application.Status),
		ClientMessage:     application.ClientMessage,
		AppliedAt:         application.AppliedAt.Format(time.RFC3339),
		CreatedAt:         application.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         application.UpdatedAt.Format(time.RFC3339),
	}

	if application.ApplicantID != nil {
		response.ApplicantID = application.ApplicantID.String()
	}

	if application.TeamID != nil {
		response.TeamID = application.TeamID.String()
	}

	if application.ReviewedAt != nil {
		response.ReviewedAt = application.ReviewedAt.Format(time.RFC3339)
	}

	return response
}