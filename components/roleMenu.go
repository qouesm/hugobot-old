package components

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func roleMenu(s *discordgo.Session) {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		roleMessages, err := ioutil.ReadDir("db")
		if err != nil {
			log.Panicln("could not read 'db/', ", err)
		}

		// future option: build a map and then just handle that map
		for _, messageInfo := range roleMessages {
			roleMessageJSON, err := ioutil.ReadFile(fmt.Sprintf("db/%s", messageInfo.Name()))
			if err != nil {
				log.Panicln("could not read file, ", err)
			}

			var roleMessage *discordgo.Message
			err = json.Unmarshal(roleMessageJSON, &roleMessage)
			if err != nil {
				log.Panicln("could not unmarshal, ", err)
			}

			// add role
			data := i.MessageComponentData()
			log.Printf("%#v", data)

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "interaction received\nstill WIP",
					Flags:   1 << 6,
				},
			})
			if err != nil {
				log.Printf("could not ACK roleMenu interaction on message %s, %v", data.CustomID, err)
				continue
			}

			time.Sleep(time.Second * 10)

			err = s.InteractionResponseDelete(s.State.User.ID, i.Interaction)
			if err != nil {
				log.Println("could not delete roleMenu ACK interaction ", err)
			}
		}
	})
}

func RoleMenuOnce(s *discordgo.Session, m *discordgo.Message) {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// add role
		data := i.MessageComponentData()
		log.Printf("%#v", data)

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "interaction received\nstill WIP",
				Flags:   1 << 6,
			},
		})
		if err != nil {
			log.Printf("could not ACK roleMenu interaction on message %s, %v", data.CustomID, err)
			return
		}

		time.Sleep(time.Second * 10)

		err = s.InteractionResponseDelete(s.State.User.ID, i.Interaction)
		if err != nil {
			log.Println("could not delete roleMenu ACK interaction ", err)
		}
	})
}
