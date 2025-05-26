package model

import (
	invoicepb "github.com/Prototype-1/freelanceX_invoice.payment_service/proto/invoice_service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToProto(m *Invoice) *invoicepb.Invoice {
	pb := &invoicepb.Invoice{
		InvoiceId:     m.ID.String(),
		FreelancerId:  m.FreelancerID.String(),
		ClientId:      m.ClientID.String(),
		ProjectId:     m.ProjectID.String(),
		Type:          invoicepb.InvoiceType(invoicepb.InvoiceType_value[m.Type]),
		Amount:        m.Amount,
		HourlyRate:    m.HourlyRate,
		HoursWorked:   m.HoursWorked,
		Status:        invoicepb.InvoiceStatus(invoicepb.InvoiceStatus_value[m.Status]),
		MilestonePhase: m.MilestonePhase,
	}
if m.DueDate != nil {
		pb.DueDate = timestamppb.New(*m.DueDate)
	}
	return pb
}
