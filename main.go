package main

import (
	"fmt"
)

var uID int

type User struct {
	ID      int
	Name    string
	Balance float32
}

func (u *User) Deposit(amount float32) {
	u.Balance += amount
	fmt.Printf("Баланс пользователя %s равен: %.2f рублей.\n", u.Name, u.Balance)
}

func (u *User) Withdraw(amount float32) {
	if amount > u.Balance {
		fmt.Printf("На счету пользователя %s недостаточно средств для списания %.2f рублей.\n", u.Name, u.Balance)
		return
	}
	u.Balance -= amount
	fmt.Printf("Баланс пользователя %s равен: %.2f рублей.\n", u.Name, u.Balance)
}

func main() {
	user1 := &User{1, "Kate", 100}
	user2 := &User{2, "Dung", 100}

	user2.Withdraw(150)
	user1.Deposit(200)
	user1.Withdraw(135)
	user2.Deposit(200)
	user1.Withdraw(237)
	user2.Deposit(1200)
	user2.Withdraw(400)
	user1.Deposit(22)
	user2.Withdraw(1000)
}
