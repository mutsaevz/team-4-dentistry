package models


type PatientRecord struct {
	Base
	PatientID uint   `json:"patient_id" gorm:"not null"`
	//Patient  *User  `json:"patient,omitempty" gorm:"foreignKey:PatientID"`
	DoctorID  uint   `json:"doctor_id" gorm:"not null"`
	Doctor   *Doctor `json:"doctor,omitempty" gorm:"foreignKey:DoctorID"`
	Diagnosis string `json:"diagnosis,omitempty"`
}
