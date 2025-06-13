package models

type Product struct {
	ID       int     `mapstructure:"id" json:"id,omitempty"`
	Name     string  `mapstructure:"name" json:"name"`
	Quantity float64 `mapstructure:"quantity" json:"quantity"`
	UserId   int     `mapstructure:"userId" json:"userId"`
	Namehash string  `mapstructure:"namehash" json:"namehash,omitempty"`
}
