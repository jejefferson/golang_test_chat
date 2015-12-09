package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
)

// передаём функции наше соединение с сервером
// читаем соединение и печатаем все сообщения
func readSocket(connection net.Conn) {
	reader := bufio.NewReader(connection)
	for { //в бескончном цикле
		temp, _ := reader.ReadString('\n') //читаем по одной строчке
		fmt.Print(temp) // и пишем в консоль
	}
}

func main() {
    conn, err := net.Dial("tcp", "127.0.0.1:3333") // соединяемся с нашим сервером, который должен быть запущен
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
	go readSocket(conn) // запускаем в потоке получение сообщений с сервера и печать
	reader := bufio.NewReader(os.Stdin) // создаём читальщика из консоли, os.Stdin обозначает консоль, в Си это файловый дскприптор 0
	for { // в бесконечном цикле
		temp, _ := reader.ReadString('\n') // читаем строчку из консоли
		conn.Write([]byte(temp)) // и отправляем на сервер

	}
}
