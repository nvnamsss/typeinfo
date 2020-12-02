package recursive

type PromotionItemFact struct {
	UserID      string                 `json:"user_id" valid:"required~user_id is required"`
	InvoiceNo   string                 `json:"invoice_no" valid:"optional,length(6|20)"`
	ServiceType string                 `json:"service_type" valid:"optional,in(DIRECT_PAY|TICKET|BILL|VOUCHER|TOPUP|TOPUP_CARD)"`
	ChannelCode string                 `json:"channel_code" valid:"optional,length(1|32)~channel_code max length 32 characters"`
	Instrument  string                 `json:"instrument" valid:"optional,in(MQR|TQR|VNQR|INAPP|CQR|ATA|AIA|WEB|SUB)"`
	Amount      int64                  `json:"amount" valid:"nullableNumeric~Amount must be numeric"`
	XDeviceID   string                 `json:"x_device_id" valid:"optional"`
	ClientID    string                 `json:"client_id" valid:"required~ClientID is required"`
	ExtraData   map[string]interface{} `json:"extra_data"`
}

func (f *PromotionItemFact) Call() {

}
