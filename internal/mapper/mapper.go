package mapper

import (
	"time"

	"github.com/imkarthi24/sf-backend/internal/entities"
	entitiy_types "github.com/imkarthi24/sf-backend/internal/entities/types"
	requestModel "github.com/imkarthi24/sf-backend/internal/model/request"
	"github.com/loop-kar/pixie/util"
)

type Mapper interface {
	User(requestModel.User) (*entities.User, error)
	Channel(requestModel.Channel) (*entities.Channel, error)
	Enquiry(e requestModel.Enquiry) (*entities.Enquiry, error)
	EnquiryHistory(e requestModel.EnquiryHistory) (*entities.EnquiryHistory, error)
	MasterConfig(e requestModel.MasterConfig) (*entities.MasterConfig, error)
	UserChannelDetail(e requestModel.UserChannelDetail) (*entities.UserChannelDetail, error)
	UserChannelDetails(items []requestModel.UserChannelDetail) ([]entities.UserChannelDetail, error)
	Customer(e requestModel.Customer) (*entities.Customer, error)
	Person(e requestModel.Person) (*entities.Person, error)
	DressType(e requestModel.DressType) (*entities.DressType, error)
	Measurement(e requestModel.Measurement) (*entities.Measurement, error)
	Order(e requestModel.Order) (*entities.Order, error)
	OrderItem(e requestModel.OrderItem) (*entities.OrderItem, error)
	OrderItems(items []requestModel.OrderItem) ([]entities.OrderItem, error)
	OrderHistory(e requestModel.OrderHistory) (*entities.OrderHistory, error)
	MeasurementHistory(e requestModel.MeasurementHistory) (*entities.MeasurementHistory, error)
	ExpenseTracker(e requestModel.ExpenseTracker) (*entities.Expense, error)
	ExpenseDetail(e requestModel.ExpenseDetail) (*entities.ExpenseDetail, error)
	ExpenseDetails(items []requestModel.ExpenseDetail) ([]entities.ExpenseDetail, error)
	Task(e requestModel.Task) (*entities.Task, error)
	Category(e requestModel.Category) (*entities.Category, error)
	Product(e requestModel.Product) (*entities.Product, error)
	Inventory(e requestModel.Inventory) (*entities.Inventory, error)
	InventoryLog(e requestModel.InventoryLog) (*entities.InventoryLog, error)
}

type mapper struct{}

func ProvideMapper() Mapper {
	return &mapper{}
}

func (m mapper) User(e requestModel.User) (*entities.User, error) {
	loginTime, err := util.GenerateDateTimeFromString(&e.LastLoginTime)
	if err != nil {
		return nil, err
	}
	userChannelDetails, err := m.UserChannelDetails(e.UserChannelDetails)
	if err != nil {
		return nil, err
	}

	return &entities.User{
		Model:               &entities.Model{ID: e.ID, IsActive: e.IsActive},
		FirstName:           e.FirstName,
		LastName:            e.LastName,
		Extension:           e.Extension,
		PhoneNumber:         e.PhoneNumber,
		Email:               e.Email,
		Password:            e.Password,
		Role:                entities.RoleType(e.Role),
		IsLoginDisabled:     e.IsLoginDisabled,
		IsLoggedIn:          e.IsLoggedIn,
		LastLoginTime:       loginTime,
		LoginFailureCounter: int16(e.LoginFailureCounter),
		ResetPasswordString: e.ResetPasswordString,
		Experience:          e.Experience,
		Department:          e.Department,
		UserChannelDetails:  userChannelDetails,
	}, nil
}

func (*mapper) Channel(chnl requestModel.Channel) (*entities.Channel, error) {
	return &entities.Channel{
		Model:       &entities.Model{ID: chnl.ID, IsActive: chnl.IsActive},
		Name:        chnl.Name,
		Status:      entities.ChannelStatus(chnl.Status),
		OwnerUserID: chnl.OwnerUserId,
	}, nil
}

