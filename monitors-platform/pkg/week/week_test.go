package week_test

import (
    "testing"
    "time"

    "github.com/org/monitors-platform/pkg/week"
    "github.com/stretchr/testify/assert"
)

func TestGetWeekStart_Lunes(t *testing.T) {
    lunes := time.Date(2026, 3, 30, 10, 0, 0, 0, time.UTC)
    result := week.GetWeekStart(lunes)
    assert.Equal(t, time.Date(2026, 3, 30, 0, 0, 0, 0, time.UTC), result)
}

func TestGetWeekStart_Miercoles(t *testing.T) {
    miercoles := time.Date(2026, 4, 1, 10, 0, 0, 0, time.UTC)
    result := week.GetWeekStart(miercoles)
    assert.Equal(t, time.Date(2026, 3, 30, 0, 0, 0, 0, time.UTC), result)
}

func TestGetWeekStart_Domingo(t *testing.T) {
    domingo := time.Date(2026, 4, 5, 23, 0, 0, 0, time.UTC)
    result := week.GetWeekStart(domingo)
    assert.Equal(t, time.Date(2026, 3, 30, 0, 0, 0, 0, time.UTC), result)
}

func TestIsLateReport_SemanaPasada(t *testing.T) {
    pasada := week.GetWeekStart(time.Now()).AddDate(0, 0, -7)
    assert.True(t, week.IsLateReport(pasada))
}

func TestIsLateReport_SemanaActual(t *testing.T) {
    actual := week.GetWeekStart(time.Now())
    assert.False(t, week.IsLateReport(actual))
}

func TestIsLateReport_SemanaFutura(t *testing.T) {
    futura := week.GetWeekStart(time.Now()).AddDate(0, 0, 7)
    assert.False(t, week.IsLateReport(futura))
}
