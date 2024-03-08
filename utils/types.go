package utils

type TCoordinates struct {
	X *int `json:"x,omitempty" validate:"required,gte=0,lte=2"`
	Y *int `json:"y,omitempty" validate:"required,gte=0,lte=2"`
}

type TFieldsData struct {
	Type *string `json:"type,omitempty" validate:"required"`
	TCoordinates
}