func (m mapper) Enquiry(e requestModel.Enquiry) (*entities.Enquiry, error) {
	return &entities.Enquiry{
		Model:               &entities.Model{ID: e.ID, IsActive: e.IsActive},
		Subject:             e.Subject,
		Notes:               e.Notes,
		Status:              entities.EnquiryStatus(e.Status),
		CustomerId:          e.CustomerId,
		Source:              e.Source,
		ReferredBy:          e.ReferredBy,
		ReferrerPhoneNumber: e.ReferrerPhoneNumber,
	}, nil
}

func (mapper) EnquiryHistory(e requestModel.EnquiryHistory) (*entities.EnquiryHistory, error) {
	var visitingDate *time.Time
	var callBackDate *time.Time
	var err error

	if e.VisitingDate != nil {
		visitingDate, err = util.GenerateDateTimeFromString(e.VisitingDate)
		if err != nil {
			return nil, err
		}
	}

	if e.CallBackDate != nil {
		callBackDate, err = util.GenerateDateTimeFromString(e.CallBackDate)
		if err != nil {
			return nil, err
		}
	}

	var enquiryDate *time.Time
	if e.EnquiryDate != "" {
		date, err := util.GenerateDateTimeFromString(&e.EnquiryDate)
		if err != nil {
			return nil, err
		}
		enquiryDate = date
	}

	var status *entities.EnquiryStatus
	if e.Status != nil {
		s := entities.EnquiryStatus(*e.Status)
		status = &s
	}

	performedAt := util.GetLocalTime()
	if e.PerformedAt != "" {
		date, err := util.GenerateDateTimeFromString(&e.PerformedAt)
		if err != nil {
			return nil, err
		}
		performedAt = *date
	}

	return &entities.EnquiryHistory{
		Model:           &entities.Model{ID: e.ID, IsActive: e.IsActive},
		Status:          status,
		EmployeeComment: e.EmployeeComment,
		CustomerComment: e.CustomerComment,
		VisitingDate:    visitingDate,
		CallBackDate:    callBackDate,
		EnquiryDate:     enquiryDate,
		ResponseStatus:  entities.ResponseStatus(e.ResponseStatus),
		EnquiryId:       e.EnquiryId,
		EmployeeId:      e.EmployeeId,
		PerformedAt:     performedAt,
		PerformedById:   e.PerformedById,
	}, nil
}

func (m *mapper) MasterConfig(e requestModel.MasterConfig) (*entities.MasterConfig, error) {
	return &entities.MasterConfig{
		Model: &entities.Model{ID: e.ID, IsActive: e.IsActive},
		Name:  e.Name,
		Type:  e.Type,

		CurrentValue:  e.CurrentValue,
		PreviousValue: e.PreviousValue,
		DefaultValue:  e.DefaultValue,
		UseDefault:    e.UseDefault,

		Description: e.Description,

		Format: e.Format,
	}, nil
}

// UserChannelDetail implements Mapper.
func (m *mapper) UserChannelDetail(e requestModel.UserChannelDetail) (*entities.UserChannelDetail, error) {
	return &entities.UserChannelDetail{
		Model:         &entities.Model{ID: e.ID, IsActive: e.IsActive},
		UserID:        e.UserID,
		UserChannelID: e.ChannelId,
	}, nil
}

// UserChannelDetails implements Mapper.
func (m *mapper) UserChannelDetails(items []requestModel.UserChannelDetail) ([]entities.UserChannelDetail, error) {
	mappedItems := make([]entities.UserChannelDetail, len(items))
	for i, item := range items {
		mapped, err := m.UserChannelDetail(item)
		if err != nil {
			return nil, err
		}
		mappedItems[i] = *mapped
	}
	return mappedItems, nil
}

func (m *mapper) Customer(e requestModel.Customer) (*entities.Customer, error) {
	return &entities.Customer{
		Model:          &entities.Model{ID: e.ID, IsActive: e.IsActive},
		FirstName:      e.FirstName,
		LastName:       e.LastName,
		Email:          e.Email,
		PhoneNumber:    e.PhoneNumber,
		WhatsappNumber: e.WhatsappNumber,
		Address:        e.Address,
	}, nil
}

