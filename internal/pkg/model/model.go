/* Модель пользовательских данных для хранения в базе */

package model

type User struct {
	Id          uint   `json:"id"`
	Name        string `json:"name,omitempty"`
	Surname     string `json:"surname,omitempty"`
	MiddleName  string `json:"fatherName,omitempty"`
	DOB         string `json:"date,omitempty"`
	Address     string `json:"address,omitempty"`
	Email       string `json:"email,omitempty"`
	Password    string `json:"password"`
	Login       string `json:"login"`
	ITN         string `json:"itn,omitempty"`
	ImageFolder string `json:"image_folder,omitempty"`
}

type SchoolOrder struct {
	Name     string `json:"name,omitempty"`
	Surname  string `json:"surname,omitempty"`
	Age      int    `json:"age,omitempty"`
	Category string `json:"category,omitempty"`
	Shift    string `json:"shift,omitempty"`
	Day      string `json:"day,omitempty"`
}
