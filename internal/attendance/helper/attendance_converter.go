package helper

import (
	"time"

	"hris/internal/attendance/dto"
	"hris/internal/attendance/entity"
)

func ToAttendanceResponse(attendance *entity.Attendance) *dto.AttendanceResponse {
	resp := &dto.AttendanceResponse{
		ID:         attendance.ID,
		EmployeeID: attendance.EmployeeID,
		Date:       attendance.Date.Format(time.RFC3339),
		Status:     attendance.Status,
		Notes:      attendance.Notes,
	}

	if attendance.ClockInTime != nil {
		clockInTime := attendance.ClockInTime.Format(time.RFC3339)
		resp.ClockInTime = &clockInTime
	}

	if attendance.ClockOutTime != nil {
		clockOutTime := attendance.ClockOutTime.Format(time.RFC3339)
		resp.ClockOutTime = &clockOutTime
	}

	return resp
}

func ToAttendanceResponseWithSchedule(attendance *entity.Attendance, scheduleName string) *dto.AttendanceResponse {
	resp := ToAttendanceResponse(attendance)
	resp.ScheduleName = &scheduleName
	return resp
}

func ToAttendanceResponses(attendances []*entity.Attendance) []*dto.AttendanceResponse {
	responses := make([]*dto.AttendanceResponse, len(attendances))
	for i, attendance := range attendances {
		responses[i] = ToAttendanceResponse(attendance)
	}
	return responses
}
