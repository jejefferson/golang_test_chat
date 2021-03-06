package main
 
import (
    "net"
    "bufio"
    "fmt"
    "strings"
    "time"
    "bytes"
    "strconv"
    "sort"
)

type User struct {
	socket net.Conn 
	name string
	addr string
}

const (
	NORMAL int = iota // для сообщений обычного характера
	SYSTEM // для системных сообщений
)

var users []User

func main() {
    server, err := net.Listen("tcp", ":3333")
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
 
func clientConns(listener net.Listener) chan User { // приём клинетов-юзеров, выделение в отдельный поток каждого
    ch := make(chan User)
    i := 0
    var T User
    go func() {
        for {
            client, err := listener.Accept()
            if client == nil {
                fmt.Printf("couldn't accept: " + err.Error())
                continue
            }
            T.socket = client
            T.addr = fmt.Sprintf("%v", client.RemoteAddr())
            i++
            fmt.Printf("%d: %v <-> %v\n", i, client.LocalAddr(), client.RemoteAddr())
            ch <- T
        }
    }()
    return ch
}
 
func handleConn(sender User) {
    b := bufio.NewReader(sender.socket)
    sender.socket.Write([]byte("Привет, назови себя \n"))
    answer,_ := b.ReadBytes('\n')
	if len(answer) == 0 { // если пользователь вышел на этом месте из программы
		sender.goodbye()
		return
	}
    ind := strings.Index(string(answer), "::")
	name := string(answer)[ind+2:len(answer)-1]
	sender.changeName([]byte(name))
	sender.socket.Write([]byte("Напиши /info чтобы узнать команды\n"))
	nick := name
	varnick := []string {"nick", "name", "nck"}
	sort.Strings(varnick)
	if len(nick) == 1 {
		nick = sender.addr
	}
	line := []byte(" зашёл в чат\n")
	sendAll(sender, line, SYSTEM)
    for { //бесконечный цикл, приём сообщений на сервер
        line, err := b.ReadBytes('\n')
        if err != nil {
            break
        }
		if len(line) == 1 {
			continue
		}
		ind := strings.Index(string(line), "::")
		linekine := strings.TrimSpace(string(line)[ind+2:])
		if strings.HasPrefix(linekine, "/") { // проверка служебного сообщения
			nomer := strings.Index(linekine, " ")
			comand := linekine
			if nomer > 0 {
				comand = linekine[1:nomer]
			} else {
				comand = linekine[1:]
			}
			switch comand {
				case "ping", "pin", "pi", "p", "pong", "png":
					ind := bytes.Index(line, []byte("::"))
					L := strings.TrimSpace(string(line[:ind]))
					pin, _:= strconv.ParseInt(L, 10, 64)
					ping := time.Now().UnixNano() - pin
					timeping := fmt.Sprintf("Pong, твоё время %v наносекунд\n", ping)
					sender.socket.Write([]byte(timeping))
					continue
				case "users", "user", "usrs", "usr", "list", "userlist":
					for _, user := range(users) {
						sender.socket.Write([]byte(fmt.Sprintf("%v ", user.name)))
					}
					sender.socket.Write([]byte("\n"))
					continue
				case "info", "inf", "ino", "ifo":
					text := "/info /ping /nick /users \n"
					sender.socket.Write([]byte(text))
					continue
				case "nick", "nik", "name", "nickname", "nam", "nme":
					indkine := strings.Index(linekine, " ")
					kine := linekine[1:indkine]
					dik := sort.SearchStrings(varnick,kine)
					if dik > 0 && dik <len(varnick) {
						ind := strings.LastIndex(string(line), " ")
						name := string(line)[ind+1:]
						sender.changeName([]byte(name))
						fmt.Printf(name)
					}
					continue
				default:
					text := "неверная команда, /info чтобы просмотреть все команды \n"
					sender.socket.Write([]byte(text))
				continue
			}
		}
		sendAll(sender, line, NORMAL)
    } 
   defer sender.goodbye()
}

func (sender *User) goodbye() { // удаление отключенного юзера
	line := []byte(" вышел из чата\n")
	sendAll(*sender, line, SYSTEM)
	for index, user := range(users) {
		if user==*sender {
			users = append(users[:index], users[index+1:]...)
			(*sender).socket.Close()
			fmt.Println("нас осталось:", len(users))
			nickk := sender.name
			if len(nickk) == 1 {
				nickk = sender.addr
			}
		}
	}
}

func (sender *User) changeName(answer []byte) { // смена имени юзера 
	answer1 := strings.TrimSpace(string(answer))
	if len(answer1) > 1 {
		for index, user := range(users){
			if user==*sender {
				(*sender).name = string(answer1)[:len(answer1)]
				users[index] = *sender
			}
		}
	}
}

func sendAll (sender User, line []byte, msg_mode int) { //рассылка сообщений
	for _, user := range(users) { 
		if user.socket != sender.socket {
			var name []byte
			if len(sender.name) == 1 || len(sender.name) > 1 {
				name = []byte(sender.name)
			} else {
				name = []byte(sender.addr)
			}
			ind := strings.Index(string(line), "::")
			h, m, s := time.Now().Clock()
			ttime := fmt.Sprint("[", h, ":", m, ":", s, "]")
			var mess string
			if msg_mode == NORMAL {
				mess = fmt.Sprint(ttime, " <", string(name), "> ", string(line)[ind+2:])
			} else {
				mess = fmt.Sprint(ttime, " *", string(name), " ", string(line)[ind+2:])
			}
			user.socket.Write([]byte(mess))
		}
	}
}
