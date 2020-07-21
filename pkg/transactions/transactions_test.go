package transactions

import (
	"reflect"
	"testing"
)

func makeTransactions() ([]*Transaction) {
	const users = 10_00
	const transactionsPerUser = 10_00
	transactions := make([]*Transaction, users * transactionsPerUser)
	for index := range transactions {
		switch index % 100 {
		case 0:
			transactions[index] = &Transaction{Id: 1, SumOfTransaction: 5_352_00, MCC: ""}// Например, каждая 100-ая транзакция в банке от нашего юзера в категории такой-то
		case 20:
			transactions[index] = &Transaction{Id: 3, SumOfTransaction: 1_362_00, MCC: "5455"}// Например, каждая 120-ая транзакция в банке от нашего юзера в категории такой-то
		case 40:
			transactions[index] = &Transaction{Id: 1, SumOfTransaction: 9_352_00, MCC: "5490"}
		default:
			transactions[index] = &Transaction{Id: 1, SumOfTransaction: 5_352_00, MCC: "5401"}// Транзакции других юзеров, нужны для "общей" массы
		}
	}
	return transactions
}

func TestSortTransactions(t *testing.T) {
	type args struct {
		slice []*Transaction
		Id    int64
	}
	tests := []struct {
		name string
		args args
		want map[string]int64
	}{
		{"first", args{makeTransactions(), 1}, map[string]int64{"Мобильная связь": 0, "Остальное": 5352000000, "Рестораны": 9352000000, "Супермаркеты": 519144000000, "Финансы": 0}},
	}
	for _, tt := range tests {
			if got := SortTransactions(tt.args.slice, tt.args.Id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortTransactions() = %v, want %v", got, tt.want)
			}
	}
}

func TestSortTransactionsByChanels(t *testing.T) {
	type args struct {
		slice      []*Transaction
		Id         int64
		partsCount int
	}
	tests := []struct {
		name string
		args args
		want map[string]int64
	}{
		{"first", args{makeTransactions(), 1, 500000}, map[string]int64{"Мобильная связь": 0, "Остальное": 5352000000, "Рестораны": 9352000000, "Супермаркеты": 519144000000, "Финансы": 0}},
	}
	for _, tt := range tests {
			if got := SortTransactionsByChanels(tt.args.slice, tt.args.Id, tt.args.partsCount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortTransactionsByChanels() = %v, want %v", got, tt.want)
			}
	}
}

func TestSortTransactionsByMutex(t *testing.T) {
	type args struct {
		slice      []*Transaction
		Id         int64
		goroutines int
	}
	tests := []struct {
		name string
		args args
		want map[string]int64
	}{
		{"first", args{makeTransactions(), 1, 500000}, map[string]int64{"Мобильная связь": 0, "Остальное": 5352000000, "Рестораны": 9352000000, "Супермаркеты": 519144000000, "Финансы": 0}},
	}
	for _, tt := range tests {
			if got := SortTransactionsByMutex(tt.args.slice, tt.args.Id, tt.args.goroutines); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortTransactionsByMutex() = %v, want %v", got, tt.want)
			}
	}
}

func TestSortTransactionsByOtherMutex(t *testing.T) {
	type args struct {
		slice      []*Transaction
		Id         int64
		goroutines int
	}
	tests := []struct {
		name string
		args args
		want map[string]int64
	}{
		{"first", args{makeTransactions(), 1, 500000}, map[string]int64{"Мобильная связь": 0, "Остальное": 5352000000, "Рестораны": 9352000000, "Супермаркеты": 519144000000, "Финансы": 0}},
	}
	for _, tt := range tests {
			if got := SortTransactionsByOtherMutex(tt.args.slice, tt.args.Id, tt.args.goroutines); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortTransactionsByOtherMutex() = %v, want %v", got, tt.want)
			}
	}
}

func BenchmarkCategorization(b *testing.B) {
	transactions := makeTransactions()
	want := map[string]int64{
		"Мобильная связь": 0,
		"Остальное": 5352000000,
		"Рестораны": 9352000000,
		"Супермаркеты": 519144000000,
		"Финансы": 0,
	}
	b.ResetTimer() // сбрасываем таймер, т.к. сама генерация транзакций достаточно ресурсоёмка
	for i := 0; i < b.N; i++ {
		result := SortTransactions(transactions,1)
		b.StopTimer() // останавливаем таймер, чтобы время сравнения не учитывалось
		if !reflect.DeepEqual(result, want) {
			b.Fatalf("invalid result, got %v, want %v", result, want)
		}
		b.StartTimer() // продолжаем работу таймера
	}
}

func BenchmarkCategorizationWithMutex(b *testing.B) {
	transactions := makeTransactions()
	want := map[string]int64{
		"Мобильная связь": 0,
		"Остальное": 5352000000,
		"Рестораны": 9352000000,
		"Супермаркеты": 519144000000,
		"Финансы": 0,
	}
	b.ResetTimer() // сбрасываем таймер, т.к. сама генерация транзакций достаточно ресурсоёмка
	for i := 0; i < b.N; i++ {
		result := SortTransactionsByMutex(transactions,1,500000)
		b.StopTimer() // останавливаем таймер, чтобы время сравнения не учитывалось
		if !reflect.DeepEqual(result, want) {
			b.Fatalf("invalid result, got %v, want %v", result, want)
		}
		b.StartTimer() // продолжаем работу таймера
	}
}

func BenchmarkCategorizationWithChanels(b *testing.B) {
	transactions := makeTransactions()
	want := map[string]int64{
		"Мобильная связь": 0,
		"Остальное": 5352000000,
		"Рестораны": 9352000000,
		"Супермаркеты": 519144000000,
		"Финансы": 0,
	}
	b.ResetTimer() // сбрасываем таймер, т.к. сама генерация транзакций достаточно ресурсоёмка
	for i := 0; i < b.N; i++ {
		result := SortTransactionsByChanels(transactions,1,500000)
		b.StopTimer() // останавливаем таймер, чтобы время сравнения не учитывалось
		if !reflect.DeepEqual(result, want) {
			b.Fatalf("invalid result, got %v, want %v", result, want)
		}
		b.StartTimer() // продолжаем работу таймера
	}
}

func BenchmarkCategorizationWithOtherMutex(b *testing.B) {
	transactions := makeTransactions()
	want := map[string]int64{
		"Мобильная связь": 0,
		"Остальное": 5352000000,
		"Рестораны": 9352000000,
		"Супермаркеты": 519144000000,
		"Финансы": 0,
	}
	b.ResetTimer() // сбрасываем таймер, т.к. сама генерация транзакций достаточно ресурсоёмка
	for i := 0; i < b.N; i++ {
		result := SortTransactionsByOtherMutex(transactions,1,500000)
		b.StopTimer() // останавливаем таймер, чтобы время сравнения не учитывалось
		if !reflect.DeepEqual(result, want) {
			b.Fatalf("invalid result, got %v, want %v", result, want)
		}
		b.StartTimer() // продолжаем работу таймера
	}
}