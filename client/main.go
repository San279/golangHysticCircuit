package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/api", api)
	app.Listen(":5001")
}

func init() {
	hystrix.ConfigureCommand("api", hystrix.CommandConfig{ //middleware config
		Timeout:                500, //timeout delay
		RequestVolumeThreshold: 1,
		ErrorPercentThreshold:  100, //the timeout request threshold
		SleepWindow:            10000,
	})

	hystrixStreamHandler := *hystrix.NewStreamHandler() //fashboard
	hystrixStreamHandler.Start()

}

func api(c *fiber.Ctx) error {
	hystrix.Go("api", func() error { //monitor given api
		res, err := http.Get("http://localhost:5000/api")
		if err != nil {
			return err
		}
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		msg := string(data)
		fmt.Println(msg)
		return nil
	}, func(err error) error {
		fmt.Println(err)
		return err
	})
	return nil
}
