package models

type Recommendation struct {
	Base

	PatientID uint     `json:"patient_id" gorm:"not null"`
	Patient   *User    `json:"patient,omitempty" gorm:"foreignKey:PatientID"`
	ServiceID uint     `json:"service_id" gorm:"not null"`
	Service   *Service `json:"service,omitempty" gorm:"foreignKey:ServiceID"`
	DoctorID  uint     `json:"doctor_id,omitempty"`
	Doctor    *Doctor  `json:"doctor,omitempty" gorm:"foreignKey:DoctorID"`
	Note      string   `json:"note,omitempty"`
}

type RecommendationCreateRequest struct {
	PatientID uint   `json:"patient_id"`
	ServiceID uint   `json:"service_id"`
	Note      string `json:"note,omitempty"`
}
