package handler

import (
	"fmt"
	"nailzbydardo/internal/service"
	"net/http"
	"strings"
	"time"
)

type CalendarHandler struct {
	appointmentService *service.AppointmentService
	calendarSecret string
	baseURL string
}

func NewCalendarHandler(appointmentService *service.AppointmentService, calendarSecret string, baseURL string) *CalendarHandler {
	return &CalendarHandler{appointmentService: appointmentService, calendarSecret: calendarSecret, baseURL: baseURL}
} 

func (h *CalendarHandler) GetCalendarFeed(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if h.calendarSecret == "" || token != h.calendarSecret {
    writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "unauthorized"})
    return
}

	appts, err := h.appointmentService.ListAppointmentsForCalendar(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}

	now := time.Now().UTC()
	salonName := "NailzByDardo"

	events := make([]string, 0, len(appts))
	for _, appt := range appts {
		start := appt.ApptDate.UTC()
		end := start.Add(120 * time.Minute)

		lines := []string{
			"BEGIN:VEVENT",
			fmt.Sprintf("UID:appointment-%s@nailzbydardo", appt.ID),
			fmt.Sprintf("DTSTAMP:%s", toICalDate(now)),
			fmt.Sprintf("DTSTART:%s", toICalDate(start)),
			fmt.Sprintf("DTEND:%s", toICalDate(end)),
			fmt.Sprintf("SUMMARY:%s", escapeText(orDefault(appt.ClientName, "Appointment"))),
		}
		if appt.Notes != nil && *appt.Notes != "" {
			lines = append(lines, fmt.Sprintf("DESCRIPTION:%s", escapeText(*appt.Notes)))
		}
		if h.baseURL != "" {
			lines = append(lines, fmt.Sprintf("URL:%s/appointments/%s", h.baseURL, appt.ID))
		}
		lines = append(lines, "STATUS:CONFIRMED", "END:VEVENT")
		events = append(events, strings.Join(lines, "\r\n"))
	}

	cal := []string{
		"BEGIN:VCALENDAR",
		"VERSION:2.0",
		fmt.Sprintf("PRODID:-//%s//Appointments//EN", salonName),
		fmt.Sprintf("X-WR-CALNAME:%s", salonName),
		"X-WR-TIMEZONE:America/Chicago",
		"REFRESH-INTERVAL;VALUE=DURATION:PT1H",
		"X-PUBLISHED-TTL:PT1H",
	}
	cal = append(cal, events...)
	cal = append(cal, "END:VCALENDAR")

	w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
	w.Header().Set("Content-Disposition", `inline; filename="nailzbydardo.ics"`)
	w.Header().Set("Cache-Control", "no-cache")
	w.Write([]byte(strings.Join(cal, "\r\n")))
}

func toICalDate(t time.Time) string {
	return t.Format("20060102T150405Z")
}

func escapeText(s string) string {
	r := strings.NewReplacer(`\`, `\\`, `;`, `\;`, `,`, `\,`, "\n", `\n`)
	return r.Replace(s)
}

func orDefault(s, def string) string {
	if s == "" {
		return def
	}
	return s
}
