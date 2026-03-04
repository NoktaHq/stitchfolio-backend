package entities

type EntityName string

const (
	Entity_BaseModel            EntityName = "BaseModel"
	Entity_Channel              EntityName = "Channel"
	Entity_Enquiry              EntityName = "Enquiry"
	Entity_EnquiryHistory       EntityName = "EnquiryHistory"
	Entity_Notification         EntityName = "Notification"
	Entity_User                 EntityName = "User"
	Entity_UserConfig           EntityName = "UserConfig"
	Entity_WhatsappNotification EntityName = "WhatsappNotification"
	Entity_MasterConfig         EntityName = "MasterConfig"
	Entity_Customer             EntityName = "Customer"
	Entity_Measurement          EntityName = "Measurement"
	Entity_Order                EntityName = "Order"
	Entity_OrderItem            EntityName = "OrderItem"
	Entity_ExpenseDetail        EntityName = "ExpenseDetail"
	Entity_Expense              EntityName = "Expense"
	Entity_EntityDocuments      EntityName = "EntityDocuments"
)

// string to entity name
func ToEntityName(s string) EntityName {
	return EntityName(s)
}
