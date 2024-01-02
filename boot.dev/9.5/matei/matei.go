package main

import (
	"crypto/md5"
	"fmt"
	"io"
)

func getCounts(userIDs []string) map[string]int {

	counts := make(map[string]int)

	for _, id := range userIDs {

		_, exists := counts[id]

		if !exists {
				// if the key doesn't exist yet,
				// initialize its value to 0
				counts[id] = 0
		} 

		// if the key does exist, increment its value
		counts[id]++
	}

	return counts
}

// don't edit below this line

func test(userIDs []string, ids []string) {
	fmt.Printf("Generating counts for %v user IDs...\n", len(userIDs))

	counts := getCounts(userIDs)
	fmt.Println("Counts from select IDs:")
	for _, k := range ids {
		v := counts[k]
		fmt.Printf(" - %s: %d\n", k, v)
	}
	fmt.Println("=====================================")
}

func main() {
	userIDs := []string{}
	for i := 0; i < 10000; i++ {
		h := md5.New()
		io.WriteString(h, fmt.Sprint(i))
		key := fmt.Sprintf("%x", h.Sum(nil))
		userIDs = append(userIDs, key[:2])
	}

	test(userIDs, []string{"00", "ff", "dd"})
	test(userIDs, []string{"aa", "12", "32"})
	test(userIDs, []string{"bb", "33"})
}
