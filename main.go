package main

import (
	"fmt"
	"sync"
	"time"
)

func logInfo(format string, a ...any) {
	fmt.Printf("\033[34m(INFO) %s\033[0m\n", fmt.Sprintf(format, a...))
}

func logError(format string, a ...any) {
	fmt.Printf("\033[31m(ERROR) %s\033[0m\n", fmt.Sprintf(format, a...))
}

type User struct {
	ID      string
	Name    string
	Balance float64
	mu      sync.Mutex
}

func (u *User) Deposit(amount float64) {
	u.mu.Lock()
	u.Balance += amount

	// Симуляция долгой операции
	time.Sleep(time.Second)
	u.mu.Unlock()
}

func (u *User) Withdraw(amount float64) error {
	u.mu.Lock()
	if u.Balance < amount {
		u.mu.Unlock()
		return fmt.Errorf("balance is less than amount to withdraw")
	}

	u.Balance -= amount

	// Симуляция долгой операции
	time.Sleep(time.Second)

	u.mu.Unlock()
	return nil
}

func (u *User) Print() {
	logInfo("Пользователь %s имеет на балансе %f", u.Name, u.Balance)
}

func NewUser(id string, name string, balance float64) *User {
	return &User{ID: id, Name: name, Balance: balance}
}

type Transaction struct {
	FromID string
	ToID   string
	Amount float64
}

func (t *Transaction) Name() string {
	return fmt.Sprintf("%s -> %s (%f)", t.FromID, t.ToID, t.Amount)
}

type PaymentSystem struct {
	Users            map[string]*User
	TransactionQueue []Transaction
}

func (ps *PaymentSystem) AddUser(u *User) {
	ps.Users[u.ID] = u
}

func (ps *PaymentSystem) AddTransaction(t Transaction) {
	ps.TransactionQueue = append(ps.TransactionQueue, t)
}

func (ps *PaymentSystem) GetUser(uid string) *User {
	for _, u := range ps.Users {
		if u.ID == uid {
			return u
		}
	}

	// Симуляция долгой операции
	time.Sleep(time.Second)

	return nil
}

func (ps *PaymentSystem) ProcessTransaction(t Transaction) error {
	fromUser := ps.GetUser(t.FromID)
	toUser := ps.GetUser(t.ToID)

	logInfo("Обработка транзакции [%s] в процессе...", t.Name())

	if fromUser == nil {
		return fmt.Errorf("Пользователь FromID (%s) не найден!", t.FromID)
	}

	if toUser == nil {
		return fmt.Errorf("Пользователь ToID (%s) не найден!", t.ToID)
	}

	if err := fromUser.Withdraw(t.Amount); err != nil {
		return err
	}

	toUser.Deposit(t.Amount)

	logInfo("Завершилась обработка транзакции [%s]", t.Name())

	return nil
}

func (ps *PaymentSystem) Worker(ch <-chan Transaction, wg *sync.WaitGroup) {
	defer wg.Done()

	for t := range ch {
		if err := ps.ProcessTransaction(t); err != nil {
			logError("При попытки обработки транзакции [%s] произошла ошибка: %s", t.Name(), err.Error())
		}

	}

}

func NewPaymentSystem() *PaymentSystem {
	return &PaymentSystem{Users: make(map[string]*User)}
}

func main() {
	user1 := NewUser("1", "Bob", 1000)
	user2 := NewUser("2", "Bib", 500)

	ps := NewPaymentSystem()
	ps.AddUser(user1)
	ps.AddUser(user2)

	user1.Print()
	user2.Print()

	ps.AddTransaction(Transaction{FromID: user1.ID, ToID: user2.ID, Amount: 200})
	ps.AddTransaction(Transaction{FromID: user2.ID, ToID: user1.ID, Amount: 50})
	ps.AddTransaction(Transaction{FromID: user1.ID, ToID: user2.ID, Amount: 30})
	ps.AddTransaction(Transaction{FromID: user2.ID, ToID: user1.ID, Amount: 80})
	ps.AddTransaction(Transaction{FromID: user1.ID, ToID: user2.ID, Amount: 5500})
	ps.AddTransaction(Transaction{FromID: user2.ID, ToID: user1.ID, Amount: 50})

	var wg sync.WaitGroup
	ch := make(chan Transaction, len(ps.TransactionQueue))

	for range 3 {
		wg.Add(1)
		go ps.Worker(ch, &wg)
	}

	for _, t := range ps.TransactionQueue {
		ch <- t
	}

	close(ch)

	wg.Wait()

	user1.Print()
	user2.Print()
}
