package repository

import (
	"errors"
	"log/slog"
	"time"

	"github.com/mutsaevz/team-4-dentistry/internal/constants"
	"github.com/mutsaevz/team-4-dentistry/internal/models"
	"gorm.io/gorm"
)

type AppointmentRepository interface {
	Delete(uint) error
	GetByID(uint) (*models.Appointment, error)
	Get() ([]models.Appointment, error)
	Transaction(func(tx *gorm.DB) error) error
	CreateTx(tx *gorm.DB, appointment *models.Appointment) error
	UpdateTx(tx *gorm.DB, appointment *models.Appointment) error
	GetByPatientID(patientID uint) ([]models.Appointment, error)
	Update(appointment *models.Appointment) error
}
type gormAppointmentRepository struct {
	DB     *gorm.DB
	logger *slog.Logger
}

func NewAppointmentRepository(db *gorm.DB, logger *slog.Logger) AppointmentRepository {
	return &gormAppointmentRepository{DB: db, logger: logger}
}

func (r *gormAppointmentRepository) Delete(id uint) error {
	r.logger.Info("Удаление записи о приеме с ID", "id", id)

	if err := r.DB.Delete(&models.Appointment{}, id).Error; err != nil {
		r.logger.Error("ошибка при удалинии appointments", "ошибка", err, "appointments_id", id)
		return err
	}

	return nil

}

func (r *gormAppointmentRepository) GetByID(id uint) (*models.Appointment, error) {
	r.logger.Debug("получение appointment блягодаря ID", "appointments_id", id)
	var appointment models.Appointment

	if err := r.DB.First(&appointment, id).Error; err != nil {
		r.logger.Error("ошибка при получении appointments по ID", "ошибка", err, "appointments_id", id)
		return nil, err
	}

	r.logger.Info("успешное получение appointments по ID", "appointments_id", id)

	return &appointment, nil
}

func (r *gormAppointmentRepository) Get() ([]models.Appointment, error) {
	r.logger.Debug("получение всех appointments")
	var appointments []models.Appointment

	if err := r.DB.Find(&appointments).Error; err != nil {
		r.logger.Error("ошибка при получении всех appointments", "ошибка", err)
		return nil, err
	}

	r.logger.Info("успешное получение всех appointments")
	return appointments, nil
}

func (r *gormAppointmentRepository) Transaction(fn func(tx *gorm.DB) error) error {
	err := r.DB.Transaction(fn)
	if err != nil {
		r.logger.Error("ошибка при выполнении транзакции", "ошибка", err)
		return err
	}

	r.logger.Info("транзакция выполнена успешно")
	return nil
}

func (r *gormAppointmentRepository) CreateTx(tx *gorm.DB, appointment *models.Appointment) error {
	if appointment == nil {
		r.logger.Warn("попытка создать nil appointment")
		return constants.Appointments_IS_nil
	}

	r.logger.Debug("создание нового appointment", "appointment", appointment)

	dateOnly := time.Date(appointment.StartAt.Year(), appointment.StartAt.Month(), appointment.StartAt.Day(), 0, 0, 0, 0, appointment.StartAt.Location())
	var schedule models.Schedule
	if err := tx.Where("doctor_id = ? AND date = ? AND start_time <= ? AND end_time >= ?", appointment.DoctorID, dateOnly, appointment.StartAt, appointment.EndAt).First(&schedule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Warn("время appointment не входит в расписание врача", "doctor_id", appointment.DoctorID, "start_at", appointment.StartAt, "end_at", appointment.EndAt)
			return constants.ErrTimeNotInSchedule
		}
		r.logger.Error("ошибка при проверке расписания врача для нового appointment", "ошибка", err)
		return err
	}

	var count int64
	if err := tx.Model(&models.Appointment{}).
		Where("doctor_id = ? AND cancelled_at IS NULL AND start_at < ? AND end_at > ?", appointment.DoctorID, appointment.EndAt, appointment.StartAt).
		Count(&count).Error; err != nil {
		r.logger.Error("ошибка при проверке конфликтов по времени для нового appointment", "ошибка", err)
		return err
	}
	if count > 0 {
		r.logger.Warn("обнаружен конфликт по времени для нового appointment с другим appointment врача", "doctor_id", appointment.DoctorID, "start_at", appointment.StartAt, "end_at", appointment.EndAt)
		return constants.ErrTimeConflict
	}

	if err := tx.Model(&models.Appointment{}).
		Where("patient_id = ? AND cancelled_at IS NULL AND start_at < ? AND end_at > ?", appointment.PatientID, appointment.EndAt, appointment.StartAt).
		Count(&count).Error; err != nil {
		r.logger.Error("ошибка при проверке конфликтов по времени для нового appointment пациента", "ошибка", err)
		return err
	}

	if count > 0 {
		r.logger.Warn("обнаружен конфликт по времени для нового appointment с другим appointment пациента", "patient_id", appointment.PatientID, "start_at", appointment.StartAt, "end_at", appointment.EndAt)
		return constants.ErrTimeConflict
	}

	if err := tx.Create(appointment).Error; err != nil {
		r.logger.Error("ошибка при создании нового appointment", "ошибка", err)
		return err
	}

	r.logger.Info("успешное создание нового appointment", "appointment_id", appointment.ID)
	return nil
}

