package main

import (
	"github.com/Geniuskaa/task6.1/pkg/transactions"
	"fmt"
)

func main() {

	const users = 10_00
	const transactionsPerUser = 10_00

	/*b := []*transactions.Transaction{{Id: 1, SumOfTransaction: 135_352_00, MCC: "5401"},
		{Id: 1, SumOfTransaction: 135_00, MCC: "5455"},
		{Id: 3, SumOfTransaction: 1_362_00, MCC: "5455"},
		{Id: 1, SumOfTransaction: 35_352_00, MCC: ""},
		{Id: 1, SumOfTransaction: 85_352_00, MCC: "5455"},
		{Id: 1, SumOfTransaction: 9_352_00, MCC: "5490"},
		{Id: 5, SumOfTransaction: 21_362_00, MCC: "5401"},
		{Id: 1, SumOfTransaction: 5_352_00, MCC: "5401"},}*/
	b := make([]*transactions.Transaction, users*transactionsPerUser)
	for index := range b {
		switch index % 100 {
		case 0:
			b[index] = &transactions.Transaction{Id: 1, SumOfTransaction: 5_352_00, MCC: ""}// Например, каждая 100-ая транзакция в банке от нашего юзера в категории такой-то
		case 20:
			b[index] = &transactions.Transaction{Id: 3, SumOfTransaction: 1_362_00, MCC: "5455"}// Например, каждая 120-ая транзакция в банке от нашего юзера в категории такой-то
		case 40:
			b[index] = &transactions.Transaction{Id: 1, SumOfTransaction: 9_352_00, MCC: "5490"}
		default:
			b[index] = &transactions.Transaction{Id: 1, SumOfTransaction: 5_352_00, MCC: "5401"}// Транзакции других юзеров, нужны для "общей" массы
		}
	}

	tinkoff := transactions.Card{
		Id:           01,
		Issuer:       "VISA",
		Balance:      194_125_00,
		Currency:     "RUB",
		Number:       "4205 1840 2045 2902",
		Icon:         "",
		Transactions: b,
	}


	fmt.Println(transactions.SortTransactions(tinkoff.Transactions, 1))
	fmt.Println(transactions.SortTransactionsByMutex(tinkoff.Transactions,1,500000))
	fmt.Println(transactions.SortTransactionsByChanels(tinkoff.Transactions,1,500000))
	fmt.Println(transactions.SortTransactionsByOtherMutex(tinkoff.Transactions,1,500000))


}


