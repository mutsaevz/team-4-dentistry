package models

type PatientRecord struct {
	Base
	PatientID uint `json:"patient_id" gorm:"not null"`
	//Patient  *User  `json:"patient,omitempty" gorm:"foreignKey:PatientID"`
	DoctorID  uint    `json:"doctor_id" gorm:"not null"`
	Doctor    *Doctor `json:"doctor,omitempty" gorm:"foreignKey:DoctorID"`
	Diagnosis string  `json:"diagnosis,omitempty"`
}

type PatientRecordCreate struct {
	PatientID uint    `json:"patient_id" validate:"required"`
	Patient   *User   `json:"patient,omitempty" validate:"omitempty"`
	DoctorID  uint    `json:"doctor_id" validate:"required"`
	Doctor    *Doctor `json:"doctor,omitempty" validate:"omitempty"`
	Diagnosis string  `json:"diagnosis,omitempty" validate:"omitempty"`
}

type PatientRecordUpdate struct {
	DoctorID  *uint   `json:"doctor_id,omitempty" validate:"omitempty"`
	Diagnosis *string `json:"diagnosis,omitempty" validate:"omitempty"`
}
