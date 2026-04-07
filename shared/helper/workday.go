package helper

import "time"

// CountWorkingDays menghitung jumlah hari kerja (Senin-Jumat) antara dua tanggal (inklusif).
func CountWorkingDays(startDate, endDate time.Time) int {
	return CountWorkingDaysExcluding(startDate, endDate, nil)
}

// CountWorkingDaysExcluding menghitung jumlah hari kerja (Senin-Jumat) yang bukan hari libur.
// holidays adalah map dengan format "YYYY-MM-DD" -> true
func CountWorkingDaysExcluding(startDate, endDate time.Time, holidays map[string]bool) int {
	if startDate.After(endDate) {
		return 0
	}

	workingDays := 0
	current := startDate

	for !current.After(endDate) {
		weekday := current.Weekday()
		// Skip weekends
		if weekday != time.Saturday && weekday != time.Sunday {
			// Skip holidays if provided
			dateStr := current.Format("2006-01-02")
			if holidays == nil || !holidays[dateStr] {
				workingDays++
			}
		}
		current = current.AddDate(0, 0, 1)
	}

	return workingDays
}
