package types

type RegisterPayload struct {
	FirstName string `json:"first_name" validate:"required,min=5,max=20"`
	LastName  string `json:"last_name" validate:"required,min=5,max=20"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
}

type LoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UpdateUserByIDPayload struct {
	FirstName string `json:"first_name" validate:"omitempty,min=5,max=20"`
	LastName  string `json:"last_name" validate:"omitempty,min=5,max=20"`
	Email     string `json:"email" validate:"omitempty,email"`
}

type CreateCategoryPayload struct {
	Name string `json:"name" validate:"required,min=3,max=50"`
	Slug string `json:"slug" validate:"required,min=3,max=50"`
}

type CreateProductPayload struct {
	CategoryID    string `json:"category_id" validate:"required"`
	Name          string `json:"name" validate:"required,min=3,max=50"`
	Slug          string `json:"slug" validate:"required,min=3,max=50"`
	Description   string `json:"description" validate:"required,min=10,max=200"`
	Price         int    `json:"price" validate:"required"`
	StockQuantity int    `json:"stock_quantity" validate:"required"`
}

type UpdateProductPayload struct {
	CategoryID    string `json:"category_id" validate:"required"`
	Name          string `json:"name" validate:"omitempty,min=3,max=50"`
	Slug          string `json:"slug" validate:"omitempty,min=3,max=50"`
	Description   string `json:"description" validate:"omitempty,min=10,max=200"`
	Price         int    `json:"price" validate:"omitempty"`
	StockQuantity int    `json:"stock_quantity" validate:"omitempty"`
}

type UpdateProductStockQuantityPayload struct {
	StockQuantity int `json:"stock_quantity" validate:"required"`
}
