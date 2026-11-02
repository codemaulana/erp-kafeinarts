package routes

import (
	"erp-kafeinart/handler"

	"github.com/gofiber/fiber/v3"
)

// RegisterHrRoutes mendaftarkan seluruh endpoint modul Sumber Daya Manusia (HR).
func RegisterHrRoutes(api fiber.Router, hrHandler *handler.HrHandler, jwtAuth fiber.Handler) {
	hr := api.Group("/hr", jwtAuth)

	// ── A. Data Master ──────────────────────────────────────────────────────────
	hr.Get("/departments", hrHandler.GetDepartments)
	hr.Get("/departments/:id", hrHandler.GetDepartmentByID)
	hr.Post("/departments", hrHandler.CreateDepartment)
	hr.Put("/departments/:id", hrHandler.UpdateDepartment)
	hr.Delete("/departments/:id", hrHandler.DeleteDepartment)

	hr.Get("/positions", hrHandler.GetPositions)
	hr.Get("/positions/:id", hrHandler.GetPositionByID)
	hr.Post("/positions", hrHandler.CreatePosition)
	hr.Put("/positions/:id", hrHandler.UpdatePosition)
	hr.Delete("/positions/:id", hrHandler.DeletePosition)

	hr.Get("/employees", hrHandler.GetEmployees)
	hr.Get("/employees/:id", hrHandler.GetEmployeeByID)
	hr.Post("/employees", hrHandler.CreateEmployee)
	hr.Put("/employees/:id", hrHandler.UpdateEmployee)
	hr.Delete("/employees/:id", hrHandler.DeleteEmployee)

	hr.Get("/leave-types", hrHandler.GetLeaveTypes)
	hr.Get("/leave-types/:id", hrHandler.GetLeaveTypeByID)
	hr.Post("/leave-types", hrHandler.CreateLeaveType)
	hr.Put("/leave-types/:id", hrHandler.UpdateLeaveType)
	hr.Delete("/leave-types/:id", hrHandler.DeleteLeaveType)

	hr.Get("/salary-components", hrHandler.GetSalaryComponents)
	hr.Get("/salary-components/:id", hrHandler.GetSalaryComponentByID)
	hr.Post("/salary-components", hrHandler.CreateSalaryComponent)
	hr.Put("/salary-components/:id", hrHandler.UpdateSalaryComponent)
	hr.Delete("/salary-components/:id", hrHandler.DeleteSalaryComponent)

	hr.Get("/work-shifts", hrHandler.GetWorkShifts)
	hr.Get("/work-shifts/:id", hrHandler.GetWorkShiftByID)
	hr.Post("/work-shifts", hrHandler.CreateWorkShift)
	hr.Put("/work-shifts/:id", hrHandler.UpdateWorkShift)
	hr.Delete("/work-shifts/:id", hrHandler.DeleteWorkShift)

	// ── B. Data Transaksi ───────────────────────────────────────────────────────
	// Presensi / Attendance
	hr.Get("/attendances", hrHandler.GetAttendances)
	hr.Get("/attendances/:id", hrHandler.GetAttendanceByID)
	hr.Post("/attendances", hrHandler.CreateAttendance)
	hr.Put("/attendances/:id", hrHandler.UpdateAttendance)
	hr.Delete("/attendances/:id", hrHandler.DeleteAttendance)
	// Presensi Aliases
	hr.Get("/presensi", hrHandler.GetAttendances)
	hr.Get("/presensi/:id", hrHandler.GetAttendanceByID)
	hr.Post("/presensi", hrHandler.CreateAttendance)
	hr.Post("/presensi/clock-in", hrHandler.CreateAttendance)
	hr.Post("/presensi/clock-out", hrHandler.CreateAttendance)
	hr.Put("/presensi/:id", hrHandler.UpdateAttendance)
	hr.Delete("/presensi/:id", hrHandler.DeleteAttendance)

	// Cuti / Leave Requests
	hr.Get("/leave-requests", hrHandler.GetLeaveRequests)
	hr.Get("/leave-requests/:id", hrHandler.GetLeaveRequestByID)
	hr.Post("/leave-requests", hrHandler.CreateLeaveRequest)
	hr.Put("/leave-requests/:id", hrHandler.UpdateLeaveRequest)
	hr.Patch("/leave-requests/:id/status", hrHandler.UpdateLeaveRequestStatus)
	hr.Delete("/leave-requests/:id", hrHandler.DeleteLeaveRequest)
	// Cuti Aliases
	hr.Get("/cuti", hrHandler.GetLeaveRequests)
	hr.Get("/cuti/:id", hrHandler.GetLeaveRequestByID)
	hr.Post("/cuti", hrHandler.CreateLeaveRequest)
	hr.Put("/cuti/:id", hrHandler.UpdateLeaveRequest)
	hr.Patch("/cuti/:id/status", hrHandler.UpdateLeaveRequestStatus)
	hr.Delete("/cuti/:id", hrHandler.DeleteLeaveRequest)

	// Payrolls
	hr.Get("/payrolls", hrHandler.GetPayrolls)
	hr.Get("/payrolls/:id", hrHandler.GetPayrollByID)
	hr.Post("/payrolls", hrHandler.CreatePayroll)
	hr.Put("/payrolls/:id", hrHandler.UpdatePayroll)
	hr.Patch("/payrolls/:id/status", hrHandler.UpdatePayrollStatus)
	hr.Delete("/payrolls/:id", hrHandler.DeletePayroll)
	// Payroll Aliases
	hr.Get("/payroll", hrHandler.GetPayrolls)
	hr.Get("/payroll/:id", hrHandler.GetPayrollByID)
	hr.Post("/payroll", hrHandler.CreatePayroll)
	hr.Put("/payroll/:id", hrHandler.UpdatePayroll)
	hr.Patch("/payroll/:id/status", hrHandler.UpdatePayrollStatus)
	hr.Delete("/payroll/:id", hrHandler.DeletePayroll)

	// Reimbursements
	hr.Get("/reimbursements", hrHandler.GetReimbursements)
	hr.Get("/reimbursements/:id", hrHandler.GetReimbursementByID)
	hr.Post("/reimbursements", hrHandler.CreateReimbursement)
	hr.Put("/reimbursements/:id", hrHandler.UpdateReimbursement)
	hr.Patch("/reimbursements/:id/status", hrHandler.UpdateReimbursementStatus)
	hr.Delete("/reimbursements/:id", hrHandler.DeleteReimbursement)
	// Reimbursement Aliases
	hr.Get("/reimbursement", hrHandler.GetReimbursements)
	hr.Get("/reimbursement/:id", hrHandler.GetReimbursementByID)
	hr.Post("/reimbursement", hrHandler.CreateReimbursement)
	hr.Put("/reimbursement/:id", hrHandler.UpdateReimbursement)
	hr.Patch("/reimbursement/:id/status", hrHandler.UpdateReimbursementStatus)
	hr.Delete("/reimbursement/:id", hrHandler.DeleteReimbursement)

	// Performance Reviews
	hr.Get("/performance-reviews", hrHandler.GetPerformanceReviews)
	hr.Get("/performance-reviews/:id", hrHandler.GetPerformanceReviewByID)
	hr.Post("/performance-reviews", hrHandler.CreatePerformanceReview)
	hr.Put("/performance-reviews/:id", hrHandler.UpdatePerformanceReview)
	hr.Patch("/performance-reviews/:id/status", hrHandler.UpdatePerformanceReviewStatus)
	hr.Delete("/performance-reviews/:id", hrHandler.DeletePerformanceReview)

	// ── C. Upload File Bukti / Struk HR ─────────────────────────────────────────
	hr.Post("/upload", hrHandler.UploadFile)

	// ── D. Laporan HR ───────────────────────────────────────────────────────────
	hr.Get("/laporan/rekap-absensi", hrHandler.GetRekapAbsensiReport)
	hr.Get("/laporan/slip-gaji", hrHandler.GetSlipGajiReport)
	hr.Get("/laporan/slip-gaji/:id", hrHandler.GetSlipGajiReport)
	hr.Get("/laporan/saldo-cuti", hrHandler.GetSaldoCutiReport)
	hr.Get("/laporan/turnover", hrHandler.GetTurnoverReport)
	// English alias
	hr.Get("/reports/rekap-absensi", hrHandler.GetRekapAbsensiReport)
	hr.Get("/reports/slip-gaji", hrHandler.GetSlipGajiReport)
	hr.Get("/reports/slip-gaji/:id", hrHandler.GetSlipGajiReport)
	hr.Get("/reports/saldo-cuti", hrHandler.GetSaldoCutiReport)
	hr.Get("/reports/turnover", hrHandler.GetTurnoverReport)
}
