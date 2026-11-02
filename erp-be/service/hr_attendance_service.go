package service

import (
	"errors"
	"time"

	"erp-kafeinart/model"

	"github.com/google/uuid"
)

// ── Transaksi Presensi ──────────────────────────────────────────────

func (s *hrService) GetAttendancesPaginated(page, perPage int, employeeID, date, status, search string) ([]model.TrxHrAttendance, int64, error) {
	return s.repo.FindAttendancesPaginated(page, perPage, employeeID, date, status, search)
}

func (s *hrService) GetAttendanceByID(id string) (*model.TrxHrAttendance, error) {
	return s.repo.FindAttendanceByID(id)
}

func (s *hrService) CreateAttendance(dto *AttendanceDTO) (*model.TrxHrAttendance, error) {
	if dto.EmployeeID == "" {
		return nil, errors.New("karyawan wajib dipilih")
	}

	attDate := dto.AttendanceDate
	if attDate == "" {
		attDate = dto.Date
	}
	if attDate == "" {
		attDate = time.Now().Format("2006-01-02")
	}

	checkIn := dto.CheckInTime
	if checkIn == "" {
		checkIn = dto.CheckIn
	}
	if checkIn == "" {
		checkIn = dto.TimeIn
	}
	if checkIn == "" {
		checkIn = time.Now().Format("15:04:05")
	}

	checkOut := dto.CheckOutTime
	if checkOut == "" {
		checkOut = dto.CheckOut
	}
	if checkOut == "" {
		checkOut = dto.TimeOut
	}

	// If attendance already exists for today and clocking out, update existing
	existing, _ := s.repo.FindAttendanceByEmployeeAndDate(dto.EmployeeID, attDate)
	if existing != nil {
		if checkOut != "" {
			existing.CheckOutTime = checkOut
		}
		if dto.Notes != "" {
			existing.Notes = dto.Notes
		}
		existing.UpdatedAt = time.Now()
		if err := s.repo.UpdateAttendance(existing); err != nil {
			return nil, err
		}
		return s.repo.FindAttendanceByID(existing.ID)
	}

	lat := dto.CheckInLat
	if lat == 0 {
		lat = dto.Latitude
	}

	lng := dto.CheckInLng
	if lng == 0 {
		lng = dto.Longitude
	}

	status := dto.AttendanceStatus
	if status == "" {
		status = dto.Status
	}
	if status == "" {
		status = "PRESENT"
	}

	att := &model.TrxHrAttendance{
		ID:               uuid.NewString(),
		EmployeeID:       dto.EmployeeID,
		AttendanceDate:   attDate,
		CheckInTime:      checkIn,
		CheckOutTime:     checkOut,
		CheckInLat:       lat,
		CheckInLng:       lng,
		AttendanceStatus: status,
		Notes:            dto.Notes,
		SelfiePhotoURL:   dto.SelfiePhotoURL,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	if err := s.repo.CreateAttendance(att); err != nil {
		return nil, err
	}
	return s.repo.FindAttendanceByID(att.ID)
}

func (s *hrService) UpdateAttendance(id string, dto *AttendanceDTO) (*model.TrxHrAttendance, error) {
	existing, err := s.repo.FindAttendanceByID(id)
	if err != nil {
		return nil, errors.New("presensi tidak ditemukan")
	}

	if dto.EmployeeID != "" {
		existing.EmployeeID = dto.EmployeeID
	}
	attDate := dto.AttendanceDate
	if attDate == "" {
		attDate = dto.Date
	}
	if attDate != "" {
		existing.AttendanceDate = attDate
	}
	checkIn := dto.CheckInTime
	if checkIn == "" {
		checkIn = dto.CheckIn
	}
	if checkIn != "" {
		existing.CheckInTime = checkIn
	}
	checkOut := dto.CheckOutTime
	if checkOut == "" {
		checkOut = dto.CheckOut
	}
	if checkOut != "" {
		existing.CheckOutTime = checkOut
	}
	lat := dto.CheckInLat
	if lat == 0 {
		lat = dto.Latitude
	}
	if lat != 0 {
		existing.CheckInLat = lat
	}
	lng := dto.CheckInLng
	if lng == 0 {
		lng = dto.Longitude
	}
	if lng != 0 {
		existing.CheckInLng = lng
	}
	status := dto.AttendanceStatus
	if status == "" {
		status = dto.Status
	}
	if status != "" {
		existing.AttendanceStatus = status
	}
	if dto.Notes != "" {
		existing.Notes = dto.Notes
	}
	if dto.SelfiePhotoURL != "" {
		existing.SelfiePhotoURL = dto.SelfiePhotoURL
	}
	existing.UpdatedAt = time.Now()

	if err := s.repo.UpdateAttendance(existing); err != nil {
		return nil, err
	}
	return s.repo.FindAttendanceByID(id)
}

func (s *hrService) DeleteAttendance(id string) error {
	return s.repo.DeleteAttendance(id)
}
