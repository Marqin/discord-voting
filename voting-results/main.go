package main

import (
	"flag"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/boltdb/bolt"
)

var topForNextRound *int

func main() {

	topForNextRound = flag.Int("top", 0, "make a line to separate top `n` places")
	flag.Parse()

	db, err := bolt.Open("results.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	results, pplVoted, err := getResults(db)
	if err != nil {
		log.Fatal(err)
	}

	if pplVoted > 0 {
		fmt.Printf("%d people have voted.\n\n", pplVoted)
		printResults(results)
	} else {
		fmt.Println("No votes casted!")
	}

}

func getResults(db *bolt.DB) (map[int][]string, int, error) {
	results := make(map[int][]string)

	votesPerOutfit := make(map[string]int)

	pplVoted := 0

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("votes"))

		if b == nil {
			return nil
		}

		return b.ForEach(func(k, v []byte) error {
			pplVoted++

			uv, err := decode(v)
			if err != nil {
				return err
			}

			for _, outfit := range uv.Votes {
				votesPerOutfit[outfit]++
			}

			return nil
		})
	})

	for outfit, votes := range votesPerOutfit {
		results[votes] = append(results[votes], outfit)
	}

	return results, pplVoted, err
}

func printResults(results map[int][]string) {

	if len(results) > 0 {
		fmt.Println("Results:")

		keys := make([]int, 0, len(results))
		for key := range results {
			keys = append(keys, key)
		}

		sort.Sort(sort.Reverse(sort.IntSlice(keys)))

		place := 1
		line := false
		for _, key := range keys {

			if *topForNextRound > 0 && !line && place > *topForNextRound {
				fmt.Println("------------------------------------------------------")
				line = true
			}

			for i, outfit := range results[key] {
				if i == 0 {
					fmt.Printf("%2d. %s (%d votes)\n", place, outfit, key)
				} else {
					fmt.Printf("    %s (%d votes)\n", outfit, key)
				}
			}
			place += len(results[key])
		}
	}

}
