package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"erp-kafeinart/model"

	"github.com/google/uuid"
)

// ── Transaksi Reimbursement ─────────────────────────────────────────

func (s *hrService) GetReimbursementsPaginated(page, perPage int, employeeID, status, search string) ([]model.TrxHrReimbursement, int64, error) {
	return s.repo.FindReimbursementsPaginated(page, perPage, employeeID, status, search)
}

func (s *hrService) GetReimbursementByID(id string) (*model.TrxHrReimbursement, error) {
	return s.repo.FindReimbursementByID(id)
}

func (s *hrService) CreateReimbursement(dto *ReimbursementDTO) (*model.TrxHrReimbursement, error) {
	if dto.EmployeeID == "" {
		return nil, errors.New("karyawan wajib dipilih")
	}

	claimNum := dto.ClaimNumber
	if claimNum == "" {
		claimNum = fmt.Sprintf("CLM-%s-%s", time.Now().Format("200601"), strings.ToUpper(uuid.NewString()[:5]))
	}

	claimDate := dto.ClaimDate
	if claimDate == "" {
		claimDate = time.Now().Format("2006-01-02")
	}

	cat := dto.CategoryID
	if cat == "" {
		cat = dto.Category
	}

	status := dto.Status
	if status == "" {
		status = "SUBMITTED"
	}

	rm := &model.TrxHrReimbursement{
		ID:            uuid.NewString(),
		ClaimNumber:   claimNum,
		EmployeeID:    dto.EmployeeID,
		ClaimDate:     claimDate,
		CategoryID:    cat,
		Amount:        dto.Amount,
		Description:   dto.Description,
		AttachmentURL: dto.AttachmentURL,
		Status:        status,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := s.repo.CreateReimbursement(rm); err != nil {
		return nil, err
	}
	return s.repo.FindReimbursementByID(rm.ID)
}

func (s *hrService) UpdateReimbursement(id string, dto *ReimbursementDTO) (*model.TrxHrReimbursement, error) {
	existing, err := s.repo.FindReimbursementByID(id)
	if err != nil {
		return nil, errors.New("reimbursement tidak ditemukan")
	}

	if dto.EmployeeID != "" {
		existing.EmployeeID = dto.EmployeeID
	}
	if dto.ClaimDate != "" {
		existing.ClaimDate = dto.ClaimDate
	}
	cat := dto.CategoryID
	if cat == "" {
		cat = dto.Category
	}
	if cat != "" {
		existing.CategoryID = cat
	}
	if dto.Amount > 0 {
		existing.Amount = dto.Amount
	}
	if dto.Description != "" {
		existing.Description = dto.Description
	}
	if dto.AttachmentURL != "" {
		existing.AttachmentURL = dto.AttachmentURL
	}
	if dto.Status != "" {
		existing.Status = dto.Status
	}
	existing.UpdatedAt = time.Now()

	if err := s.repo.UpdateReimbursement(existing); err != nil {
		return nil, err
	}
	return s.repo.FindReimbursementByID(id)
}

func (s *hrService) UpdateReimbursementStatus(id, status string) error {
	return s.repo.UpdateReimbursementStatus(id, status)
}

func (s *hrService) DeleteReimbursement(id string) error {
	return s.repo.DeleteReimbursement(id)
}

// ── Transaksi Penilaian Kinerja ─────────────────────────────────────

func (s *hrService) GetPerformanceReviewsPaginated(page, perPage int, employeeID, period, status, search string) ([]model.TrxHrPerformanceReview, int64, error) {
	return s.repo.FindPerformanceReviewsPaginated(page, perPage, employeeID, period, status, search)
}

func (s *hrService) GetPerformanceReviewByID(id string) (*model.TrxHrPerformanceReview, error) {
	return s.repo.FindPerformanceReviewByID(id)
}

func (s *hrService) CreatePerformanceReview(dto *PerformanceReviewDTO) (*model.TrxHrPerformanceReview, error) {
	if dto.EmployeeID == "" {
		return nil, errors.New("karyawan wajib dipilih")
	}
	if dto.ReviewerID == "" {
		dto.ReviewerID = "admin-reviewer"
	}
	if dto.OverallFeedback == "" && dto.Feedback != "" {
		dto.OverallFeedback = dto.Feedback
	}

	reviewNum := dto.ReviewNumber
	if reviewNum == "" {
		reviewNum = fmt.Sprintf("REV-%s-%s", time.Now().Format("200601"), strings.ToUpper(uuid.NewString()[:5]))
	}

	period := dto.ReviewPeriod
	if period == "" {
		period = dto.Period
	}
	if period == "" {
		period = fmt.Sprintf("%d-Q%d", time.Now().Year(), (int(time.Now().Month())-1)/3+1)
	}

	status := dto.Status
	if status == "" {
		status = "DRAFT"
	}

	reviewDate := dto.ReviewDate
	if reviewDate == "" {
		reviewDate = time.Now().Format("2006-01-02")
	}

	reviewID := uuid.NewString()
	var items []model.TrxHrPerformanceReviewItem
	var weightedSum float64 = 0
	var totalWeight float64 = 0

	for _, it := range dto.Items {
		crit := it.Criterion
		if crit == "" {
			crit = it.CriteriaName
		}
		comm := it.Comments
		if comm == "" {
			comm = it.Notes
		}
		itemObj := model.TrxHrPerformanceReviewItem{
			ID:        uuid.NewString(),
			ReviewID:  reviewID,
			Criterion: crit,
			Score:     it.Score,
			Weight:    it.Weight,
			Comments:  comm,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		items = append(items, itemObj)
		w := it.Weight
		if w <= 0 {
			w = 1
		}
		weightedSum += it.Score * w
		totalWeight += w
	}

	score := dto.TotalScore
	if score == 0 && totalWeight > 0 {
		score = weightedSum / totalWeight
	}

	pr := &model.TrxHrPerformanceReview{
		ID:              reviewID,
		ReviewNumber:    reviewNum,
		EmployeeID:      dto.EmployeeID,
		ReviewerID:      dto.ReviewerID,
		ReviewPeriod:    period,
		ReviewDate:      reviewDate,
		TotalScore:      score,
		OverallFeedback: dto.OverallFeedback,
		Status:          status,
		Items:           items,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	if err := s.repo.CreatePerformanceReview(pr); err != nil {
		return nil, err
	}
	return s.repo.FindPerformanceReviewByID(pr.ID)
}

func (s *hrService) UpdatePerformanceReview(id string, dto *PerformanceReviewDTO) (*model.TrxHrPerformanceReview, error) {
	existing, err := s.repo.FindPerformanceReviewByID(id)
	if err != nil {
		return nil, errors.New("penilaian kinerja tidak ditemukan")
	}

	if dto.EmployeeID != "" {
		existing.EmployeeID = dto.EmployeeID
	}
	if dto.ReviewerID != "" {
		existing.ReviewerID = dto.ReviewerID
	}
	period := dto.ReviewPeriod
	if period == "" {
		period = dto.Period
	}
	if period != "" {
		existing.ReviewPeriod = period
	}
	if dto.ReviewDate != "" {
		existing.ReviewDate = dto.ReviewDate
	}
	if dto.OverallFeedback != "" {
		existing.OverallFeedback = dto.OverallFeedback
	}
	if dto.Status != "" {
		existing.Status = dto.Status
	}
	existing.UpdatedAt = time.Now()

	var items []model.TrxHrPerformanceReviewItem
	var weightedSum float64 = 0
	var totalWeight float64 = 0

	for _, it := range dto.Items {
		itemObj := model.TrxHrPerformanceReviewItem{
			ID:        uuid.NewString(),
			ReviewID:  id,
			Criterion: it.Criterion,
			Score:     it.Score,
			Weight:    it.Weight,
			Comments:  it.Comments,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		items = append(items, itemObj)
		w := it.Weight
		if w <= 0 {
			w = 1
		}
		weightedSum += it.Score * w
		totalWeight += w
	}

	existing.Items = items
	if dto.TotalScore > 0 {
		existing.TotalScore = dto.TotalScore
	} else if totalWeight > 0 {
		existing.TotalScore = weightedSum / totalWeight
	}

	if err := s.repo.UpdatePerformanceReview(existing); err != nil {
		return nil, err
	}
	return s.repo.FindPerformanceReviewByID(id)
}

func (s *hrService) UpdatePerformanceReviewStatus(id, status string) error {
	return s.repo.UpdatePerformanceReviewStatus(id, status)
}

func (s *hrService) DeletePerformanceReview(id string) error {
	return s.repo.DeletePerformanceReview(id)
}
