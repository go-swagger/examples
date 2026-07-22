//go:build ignore
// +build ignore

package main

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/swag/conv"
	"github.com/go-swagger/examples/stream-server/client"
	"github.com/go-swagger/examples/stream-server/client/operations"
	"github.com/go-swagger/examples/stream-server/models"
)

func main() {
	n := int64(5)
	if len(os.Args) > 1 {
		var err error
		n, err = conv.ConvertInt64(os.Args[1])
		if err != nil {
			log.Fatalln("pass an integer as argument")
			return
		}
	}
	log.Printf("asking server for countdown timings: %d", n)

	if err := ask(n); err != nil {
		log.Printf("failure: %v", err)
	}
}

func ask(n int64) error {
	// snippet:consumer
	customized := httptransport.New("localhost:8000", "/", []string{"http"})
	customized.Consumers[runtime.JSONMime] = runtime.ByteStreamConsumer()

	countdowns := client.New(customized, nil)

	reader, writer := io.Pipe()

	scanner := bufio.NewScanner(reader)
	// endsnippet:consumer

	ctx, cancel := context.WithCancel(context.Background())

	// consumes asynchronously the response buffer
	var wg sync.WaitGroup

	// snippet:scan
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		defer cancel()

		// read response items line by line
		for scanner.Scan() {
			// each response item is JSON
			txt := scanner.Text()
			log.Printf("received countdown mark - raw: %s", txt)

			var mark models.Mark

			err := json.Unmarshal([]byte(txt), &mark)
			if err != nil {
				log.Printf("unmarshal error: %v", err)
				return
			}

			log.Printf("received countdown mark - remaining: %d", conv.Value(mark.Remains))
		}

		if err := scanner.Err(); err != nil {
			log.Printf("scanner err: %v", err)
		}

		log.Println("EOF")
	}(&wg)
	// endsnippet:scan

	queryCtx, timedOut := context.WithTimeout(ctx, 7*time.Second)
	defer timedOut()

	// snippet:request
	elapsed := operations.NewElapseParamsWithContext(queryCtx).WithLength(n)

	_, err := countdowns.Operations.Elapse(elapsed, writer)
	// endsnippet:request

	if err == nil {
		log.Printf("response complete")
	} else {
		log.Printf("got an error")
	}

	_ = writer.Close()

	wg.Wait()

	return err
}
