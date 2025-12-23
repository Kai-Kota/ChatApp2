package types

type Auth struct {
	UserName string `dynamodbav:"user_name" json:"user_name"`
	Password string `dynamodbav:"password" json:"password"`
}
