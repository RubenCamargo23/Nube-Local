package week

import "time"

// GetWeekStart retorna el lunes de la semana de la fecha dada.
func GetWeekStart(t time.Time) time.Time {
    t = t.UTC().Truncate(24 * time.Hour)
    weekday := int(t.Weekday())
    if weekday == 0 {
        weekday = 7 // domingo = 7
    }
    return t.AddDate(0, 0, -(weekday - 1))
}

// IsLateReport retorna true si la semana ya cerró (es anterior a la semana actual).
func IsLateReport(weekStart time.Time) bool {
    currentWeekStart := GetWeekStart(time.Now())
    return weekStart.Before(currentWeekStart)
}
