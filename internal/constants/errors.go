package constants

import "errors"

var (
	Appointments_IS_nil             = errors.New("структура записи о приёме не задана: передайте непустой объект appointment для создания или обновления записи")
	PatientRecord_IS_nil            = errors.New("patientRecord is nil ")
	User_IS_nil                     = errors.New("user is nil")
	Service_IS_nil                  = errors.New("service is nil")
	ErrTimeNotInSchedule            = errors.New("выбранное время не входит в рабочее расписание врача: проверьте дату и время, или обновите расписание врача перед созданием приёма")
	ErrTimeConflict                 = errors.New("конфликт времени: выбранный интервал пересекается с существующей записью (у врача или у пациента). Выберите другое время или отмените/перенесите существующий приём")
	AppointmentCreateRequest_IS_nil = errors.New("appointment create request is nil")
	DoctorIDIsIncorrect             = errors.New("doctor id is incorrect")
	PatientIDIsIncorrect            = errors.New("patient id is incorrect")
	ServiceIDIsIncorrect            = errors.New("service id is incorrect")
	ErrInvalidAppointmentTime       = errors.New("invalid appointment time")
	ErrInvalidPrice                 = errors.New("invalid price")
	ErrInvalidAppointmentID         = errors.New("invalid appointment id")
	ErrUpdateAppointments           = errors.New("update error")
	ErrDeleteAppointments           = errors.New("delete error")
	ErrGetByIDAppointments          = errors.New("appointment not found ")
	ErrGetAppointments              = errors.New("error getting appointments")
	Rec_IS_nil                      = errors.New("recommendation is nil")
	PatientID_IS_incorrect          = errors.New("patientID is incorrect")
	Diagnosis_IS_empty              = errors.New("diagnosis is empty")
	DoctorID_IS_incorrect           = errors.New("doctorID is incorrect")
	Parse_ID_Error                  = errors.New("parse id  error")
	Invalid_JSON_Error              = errors.New("invalid json error")
	ErrCreateAppointment            = errors.New("error creating appointment")
	User_appointments_Not_Found     = errors.New("записи о приёмах для указанного пациента не найдены: проверьте корректность patientID или создайте приёмы для этого пациента")
)

// Schedule errors
var Schedule_IS_nil = errors.New("schedule is nil")
var Schedule_Conflict = errors.New("schedule conflict detected")
var Schedule_Not_Found = errors.New("schedule not found")
var Schedule_Invalid_Date_Range = errors.New("invalid date range")
var Schedule_No_Available_Slots = errors.New("no available slots found")
var Schedule_Doctor_Not_Found = errors.New("doctor not found")
var Schedule_Invalid_Input = errors.New("invalid schedule input")
var Schedule_Creation_Failed = errors.New("failed to create schedule")
var Schedule_Update_Failed = errors.New("failed to update schedule")
var Schedule_Deletion_Failed = errors.New("failed to delete schedule")
var ErrInvalidRoomNumber = errors.New("invalid room number")
var ErrInvalidTimeRange = errors.New("invalid time range")
var ErrInvalidDoctorID = errors.New("invalid doctor ID")

// Review errors
var Review_IS_nil = errors.New("review is nil")
