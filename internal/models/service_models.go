package models

type Service struct {
	Base
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Duration    int     `json:"duration"`
	Price       float64 `json:"price"`
}

type ServiceCreateRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Duration    int     `json:"duration"`
	Price       float64 `json:"price"`
}

type ServiceUpdateRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Category    *string  `json:"category"`
	Duration    *int     `json:"duration"`
	Price       *float64 `json:"price"`
}
