package service

import (
	"context"
	"nailzbydardo/internal/model"
	"time"
)

type DashboardSummary struct {
	UpcomingAppointments    []model.Appointment `json:"upcoming_appointments"`
	MonthlyRevenue          int64               `json:"monthly_revenue"`
	MonthlyTips             int64               `json:"monthly_tips"`
	MonthlyAppointmentCount int64               `json:"monthly_appointment_count"`
	MonthlyExpenses         int64               `json:"monthly_expenses"`
}

type FinancialsSummary struct {
	Revenue          int64 `json:"revenue"`
	Expenses         int64 `json:"expenses"`
	AppointmentCount int64 `json:"appointment_count"`
	Tips             int64 `json:"tips"`
}

type DashboardService struct {
	appointmentService *AppointmentService
	expenseService     *ExpenseService
}

func NewDashboardService(appointmentService *AppointmentService, expenseService *ExpenseService) *DashboardService {
	return &DashboardService{appointmentService: appointmentService, expenseService: expenseService}
}

func (s *DashboardService) GetDashboard(ctx context.Context) (DashboardSummary, error) {
	now := time.Now()
	startOfMonth := time.Date(
		now.Year(),
		now.Month(),
		1,
		0, 0, 0, 0,
		now.Location(),
	)
	startOfNextMonth := startOfMonth.AddDate(0, 1, 0)

	upcoming, err := s.appointmentService.ListUpcomingAppointments(ctx)
	if err != nil {
		return DashboardSummary{}, err
	}
	completeAppts, err := s.appointmentService.ListCompleteAppointmentsForPeriod(ctx, startOfMonth, startOfNextMonth)
	if err != nil {
		return DashboardSummary{}, err
	}
	var revenue int64
	var tips int64
	for _, completeAppt := range completeAppts {
		total, err := s.appointmentService.CalculateAppointmentTotal(ctx, completeAppt.ID)
		if err != nil {
			return DashboardSummary{}, err
		}
		revenue += total.ServiceTotal
		tips += total.Tip
	}
	count, err := s.appointmentService.GetAppointmentCountForPeriod(ctx, startOfMonth, startOfNextMonth)
	if err != nil {
		return DashboardSummary{}, err
	}
	expenses, err := s.expenseService.GetExpensesForPeriod(ctx, startOfMonth, startOfNextMonth)
	if err != nil {
		return DashboardSummary{}, err
	}
	dashboardSummary := DashboardSummary{
		UpcomingAppointments:    upcoming,
		MonthlyRevenue:          revenue,
		MonthlyTips:             tips,
		MonthlyAppointmentCount: count,
		MonthlyExpenses:         expenses,
	}
	return dashboardSummary, nil
}


func (s *DashboardService) GetFinancials(ctx context.Context, dateOne time.Time, dateTwo time.Time) (FinancialsSummary, error) {
    completeAppts, err := s.appointmentService.ListCompleteAppointmentsForPeriod(ctx, dateOne, dateTwo)
    if err != nil {
        return FinancialsSummary{}, err
    }

    var revenue, tips int64
    for _, appt := range completeAppts {
        total, err := s.appointmentService.CalculateAppointmentTotal(ctx, appt.ID)
        if err != nil {
            return FinancialsSummary{}, err
        }
        revenue += total.ServiceTotal
        tips += total.Tip
    }

    count, err := s.appointmentService.GetAppointmentCountForPeriod(ctx, dateOne, dateTwo)
    if err != nil {
        return FinancialsSummary{}, err
    }

    expenses, err := s.expenseService.GetExpensesForPeriod(ctx, dateOne, dateTwo)
    if err != nil {
        return FinancialsSummary{}, err
    }

    return FinancialsSummary{
        Revenue:          revenue,
        Expenses:         expenses,
        AppointmentCount: count,
        Tips:             tips,
    }, nil
}
