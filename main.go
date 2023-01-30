package main

import (
	"fmt"
	"sync"
)

func main() {
	ch1, err := read("file1.csv")
	if err != nil {
		panic(fmt.Errorf("could not read file1 %v", err))
	}

	ch2, err := read("file2.csv")
	if err != nil {
		panic(fmt.Errorf("could not read file2 %v", err))
	}

	//-

	exit := make(chan struct{})

	chM := merge(ch1, ch2)

	go func() {
		for v := range chM {
			fmt.Println(v)
		}

		close(exit)
	}()

	<-exit

	fmt.Println("All completed, exiting")
}

func merge(cs ...<-chan string) <-chan string {
	var wg sync.WaitGroup

	out := make(chan string)

	send := func(c <-chan string) {
		for n := range c {
			out <- n
		}

		wg.Done()
	}

	wg.Add(len(cs))

	for _, c := range cs {
		go send(c)
	}

	go func() {
		wg.Wait()

		close(out)
	}()

	return out
}

func read(fileName string) (<-chan string, error) {
	ch := make(chan string, 10)

	go func(ch chan string) {
		for i := 1; i <= 100; i++ {
			ch <- fmt.Sprintf("fileName %s val: %d", fileName, i)
		}
		close(ch)
	}(ch)

	return ch, nil
}
