package reports

import (
    "fmt"
    "time"

    "github.com/jung-kurt/gofpdf"
)

func GeneratePDF(outputPath, nombrePersona, semana, resumen string) error {
    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.AddPage()
    pdf.SetFont("Arial", "B", 16)

    pdf.Cell(190, 10, "Reporte Semanal de Actividades")
    pdf.Ln(12)

    pdf.SetFont("Arial", "", 12)
    pdf.Cell(190, 8, fmt.Sprintf("Persona: %s", nombrePersona))
    pdf.Ln(8)
    pdf.Cell(190, 8, fmt.Sprintf("Semana: %s", semana))
    pdf.Ln(8)
    pdf.Cell(190, 8, fmt.Sprintf("Generado: %s", time.Now().Format("02/01/2006 15:04")))
    pdf.Ln(12)

    pdf.SetFont("Arial", "B", 12)
    pdf.Cell(190, 8, "Resumen generado por IA:")
    pdf.Ln(10)

    pdf.SetFont("Arial", "", 11)
    pdf.MultiCell(190, 6, resumen, "", "", false)

    return pdf.OutputFileAndClose(outputPath)
}
