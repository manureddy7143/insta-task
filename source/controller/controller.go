package controller

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"net/http"

	"github.com/go-chassis/openlog"
	dto "github.com/manureddy7143/insta-task/source/dto"

	"github.com/gin-gonic/gin"
)

type Transaction struct{}

var TransactionsSorage []map[string]interface{}

type AddTransactionInput struct {
	Data map[string]interface{}
}

// LoginPath - URL Path for Login
const Transactions = "/transactions"
const GetStat = "/statstics"
const Delete = "/delete"

const secretKey = "secret"

// PostTransactions will Store valid Tranasactins
func (a Transaction) PostTransactions(c *gin.Context) {
	openlog.Info("Got a request to add Transaction")

	data := make(map[string]interface{})
	json.NewDecoder(c.Request.Body).Decode(&data)
	input := AddTransactionInput{Data: data}
	res := AddTransaction(input)
	c.JSON(http.StatusAccepted, res)
}

// FetchAllTransactions will list Trasactions Statistics
func (a Transaction) GetStatstics(c *gin.Context) {
	res := GetAllStastics()
	c.JSON(http.StatusAccepted, res)
}

// DeleteAllTransactions delete all transactions
func (a Transaction) DeleteAllTransactions(c *gin.Context) {
	openlog.Info("Got a request to Delete All Transactions")
	res := DeleteAllTransactions()
	c.JSON(http.StatusAccepted, res)
}

func AddTransaction(input AddTransactionInput) dto.Response {
	timestamp := input.Data["timestamp"].(string)
	transactiontime, err := time.Parse("2006-01-02T15:04:05.999Z", timestamp)
	if err != nil {
		fmt.Println(err)
		return dto.Response{Msg: "Internal Server Error", Data: nil, Status: 500}
	}
	timeNow := time.Now()
	diff := timeNow.Sub(transactiontime)

	if diff > time.Duration(time.Second*60) {
		fmt.Println("duration :", diff)
		return dto.Response{Msg: "Timestamp shoud be with 60 seconds", Data: input.Data, Status: 204}
	} else if diff < 0 {
		return dto.Response{Msg: "Timestamp shoud not be in future", Data: input.Data, Status: 422}
	}

	TransactionsSorage = append(TransactionsSorage, input.Data)
	return dto.Response{Msg: "Transaction added successfully", Data: input.Data, Status: 201}
}

func GetAllStastics() dto.Response {

	stats, err := GetStatstics(time.Now())
	if err != nil {
		fmt.Println(err)
		return dto.Response{Msg: "Internal Server Error", Data: nil, Status: 500}
	}
	return dto.Response{Msg: "All Transaction Fetched successfully", Data: stats, Status: 200}
}

func DeleteAllTransactions() dto.Response {
	TransactionsSorage = nil
	return dto.Response{Msg: "All Transactions deleted successfully", Data: nil, Status: 204}
}
func GetStatstics(timenow time.Time) (map[string]interface{}, error) {
	var max, min, average, count, sum float64

	if len(TransactionsSorage) > 0 {
		min = math.MaxFloat64
	}
	for _, transaction := range TransactionsSorage {
		timestamp := transaction["timestamp"].(string)
		transactiontime, err := time.Parse("2006-01-02T15:04:05.999Z", timestamp)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		diff := timenow.Sub(transactiontime)
		if diff < time.Duration(time.Second*60) {
			amount := transaction["amount"].(string)
			cost, err := strconv.ParseFloat(amount, 64)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			sum += cost
			count++
			average = sum / count
			if cost > max {
				max = cost
			}
			if cost < min {
				min = cost
			}
		}
	}
	return map[string]interface{}{
		"sum":     sum,
		"average": average,
		"count":   count,
		"min":     min,
		"max":     max,
	}, nil
}
