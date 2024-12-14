package domain

type Address struct {
	Homecode    string `json:"home_code" bson:"home_code"`
	Street      string `json:"street" bson:"street"`
	District    string `json:"district" bson:"district"`
	Province    string `json:"province" bson:"province"`
	FullAddress string `json:"full_address" bson:"full_address"`
}