func (m *mapper) Person(e requestModel.Person) (*entities.Person, error) {
	var customerId uint
	if e.CustomerId != nil {
		customerId = *e.CustomerId
	}

	return &entities.Person{
		Model:      &entities.Model{ID: e.ID, IsActive: e.IsActive},
		FirstName:  e.FirstName,
		LastName:   e.LastName,
		Gender:     entities.Gender(e.Gender),
		Age:        e.Age,
		CustomerId: customerId,
	}, nil
}

func (m *mapper) DressType(e requestModel.DressType) (*entities.DressType, error) {
	return &entities.DressType{
		Model:        &entities.Model{ID: e.ID, IsActive: e.IsActive},
		Name:         e.Name,
		Description:  e.Description,
		Measurements: e.Measurements,
	}, nil
}

func (m *mapper) Measurement(e requestModel.Measurement) (*entities.Measurement, error) {
	// Convert values JSON
	var values entitiy_types.JSON
	if len(e.Values) > 0 {
		values = entitiy_types.JSON(e.Values)
	}

	var personId uint
	if e.PersonId != nil {
		personId = *e.PersonId
	}

	var dressTypeId uint
	if e.DressTypeId != nil {
		dressTypeId = *e.DressTypeId
	}

	return &entities.Measurement{
		Model:       &entities.Model{ID: e.ID, IsActive: e.IsActive},
		Value:       values,
		PersonId:    personId,
		DressTypeId: dressTypeId,
		TakenById:   e.TakenById,
	}, nil
}

func (m *mapper) Order(e requestModel.Order) (*entities.Order, error) {
	orderItems, err := m.OrderItems(e.OrderItems)
	if err != nil {
		return nil, err
	}

	var expectedDeliveryDate *time.Time
	if e.ExpectedDeliveryDate != nil {
		date, err := util.GenerateDateTimeFromString(e.ExpectedDeliveryDate)
		if err != nil {
			return nil, err
		}
		expectedDeliveryDate = date
	}

	var deliveredDate *time.Time
	if e.DeliveredDate != nil {
		date, err := util.GenerateDateTimeFromString(e.DeliveredDate)
		if err != nil {
			return nil, err
		}
		deliveredDate = date
	}

	return &entities.Order{
		Model:                &entities.Model{ID: e.ID, IsActive: e.IsActive},
		Status:               entities.OrderStatus(e.Status),
		Notes:                e.Notes,
		AdditionalCharges:    e.AdditionalCharges,
		ExpectedDeliveryDate: expectedDeliveryDate,
		DeliveredDate:        deliveredDate,
		CustomerId:           e.CustomerId,
		OrderTakenById:       e.OrderTakenById,
		OrderItems:           orderItems,
	}, nil
}

func (m *mapper) OrderItem(e requestModel.OrderItem) (*entities.OrderItem, error) {
	var expectedDeliveryDate *time.Time
	if e.ExpectedDeliveryDate != nil {
		date, err := util.GenerateDateTimeFromString(e.ExpectedDeliveryDate)
		if err != nil {
			return nil, err
		}
		expectedDeliveryDate = date
	}

	var deliveredDate *time.Time
	if e.DeliveredDate != nil {
		date, err := util.GenerateDateTimeFromString(e.DeliveredDate)
		if err != nil {
			return nil, err
		}
		deliveredDate = date
	}

	return &entities.OrderItem{
		Model:                &entities.Model{ID: e.ID, IsActive: e.IsActive},
		Description:          e.Description,
		Quantity:             e.Quantity,
		Price:                e.Price,
		Total:                e.Total,
		AdditionalCharges:    e.AdditionalCharges,
		ExpectedDeliveryDate: expectedDeliveryDate,
		DeliveredDate:        deliveredDate,
		PersonId:             e.PersonId,
		MeasurementId:        e.MeasurementId,
		OrderId:              e.OrderId,
	}, nil
}

