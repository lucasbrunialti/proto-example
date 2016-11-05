package main

import (
	"log"
	"math/rand"
	"time"
	"unsafe"

	"github.com/lucasbrunialti/proto-example/example"
	"github.com/golang/protobuf/proto"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandomTest() *example.Test {
	test := &example.Test{
		Label: proto.String(RandStringRunes(10)),
		Type:  proto.Int32(rand.Int31()),
		Reps:  []int64{rand.Int63(), rand.Int63(), rand.Int63()},
	}

	return test
}

func Run() {
	test := RandomTest()

	// log.Printf("Serializing it...")

	start := time.Now()
	data, err := proto.Marshal(test)
	elapsed := time.Since(start)

	log.Printf("Serialization in %d ns", elapsed.Nanoseconds())

	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	newTest := &example.Test{}

	// log.Printf("Retrieving it...")

	start = time.Now()
	err = proto.Unmarshal(data, newTest)
	elapsed = time.Since(start)

	log.Printf("Unmarshal in %d ns", elapsed.Nanoseconds())

	if err != nil {
		log.Fatal("unmarshaling error: ", err)
	}
	// Now test and newTest contain the same data.

	log.Printf("Test as proto object", test.String())
	if test.String() != newTest.String() {
		log.Fatalf("data mismatch %q != %q", test.String(), newTest.String())
	}

	// log.Printf("Label: %d bytes", int(unsafe.Sizeof(*test.Label))+len(*test.Label))
	// log.Printf("Type: %d bytes", int(unsafe.Sizeof(*test.Type)))
	// log.Printf("Reps: %d bytes", int(unsafe.Sizeof(test.Reps)))
	log.Printf("Size of test (proto) %d bytes vs %d raw bytes", proto.Size(newTest), (int(unsafe.Sizeof(*test.Label)) + len(*test.Label) + int(unsafe.Sizeof(*test.Type)) + int(unsafe.Sizeof(test.Reps))))
}

func main() {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 100; i++ {
		Run()
	}
	// 1st conclusion: pb spends less than raw
}
