package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"erp-kafeinart/model"

	"github.com/google/uuid"
)

// ── Transaksi Payroll ───────────────────────────────────────────────

func (s *hrService) GetPayrollsPaginated(page, perPage int, period, employeeID, status, search string) ([]model.TrxHrPayroll, int64, error) {
	return s.repo.FindPayrollsPaginated(page, perPage, period, employeeID, status, search)
}

func (s *hrService) GetPayrollByID(id string) (*model.TrxHrPayroll, error) {
	return s.repo.FindPayrollByID(id)
}

func (s *hrService) CreatePayroll(dto *PayrollDTO) (*model.TrxHrPayroll, error) {
	if dto.EmployeeID == "" {
		return nil, errors.New("karyawan wajib dipilih")
	}

	payrollNum := dto.PayrollNumber
	if payrollNum == "" {
		payrollNum = fmt.Sprintf("PAY-%s-%s", time.Now().Format("200601"), strings.ToUpper(uuid.NewString()[:5]))
	}

	period := dto.PayrollPeriod
	if period == "" {
		period = time.Now().Format("2006-01")
	}

	status := dto.PaymentStatus
	if status == "" {
		status = dto.Status
	}
	if status == "" {
		status = "DRAFT"
	}

	payrollID := uuid.NewString()
	var items []model.TrxHrPayrollItem
	var calcEarnings float64 = 0
	var calcDeductions float64 = 0

	for _, it := range dto.Items {
		itemObj := model.TrxHrPayrollItem{
			ID:          uuid.NewString(),
			PayrollID:   payrollID,
			ComponentID: it.ComponentID,
			Amount:      it.Amount,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		items = append(items, itemObj)

		isDeduction := false
		if strings.EqualFold(it.ItemType, "DEDUCTION") || strings.EqualFold(it.ComponentType, "DEDUCTION") {
			isDeduction = true
		} else if scObj, err := s.repo.FindSalaryComponentByID(it.ComponentID); err == nil && scObj != nil && strings.EqualFold(scObj.ComponentType, "DEDUCTION") {
			isDeduction = true
		}

		if isDeduction {
			calcDeductions += it.Amount
		} else {
			calcEarnings += it.Amount
		}
	}

	// Add top-level fields if provided
	if dto.OvertimePay > 0 {
		calcEarnings += dto.OvertimePay
	}
	if dto.Bonus > 0 {
		calcEarnings += dto.Bonus
	}
	if dto.Deductions > 0 && calcDeductions == 0 {
		calcDeductions += dto.Deductions
	}

	earnings := dto.TotalEarnings
	if earnings == 0 {
		earnings = calcEarnings
	}
	deductions := dto.TotalDeductions
	if deductions == 0 {
		deductions = calcDeductions
	}
	netSalary := dto.NetSalary
	if netSalary == 0 {
		netSalary = earnings - deductions
	}

	p := &model.TrxHrPayroll{
		ID:              payrollID,
		PayrollNumber:   payrollNum,
		PayrollPeriod:   period,
		PeriodMonth:     dto.PeriodMonth,
		PeriodYear:      dto.PeriodYear,
		EmployeeID:      dto.EmployeeID,
		TotalEarnings:   earnings,
		TotalDeductions: deductions,
		NetSalary:       netSalary,
		PaymentStatus:   status,
		PaymentDate:     dto.PaymentDate,
		Notes:           dto.Notes,
		Items:           items,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	if err := s.repo.CreatePayroll(p); err != nil {
		return nil, err
	}
	return s.repo.FindPayrollByID(p.ID)
}

func (s *hrService) UpdatePayroll(id string, dto *PayrollDTO) (*model.TrxHrPayroll, error) {
	existing, err := s.repo.FindPayrollByID(id)
	if err != nil {
		return nil, errors.New("payroll tidak ditemukan")
	}

	if dto.EmployeeID != "" {
		existing.EmployeeID = dto.EmployeeID
	}
	if dto.PayrollPeriod != "" {
		existing.PayrollPeriod = dto.PayrollPeriod
	}
	status := dto.PaymentStatus
	if status == "" {
		status = dto.Status
	}
	if status != "" {
		existing.PaymentStatus = status
	}
	if dto.PaymentDate != "" {
		existing.PaymentDate = dto.PaymentDate
	}
	existing.Notes = dto.Notes
	existing.UpdatedAt = time.Now()

	var items []model.TrxHrPayrollItem
	var calcEarnings float64 = 0
	var calcDeductions float64 = 0

	for _, it := range dto.Items {
		itemObj := model.TrxHrPayrollItem{
			ID:          uuid.NewString(),
			PayrollID:   id,
			ComponentID: it.ComponentID,
			Amount:      it.Amount,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		items = append(items, itemObj)

		scObj, err := s.repo.FindSalaryComponentByID(it.ComponentID)
		if err == nil && scObj.ComponentType == "DEDUCTION" {
			calcDeductions += it.Amount
		} else {
			calcEarnings += it.Amount
		}
	}

	existing.Items = items
	existing.TotalEarnings = dto.TotalEarnings
	if existing.TotalEarnings == 0 {
		existing.TotalEarnings = calcEarnings
	}
	existing.TotalDeductions = dto.TotalDeductions
	if existing.TotalDeductions == 0 {
		existing.TotalDeductions = calcDeductions
	}
	existing.NetSalary = dto.NetSalary
	if existing.NetSalary == 0 {
		existing.NetSalary = existing.TotalEarnings - existing.TotalDeductions
	}

	if err := s.repo.UpdatePayroll(existing); err != nil {
		return nil, err
	}
	return s.repo.FindPayrollByID(id)
}

func (s *hrService) UpdatePayrollStatus(id, status string) error {
	return s.repo.UpdatePayrollStatus(id, status)
}

func (s *hrService) DeletePayroll(id string) error {
	return s.repo.DeletePayroll(id)
}
