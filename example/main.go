package main

import (
	"context"
	"log"
	"time"

	redis "github.com/go-redis/redis/v8"
	redlock "github.com/kecci/go-redis-lock"
)

func main() {
	log.SetFlags(log.Ltime)
	rc1 := redis.NewClient(&redis.Options{Addr: "0.0.0.0:6379"})
	// rc2 := redis.NewClient(&redis.Options{Addr: "0.0.0.0:7002"})
	// rc3 := redis.NewClient(&redis.Options{Addr: "0.0.0.0:7003"})

	dlm := redlock.NewDLM([]*redis.Client{rc1}, 10*time.Second, 2*time.Second)

	// With Lock Only
	// withLockOnly(dlm)

	var messages = make(chan string)

	go func() {
		// With Lock Unlock
		withLockAndUnlock(dlm, "1")
		messages <- "1"
	}()

	go func() {
		// With Lock Unlock
		withLockAndUnlock(dlm, "1")
		messages <- "2"
	}()

	var message = <-messages
	println(message)
}

func withLockAndUnlock(dlm *redlock.DLM, orderID string) {
	ctx := context.Background()
	locker := dlm.NewLocker("this-is-a-key-00" + orderID)

	if err := locker.Lock(ctx); err != nil {
		log.Println("[main] Failed when locking, err:", err)
	}

	// Perform operation.
	someOperation()

	if err := locker.Unlock(ctx); err != nil {
		log.Println("[main] Failed when unlocking, err:", err)
	}

	log.Println("[main] Done")
}

// func withLockOnly(dlm *redlock.DLM) {
// 	ctx := context.Background()
// 	locker := dlm.NewLocker("this-is-a-key-002")

// 	if err := locker.Lock(ctx); err != nil {
// 		log.Fatal("[main] Failed when locking, err:", err)
// 	}

// 	// Perform operation.
// 	someOperation()

// 	// Don't unlock

// 	log.Println("[main] Done")
// }

func someOperation() {
	log.Println("[someOperation] Process has been started")
	time.Sleep(2 * time.Second)
	log.Println("[someOperation] Process has been finished")
}