func (m *mapper) OrderItems(items []requestModel.OrderItem) ([]entities.OrderItem, error) {
	mappedItems := make([]entities.OrderItem, len(items))
	for i, item := range items {
		mapped, err := m.OrderItem(item)
		if err != nil {
			return nil, err
		}
		mappedItems[i] = *mapped
	}
	return mappedItems, nil
}

func (m *mapper) OrderHistory(e requestModel.OrderHistory) (*entities.OrderHistory, error) {
	var expectedDeliveryDate *time.Time
	if e.ExpectedDeliveryDate != nil {
		date, err := util.GenerateDateTimeFromString(e.ExpectedDeliveryDate)
		if err != nil {
			return nil, err
		}
		expectedDeliveryDate = date
	}

	var deliveredDate *time.Time
	if e.DeliveredDate != nil {
		date, err := util.GenerateDateTimeFromString(e.DeliveredDate)
		if err != nil {
			return nil, err
		}
		deliveredDate = date
	}

	var status *entities.OrderStatus
	if e.Status != nil {
		s := entities.OrderStatus(*e.Status)
		status = &s
	}

	var orderItemData entitiy_types.JSON
	if e.OrderItemData != "" {
		orderItemData = entitiy_types.JSON(e.OrderItemData)
	}

	performedAt := util.GetLocalTime()
	if e.PerformedAt != "" {
		date, err := util.GenerateDateTimeFromString(&e.PerformedAt)
		if err != nil {
			return nil, err
		}
		performedAt = *date
	}

	return &entities.OrderHistory{
		Model:                &entities.Model{ID: e.ID, IsActive: e.IsActive},
		Action:               entities.OrderHistoryAction(e.Action),
		ChangedFields:        e.ChangedFields,
		Status:               status,
		ExpectedDeliveryDate: expectedDeliveryDate,
		DeliveredDate:        deliveredDate,
		OrderItemId:          e.OrderItemId,
		OrderItemData:        &orderItemData,
		OrderId:              e.OrderId,
		PerformedAt:          performedAt,
		PerformedById:        e.PerformedById,
	}, nil
}

func (m *mapper) MeasurementHistory(e requestModel.MeasurementHistory) (*entities.MeasurementHistory, error) {
	var oldValues entitiy_types.JSON
	if e.OldValues != "" {
		oldValues = entitiy_types.JSON(e.OldValues)
	}

	performedAt := util.GetLocalTime()
	if e.PerformedAt != "" {
		date, err := util.GenerateDateTimeFromString(&e.PerformedAt)
		if err != nil {
			return nil, err
		}
		performedAt = *date
	}

	return &entities.MeasurementHistory{
		Model:         &entities.Model{ID: e.ID, IsActive: e.IsActive},
		Action:        entities.MeasurementHistoryAction(e.Action),
		OldValues:     oldValues,
		MeasurementId: e.MeasurementId,
		PerformedAt:   performedAt,
		PerformedById: e.PerformedById,
	}, nil
}

func (m *mapper) ExpenseTracker(e requestModel.ExpenseTracker) (*entities.Expense, error) {
	var purchaseDate *time.Time
	if e.PurchaseDate != nil {
		date, err := util.GenerateDateTimeFromString(e.PurchaseDate)
		if err != nil {
			return nil, err
		}
		purchaseDate = date
	}

	var isActive bool = true
	if e.IsActive != nil {
		isActive = *e.IsActive
	}

	expenseDetails, err := m.ExpenseDetails(e.ExpenseDetails)
	if err != nil {
		return nil, err
	}

	return &entities.Expense{
		Model:          &entities.Model{ID: e.ID, IsActive: isActive},
		PurchaseDate:   purchaseDate,
		BillNumber:     e.BillNumber,
		CompanyName:    e.CompanyName,
		Material:       e.Material,
		Price:          e.Price,
		Balance:        e.Balance,
		Location:       e.Location,
		Notes:          e.Notes,
		ExpenseDetails: expenseDetails,
	}, nil
}

