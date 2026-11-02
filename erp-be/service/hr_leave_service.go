package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"erp-kafeinart/model"

	"github.com/google/uuid"
)

// ── Transaksi Pengajuan Cuti ────────────────────────────────────────

func (s *hrService) GetLeaveRequestsPaginated(page, perPage int, employeeID, status, search string) ([]model.TrxHrLeaveRequest, int64, error) {
	return s.repo.FindLeaveRequestsPaginated(page, perPage, employeeID, status, search)
}

func (s *hrService) GetLeaveRequestByID(id string) (*model.TrxHrLeaveRequest, error) {
	return s.repo.FindLeaveRequestByID(id)
}

func (s *hrService) CreateLeaveRequest(dto *LeaveRequestDTO) (*model.TrxHrLeaveRequest, error) {
	if dto.EmployeeID == "" {
		return nil, errors.New("karyawan wajib dipilih")
	}
	if dto.LeaveTypeID == "" {
		return nil, errors.New("jenis cuti wajib dipilih")
	}

	reqNum := dto.RequestNumber
	if reqNum == "" {
		reqNum = fmt.Sprintf("CUTI-%s-%s", time.Now().Format("200601"), strings.ToUpper(uuid.NewString()[:5]))
	}

	status := dto.ApprovalStatus
	if status == "" {
		status = dto.Status
	}
	if status == "" {
		status = "SUBMITTED"
	}

	totalDays := dto.TotalDays
	if totalDays <= 0 {
		totalDays = 1
	}

	lr := &model.TrxHrLeaveRequest{
		ID:              uuid.NewString(),
		RequestNumber:   reqNum,
		EmployeeID:      dto.EmployeeID,
		LeaveTypeID:     dto.LeaveTypeID,
		StartDate:       dto.StartDate,
		EndDate:         dto.EndDate,
		TotalDays:       totalDays,
		Reason:          dto.Reason,
		AttachmentURL:   dto.AttachmentURL,
		ApprovalStatus:  status,
		Status:          status,
		RejectionReason: dto.RejectionReason,
		ApprovedBy:      dto.ApprovedBy,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	if err := s.repo.CreateLeaveRequest(lr); err != nil {
		return nil, err
	}
	return s.repo.FindLeaveRequestByID(lr.ID)
}

func (s *hrService) UpdateLeaveRequest(id string, dto *LeaveRequestDTO) (*model.TrxHrLeaveRequest, error) {
	existing, err := s.repo.FindLeaveRequestByID(id)
	if err != nil {
		return nil, errors.New("pengajuan cuti tidak ditemukan")
	}

	if dto.EmployeeID != "" {
		existing.EmployeeID = dto.EmployeeID
	}
	if dto.LeaveTypeID != "" {
		existing.LeaveTypeID = dto.LeaveTypeID
	}
	if dto.StartDate != "" {
		existing.StartDate = dto.StartDate
	}
	if dto.EndDate != "" {
		existing.EndDate = dto.EndDate
	}
	if dto.TotalDays > 0 {
		existing.TotalDays = dto.TotalDays
	}
	if dto.Reason != "" {
		existing.Reason = dto.Reason
	}
	if dto.AttachmentURL != "" {
		existing.AttachmentURL = dto.AttachmentURL
	}
	status := dto.ApprovalStatus
	if status == "" {
		status = dto.Status
	}
	if status != "" {
		existing.ApprovalStatus = status
		existing.Status = status
	}
	if dto.RejectionReason != "" {
		existing.RejectionReason = dto.RejectionReason
	}
	if dto.ApprovedBy != "" {
		existing.ApprovedBy = dto.ApprovedBy
	}
	existing.UpdatedAt = time.Now()

	if err := s.repo.UpdateLeaveRequest(existing); err != nil {
		return nil, err
	}
	return s.repo.FindLeaveRequestByID(id)
}

func (s *hrService) UpdateLeaveRequestStatus(id, status, reason, approvedBy string) error {
	return s.repo.UpdateLeaveRequestStatus(id, status, reason, approvedBy)
}

func (s *hrService) DeleteLeaveRequest(id string) error {
	return s.repo.DeleteLeaveRequest(id)
}
