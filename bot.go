package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

const nome = "Inversor"


func inverter(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func main() {
	conn, err := net.Dial("tcp", "localhost:4444")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	
	fmt.Fprintln(conn, `\bot `+nome)
	fmt.Println("Entrando como Bot @" + nome + ". Aguardando mensagens privadas...")

	respondeu := false
	servidor := bufio.NewScanner(conn)
	for servidor.Scan() {
		linha := servidor.Text()

		
		
		
		if strings.HasPrefix(linha, "Erro:") {
			fmt.Println(linha)
			if !respondeu {
				return
			}
			continue
		}

		
		de, texto, ok := strings.Cut(linha, " disse em privado: ")
		if !ok {
			continue
		}

		resposta := inverter(texto)
		fmt.Println("Mensagem de: " + de)
		fmt.Println("Recebi: " + texto)
		fmt.Println("Resposta: " + resposta)
		fmt.Println()

		fmt.Fprintln(conn, `\msg `+de+" "+resposta)
		respondeu = true
	}
	fmt.Println("Conexão encerrada.")
}