package dto

type PurchaseTicketRequest struct {
	ProductID      int64  `json:"product_id" validate:"required"`
	SelectedDate   string `json:"selected_date" validate:"required"`
	SelectedTime   string `json:"selected_time" validate:"required"`
	Category       string `json:"category" validate:"required"`
	AdditionalInfo string `json:"additional_info"` // Optional fields
}

type PurchaseTicketResponse struct {
	Message      string      `json:"message"`
	PaymentLink  string      `json:"payment_link,omitempty"`
	TicketDetail interface{} `json:"ticket_detail,omitempty"`
}
