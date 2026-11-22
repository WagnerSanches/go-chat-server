package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:4444")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	servidor := bufio.NewScanner(conn)    
	teclado := bufio.NewScanner(os.Stdin) 

	
	
	
	for {
		fmt.Print("Escolha seu apelido: ")
		if !teclado.Scan() {
			return
		}
		fmt.Fprintln(conn, `\user `+teclado.Text())

		if !servidor.Scan() {
			log.Fatal("o servidor encerrou a conexão")
		}
		resposta := servidor.Text()
		fmt.Println(resposta)
		if !strings.HasPrefix(resposta, "Erro:") {
			break
		}
	}
	fmt.Println(`Comandos: \msg <texto> | \msg @<apelido> <texto> | \changenick <apelido> | \exit`)

	
	go func() {
		for teclado.Scan() {
			linha := teclado.Text()
			if linha == `\exit` {
				break
			}
			fmt.Fprintln(conn, linha)
		}
		fmt.Fprintln(conn, `\exit`) 
	}()

	
	for servidor.Scan() {
		fmt.Println(servidor.Text())
	}
	fmt.Println("Conexão encerrada.")
}