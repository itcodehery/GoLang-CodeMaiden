package main

import (
	"errors"
	"fmt"
)

// Q1: Divide by Zero
func divide(a int, b int) (res int, err error) {
	if b == 0 {
		res = 0
		err = errors.New("Error: cannot divide by zero")
	} else {
		res = a / b
		err = nil

	}

	return res, err
}

// Q2: Withdrawal
func withdraw(balance float64, amount float64) (bal float64, err error) {
	bal = balance
	if amount > balance {
		bal = balance
		err = errors.New("Insufficient Balance")
	} else if amount <= 0 {
		bal = balance
		err = errors.New("Invalid Amount")
	} else {
		bal = bal - amount
		err = nil
	}
	return bal, err

}

// Q3: Custom Error
type ScoreError struct {
	Score int
}

func (e ScoreError) Error() string {
	return fmt.Sprintf("invalid score: ", e.Score)
}

func validateScore(score int) error {
	if score < 0 {
		errscore := ScoreError{Score: score}
		return errscore
	} else {
		return nil
	}
}

func main() {
	result, err := divide(10, 2)
	if err != nil {
		fmt.Println("Brother you have sinned: ", err)
	} else {
		fmt.Println("Result: ", result)
	}

	balance := 1000.0
	newBalance, err := withdraw(balance, 400)
	if err != nil {
		fmt.Println("Transaction failed: ", err)
	} else {
		fmt.Println("New Balance:", newBalance)
	}

	score := 23
	verr := validateScore(score)
	if verr != nil {
		fmt.Println("Error occurred: ", verr.Error())
	} else {
		fmt.Println("Score is valid!")
	}

}
