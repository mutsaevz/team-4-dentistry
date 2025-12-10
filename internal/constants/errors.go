package constants

import "errors"

var (
	Appointments_IS_nil             = errors.New("appointments is nil ")
	PatientRecord_IS_nil            = errors.New("patientRecord is nil ")
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
	Rec_IS_nil                      = errors.New("recommendation is nil")
	PatientID_IS_incorrect          = errors.New("patientID is incorrect")
	Diagnosis_IS_empty              = errors.New("diagnosis is empty")
	DoctorID_IS_incorrect           = errors.New("doctorID is incorrect")
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
