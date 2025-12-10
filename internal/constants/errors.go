package constants

import "errors"

var Appointments_IS_nil = errors.New("appointments is nil ")
var PatientRecord_IS_nil = errors.New("PatientRecord is nil ")
var User_IS_nil = errors.New("user is nil")
var Service_IS_nil = errors.New("service is nil")

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
