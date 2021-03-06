package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// type Blockchain struct{}

type TicketProperty struct {
	Fifteen float32 `json:"15m"`
	Last    float32 `json:"last"`
	Buy     float32 `json:"buy"`
	Sell    float32 `json:"sell"`
	Symbol  string  `json:"symbol"`
}

func GetPrices() (map[string]TicketProperty, error) {

	resp, err := http.Get("https://blockchain.info/ticker")
	if err != nil {
		return nil, err
	}
	var tickets = make(map[string]TicketProperty)
	if err := json.NewDecoder(resp.Body).Decode(&tickets); err != nil {
		return nil, err
	}
	return tickets, nil
}

func GetAddress() (string, error) {
	req, err := http.NewRequest("POST", "https://www.blockonomics.co/api/new_address", nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", key))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	var address = make(map[string]string)
	err = json.NewDecoder(resp.Body).Decode(&address)
	return address["address"], err
}

func handleDeposit(w http.ResponseWriter, r *http.Request) {
	fmt.Println(channelId)
	status := r.URL.Query().Get("status")
	address := r.URL.Query().Get("addr")
	value := r.URL.Query().Get("value")
	txid := r.URL.Query().Get("txid")
	switch status {
	case "0":
		u, err := GetUserByDepositAddress(address)
		if err != nil {
			logrus.Error(err)
		}
		logrus.Info(
			fmt.Sprintf("Se ha detectado una tranccion "+
				"no confirmada por parte del usuario %s.\n"+
				"Transaction ID:%s", u.GetName(), txid))
		break
	case "1":
		u, err := GetUserByDepositAddress(address)
		if err != nil {
			logrus.Error(err)
		}
		logrus.Info(
			fmt.Sprintf("Se ha detectado una tranccion "+
				"parcialmente confirmada por parte del usuario %s.\n"+
				"Transaction ID:%s", u.GetName(), txid))
		break
	case "2":
		u, err := GetUserByDepositAddress(address)
		if err != nil {
			logrus.Error(err)
			return
		}

		if err := AddInvestToUser(value, u.GetID()); err != nil {
			logrus.Error(err)
			return
		}

		if err := AddTransactionToUser(u.GetID(), true, txid, value); err != nil {
			logrus.Error(err)
			return
		}

		intValue, err := strconv.ParseInt(value, 10, 0)
		if err != nil {
			logrus.Error(err)
			return
		}
		invert := decimal.New(intValue, -Exponent)
		msg := tgbotapi.NewMessageToChannel(channelId,
			fmt.Sprintf("Nueva inversion:\n "+
				"%s ha invertido %s BTC!\n"+
				"Transaction ID:\n"+
				"<a href=\"https://blockchain.info/tx/%s\">%s</a>", u.GetName(),
				invert.StringFixed(Exponent), txid, txid))

		msg.ParseMode = "html"
		if _, err := bot.Send(msg); err != nil {
			logrus.Error(err)
		}
		break
	}
}
