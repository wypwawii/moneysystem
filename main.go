package main

import "fmt"

type User struct {
	ID      string
	Name    string
	Balance float64
}

func (u *User) Deposit(amount float64) {
	u.Balance += amount
}

func (u *User) Withdraw(amount float64) error {
	if u.Balance < amount {
		return fmt.Errorf("Баланс меньше чем сумма перевода")
	}

	u.Balance -= amount

	return nil
}

func (u *User) Print() {
	fmt.Printf("Пользователь %s имеет на балансе %f\n", u.Name, u.Balance)
}

func NewUser(id string, name string, balance float64) *User {
	return &User{ID: id, Name: name, Balance: balance}
}

type Transaction struct {
	FromID string
	ToID   string
	Amount float64
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

	return nil
}

func (ps *PaymentSystem) ProcessTransaction(t Transaction) error {
	fromUser := ps.GetUser(t.FromID)
	toUser := ps.GetUser(t.ToID)

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

	return nil
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

	for _, t := range ps.TransactionQueue {
		if err := ps.ProcessTransaction(t); err != nil {
			fmt.Println(err)
		}
	}

	user1.Print()
	user2.Print()
}
