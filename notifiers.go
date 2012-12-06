package main

import (
	"fmt"
)

func sendNotifier(url Url) {
	// In the future, it'd be cool to make this good.
	subject := fmt.Sprintf("DOWN: Failure for %s", url.Url)
	body := fmt.Sprintf("This is a notification that %s failed a check, with status code %d on %s.", url.Url, url.LastCheckStatus, url.LastCheck)
	mail := Mail{Recipient: "jackpearkes@gmail.com", Subject: subject, Body: body, From: "postmaster@jack.mailgun.org"}
	sendMail(mail)
}
