package helper

import (
	"time"

	"hris/internal/schedule/dto"
	"hris/internal/schedule/entity"
)

func ToScheduleResponse(schedule *entity.Schedule) *dto.ScheduleResponse {
	return &dto.ScheduleResponse{
		ID:                 schedule.ID,
		Name:               schedule.Name,
		TimeIn:             schedule.TimeIn,
		TimeOut:            schedule.TimeOut,
		AllowedLateMinutes: schedule.AllowedLateMinutes,
		Description:        schedule.Description,
		CreatedAt:          schedule.CreatedAt.Format(time.RFC3339),
		UpdatedAt:          schedule.UpdatedAt.Format(time.RFC3339),
	}
}

func ToScheduleResponses(schedules []*entity.Schedule) []*dto.ScheduleResponse {
	responses := make([]*dto.ScheduleResponse, len(schedules))
	for i, schedule := range schedules {
		responses[i] = ToScheduleResponse(schedule)
	}
	return responses
}
