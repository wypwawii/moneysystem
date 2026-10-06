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

func (u *User) Withdraw(amount float64) {
	if u.Balance < amount {
		return
	}

	u.Balance -= amount
}

func (u *User) Print() {
	fmt.Printf("Пользователь %s имеет на балансе %f\n", u.Name, u.Balance)
}

func main() {
	user1 := &User{ID: "1", Name: "Bob", Balance: 1000}
	user2 := &User{ID: "2", Name: "Bib", Balance: 2000}

	user1.Print()
	user2.Print()

	user1.Deposit(500)
	user2.Withdraw(1500)

	user1.Print()
	user2.Print()
}
