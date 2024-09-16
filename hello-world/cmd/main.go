package main

import (
	"log"

	"hello-world/db"
)

func main() {
	op1 := db.GetOp{K: "key1"}
	op2 := db.GetOp{K: "key2"}
	db.AppendToBatch(op1)
	db.AppendToBatch(op2)
	log.Printf("batch => %#v", db.GetBatch())
}