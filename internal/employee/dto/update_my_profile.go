package dto

// UpdateMyProfileRequest — hanya field yang boleh diubah sendiri oleh karyawan.
// Tidak termasuk: position, salaryBase, jobLevel, division, scheduleId, employmentStatus.
type UpdateMyProfileRequest struct {
	PhoneNumber       *string `json:"phoneNumber" validate:"omitempty,min=10,max=15"`
	Address           *string `json:"address" validate:"omitempty,min=5"`
	BankName          *string `json:"bankName" validate:"omitempty"`
	BankAccountNumber *string `json:"bankAccountNumber" validate:"omitempty"`
	BankAccountHolder *string `json:"bankAccountHolder" validate:"omitempty"`
}
