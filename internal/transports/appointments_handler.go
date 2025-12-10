package transports

import "github.com/mutsaevz/team-4-dentistry/internal/services"

type AppointmentsHandler struct {
	Service services.AppointmentService
}

func NewAppointmentsHandler(appointmentService services.AppointmentService) *AppointmentsHandler {
	return &AppointmentsHandler{Service: appointmentService}
}

