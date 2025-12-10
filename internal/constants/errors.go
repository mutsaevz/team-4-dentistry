package constants

import "errors"

var (
	Appointments_IS_nil             = errors.New("appointments is nil ")
	PatientRecord_IS_nil            = errors.New("PatientRecord is nil ")
	User_IS_nil                     = errors.New("user is nil")
	Service_IS_nil                  = errors.New("service is nil")
	ErrTimeNotInSchedule            = errors.New("выбранное время не входит в расписание врача")
	ErrTimeConflict                 = errors.New("выбранное время уже занято")
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
)
