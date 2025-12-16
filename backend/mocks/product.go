package types

type Product struct {
	Id    string `dynamodbv:"id" json:"id"`
	Name  string `dynamodbv:"name" json:"name"`
	Price uint   `dynamodbv:"price" json:"price"`
}

type ProductRange struct {
	Products []Product `json:"products"`
	Next     *string   `json:"next,omitempty"`
}