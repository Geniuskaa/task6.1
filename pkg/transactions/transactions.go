package transactions

import (
	"sync"
)

type Transaction struct {
	Id               int64
	SumOfTransaction int64
	Date             int64 // в формате Unix Timeship
	MCC              string
	Status           string
}

type Card struct {
	Id int64
	Issuer string
	Balance int64
	Currency string
	Number string
	Icon string
	Transactions []*Transaction
}

func SortTransactions(slice[]*Transaction, Id int64) map[string]int64{ // Нет фильтрации по Id пользователя
	mapOfTransactions := make(map[string]int64)
	var sumOfSupermarkets, sumOfRestaurants, sumOfMobile, sumOfFinance, sumOfOtherThings int64
	for m := range slice {
		if slice[m].Id == Id { // Проверка по Id
			switch slice[m].MCC {
			case "5401": // Супермаркеты
				sumOfSupermarkets += slice[m].SumOfTransaction
			case "5490": // Рестораны
				sumOfRestaurants += slice[m].SumOfTransaction
			case "5500": // Мобильная связь
				sumOfMobile += slice[m].SumOfTransaction
			case "5455": // Финансы
				sumOfFinance += slice[m].SumOfTransaction
			default: // Остальное
				sumOfOtherThings += slice[m].SumOfTransaction
			}
		}
	}
	mapOfTransactions["Супермаркеты"] = sumOfSupermarkets
	mapOfTransactions["Рестораны"] = sumOfRestaurants
	mapOfTransactions["Мобильная связь"] = sumOfMobile
	mapOfTransactions["Финансы"] = sumOfFinance
	mapOfTransactions["Остальное"] = sumOfOtherThings

	return mapOfTransactions
}

func SortTransactionsByMutex(slice[]*Transaction, Id int64, goroutines int) map[string]int64 {
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	result := make(map[string]int64)

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		partSize := 2
		part := slice[i * partSize : (i+1) * partSize]
		go func() {
			m := SortTransactions(part, Id)

			mu.Lock()
			for key, value:= range m {
				switch key {
				case "Мобильная связь":
					result["Мобильная связь"] += value
				case "Рестораны":
					result["Рестораны"] += value
				case "Супермаркеты":
					result["Супермаркеты"] += value
				case "Финансы":
					result["Финансы"] += value
				default:
					result["Остальное"] += value
				}
			}

			mu.Unlock()
			wg.Done()
		}()

	}
	wg.Wait()

	return result
}

func SortTransactionsByChanels(slice[]*Transaction, Id int64, partsCount int) map[string]int64 {
	result := make(map[string]int64)
	ch := make(chan map[string]int64)

	for i := 0; i < partsCount; i++ { // TODO здесь ваши условия разделения
		partSize := 2
		part := slice[i * partSize : (i+1) * partSize]
		go func(ch chan <- map[string]int64) {
			ch <- SortTransactions(part, Id)
		}(ch)
	}

	finished := 0
	for valueMap:= range ch {
		for key, value := range valueMap {
			switch key {
			case "Мобильная связь":
				result["Мобильная связь"] += value
			case "Рестораны":
				result["Рестораны"] += value
			case "Супермаркеты":
				result["Супермаркеты"] += value
			case "Финансы":
				result["Финансы"] += value
			default:
				result["Остальное"] += value
			}
		}
		finished++
		if finished == partsCount {
			break
		}
	}

	return result
}

func SortTransactionsByOtherMutex(slice[]*Transaction, Id int64, goroutines int) map[string]int64 {
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	result := make(map[string]int64)
	var sumOfSupermarkets, sumOfRestaurants, sumOfMobile, sumOfFinance, sumOfOtherThings int64

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		partSize := 2
		part := slice[i * partSize : (i+1) * partSize]
		go func() {
			for _, t := range part {
				if t.Id == Id {
					mu.Lock()

					switch t.MCC {
					case "5401": // Супермаркеты
						sumOfSupermarkets += t.SumOfTransaction
					case "5490": // Рестораны
						sumOfRestaurants += t.SumOfTransaction
					case "5500": // Мобильная связь
						sumOfMobile += t.SumOfTransaction
					case "5455": // Финансы
						sumOfFinance += t.SumOfTransaction
					default: // Остальное
						sumOfOtherThings += t.SumOfTransaction
					}

					result["Супермаркеты"] = sumOfSupermarkets
					result["Рестораны"] = sumOfRestaurants
					result["Мобильная связь"] = sumOfMobile
					result["Финансы"] = sumOfFinance
					result["Остальное"] = sumOfOtherThings
					mu.Unlock()
				}


			}
			wg.Done()
		}()
	}
	wg.Wait()

	return result
}
