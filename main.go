package main

import (
	"fmt"
	"net"
	"sort"
)

// Поиск открытого порта, в случае когда открыт передаем его, если закрыт 0
func worker(ports chan int, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("scanme.nmap.org:%d", p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		_ = conn.Close()
		results <- p
	}
}

func main() {
	ports := make(chan int, 100)
	results := make(chan int)

	var openports []int

	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}
	//Заполняем канал для дальнейшей проверки в worker
	go func() {
		for i := 0; i <= 1024; i++ {
			ports <- i
		}
	}()

	//Проверяем полученные результаты и открытые порты записываются в срез
	for i := 0; i <= 1024; i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}
	//Закрываем каналы, сортируем полученный срез с открытыми портами, выводим их в консоль
	close(ports)
	close(results)
	sort.Ints(openports)
	for _, port := range openports {
		fmt.Printf("Port %d open\n", port)
	}
}
