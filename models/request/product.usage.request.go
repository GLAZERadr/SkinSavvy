package request

import (
	"time"
)

type UserProductUsage struct {
	ID					string 		`json:"Id"`
	UserID				string 		`json:"userId"`
	ProductBrand		string 		`json:"productBrand"`
	ProductName			string 		`json:"productName"`
	ProductImage		string 		`json:"productImage"`
	CreatedAt			time.Time	`json:"created_at"`
}