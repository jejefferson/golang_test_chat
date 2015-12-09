package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"flag"
	"strconv"
)

type User struct {
	socket net.Conn
	name   string
	address string
}

func (s *User) setName(name string) {
	for index, user := range users {
		if user == *s {
			(*s).name = name
			users[index] = *s
		}
	}
	return
}

func (s *User) Close() {
	for index, user := range users {
		if user == *s {
			(*s).socket.Close()
			users = append(users[:index], users[index+1:]...)
		}
	}
}

var users []User

func main() {
	port := flag.Int("port", 3333, "You should specify address")
	flag.Parse()
	server, err := net.Listen("tcp", ":"+strconv.Itoa(*port))
	if server == nil {
		panic("couldn't start listening: " + err.Error())
	}
	conns := clientConns(server)
	for {
		client := <-conns
		users = append(users, client)
		go handleConn(client)
	}
}

func clientConns(listener net.Listener) chan User {
	ch := make(chan User)
	i := 0
	go func() {
		for {
			client, err := listener.Accept()
			if client == nil {
				fmt.Printf("couldn't accept: " + err.Error())
				continue
			}
			i++
			addr := fmt.Sprintf("%v", client.RemoteAddr())
			fmt.Printf("%d %v <=> %v\n", i, client.LocalAddr(), addr)
			ch <- User{socket: client, address: addr}
		}
	}()
	return ch
}

func handleConn(sender User) {
	defer sender.Close()
	b := bufio.NewReader(sender.socket)
	sender.socket.Write([]byte("You entred to chat. Hello.\nPlease enter you name: "))
	username, _ := b.ReadBytes('\n')
	if len(username) > 1 {
		sender.setName(string(username)[:len(username)-1])
	} else {
		sender.socket.Write([]byte("Using your ip:port as name, you can type \"nick: yournick\" to change it\n"))
	}
	for {
		line, err := b.ReadBytes('\n')
		if err != nil {
			break
		}
		stringline := string(line)
		if strings.HasPrefix(stringline, "nick:") {
			ind := strings.Index(stringline, ":")
			name := strings.TrimSpace(strings.Trim(stringline[ind+1:], " "))
			sender.setName(name)
			continue
		}
		for _, user := range users {
			if user.socket != sender.socket {
				if sender.name != "" {
					user.socket.Write([]byte(sender.name))
				} else {
					user.socket.Write([]byte(sender.address))
				}
				user.socket.Write([]byte("> "))
				user.socket.Write(line)
			}
		}
	}
}
