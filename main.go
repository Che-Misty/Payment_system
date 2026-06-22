package main

import (
	"fmt"
	"math/rand"
	"sync"
)

var uID int

type User struct {
	ID      int
	Name    string
	Balance float64
	mu      sync.Mutex
}

func NewUser(name string, balance float64) *User {
	uID += 1

	return &User{ID: uID, Name: name, Balance: balance}
}

type Transaction struct {
	FromID int
	ToID   int
	Amount float64
}

type PaymentSystem struct {
	Users        map[int]*User
	Transactions []Transaction
	mu           sync.Mutex
}

func NewPaymentSystem() *PaymentSystem {
	return &PaymentSystem{Users: make(map[int]*User), Transactions: make([]Transaction, 0)}
}

func (ps *PaymentSystem) AddUser(name string, balance float64) {
	user := NewUser(name, balance)
	id := user.ID
	ps.Users[id] = user
}

func (ps *PaymentSystem) AddTransaction(id1 int, id2 int, amount float64) {
	trx := Transaction{FromID: id1, ToID: id2, Amount: amount}
	ps.Transactions = append(ps.Transactions, trx)
}

func (u *User) Deposit(amount float64) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.Balance += amount
}

func (u *User) Withdraw(amount float64) (bool, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	if amount > u.Balance {
		return false, fmt.Errorf("User %s haven't enought money", u.Name)
	}
	u.Balance -= amount
	return true, nil
}

func (ps *PaymentSystem) ProcessingTransaction(i int) error {
	trx := ps.Transactions[i]
	sender := ps.Users[trx.FromID]
	recipient := ps.Users[trx.ToID]

	if len(sender.Name) == 0 {
		return fmt.Errorf("User %s not found!", sender.Name)
	}

	if len(recipient.Name) == 0 {
		return fmt.Errorf("User %s not found!", recipient.Name)
	}

	if ok, err := sender.Withdraw(trx.Amount); !ok {
		return err
	}
	recipient.Deposit(trx.Amount)
	return nil
}

func main() {
	// 	accounts := make(map[string]*User)
	// 	username := []string{"Kate", "Dung", "Son", "Ngoc"}

	ps := NewPaymentSystem()

	ps.AddUser("Khena", 400)
	ps.AddUser("Shiro", 200)

	ps.AddTransaction(1, 2, 100)
	ps.AddTransaction(1, 2, 91.23)
	ps.AddTransaction(2, 1, 44.41)
	ps.AddTransaction(2, 1, 250.12)
	ps.AddTransaction(1, 2, 121.17)
	
	

	for idx := range ps.Transactions {
		ps.ProcessingTransaction(idx)
	}

	for _, user := range ps.Users {
		fmt.Printf("Баланс пользователя %s на текущий момент: %.2f$\n", user.Name, user.Balance)
	}
}
