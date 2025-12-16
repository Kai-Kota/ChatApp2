package handlers

type APIGatewayV2Handler struct {
	products *domain.Products
}

func NewAPIGatewayV2Hander(d *domain.Products) *APIGatewayV2Handler {
	return &APIGatewayV2Handler{
		products: d,
	}
}
