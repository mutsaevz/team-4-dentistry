package services

import (
	"log/slog"

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
	repo   repository.PatientRecordRepo
	logger *slog.Logger
}

func NewPatientRecordService(repo repository.PatientRecordRepo, logger *slog.Logger) PatientRecordService {
	return &patientRecord{repo: repo, logger: logger}
}

func (s *patientRecord) Create(req *models.PatientRecordCreate) (*models.PatientRecord, error) {
	s.logger.Debug("Create PatientRecord вызван", "patient_id", req.PatientID)

	if req == nil {
		s.logger.Warn("передан nil PatientRecordCreate")
		return nil, constants.PatientRecord_IS_nil
	}

	if req.PatientID == 0 {
		s.logger.Warn("некорректный PatientID при создании patient record")
		return nil, constants.PatientID_IS_incorrect
	}

	if req.Diagnosis == "" {
		s.logger.Warn("пустая диагноза при создании patient record", "patient_id", req.PatientID)
		return nil, constants.Diagnosis_IS_empty
	}

	if req.DoctorID == 0 {
		s.logger.Warn("некорректный DoctorID при создании patient record", "patient_id", req.PatientID)
		return nil, constants.DoctorID_IS_incorrect
	}

	patientRecord := &models.PatientRecord{
		PatientID: req.PatientID,
		Diagnosis: req.Diagnosis,
		DoctorID:  req.DoctorID,
	}

	if err := s.repo.Create(patientRecord); err != nil {
		s.logger.Error("ошибка при создании patient record", "error", err)
		return nil, err
	}

	s.logger.Info("patient record создан", "id", patientRecord.ID, "patient_id", patientRecord.PatientID)
	return patientRecord, nil
}

func (s *patientRecord) GetByID(ID uint) (*models.PatientRecord, error) {
	s.logger.Debug("GetByID PatientRecord вызван", "id", ID)

	if ID <= 0 {
		s.logger.Warn("некорректный ID при GetByID patient record", "id", ID)
		return nil, constants.PatientID_IS_incorrect
	}

	patientRecord, err := s.repo.GetID(ID)
	if err != nil {
		s.logger.Error("ошибка при получении patient record по ID", "error", err, "id", ID)
		return nil, err
	}

	s.logger.Info("patient record получен по ID", "id", ID)
	return patientRecord, nil
}

func (s *patientRecord) GetAll() ([]models.PatientRecord, error) {
	s.logger.Debug("GetAll PatientRecords вызван")
	patientRecords, err := s.repo.Get()
	if err != nil {
		s.logger.Error("ошибка при получении всех patient records", "error", err)
		return nil, err
	}

	s.logger.Info("patient records получены", "count", len(patientRecords))
	return patientRecords, nil
}

func (s *patientRecord) Update(id uint, req *models.PatientRecordUpdate) error {
	s.logger.Debug("Update PatientRecord вызван", "id", id)
	if req == nil {
		s.logger.Warn("передан nil PatientRecordUpdate", "id", id)
		return constants.PatientRecord_IS_nil
	}

	patientRecord, err := s.repo.GetID(id)
	if err != nil {
		s.logger.Error("ошибка при получении patient record для обновления", "error", err, "id", id)
		return err
	}

	if req.Diagnosis != nil && *req.Diagnosis == "" {
		s.logger.Warn("пустая диагноза в Update", "id", id)
		return constants.Diagnosis_IS_empty
	}

	if req.DoctorID != nil && *req.DoctorID <= 0 {
		s.logger.Warn("некорректный DoctorID в Update", "id", id)
		return constants.DoctorID_IS_incorrect
	}

	if req.Diagnosis != nil {
		patientRecord.Diagnosis = *req.Diagnosis
	}

	if req.DoctorID != nil {
		patientRecord.DoctorID = *req.DoctorID
	}

	if err := s.repo.Update(patientRecord); err != nil {
		s.logger.Error("ошибка при обновлении patient record", "error", err, "id", id)
		return err
	}

	s.logger.Info("patient record успешно обновлен", "id", id)
	return nil
}

func (r *patientRecord) Delete(ID uint) error {
	r.logger.Debug("Delete PatientRecord вызван", "id", ID)
	if ID <= 0 {
		r.logger.Warn("некорректный ID при Delete patient record", "id", ID)
		return constants.PatientID_IS_incorrect
	}

	if err := r.repo.Delete(ID); err != nil {
		r.logger.Error("ошибка при удалении patient record", "error", err, "id", ID)
		return err
	}

	r.logger.Info("patient record удален", "id", ID)
	return nil
}
