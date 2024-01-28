package main

import (
	"flag"
	"fmt"
	"gelato/main/message"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	RemoveCommands           = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
	dg                       *discordgo.Session
	Token                    string
	dmPermission                   = false
	defaultMemberPermissions int64 = discordgo.PermissionManageServer
	commands                       = []*discordgo.ApplicationCommand{
		{
			Name:                     "t",
			Description:              "Manage",
			DefaultMemberPermissions: &defaultMemberPermissions,
		},
		{

			Name:                     "options",
			Description:              "옵션",
			DefaultMemberPermissions: &defaultMemberPermissions,
		},
		{
			Name:        "options",
			Description: "사용가능한 커맨드 목록",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "t",
					Description: "시간확인",
					Required:    true,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"t": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "시간확인커맨드",
				},
			})
		},
		"options": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options

			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			margs := make([]interface{}, 0, len(options))
			msgFormat := "젤라또 봇 사용 명령어입니다\n"

			if option, ok := optionMap["t"]; ok {
				margs = append(margs, option.StringValue())
				msgFormat += "> t: %s\n"
			}
			fmt.Println(msgFormat, margs, options)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(msgFormat, margs...),
				},
			})
		},
	}
)

func loadENV() {
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

func createSession() {
	var err error
	dg, err = discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("디스코드 세션 생성 중 에러 발생", err)
		return
	}
}

func addHandlers() {
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func init() {
	loadENV()
	createSession()
	addHandlers()
}
func main() {
	dg.AddHandler(message.Create)
	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := dg.Open()
	if err != nil {
		fmt.Println("연결 중 에러 발생", err)
		return
	}
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := dg.ApplicationCommandCreate(dg.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer dg.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if *RemoveCommands {
		log.Println("Removing commands...")
		// We need to fetch the commands, since deleting requires the command ID.
		// We are doing this from the returned commands on line 375, because using
		// this will delete all the commands, which might not be desirable, so we
		// are deleting only the commands that we added.
		registeredCommands, err := dg.ApplicationCommands(dg.State.User.ID, "")
		if err != nil {
			log.Fatalf("Could not fetch registered commands: %v", err)
		}

		for _, v := range registeredCommands {
			err := dg.ApplicationCommandDelete(dg.State.User.ID, "", v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Gracefully shutting down.")

}
