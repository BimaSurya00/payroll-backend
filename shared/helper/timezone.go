package helper

import (
	"time"
	"sync"
)

var (
	appTimezone *time.Location
	once        sync.Once
)

// InitTimezone menginisialisasi timezone aplikasi. Dipanggil sekali saat startup.
func InitTimezone(tz string) error {
	var err error
	once.Do(func() {
		if tz == "" {
			tz = "Asia/Jakarta" // default WIB
		}
		appTimezone, err = time.LoadLocation(tz)
	})
	return err
}

// Now mengembalikan waktu sekarang dalam timezone aplikasi.
func Now() time.Time {
	if appTimezone == nil {
		return time.Now() // fallback
	}
	return time.Now().In(appTimezone)
}

// Today mengembalikan awal hari ini (00:00:00) dalam timezone aplikasi.
func Today() time.Time {
	now := Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, appTimezone)
}

// GetLocation mengembalikan timezone location.
func GetLocation() *time.Location {
	if appTimezone == nil {
		loc, _ := time.LoadLocation("Asia/Jakarta")
		return loc
	}
	return appTimezone
}
