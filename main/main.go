package main

import (
	"fmt"
	"gelato/main/message"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"syscall"
)

var (
	Token string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env 파일 로딩 중 에러 발생", err)
		os.Exit(1)
	}
	Token = os.Getenv("TOKEN")
	if Token == "" {
		fmt.Println("TOKEN을 읽지 못함")
		os.Exit(1)
	}
}

func main() {
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("디스코드 세션 생성 중 에러 발생", err)
		return
	}
	dg.AddHandler(message.Create)
	err = dg.Open()
	if err != nil {
		fmt.Println("연결 중 에러 발생", err)
		return
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	_ = dg.Close()
}
