package main

import (
	"fmt"
	"sync"
)

var uID int

type User struct {
	ID      int
	Name    string
	Balance float32
	mu      sync.Mutex
}

func (u *User) NewUser(name string) *User {
	u.mu.Lock()
	uID += 1
	u.mu.Unlock()
	var balance float32

	return &User{ID: uID, Name: name, Balance: balance}
}

func (u *User) Deposit(amount float32) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.Balance += amount
	fmt.Printf("Баланс пользователя %s равен: %.2f рублей.\n", u.Name, u.Balance)
}

func (u *User) Withdraw(amount float32) {
	u.mu.Lock()
	defer u.mu.Unlock()

	if amount > u.Balance {
		fmt.Printf("На счету пользователя %s недостаточно средств для списания %.2f рублей.\n", u.Name, u.Balance)
		return
	}
	u.Balance -= amount
	fmt.Printf("Баланс пользователя %s равен: %.2f рублей.\n", u.Name, u.Balance)
}

func main() {
	var u User
	accounts := make(map[string]*User)
	username := []string{"Kate", "Dung", "Son", "Ngoc"}
	for _, val := range username {
		v := val
		accounts[v] = u.NewUser(v)
	}

	for _, usr := range accounts {
		fmt.Printf("id=%d name=%s balance=%.2f\n", usr.ID, usr.Name, usr.Balance)
	}
}
