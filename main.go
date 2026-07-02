package main

import (
	"fmt"
	"sync"
)

type User struct {
	ID      string
	Name    string
	Balance float64
	mu      sync.Mutex
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

func NewPaymentSystem() *PaymentSystem {
	return &PaymentSystem{Users: make(map[string]*User), TransactionQueue: make([]Transaction, 0)}
}

func (ps *PaymentSystem) AddUser(id string, name string, balance float64) {
	user := NewUser(id, name, balance)
	ps.Users[id] = user
}

func (ps *PaymentSystem) AddTransaction(id1 string, id2 string, amount float64) {
	trx := Transaction{FromID: id1, ToID: id2, Amount: amount}
	ps.TransactionQueue = append(ps.TransactionQueue, trx)
}

func (u *User) Deposit(amount float64) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.Balance += amount
}

func (u *User) Withdraw(amount float64) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if amount > u.Balance {
		return fmt.Errorf("User %s haven't enought money", u.Name)
	}
	u.Balance -= amount
	return nil
}

func (ps *PaymentSystem) ProcessingTransaction(trx Transaction) error {
	sender, ok := ps.Users[trx.FromID]
	if !ok {
		return fmt.Errorf("User %s not found!", trx.FromID)
	}
	recipient, ok := ps.Users[trx.ToID]
	if !ok {
		return fmt.Errorf("User %s not found!", trx.ToID)
	}

	if err := sender.Withdraw(trx.Amount); err != nil {
		return err
	}
	recipient.Deposit(trx.Amount)
	return nil
}

func (ps *PaymentSystem) Worker(ch <-chan Transaction) {
	for trx := range ch {
		if err := ps.ProcessingTransaction(trx); err != nil {
			fmt.Println(err)
		}
	}
}

func main() {
	var wg sync.WaitGroup

	ps := NewPaymentSystem()

	ps.AddUser("1", "Khena", 400)
	ps.AddUser("2", "Shiro", 200)

	ps.AddTransaction("1", "2", 100)
	ps.AddTransaction("1", "2", 91.23)
	ps.AddTransaction("2", "1", 44.41)
	ps.AddTransaction("2", "1", 250.12)
	ps.AddTransaction("1", "2", 121.17)
	ps.AddTransaction("1", "2", 500)
	ps.AddTransaction("1", "2", 11.11)
	ps.AddTransaction("1", "2", 10.01)
	ps.AddTransaction("2", "1", 250.12)
	ps.AddTransaction("2", "1", 54.12)

	ch := make(chan Transaction, len(ps.TransactionQueue))

	for _, t := range ps.TransactionQueue {
		ch <- t
	}

	for i := 0; i < 3; i++ {
		wg.Go(func() { ps.Worker(ch) })
	}

	close(ch)
	wg.Wait()

	for _, user := range ps.Users {
		fmt.Printf("Баланс пользователя %s на текущий момент: %.2f$\n", user.Name, user.Balance)
	}
}
