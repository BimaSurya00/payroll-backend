package dto

type DashboardSummary struct {
	Attendance AttendanceSummary `json:"attendance"`
	Leave      LeaveSummary      `json:"leave"`
	Payroll    PayrollSummary    `json:"payroll"`
	Employee   EmployeeSummary   `json:"employee"`
}

type AttendanceSummary struct {
	TodayPresent   int `json:"todayPresent"`
	TodayLate      int `json:"todayLate"`
	TodayAbsent    int `json:"todayAbsent"`
	TodayLeave     int `json:"todayLeave"`
	TotalEmployees int `json:"totalEmployees"`
}

type LeaveSummary struct {
	PendingRequests   int `json:"pendingRequests"`
	ApprovedThisMonth int `json:"approvedThisMonth"`
	RejectedThisMonth int `json:"rejectedThisMonth"`
}

type PayrollSummary struct {
	DraftCount     int     `json:"draftCount"`
	ApprovedCount  int     `json:"approvedCount"`
	PaidCount      int     `json:"paidCount"`
	TotalNetSalary float64 `json:"totalNetSalary"`
	CurrentPeriod  string  `json:"currentPeriod"`
}

type EmployeeSummary struct {
	TotalActive   int `json:"totalActive"`
	TotalInactive int `json:"totalInactive"`
	NewThisMonth  int `json:"newThisMonth"`
}

// AttendanceStats - Detailed attendance statistics
type AttendanceStats struct {
	Period         string                  `json:"period"`
	DailyBreakdown []DailyAttendanceStat   `json:"dailyBreakdown"`
	Summary        AttendancePeriodSummary `json:"summary"`
}

type DailyAttendanceStat struct {
	Date    string `json:"date"`
	Present int    `json:"present"`
	Late    int    `json:"late"`
	Absent  int    `json:"absent"`
	Leave   int    `json:"leave"`
	Total   int    `json:"total"`
}

type AttendancePeriodSummary struct {
	TotalPresent int     `json:"totalPresent"`
	TotalLate    int     `json:"totalLate"`
	TotalAbsent  int     `json:"totalAbsent"`
	TotalLeave   int     `json:"totalLeave"`
	AvgPresent   float64 `json:"avgPresent"`
	AvgLate      float64 `json:"avgLate"`
}

// PayrollStats - Detailed payroll statistics
type PayrollStats struct {
	Period          string                  `json:"period"`
	TotalPayrolls   int                     `json:"totalPayrolls"`
	TotalAmount     float64                 `json:"totalAmount"`
	AverageSalary   float64                 `json:"averageSalary"`
	StatusBreakdown PayrollStatusBreakdown  `json:"statusBreakdown"`
	DepartmentStats []DepartmentPayrollStat `json:"departmentStats"`
}

type PayrollStatusBreakdown struct {
	DraftCount     int     `json:"draftCount"`
	ApprovedCount  int     `json:"approvedCount"`
	PaidCount      int     `json:"paidCount"`
	DraftAmount    float64 `json:"draftAmount"`
	ApprovedAmount float64 `json:"approvedAmount"`
	PaidAmount     float64 `json:"paidAmount"`
}

type DepartmentPayrollStat struct {
	DepartmentID   string  `json:"departmentId"`
	DepartmentName string  `json:"departmentName"`
	EmployeeCount  int     `json:"employeeCount"`
	TotalPayroll   float64 `json:"totalPayroll"`
}

// EmployeeStats - Detailed employee statistics
type EmployeeStats struct {
	TotalCount      int                      `json:"totalCount"`
	StatusBreakdown EmployeeStatusBreakdown  `json:"statusBreakdown"`
	DepartmentStats []DepartmentEmployeeStat `json:"departmentStats"`
	JobLevelStats   []JobLevelStat           `json:"jobLevelStats"`
	RecentHires     []RecentHire             `json:"recentHires"`
}

type EmployeeStatusBreakdown struct {
	Permanent int `json:"permanent"`
	Contract  int `json:"contract"`
	Probation int `json:"probation"`
	Intern    int `json:"intern"`
	Resigned  int `json:"resigned"`
}

type DepartmentEmployeeStat struct {
	DepartmentID   string `json:"departmentId"`
	DepartmentName string `json:"departmentName"`
	EmployeeCount  int    `json:"employeeCount"`
	ActiveCount    int    `json:"activeCount"`
}

type JobLevelStat struct {
	Level         string `json:"level"`
	EmployeeCount int    `json:"employeeCount"`
}

type RecentHire struct {
	EmployeeID   string `json:"employeeId"`
	EmployeeName string `json:"employeeName"`
	Position     string `json:"position"`
	JoinDate     string `json:"joinDate"`
}

// RecentActivity - Activity log entry
type RecentActivity struct {
	ID           string `json:"id"`
	Timestamp    string `json:"timestamp"`
	Action       string `json:"action"`
	ResourceType string `json:"resourceType"`
	ResourceID   string `json:"resourceId"`
	UserID       string `json:"userId"`
	UserName     string `json:"userName"`
	Description  string `json:"description"`
}

type RecentActivitiesResponse struct {
	Activities []RecentActivity `json:"activities"`
	Total      int              `json:"total"`
}

// SuperUser Dashboard DTOs

type SuperUserSummary struct {
	TotalCompanies  int           `json:"totalCompanies"`
	ActiveCompanies int           `json:"activeCompanies"`
	TotalUsers      int           `json:"totalUsers"`
	TotalEmployees  int           `json:"totalEmployees"`
	TotalAdmins     int           `json:"totalAdmins"`
	TotalSuperUsers int           `json:"totalSuperUsers"`
	CompanyStats    []CompanyStat `json:"companyStats"`
}

type CompanyStat struct {
	CompanyID     string `json:"companyId"`
	CompanyName   string `json:"companyName"`
	Plan          string `json:"plan"`
	IsActive      bool   `json:"isActive"`
	MaxEmployees  int    `json:"maxEmployees"`
	UserCount     int    `json:"userCount"`
	EmployeeCount int    `json:"employeeCount"`
	CreatedAt     string `json:"createdAt"`
}
