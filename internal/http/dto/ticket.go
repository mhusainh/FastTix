// internal/http/dto/ticket.go

package dto

type PurchaseTicketRequest2 struct {
	ProductID    int64  `json:"product_id" validate:"required"`
	SelectedDate string `json:"selected_date" validate:"required"` // e.g., "2024-12-25"
	SelectedTime string `json:"selected_time" validate:"required"` // e.g., "18:00"
	Category     string `json:"category" validate:"required"`      // e.g., "VIP", "Regular"
	eTicketInfo  string `json:"eticket_info" validate:"required"`  // Additional eTicket information
}

type PurchaseTicketResponse2 struct {
	Message     string `json:"message"`
	PaymentLink string `json:"payment_link,omitempty"`
	TicketURL   string `json:"ticket_url,omitempty"`
}
