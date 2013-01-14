/*
Provides simple interface to mysql transactions
*/
package mysql

import (
	"strings"
	"fmt"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native" // Native engine
)

type transactionError struct {
	message string
}

// Provides a simple container for the query parameters of an transaction
type TransactionQuery struct {
	Database, Table string
	Columns         []string
}

// Our class to handle transactions with
type Transaction struct {
	stmt mysql.Stmt
	transaction mysql.Transaction
	numParams int
}

// To implement the error interface
func (te transactionError) Error() string {
	return te.message
}

type Entry interface {
	GetParams() ([]interface{})
}

// Allocates a new transaction
func NewTransaction(connection mysql.Conn, query TransactionQuery) (*Transaction, error) {
	if len(query.Columns) < 1 {
		return nil, transactionError{"Too few columns"}
	}

	sss := fmt.Sprintf("INSERT INTO `%s`.`%s` (`%s`) VALUES (?%s)", query.Database, query.Table, strings.Join(query.Columns, "`, `"), strings.Repeat(", ?", len(query.Columns)-1))
	ins, err := connection.Prepare(sss)
	if err != nil {
		return nil, err
	}

	trans, err := connection.Begin()
	if err != nil {
		return nil, err
	}

	return &Transaction{trans.Do(ins), trans, len(query.Columns)}, nil
}

// inserts everything from the channel until it gets closed
func (trans *Transaction) BeginInsert(c chan Entry, mut chan int) {
	defer trans.transaction.Commit()

	for {
		entry, ok := <-c
		if !ok {
			break
		}

		entries := entry.GetParams()

		if len(entries) == trans.numParams {
			_, err := trans.stmt.Run(entries...)
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Println("Error: arguments count does not match")
		}
	}

	mut <- 0
}