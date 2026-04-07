package dto

type AttendanceReportItem struct {
	EmployeeID     string  `json:"employeeId"`
	EmployeeName   string  `json:"employeeName"`
	Position       string  `json:"position"`
	Division       string  `json:"division"`
	TotalPresent   int     `json:"totalPresent"`
	TotalLate      int     `json:"totalLate"`
	TotalAbsent    int     `json:"totalAbsent"`
	TotalLeave     int     `json:"totalLeave"`
	TotalDays      int     `json:"totalDays"`
	AttendanceRate float64 `json:"attendanceRate"` // Percentage
	Period         string  `json:"period"`
}

type MonthlyAttendanceReport struct {
	Period         string                  `json:"period"`
	Month          int                     `json:"month"`
	Year           int                     `json:"year"`
	TotalEmployees int                     `json:"totalEmployees"`
	Summary        AttendanceReportItem    `json:"summary"` // Aggregate of all employees
	Items          []AttendanceReportItem  `json:"items"`
}

type MyAttendanceSummary struct {
	Period         string  `json:"period"`
	TotalPresent   int     `json:"totalPresent"`
	TotalLate      int     `json:"totalLate"`
	TotalAbsent    int     `json:"totalAbsent"`
	TotalLeave     int     `json:"totalLeave"`
	TotalDays      int     `json:"totalDays"`
	AttendanceRate float64 `json:"attendanceRate"`
}

type AttendanceSummary struct {
	TotalPresent int `json:"totalPresent"`
	TotalLate    int `json:"totalLate"`
	TotalAbsent  int `json:"totalAbsent"`
	TotalLeave   int `json:"totalLeave"`
}

func NewAttendanceSummary(present, late, absent, leave int) AttendanceSummary {
	return AttendanceSummary{
		TotalPresent: present,
		TotalLate:    late,
		TotalAbsent:  absent,
		TotalLeave:   leave,
	}
}

func CalculateAttendanceRate(present, late, totalDays int) float64 {
	if totalDays == 0 {
		return 0.0
	}
	presentDays := present + late
	return float64(presentDays) / float64(totalDays) * 100.0
}
