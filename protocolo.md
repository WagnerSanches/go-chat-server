# Protocolo do chat

Regras de comunicação entre o servidor (`server.go`), o cliente (`client.go`) e o bot (`bot.go`).

- Conexão TCP em `localhost:4444`.
- Texto puro, **uma mensagem por linha** (terminada em `\n`).

## 1. Identificação

A primeira linha enviada depois de conectar diz quem é o programa.

| Quem conecta | Primeira linha | Anúncio do servidor |
|---|---|---|
| Cliente | `\user <apelido>` | `Usuário @<apelido> entrou!` |
| Bot | `\bot <apelido>` | `Bot @<apelido> se conectou!` |

O usuário não digita `\user`: o cliente pergunta o apelido e envia essa linha sozinho.

Se a identificação falhar (apelido em uso ou inválido), o servidor responde com o erro e espera uma nova linha de identificação.

### Regras do apelido

- Não pode ser vazio.
- Não pode ter espaço, `@` nem `\`.
- Precisa ser único entre todos os conectados (usuários e bots).

## 2. Comandos (cliente → servidor)

| Linha enviada | O que o servidor faz |
|---|---|
| `\msg <texto>` | Envia uma mensagem pública. |
| `\msg @<apelido> <texto>` | Envia uma mensagem privada para `@<apelido>`. |
| `\changenick <novo>` | Troca o apelido. |
| `\exit` | Tira o usuário do chat e fecha a conexão. |
| Texto sem comando | Nada: a linha é ignorada. |

Fechar o programa ou perder a conexão sem `\exit` é tratado igual ao `\exit`.

## 3. Mensagens do servidor

| Evento | Texto enviado | Quem recebe | Console do servidor |
|---|---|---|---|
| Usuário entra | `Usuário @abc entrou!` | Todos os usuários | Sim |
| Bot entra | `Bot @abc se conectou!` | Todos os usuários | Sim |
| Mensagem pública | `@abc: <texto>` | Todos os usuários, inclusive quem enviou | Não |
| Mensagem privada | `@abc disse em privado: <texto>` | Só o destinatário | Sim |
| Usuário sai | `Usuário @abc saiu` | Todos os usuários | Sim |
| Bot sai | `Bot @abc saiu` | Todos os usuários | Sim |
| Troca de apelido | `Usuário @abc agora é @xyz` | Todos os usuários | Sim |

**Bots recebem apenas mensagens privadas.** Anúncios e mensagens públicas vão só para os usuários.

## 4. Erros

Os erros vão só para quem causou.

| Situação | Resposta |
|---|---|
| Apelido já em uso (na entrada ou no `\changenick`) | `Erro: Este nick já está em uso` |
| Apelido inválido | `Erro: apelido inválido` |
| Primeira linha não é `\user` nem `\bot` | `Erro: identifique-se com \user <apelido>` |
| Privado para um apelido que não existe | `Erro: usuário @xyz não encontrado` |

## 5. Bot Inversor

1. Conecta e envia `\bot Inversor`.
2. Para cada linha `@abc disse em privado: <texto>`, inverte o texto e responde com `\msg @abc <texto invertido>`.
3. Ignora qualquer outra linha.
4. No próprio console, imprime:

```
Mensagem de: @abc
Recebi: abcdef
Resposta: fedcba
```

## Exemplo de sessão

```
Zezinho  → servidor : \user Zezinho
servidor → usuários : Usuário @Zezinho entrou!
Inversor → servidor : \bot Inversor
servidor → usuários : Bot @Inversor se conectou!
Zezinho  → servidor : \msg Olá pessoal
servidor → usuários : @Zezinho: Olá pessoal
Zezinho  → servidor : \msg @Inversor abcdef
servidor → Inversor : @Zezinho disse em privado: abcdef
Inversor → servidor : \msg @Zezinho fedcba
servidor → Zezinho  : @Inversor disse em privado: fedcba
Zezinho  → servidor : \exit
servidor → usuários : Usuário @Zezinho saiu
```