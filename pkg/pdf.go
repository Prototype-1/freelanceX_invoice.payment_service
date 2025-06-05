package pkg

import (
	"bytes"
	"fmt"
	"time"
	"github.com/jung-kurt/gofpdf"
	"github.com/Prototype-1/freelanceX_invoice.payment_service/internal/model" 
)

func GenerateInvoicePDF(invoice *model.Invoice) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 20)
	pdf.Cell(0, 10, "Invoice")
	pdf.Ln(15)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Invoice ID: %s", invoice.ID.String()))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Project ID: %s", invoice.ProjectID.String()))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Client ID: %s", invoice.ClientID.String()))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Freelancer ID: %s", invoice.FreelancerID.String()))
	pdf.Ln(8)

	pdf.Cell(40, 10, fmt.Sprintf("Invoice Type: %s", invoice.Type))
	pdf.Ln(8)

	pdf.Cell(40, 10, fmt.Sprintf("Status: %s", invoice.Status))
	pdf.Ln(8)

	if invoice.DueDate != nil {
		pdf.Cell(40, 10, fmt.Sprintf("Due Date: %s", invoice.DueDate.Format("2006-01-02")))
		pdf.Ln(8)
	}

	pdf.Cell(40, 10, fmt.Sprintf("Hours Worked: %.2f", invoice.HoursWorked))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Hourly Rate: $%.2f", invoice.HourlyRate))
	pdf.Ln(8)

	pdf.Cell(40, 10, fmt.Sprintf("Milestone/Phase: %s", invoice.MilestonePhase))
	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, fmt.Sprintf("Amount Due: ₹%.2f", invoice.Amount))
	pdf.Ln(15)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 10, fmt.Sprintf("Generated on %s", time.Now().Format("2006-01-02 15:04:05")))
	pdf.Ln(10)

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