func (r *gormAppointmentRepository) UpdateTx(tx *gorm.DB, appointment *models.Appointment) error {
	if appointment == nil {
		r.logger.Warn("попытка обновить nil appointment")
		return constants.Appointments_IS_nil
	}

	dateOnly := time.Date(appointment.StartAt.Year(), appointment.StartAt.Month(), appointment.StartAt.Day(), 0, 0, 0, 0, appointment.StartAt.Location())
	var schedule models.Schedule
	if err := tx.Where("doctor_id = ? AND date = ? AND start_time <= ? AND end_time >= ?", appointment.DoctorID, dateOnly, appointment.StartAt, appointment.EndAt).First(&schedule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Warn("время appointment не входит в расписание врача", "doctor_id", appointment.DoctorID, "start_at", appointment.StartAt, "end_at", appointment.EndAt)
			return constants.ErrTimeNotInSchedule
		}
		r.logger.Error("ошибка при проверке расписания врача для обновленного appointment", "ошибка", err)
		return err
	}

	var count int64
	if err := tx.Model(&models.Appointment{}).
		Where("doctor_id = ? AND id <> ? AND cancelled_at IS NULL AND start_at < ? AND end_at > ?", appointment.DoctorID, appointment.ID, appointment.EndAt, appointment.StartAt).
		Count(&count).Error; err != nil {
		r.logger.Error("ошибка при проверке конфликтов по времени для обновленного appointment", "ошибка", err)
		return err
	}
	if count > 0 {
		r.logger.Warn("обнаружен конфликт по времени для обновленного appointment с другим appointment врача", "doctor_id", appointment.DoctorID, "start_at", appointment.StartAt, "end_at", appointment.EndAt)
		return constants.ErrTimeConflict
	}

	if err := tx.Model(&models.Appointment{}).
		Where("patient_id = ? AND id <> ? AND cancelled_at IS NULL AND start_at < ? AND end_at > ?", appointment.PatientID, appointment.ID, appointment.EndAt, appointment.StartAt).
		Count(&count).Error; err != nil {
		r.logger.Error("ошибка при проверке конфликтов по времени для обновленного appointment", "ошибка", err)
		return err
	}
	if count > 0 {
		r.logger.Warn("обнаружен конфликт по времени для обновленного appointment с другим appointment пациента", "patient_id", appointment.PatientID, "start_at", appointment.StartAt, "end_at", appointment.EndAt)
		return constants.ErrTimeConflict
	}

	err := tx.Save(appointment).Error
	if err != nil {
		r.logger.Error("ошибка при обновлении appointment", "ошибка", err)
		return err
	}

	r.logger.Info("успешное обновление appointment", "appointment_id", appointment.ID)
	return nil
}

func (r *gormAppointmentRepository) GetByPatientID(patientID uint) ([]models.Appointment, error) {
	r.logger.Debug("получение appointments по patientID", "patient_id", patientID)
	var appointment []models.Appointment

	if err := r.DB.Where("patient_id = ?", patientID).Find(&appointment); err != nil {
		r.logger.Error("ошибка при получении appointments по patientID", "ошибка", err, "patient_id", patientID)
		return nil, constants.User_appointments_Not_Found
	}

	if len(appointment) == 0 {
		r.logger.Warn("appointments для данного patientID не найдены", "patient_id", patientID)
		return nil, constants.User_appointments_Not_Found
	}

	r.logger.Info("успешное получение appointments по patientID", "patient_id", patientID)

	return appointment, nil
}

func (r *gormAppointmentRepository) Update(appointment *models.Appointment) error {
	if err := r.DB.Save(appointment).Error; err != nil {
		r.logger.Error("ошибка при обновлении appointment", "ошибка", err)
		return err
	}
	r.logger.Info("успешное обновление appointment", "appointment_id", appointment.ID)
	return nil
}
