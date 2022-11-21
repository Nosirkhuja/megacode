package controllers

import (
	"MegaCode/database"
	"MegaCode/internal/pkg/model"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"time"
)

const SecretKey = "secret"

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	user := model.User{
		Name:       data["name"],
		Email:      data["email"],
		Surname:    data["surname"],
		MiddleName: data["fathername"],
		DOB:        data["date"],
		Address:    data["address"],
		Login:      data["login"],
		Password:   string(password),
	}

	var exist model.User

	database.DB.Where("email = ?", data["email"]).First(&exist)
	if exist.Email == "" {
		database.DB.Create(&user)
		return c.JSON(fiber.Map{"message": "Успешно!"})
	}

	return c.JSON(fiber.Map{"message": "Пользователь с таким email уже зарегистрирован"})
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user model.User

	database.DB.Where("email = ?", data["email"]).First(&user)

	if user.Id == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "Пользователь не найден",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Неправильный пароль",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), //1 day
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Не удалось зайти",
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "Успешно!",
	})
}

func User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Неопознанный пользователь",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user model.User

	database.DB.Where("id = ?", claims.Issuer).First(&user)

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "Успешно",
	})
}

func Order(c *fiber.Ctx) error {
	/*err := Register(c)
	if err != nil {
		return err
	}*/

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	user := model.User{
		Name:       data["firstName"],
		Email:      data["email"],
		Surname:    data["secondName"],
		MiddleName: data["Lastname"],
		DOB:        data["birthday"],
		Address:    data["address"],
		Login:      data["login"],
		Password:   string(password),
	}

	err := SendToSchool(user)
	if err != nil {
		return c.JSON(fiber.Map{"message": "Не удалось записать в автошколу("})
	}
	return c.JSON(fiber.Map{"message": "Успешно!"})

}

func SendToSchool(user model.User) error {
	httpposturl := "http://localhost:3001/messages"
	order := new(model.SchoolOrder)
	order.Name = user.Name
	order.Surname = user.Surname
	bd := user.DOB
	age := 0
	if len(bd) > 4 {
		bd = bd[:4]
		var err error
		age, err = strconv.Atoi(bd)
		if err != nil {
			age = 19
		}
		age = 2022 - age
	}
	order.Age = age
	jsonData, err := json.Marshal(order)
	if err != nil {
		return err
	}
	request, err := http.NewRequest("POST", httpposturl, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	fmt.Println("данные успешно отправлено в автошколу\n", order)
	return nil
}

func OrderCollects(c *fiber.Ctx) error {
	return nil
}
