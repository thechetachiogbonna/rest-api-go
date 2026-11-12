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
	FirstName string `json:"first_name" validate:"min=5,max=20"`
	LastName  string `json:"last_name" validate:"min=5,max=20"`
	Email     string `json:"email" validate:"email"`
}

type CreateCategoryPayload struct {
	Name string `json:"name" validate:"required,min=3,max=50"`
	Slug string `json:"slug" validate:"required,min=3,max=50"`
}
