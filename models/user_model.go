package models

type User struct {
	UserId    int    `json:"userid,omitempty"`
	FirstName string `json:"firstname,omitempty" validate:"required"`
	LastName  string `json:"lastname,omitempty" validate:"required"`
	Age       int    `json:"age,omitempty" validate:"gte=0,lte=130"`
	Gender    string `json:"gender,omitempty" validate:"oneof=male female others"`
}
