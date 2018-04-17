package jsonsideload

import (
	"encoding/json"
)

type Session struct {
	ID      float64  `jsonsideload:"attr,id"`
	Account *Account `jsonsideload:"hasone,accounts,account_id"`
	User    *User    `jsonsideload:"hasone,users,user_id"`
}

type Account struct {
	ID         float64    `jsonsideload:"attr,id"`
	FullDomain string     `jsonsideload:"attr,full_domain"`
	Name       string     `jsonsideload:"attr,name"`
	OfferForm  *OfferForm `jsonsideload:"hasone,offer_forms,offer_form_id"`
	Users      []*User    `jsonsideload:"hasmany,users,user_ids"`
	//Features     []string      `jsonsideload:"attr,features"`
	//Subscription *Subscription `jsonsideload:"hasone,subscription"`
}

type Subscription struct {
	ID     json.Number `jsonsideload:"attr,id"`
	Status int         `jsonsideload:"attr,status"`
	Plan   *Plan       `jsonsideload:"hasone,plan"`
}

type Plan struct {
	ID   json.Number `jsonsideload:"attr,id"`
	Name string      `jsonsideload:"attr,name"`
}

type OfferForm struct {
	ID   float64 `jsonsideload:"attr,id"`
	Name string  `jsonsideload:"attr,name"`
}

type User struct {
	Token       string  `jsonsideload:"attr,token"`
	EmployeeID  string  `jsonsideload:"attr,employee_id"`
	FirstName   string  `jsonsideload:"attr,first_name"`
	MiddleName  string  `jsonsideload:"attr,middle_name"`
	LastName    string  `jsonsideload:"attr,last_name"`
	Name        string  `jsonsideload:"attr,name"`
	ID          float64 `jsonsideload:"attr,id"`
	Email       string  `jsonsideload:"attr,email"`
	Designation string  `jsonsideload:"attr,designation"`
	// Avatar      *Avatar  `jsonsideload:"hasone,avatar"`
	//Abilities []string `jsonsideload:"attr,abilities"`
	//UserEmails    []UserEmail `jsonsideload:"user_emails"`
	IrisJWTSecret string `jsonsideload:"attr,iris_jwt_secret"`
}

type UserEmail struct {
	ID           json.Number `jsonsideload:"attr,id"`
	Email        string      `jsonsideload:"attr,email"`
	Confirmed    bool        `jsonsideload:"attr,confirmed"`
	PrimaryEmail bool        `jsonsideload:"attr,primary_email"`
	UserID       json.Number `jsonsideload:"attr,user_id"`
}

type Avatar struct {
	ID           json.Number  `jsonsideload:"attr,id"`
	ExpiringURLs *ExpiringURL `jsonsideload:"hasone,expiring_urls"`
}

type ExpiringURL struct {
	Medium   string `jsonsideload:"attr,medium"`
	Original string `jsonsideload:"attr,original"`
	Thumb    string `jsonsideload:"attr,thumb"`
}
