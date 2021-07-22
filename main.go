//! KROSSMAN (version 0.2)
//@ Author ivkis (t.me/iivkis)
//# Support Shop(-s): Nike.

package main

import (
	"krossman/parsers/nikeParser"
	"time"

	"github.com/gofiber/fiber/v2"
)

const addr = "localhost:8027"

/*
	---Shop links---
	Nike: http://export.admitad.com/ru/webmaster/websites/213277/products/export_adv_products/?user=millstone&code=56d7cfafe9&feed_id=18062&format=xml
*/
func main() {
	//init parsers
	nike := nikeParser.New(&nikeParser.Config{
		SaveAs:     "nike",
		Address:    "http://export.admitad.com/ru/webmaster/websites/213277/products/export_adv_products",
		User:       "millstone",
		Code:       "56d7cfafe9",
		FeedID:     18062,
		CategoryID: "c71a18083d",
	})

	//run parsers
	go func() {
		for {
			nike.Parse()
			time.Sleep(time.Minute * 50)
		}
	}()

	//server
	app := fiber.New()
	app.Static("/shops", "./static/shops")
	if err := app.Listen(addr); err != nil {
		panic(err)
	}
}
