package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/status"
	"os"
	"time"
)

import "cloud.google.com/go/firestore"

func main() {
	projectID := flag.String("project", "", "project id")
	dbName := flag.String("db", "", "database name")
	timeout := flag.Int("timeout", 30, "timeout in seconds")
	printRefs := flag.Bool("refs-only", false, "print document refs instead")
	getByRef := flag.String("get-by-ref", "", "get document by ref")
	collection := flag.String("collection", "tasks", "collection name")

	flag.Parse()

	if *projectID == "" {
		flag.Usage()
		fmt.Println("project id is required")
		os.Exit(1)
	}

	if *dbName == "" {
		flag.Usage()
		fmt.Println("database name is required")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout)*time.Second)
	defer cancel()

	var err error
	var total int
	var lastID string

	client, err := firestore.NewClientWithDatabase(ctx, *projectID, *dbName)
	if err != nil {
		fmt.Printf("error creating firestore client: %v\n", err)
		os.Exit(1)
	}

	defer func() {
		err := client.Close()
		if err != nil {
			fmt.Printf("error closing firestore client: %v\n", err)
		}
	}()

	if *getByRef != "" {
		doc, err := client.Collection(*collection).Doc(*getByRef).Get(ctx)
		if err != nil {
			fmt.Printf("error getting document by ref: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("document:")
		spew.Dump(doc.Data())
		return
	}

	start := time.Now()

	defer func() {
		fmt.Println("\n-------------------------")
		fmt.Printf("total documents processed: %d\n", total)
		fmt.Printf("last document ref: %s\n", lastID)
		fmt.Printf("elapsed time: %v\n", time.Since(start))
		fmt.Println("-------------------------")
	}()

	iter := client.Collection(*collection)

	if *printRefs {
		total, lastID, err = getRefsOnly(ctx, iter)
	} else {
		total, lastID, err = getAllDocuments(ctx, iter)
	}

	if err != nil {
		fmt.Println(err.Error())
		s, ok := status.FromError(err)
		if ok {
			fmt.Printf("status code: %d\n", s.Code())
			fmt.Printf("status message: %s\n", s.Message())
		}
	}
}

func getAllDocuments(ctx context.Context, collection *firestore.CollectionRef) (total int, lastID string, err error) {
	iter := collection.Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return total, lastID, fmt.Errorf("error getting document: %v\n", err)
		}

		fmt.Println(doc.Ref.ID)
		lastID = doc.Ref.ID
		total++
	}

	return total, lastID, nil
}

func getRefsOnly(ctx context.Context, collection *firestore.CollectionRef) (total int, lastID string, err error) {
	iter := collection.DocumentRefs(ctx)

	for {
		ref, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return total, lastID, fmt.Errorf("error getting document ref: %v\n", err)
		}

		fmt.Println(ref.ID)
		lastID = ref.ID
		total++
	}

	return total, lastID, nil
}
