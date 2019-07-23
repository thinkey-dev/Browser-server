package mongo_test

import (
	"fmt"
	"testing"

	"github.com/dvpp/go-dvpp/tests/mongo"
)

func TestMongo(t *testing.T) {
	collection, err := mongo.InitMongod("contracts", "contract")
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	res := collection.Indexes()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	id := res
	fmt.Printf("%v\n", id)
}
