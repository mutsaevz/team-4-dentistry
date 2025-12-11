package services

import (
	"github.com/mutsaevz/team-4-dentistry/internal/constants"
	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"github.com/mutsaevz/team-4-dentistry/internal/repository"
)

type PatientRecordService interface {
	Create(req *models.PatientRecordCreate) (*models.PatientRecord, error)
	GetByID(ID uint) (*models.PatientRecord, error)
	GetAll() ([]models.PatientRecord, error)
	Update(id uint, req *models.PatientRecordUpdate) error
	Delete(ID uint) error
}

type patientRecord struct {
	repo repository.PatientRecordRepo
}

func NewPatientRecordService(repo repository.PatientRecordRepo) PatientRecordService {
	return &patientRecord{repo: repo}
}

func (s *patientRecord) Create(req *models.PatientRecordCreate) (*models.PatientRecord, error) {
	if req == nil {
		return nil, constants.PatientRecord_IS_nil
	}

	if req.PatientID == 0 {
		return nil, constants.PatientID_IS_incorrect
	}

	if req.Diagnosis == "" {
		return nil, constants.Diagnosis_IS_empty
	}

	if req.DoctorID == 0 {	
		return nil, constants.DoctorID_IS_incorrect
	}

	patientRecord := &models.PatientRecord{
		PatientID: req.PatientID,
		Diagnosis: req.Diagnosis,
		DoctorID:  req.DoctorID,
	}

	if err := s.repo.Create(patientRecord); err != nil {
		return nil, err
	}

	return patientRecord, nil
}

func (s *patientRecord) GetByID(ID uint) (*models.PatientRecord, error) {

	if ID <= 0 {
		return nil, constants.PatientID_IS_incorrect
	}

	patientRecord, err := s.repo.GetID(ID)
	if err != nil {
		return nil, err
	}

	return patientRecord, nil
}

func (s *patientRecord) GetAll() ([]models.PatientRecord, error) {
	patientRecords, err := s.repo.Get()
	if err != nil {
		return nil, err
	}

	return patientRecords, nil
}	

func (s *patientRecord) Update(id uint,req *models.PatientRecordUpdate) error {
	if req == nil {
		return constants.PatientRecord_IS_nil
	}

	 patientRecord,err := s.repo.GetID(id)
	if err != nil {
		return err
	}

	if req.Diagnosis != nil && *req.Diagnosis == "" {
		return constants.Diagnosis_IS_empty
	}

	if req.DoctorID != nil && *req.DoctorID <= 0 {
		return constants.DoctorID_IS_incorrect
	}

	if req.Diagnosis != nil {
		patientRecord.Diagnosis = *req.Diagnosis
	}

	if req.DoctorID != nil {
		patientRecord.DoctorID = *req.DoctorID
	}

	if err := s.repo.Update(patientRecord);  err != nil {
		return err
	}

	return nil
}

func (r *patientRecord) Delete(ID uint) error {
	if ID <= 0 {
		return constants.PatientID_IS_incorrect
	}

	if err := r.repo.Delete(ID); err != nil {
		return err
	}

	return nil
}

