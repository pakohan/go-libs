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

//Provides a simple container for the query parameters of an transaction
type TransactionQuery struct {
	Database, Table string
	Columns         []string
}

//Our class to handle transactions with
type Transaction struct {
	stmt mysql.Stmt
	transaction mysql.Transaction
	numParams int
}

//To implement the error interface
func (te transactionError) Error() string {
	return te.message
}

func NewTransaction(connection mysql.Conn, query TransactionQuery) (*Transaction, error) {
	if len(query.Columns) < 1 {
		return nil, transactionError{"Too few columns"}
	}

	ins, err := connection.Prepare(fmt.Sprintf("INSERT INTO %s.%s (%s) VALUES (?%s)", query.Database, query.Table, strings.Join(query.Columns, ", "), strings.Repeat(", ?", len(query.Columns)-1)))
	if err != nil {
		return nil, transactionError{"Error creating prepared Statement"}
	}

	trans, err := connection.Begin()
	if err != nil {
		return nil, transactionError{"Error creating Transaction"}
	}

	return &Transaction{trans.Do(ins), trans, len(query.Columns)}, nil
}

func (trans *Transaction) BeginInsert(c chan []interface{}, mut chan int) {
	for {
		params, ok := <-c
		if !ok {
			break
		}

		if len(params) == trans.numParams {
			_, err := trans.stmt.Run(params)
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
