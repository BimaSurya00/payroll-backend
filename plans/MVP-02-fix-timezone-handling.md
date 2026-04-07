# MVP-02: Fix Timezone Handling

## Prioritas: 🔴 CRITICAL — Data Accuracy
## Estimasi: 2 jam
## Tipe: Bug Fix

---

## Deskripsi Masalah

Attendance service menggunakan `time.Now().UTC()` untuk menentukan tanggal dan waktu clock in/out.
Indonesia menggunakan WIB (UTC+7). Karyawan yang clock in jam 08:00 WIB akan tercatat sebagai 01:00 UTC.

Dampak:
- **Status LATE/PRESENT salah** — schedule time (e.g. "09:00") di-compare dengan waktu UTC
- **Tanggal bisa salah** — clock in jam 00:01-07:00 WIB = hari sebelumnya di UTC

Lokasi bug:
- `internal/attendance/service/attendance_service_impl.go` line 67 dan 132
- `internal/leave/service/leave_service.go` line beberapa tempat menggunakan `time.Now()`

## Solusi

Tambah konfigurasi timezone di config, lalu gunakan di semua service yang membutuhkan waktu lokal.

## File yang Diubah

### 1. [MODIFY] `config/config.go`

**Tambah field di `AppConfig`:**
```go
type AppConfig struct {
    Name     string
    Env      string
    Port     string
    Host     string
    Timezone string // Tambah ini
}
```

**Di `LoadConfig()`, tambah:**
```go
App: AppConfig{
    Name:     viper.GetString("APP_NAME"),
    Env:      viper.GetString("APP_ENV"),
    Port:     viper.GetString("APP_PORT"),
    Host:     viper.GetString("APP_HOST"),
    Timezone: viper.GetString("APP_TIMEZONE"), // Tambah ini
},
```

### 2. [NEW] `shared/helper/timezone.go`

Buat helper untuk memuat timezone:

```go
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
```

### 3. [MODIFY] `main.go`

**Tambah inisialisasi timezone setelah config load:**
```go
// Setelah LoadConfig(), tambahkan:
if err := sharedHelper.InitTimezone(cfg.App.Timezone); err != nil {
    zap.L().Fatal("failed to initialize timezone", zap.Error(err))
}
```

**Import yang perlu ditambah:**
```go
sharedHelper "hris/shared/helper"
```

### 4. [MODIFY] `internal/attendance/service/attendance_service_impl.go`

**Replace semua `time.Now().UTC()` dengan `sharedHelper.Now()`:**

**Line 67-68 (ClockIn):**
```go
// SEBELUM:
now := time.Now().UTC()
today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

// SESUDAH:
now := sharedHelper.Now()
today := sharedHelper.Today()
```

**Line 132-133 (ClockOut):**
```go
// SEBELUM:
now := time.Now().UTC()
today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

// SESUDAH:
now := sharedHelper.Now()
today := sharedHelper.Today()
```

**Tambah import:**
```go
sharedHelper "hris/shared/helper"
```

**Juga fix `determineStatus` (line 209-224):**
```go
func (s *attendanceService) determineStatus(clockInTime time.Time, scheduleTimeIn string, allowedLateMinutes int) string {
    scheduleTime, err := time.Parse("15:04", scheduleTimeIn)
    if err != nil {
        return "PRESENT"
    }

    loc := sharedHelper.GetLocation()
    localClockIn := clockInTime.In(loc)

    scheduledTime := time.Date(localClockIn.Year(), localClockIn.Month(), localClockIn.Day(),
        scheduleTime.Hour(), scheduleTime.Minute(), 0, 0, loc)

    deadline := scheduledTime.Add(time.Duration(allowedLateMinutes) * time.Minute)

    if localClockIn.After(deadline) {
        return "LATE"
    }

    return "PRESENT"
}
```

### 5. [MODIFY] `.env` dan `.env.example`

**Tambah:**
```env
APP_TIMEZONE=Asia/Jakarta
```

## Verifikasi

1. `go build ./...` — pastikan compile sukses
2. Test clock in:
   - Pada jam 08:00 WIB → status harus `PRESENT` (schedule 09:00)
   - Pada jam 09:30 WIB (late > 15 min) → status harus `LATE`
   - Pada jam 00:01 WIB → tanggal harus hari ini (bukan kemarin)
3. Cek log output — waktu harus menampilkan timezone WIB
