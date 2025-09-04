// package faulttolerance

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"sync"
// 	"time"
// )

// func StartHealthMonitor() bool {
// 	for {
// 		q := health_checker()
// 		counter := 0
// 		for _, val := range q {
// 			if !val {
// 				counter++
// 			}
// 		}
// 		if counter > 1 {
// 			return false
// 		}
// 		time.Sleep(10 * time.Second)
// 	}
// }

// func health_checker() map[int]bool {
// 	var wg sync.WaitGroup
// 	var m map[int]bool
// 	var mt sync.Mutex
// 	for i := 3001; i < 3004; i++ {
// 		url := fmt.Sprintf("http://localhost:%d/healthcheck", i)
// 		wg.Add(1)
// 		go func(url string, i int) {
// 			defer wg.Done()
// 			a := hit_route(url)
// 			mt.Lock()
// 			m[i] = a
// 			mt.Unlock()
// 		}(url, i)
// 	}
// 	wg.Wait()
// 	return m
// }
// func hit_route(url string) bool {
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		log.Fatal(err)
// 		return false
// 	}
// 	fmt.Println(resp)
// 	if resp.StatusCode != http.StatusAccepted {
// 		return false
// 	}
// 	return true
// }

package faulttolerance

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

func StartHealthMonitor() bool {
	for {
		q := health_checker()
		counter := 0
		for _, val := range q {
			if !val {
				counter++
			}
		}
		if counter > 1 {
			return false
		}
		time.Sleep(10 * time.Second)
	}
}

func health_checker() map[int]bool {
	var wg sync.WaitGroup
	m := make(map[int]bool) // FIX: initialize map
	var mt sync.Mutex

	for i := 3001; i < 3004; i++ {
		url := fmt.Sprintf("http://localhost:%d/healthcheck", i)
		wg.Add(1)
		go func(url string, i int) {
			defer wg.Done()
			a := hit_route(url)
			mt.Lock()
			m[i] = a
			mt.Unlock()
		}(url, i)
	}

	wg.Wait()
	return m
}

func hit_route(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("[X] Error hitting %s: %v\n", url, err)
		return false
	}
	defer resp.Body.Close()

	fmt.Println(resp.Status)

	if resp.StatusCode != http.StatusOK { // FIX: expect 200
		return false
	}
	return true
}
