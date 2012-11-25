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

	ins, err := connection.Prepare(fmt.Sprintf("INSERT INTO %s.%s (%s) VALUES (?%s)", query.Database, query.Table, strings.Join(query.Columns, ", "), strings.Repeat(", ?", len(query.Columns)-1)))
	if err != nil {
		return nil, transactionError{"Error creating prepared Statement"}
	}

	fmt.Println(ins)

	trans, err := connection.Begin()
	if err != nil {
		return nil, transactionError{"Error creating Transaction"}
	}

	return &Transaction{trans.Do(ins), trans, len(query.Columns)}, nil
}

// inserts everything from the channel until it gets closed
func (trans *Transaction) BeginInsert(c chan *Entry, mut chan int) {
	for {
		entry, ok := <-c
		if !ok {
			break
		}

		entries := entry.GetParams()
		if len(entries) == trans.numParams {
			_, err := trans.stmt.Run(entries)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Error: arguments count does not match")
		}
	}

	trans.transaction.Commit()
	mut <- 0
}
