package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

type client chan<- string // canal de mensagem

type usuario struct {
	apelido string
	ch      chan string
	bot     bool
}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

var (
	mu       sync.Mutex
	usuarios = make(map[string]*usuario)
	publicas = make(chan string) // canal de mensagens públicas
)

func broadcaster() {
	for msg := range publicas {
		mu.Lock()
		for _, u := range usuarios {
			if !u.bot { // bots nao recebem mensagens publicas
				u.ch <- msg
			}
		}
		mu.Unlock()
	}
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string)
	go clientWriter(conn, ch)
	input := bufio.NewScanner(conn)
	u := identificar(input, ch)
	if u == nil {
		close(ch)
		conn.Close()
		return
	}

	tipo := "Usuário"
	anuncio_entrada := "Usuário @" + u.apelido + " entrou!"
	if u.bot {
		tipo = "Bot"
		anuncio_entrada = "Bot @" + u.apelido + "se conectou!"
	}
	fmt.Println(anuncio_entrada)
	publicas <- anuncio_entrada

	for input.Scan() {
		cmd, resto, _ := strings.Cut(input.Text(), "")
		if cmd == "/exit" {
			break
		}

		switch cmd {
		case "/msg":

			if strings.HasPrefix(resto, "@") {

			} else if resto != "" {
				publicas <- "@" + u.apelido + ": " + resto
			}

		case "/changenick":

		}
	}

	remover(u)
	saida := tipo + " @" + u.apelido + " saiu"
	fmt.Println(saida)
	publicas <- saida
	close(ch)
	conn.Close()
}

func registrar(u *usuario) bool {
	mu.Lock()
	defer mu.Unlock()
	if _, existe := usuarios[u.apelido]; existe {
		return false
	}

	usuarios[u.apelido] = u
	return true
}

func remover(u *usuario) {
	mu.Lock()
	delete(usuarios, u.apelido)
	mu.Unlock()
}

func apelidoValido(apelido string) bool {
	return apelido != "" && !strings.ContainsAny(apelido, ` @/`)
}

func identificar(input *bufio.Scanner, ch chan string) *usuario {
	for input.Scan() {
		// primeiro pega o valor
		cmd, apelido, _ := strings.Cut(input.Text(), " ")
		if cmd != "/bot" && cmd != "/user" { // se nao for nenhuma das opcoes
			// escreva no channel
			ch <- `Erro: Identifique-se comm /user <apelido>`
			continue
		}

		// agora valida o apelido
		if !apelidoValido(apelido) {
			ch <- `Erro: apelido inválido`
			continue
		}

		// verifica se ja existe
		u := &usuario{apelido: apelido, ch: ch, bot: cmd == `/bot`}
		if registrar(u) {
			ch <- `Erro: este apelido já está em uso`
			continue
		}

		return u
	}

	return nil
}

func main() {
	fmt.Println("Iniciando servidor...")
	listener, err := net.Listen("tcp", "localhost:4444")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Portas abertas!")
	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}
