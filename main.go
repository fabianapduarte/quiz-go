package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"time"

	utils "github.com/fabianapduarte/quiz-go/utils"
)

type Question struct {
	Text    string
	Options []string
	Answer  int
}

type GameState struct {
	Name      string
	Points    int
	Questions []Question
}

func (g *GameState) Init() {
	fmt.Println("Seja bem vindo(a) ao quiz!")
	fmt.Print("Escreva seu nome: ")

	reader := bufio.NewReader(os.Stdin)
	name, err := reader.ReadString('\n')

	if err != nil {
		panic("Erro ao ler string")
	}

	g.Name = name

	fmt.Printf("Vamos ao jogo, %s", g.Name)
	fmt.Println("Você terá 1 minuto para responder todas as perguntas")
	fmt.Println()
}

func (g *GameState) ProcessCSV() {
	file, err := os.Open("data/quiz-go.csv")
	if err != nil {
		panic("Erro ao abrir arquivo CSV")
	}

	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		panic("Erro ao ler arquivo CSV")
	}

	for index, record := range records {
		if index > 0 {
			correctAnswer, _ := utils.ToInt(record[6])
			question := Question{
				Text:    record[0],
				Options: record[1:6],
				Answer:  correctAnswer,
			}

			g.Questions = append(g.Questions, question)
		}
	}
}

func (g *GameState) SetTimer() {
	timer := time.NewTimer(1 * time.Minute)
	<-timer.C

	fmt.Println("\n\n\033[31mTempo encerrado!\033[0m")
	g.PrintResult()

	os.Exit(0)
}

func (g *GameState) PrintResult() {
	fmt.Println("Fim do jogo!")

	if g.Points > 60 {
		fmt.Println("Resultado: \033[32mAprovado(a)!\033[0m")
	} else {
		fmt.Println("Resultado: \033[31mReprovado(a)!\033[0m")
	}
	fmt.Printf("Pontuação total: %d", g.Points)
}

func (g *GameState) Run() {
	go g.SetTimer()

	for index, question := range g.Questions {
		fmt.Printf("\033[34m %d. %s \033[0m\n", index+1, question.Text)

		for indexOption, option := range question.Options {
			fmt.Printf("[%d] %s\n", indexOption+1, option)
		}

		fmt.Print("\nDigite a alternativa: ")

		var answer int
		var err error

		for {
			reader := bufio.NewReader(os.Stdin)
			read, _ := reader.ReadString('\r')

			answer, err = utils.ToInt(read[:len(read)-1])
			if err != nil {
				fmt.Println(err.Error())
				fmt.Print("Digite a alternativa: ")
				continue
			}
			break
		}

		if answer == question.Answer {
			fmt.Println("\033[32mParabéns, você acertou!\033[0m")
			g.Points += 10
		} else {
			fmt.Println("\033[31mResposta incorreta :/\033[0m")
		}
		fmt.Println()
	}
}

func main() {
	game := &GameState{Points: 0}

	utils.ClearTerminal()

	go game.ProcessCSV()
	game.Init()
	game.Run()
	game.PrintResult()
}
