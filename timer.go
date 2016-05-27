package main

import (
	"fmt"
	"github.com/walltearer/gonsole"
	"math"
	"sync"
	"time"
)

// Refresh counter of remaining seconds in console
// NOTE: Windows cmd doesn't seem to support saving and restoring cursor position
func refreshRemaining(i int) {
	// \033[s - saves current cursor position
	// \033[A - moves line up
	// \r - goes to the beginning of the line
	// \033[K - clears the whole line
	// \033[u - restores saved cursor position
	fmt.Printf("\033[s\033[A\r\033[KTime remaining: %dm %ds\033[u", i/60, int(math.Mod(float64(i), 60)))
}

func main() {
	gon := gonsole.New()

	var wg sync.WaitGroup
	wg.Add(1)

	ch := make(chan int)

	fmt.Println("Timer")
	d := gon.ReadInt("Enter seconds: ")
	fmt.Println()

	go func() {
		// starting countdown for the first second
		second := time.After(time.Second)
		// looping until there are remaining seconds
		for i := d; i > 0; {
			refreshRemaining(i)
			select {
			case up := <-ch:
				i += up
			case <-second:
				// remaining seconds should be substracted only after second passes, not on each iteration
				i--
				// starting countdown for next second
				second = time.After(time.Second)
			}
		}
		wg.Done()
	}()

	go func() {
		for {
			upd := gon.ReadInt("Add/Substract seconds: ")
			ch <- upd
		}
	}()

	wg.Wait()
	fmt.Print("\nDone\n")
}