func (m *mapper) ExpenseDetail(e requestModel.ExpenseDetail) (*entities.ExpenseDetail, error) {
	var isActive bool = true
	if e.IsActive != nil {
		isActive = *e.IsActive
	}
	return &entities.ExpenseDetail{
		Model:     &entities.Model{ID: e.ID, IsActive: isActive},
		Source:    e.Source,
		Price:     e.Price,
		ExpenseId: e.ExpenseId,
	}, nil
}

func (m *mapper) ExpenseDetails(items []requestModel.ExpenseDetail) ([]entities.ExpenseDetail, error) {
	if len(items) == 0 {
		return nil, nil
	}
	mapped := make([]entities.ExpenseDetail, 0, len(items))
	for _, item := range items {
		ent, err := m.ExpenseDetail(item)
		if err != nil {
			return nil, err
		}
		mapped = append(mapped, *ent)
	}
	return mapped, nil
}

func (m *mapper) Task(e requestModel.Task) (*entities.Task, error) {
	var dueDate *time.Time
	if e.DueDate != nil {
		date, err := util.GenerateDateTimeFromString(e.DueDate)
		if err != nil {
			return nil, err
		}
		dueDate = date
	}

	var reminderDate *time.Time
	if e.ReminderDate != nil {
		date, err := util.GenerateDateTimeFromString(e.ReminderDate)
		if err != nil {
			return nil, err
		}
		reminderDate = date
	}

	var completedAt *time.Time
	if e.CompletedAt != nil {
		date, err := util.GenerateDateTimeFromString(e.CompletedAt)
		if err != nil {
			return nil, err
		}
		completedAt = date
	}

	var isActive bool = true
	if e.IsActive != nil {
		isActive = *e.IsActive
	}

	if e.Status == "" {
		e.Status = string(entities.TaskStatusPending)
	}
	return &entities.Task{
		Model:        &entities.Model{ID: e.ID, IsActive: isActive},
		Title:        e.Title,
		Description:  e.Description,
		Status:       entities.TaskStatus(e.Status),
		Priority:     e.Priority,
		DueDate:      dueDate,
		ReminderDate: reminderDate,
		CompletedAt:  completedAt,
		AssignedToId: e.AssignedToId,
	}, nil
}

func (m *mapper) Category(e requestModel.Category) (*entities.Category, error) {
	return &entities.Category{
		Model: &entities.Model{ID: e.ID, IsActive: e.IsActive},
		Name:  e.Name,
	}, nil
}

func (m *mapper) Product(e requestModel.Product) (*entities.Product, error) {
	return &entities.Product{
		Model:        &entities.Model{ID: e.ID, IsActive: e.IsActive},
		Name:         e.Name,
		SKU:          e.SKU,
		CategoryId:   e.CategoryId,
		Description:  e.Description,
		CostPrice:    e.CostPrice,
		SellingPrice: e.SellingPrice,
	}, nil
}

func (m *mapper) Inventory(e requestModel.Inventory) (*entities.Inventory, error) {
	return &entities.Inventory{
		Model:             &entities.Model{ID: e.ID, IsActive: e.IsActive},
		ProductId:         e.ProductId,
		LowStockThreshold: e.LowStockThreshold,
	}, nil
}

func (m *mapper) InventoryLog(e requestModel.InventoryLog) (*entities.InventoryLog, error) {
	var loggedAt time.Time
	if e.LoggedAt != "" {
		date, err := util.GenerateDateTimeFromString(&e.LoggedAt)
		if err != nil {
			return nil, err
		}
		loggedAt = *date
	} else {
		loggedAt = util.GetLocalTime()
	}

	return &entities.InventoryLog{
		Model:      &entities.Model{ID: e.ID, IsActive: e.IsActive},
		ProductId:  e.ProductId,
		ChangeType: entities.InventoryLogChangeType(e.ChangeType),
		Quantity:   e.Quantity,
		Reason:     e.Reason,
		Notes:      e.Notes,
		LoggedAt:   loggedAt,
	}, nil
}
