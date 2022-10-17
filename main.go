package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Container struct {
	goodsName string
}

type Harbor struct {
	c          chan Ship
	piers      []bool
	containers []Container
}

type Ship struct {
	id         int
	containers []Container
	size       int
	harbor     Harbor
}

func (receiver Ship) run(wg *sync.WaitGroup) {
	defer wg.Done()

	ship := <-receiver.harbor.c
	ship.moveContainers()
}

func (receiver Ship) moveContainers() {
	fmt.Print("Ship ")
	fmt.Print(receiver.id + 1)
	fmt.Println(" start moving containers")
	time.Sleep(time.Duration(len(receiver.containers)) * time.Second)
	fmt.Print("Ship ")
	fmt.Print(receiver.id + 1)
	fmt.Println(" finish moving containers")
}

func generatePiers(count int) []bool {
	var containers = []bool{}
	for i := 0; i < count; i++ {
		containers = append(containers, false)
	}
	return containers
}

func generateContainers(count int) []Container {
	var containers = []Container{}
	for i := 0; i < count; i++ {
		ship := Container{goodsName: "Banana"}
		containers = append(containers, ship)
	}
	return containers
}

func generateShips(harbor Harbor, harborSize int) []Ship {
	var ships = []Ship{}
	for i := 0; i < harborSize; i++ {
		maxSheepSize := 50
		ship := Ship{
			id:         i,
			containers: generateContainers(rand.Intn(maxSheepSize)),
			size:       maxSheepSize,
			harbor:     harbor,
		}
		ships = append(ships, ship)
	}
	return ships
}

func start(ships []Ship, harbor Harbor) {
	var wg sync.WaitGroup
	for i := 0; i < len(ships); i++ {
		wg.Add(1)
		harbor.c <- ships[i]
		go ships[i].run(&wg)
	}
	wg.Wait()
}

func main() {
	harborSize := 1000
	harbor := Harbor{
		c:          make(chan Ship, 10),
		piers:      generatePiers(harborSize),
		containers: generateContainers(harborSize / 2),
	}
	ships := generateShips(harbor, harborSize/5)

	start(ships, harbor)
}
