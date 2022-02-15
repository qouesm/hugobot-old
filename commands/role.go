package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var Role = Command{
	AppCommand: discordgo.ApplicationCommand{
		Name:        "role",
		Type:        discordgo.ChatApplicationCommand,
		Description: "Role management messages",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "create",
				Description: "Create a role management category",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "title",
						Description: "Title of message",
						Required:    true,
						Type:        discordgo.ApplicationCommandOptionString,
					},

					{
						Name:        "role",
						Description: "1st role (mentionable)",
						Required:    true,
						Type:        discordgo.ApplicationCommandOptionRole,
					},

					{
						Name:        "emoji",
						Description: "Emoji of 1st role",
						Required:    false,
						Type:        discordgo.ApplicationCommandOptionString,
					},

					{
						Name:        "description",
						Description: "Description of 1st role",
						Required:    false,
						Type:        discordgo.ApplicationCommandOptionString,
					},
				},
			},

			{
				Name:        "edit",
				Description: "Edit a role management message",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			},

			{
				Name:        "delete",
				Description: "Delete a role management message",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			},
		},
	},

	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}

		switch i.ApplicationCommandData().Options[0].Name {

		case "create":
			createRoleMessage(s, i)
		case "edit":
		case "delete":
		}

		// s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		// 	Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		// 	Data: &discordgo.InteractionResponseData{},
		// })
	},
}

func createRoleMessage(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// gather information
	// roleEmbed := NewEmbed().
	// 	SetTitle("Creating Role Message").
	// 	SetDescription("Enter the message title...").
	// 	// AddField("Emoji", "").
	// 	// AddField("Role", "").
	// 	SetColor(0xff0000).
	// 	InlineAllFields().MessageEmbed

	// ACK interation
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Println("could not create 'role create' interaction")
	}

	// get/create master message
	var masterMessage *discordgo.Message
	// TODO change db name
	// TODO PROBLEM This code walks the db but only checks the first *.json file in a dir for whether it is masterMessage.json or not
	// if there are other files this doesn't work
	err = filepath.WalkDir("jsondb/", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Panicln(err.Error())
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(d.Name()) != ".json" {
			log.Panicln("unexpected non-json in db", err.Error())
		}

		// log.Println(path, d.Name())

		// looking whether [ i.GuildID, i.ChannelID ] masterMessage.json exists
		if d.Name() == "masterMessage.json" {
			location := strings.Split(path, "/")
			location = location[1 : len(location)-1]

			log.Println(location)

			if reflect.DeepEqual(location, []string{i.GuildID, i.ChannelID}) {
				// master exists
				log.Println("master exist")
				// unmarshall it to masterMessage
				masterMessageJSON, err := ioutil.ReadFile(path)
				if err != nil {
					log.Panicln("could not read file", err.Error())
				}

				err = json.Unmarshal(masterMessageJSON, &masterMessage)
				if err != nil {
					log.Panicln("could not unmarshal", err.Error())
				}

			} else {
				// master does not exist
				log.Println("master nonexist")
				masterMessage = createMaster(s, i)

				// DEV TEMPORARY put this message in db so above if can check
				// masterMessageFile, err := os.Create(fmt.Sprintf("db/%s/%s/masterMessage.json", i.GuildID, i.ChannelID))

				if err := os.MkdirAll(filepath.Dir(path), 0770); err != nil {
					log.Println("problem creating file,", err)
				}
				masterMessageFile, err := os.Create(path)
				if err != nil {
					log.Panicln("couldn't save file to db", err.Error())
				}
				defer masterMessageFile.Close()

				masterMessageJSON, err := json.Marshal(masterMessage)
				if err != nil {
					log.Println("could not marshal message,", err)
				}
				_, err = masterMessageFile.Write(masterMessageJSON)
				if err != nil {
					log.Println("could not write file", err.Error())
				}

				log.Println("ok")
				return io.EOF
			}
		}
		return nil
	})
	if err == io.EOF {
		err = nil
	}

	// // probably doesn't initialize masterMessage correctly
	// if _, err := os.Stat(fmt.Sprintf("db/%s", i.GuildID)); err != nil {
	// 	if os.IsExist(err) {
	// 		// guild has master messages
	// 		log.Println("guild has master messages")
	// 		if _, err := os.Stat(fmt.Sprintf("db/%s/%s", i.GuildID, i.ChannelID)); err != nil {
	// 			if os.IsExist(err) {
	// 				if _, err := os.Stat(fmt.Sprintf("db/%s/%s/masterMessage.json", i.GuildID, i.ChannelID)); err != nil {
	// 					if os.IsNotExist(err) {
	// 						log.Panicln("channel dir exists without masterMessage, ", err)
	// 					}
	// 				}
	// 				// channel has master message
	// 				log.Println("channel has master message")
	// 				// get master message
	// 				masterMessageJSON, err := ioutil.ReadFile(fmt.Sprintf("db/%s/%s/masterMessage.json", i.GuildID, i.ChannelID))
	// 				if err != nil {
	// 					log.Panicln("could not read file, ", err)
	// 				}

	// 				err = json.Unmarshal(masterMessageJSON, &masterMessage)
	// 				if err != nil {
	// 					log.Panicln("could not unmarshal, ", err)
	// 				}
	// 			} else {
	// 				// channel does not have master message
	// 				log.Println("channel does not have a master message")
	// 				// create channel folder
	// 				err = os.Mkdir(fmt.Sprintf("db/%s/%s", i.GuildID, i.ChannelID), os.ModeDir)
	// 				if err != nil {
	// 					log.Println("could not create channel dir, ", err)
	// 				}
	// 				// create master message
	// 				masterMessage = createMaster(s, i)
	// 			}
	// 		} else {
	// 			// guild does not have any master messages
	// 			log.Println("guild does not have any master messages")
	// 			// create guild folder
	// 			err = os.Mkdir(fmt.Sprintf("db/%s", i.GuildID), os.ModeDir)
	// 			if err != nil {
	// 				log.Println("could not create guild dir, ", err)
	// 			}
	// 			// create channel folder
	// 			err = os.Mkdir(fmt.Sprintf("db/%s/%s", i.GuildID, i.ChannelID), os.ModeDir)
	// 			if err != nil {
	// 				log.Println("could not create channel dir, ", err)
	// 			}
	// 			// create master message
	// 			masterMessage = createMaster(s, i)
	// 		}
	// 	}
	// }

	// verify masterMessage's selectMenu is not full
	log.Printf("%#v", masterMessage)
	// log.Printf("%#v", &masterMessage)
	// log.Println(masterMessage.Embeds[0].Title)
	log.Println(len(masterMessage.Components))
	masterMenu, ok := masterMessage.Components[0].(discordgo.SelectMenu)
	if !ok {
		log.Panicln("unexpected masterMessage, ", err)
	}
	log.Printf("%#v", masterMenu)
	if len(masterMenu.Options) >= 25 {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Cannot create category: Menu limit reached (25) for this channel",
				Flags:   1 << 6,
			},
		})
		if err != nil {
			log.Println("could not respond to role create, ", err)
			return
		}

		time.Sleep(time.Second * 10)

		err = s.InteractionResponseDelete(s.State.User.ID, i.Interaction)
		if err != nil {
			log.Println("could not delete role create response, ", err)
		}
		return
	}

	// get command info
	createOptions := i.ApplicationCommandData().Options[0].Options

	var (
		title       string
		role        *discordgo.Role
		description string
		emoji       discordgo.ComponentEmoji
	)

	for _, option := range createOptions {
		switch option.Name {
		case "title":
			title = option.StringValue()
		case "role":
			role = option.RoleValue(s, i.GuildID)
		case "description":
			description = option.StringValue()
		case "emoji":
			emoji = discordgo.ComponentEmoji{
				Name: option.StringValue(),
			}
		}
	}

	// verify new menu's title is not pre-existing in the masterMenu
	for _, option := range masterMenu.Options {
		if strings.EqualFold(option.Label, title) {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "A category with this name already exists in this channel",
					Flags:   1 << 6,
				},
			})
			if err != nil {
				log.Println("could not respond to role create, ", err)
				return
			}
			time.Sleep(time.Second * 10)
			err = s.InteractionResponseDelete(s.State.User.ID, i.Interaction)
			if err != nil {
				log.Println("could not delete role create response, ", err)
			}
			return
		}
	}

	// create new menu
	roleMenu := discordgo.SelectMenu{
		CustomID:    i.ID,
		Placeholder: title,
		Options: []discordgo.SelectMenuOption{
			{
				Value: "r1",
				Label: role.Name,
			},
		},
		MinValues: 0,
		MaxValues: 1,
	}

	if description != "" {
		roleMenu.Options[0].Description = description
	}
	if emoji.Name != "" {
		roleMenu.Options[0].Emoji = emoji
	}

	// add new menu info to master
	// insert at len-1 (append just before last element)
	mmLen := len(masterMenu.Options)
	masterMenu.Options = append(masterMenu.Options[:mmLen], masterMenu.Options[mmLen-1:]...)
	masterMenu.Options[mmLen-1] = discordgo.SelectMenuOption{
		// Value:       fmt.Sprintf("c%d", mmLen-1), // should be c1, c2, c3...
		Label:       title,
		Value:       strings.ToLower(title),
		Description: description,
		Emoji:       emoji,
	}
	// put menu back into message
	masterMessage, err = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		ID:         masterMessage.ID,
		Channel:    masterMessage.ChannelID,
		Embeds:     []*discordgo.MessageEmbed{infoEmbed()},
		Components: []discordgo.MessageComponent{masterMenu},
	})
	if err != nil {
		log.Println("could not update masterMessage, ", err)
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Could not add category",
				Flags:   1 << 6,
			},
		})
		if err != nil {
			log.Println("could not respond to role create, ", err)
			return
		}

		time.Sleep(time.Second * 10)

		err = s.InteractionResponseDelete(s.State.User.ID, i.Interaction)
		if err != nil {
			log.Println("could not delete role create response, ", err)
		}
		return
	}

	// marshal masterMessage
	masterMessageFile, err := os.Create(fmt.Sprintf("db/%s/%s/masterMessage.json", i.GuildID, i.ChannelID))
	if err != nil {
		log.Println("problem creating file,", err)
	}
	defer masterMessageFile.Close()

	masterMessageJSON, err := json.Marshal(masterMessage)
	if err != nil {
		log.Println("could not marshal message,", err)
	}
	masterMessageFile.Write(masterMessageJSON)

	// marshal new roleMenu
	roleMenuFile, err := os.Create(fmt.Sprintf("db/%s/%s/%s%s.json", i.GuildID, i.ChannelID, title, roleMenu.CustomID))
	if err != nil {
		log.Println("problem creating file,", err)
	}
	defer roleMenuFile.Close()

	roleMenuJSON, err := json.Marshal(roleMenu)
	if err != nil {
		log.Println("could not marshal message,", err)
	}
	roleMenuFile.Write(roleMenuJSON)

	// add handler
	// components.RoleMenuOnce(s, roleMessage)

	// clean up interaction
	err = s.InteractionResponseDelete(s.State.User.ID, i.Interaction)
	if err != nil {
		log.Println("could not delete 'role create' interaction", err)
	}

}

// // role adding hook
// s.AddHandler(func(s *discordgo.Session, mr *discordgo.MessageReactionAdd) {
// 	if mr.UserID == s.State.User.ID {
// 		return
// 	}
// 	if mr.MessageID == save.Msg.ID {
// 		number := emojiNum[mr.Emoji.APIName()]
// 		role := save.Roles[number]
// 		err1 := s.GuildMemberRoleAdd(mr.GuildID, mr.UserID, role.ID)
// 		user, err2 := s.User(mr.UserID)
// 		if err1 != nil {
// 			if err2 != nil {
// 				log.Println("Couldn't get user's struct,", err1)
// 				log.Println("Couldn't add role,", err)
// 				return
// 			}
// 			log.Println("Couldn't add role:", user.Username, ",", err)
// 			return
// 		}
// 		log.Println("rr " + user.Username + " + @" + role.Name)
// 		dm, err := s.UserChannelCreate(mr.UserID)
// 		if err != nil {
// 			log.Println("could not create dm channel,", err)
// 			return
// 		}
// 		g, err := s.Guild(mr.GuildID)
// 		if err != nil {
// 			log.Println("could not get current guild,", err)
// 			return
// 		}
// 		_, err = s.ChannelMessageSend(dm.ID, g.Name+": ADDED @"+role.Name)
// 		if err != nil {
// 			log.Println("could not dm user,", err)
// 			return
// 		}
// 	}
// })

// // role deletion hook
// s.AddHandler(func(s *discordgo.Session, mr *discordgo.MessageReactionRemove) {
// 	if mr.UserID == s.State.User.ID {
// 		return
// 	}
// 	if mr.MessageID == save.Msg.ID {
// 		number := emojiNum[mr.Emoji.APIName()]
// 		role := save.Roles[number]
// 		err1 := s.GuildMemberRoleRemove(mr.GuildID, mr.UserID, role.ID)
// 		user, err2 := s.User(mr.UserID)
// 		if err1 != nil {
// 			if err2 != nil {
// 				log.Println("Couldn't get user's struct,", err1)
// 				log.Println("Couldn't add role,", err)
// 				return
// 			}
// 			log.Println("Couldn't del role:", user.Username, ",", err)
// 			return
// 		}
// 		log.Println("rr " + user.Username + " - @" + role.Name)
// 		dm, err := s.UserChannelCreate(mr.UserID)
// 		if err != nil {
// 			log.Println("could not create dm channel,", err)
// 			return
// 		}
// 		g, err := s.Guild(mr.GuildID)
// 		if err != nil {
// 			log.Println("could not get current guild,", err)
// 			return
// 		}
// 		_, err = s.ChannelMessageSend(dm.ID, g.Name+": REMOVED @"+role.Name)
// 		if err != nil {
// 			log.Println("could not dm user,", err)
// 			return
// 		}
// 	}
// })

func createMaster(s *discordgo.Session, i *discordgo.InteractionCreate) *discordgo.Message {
	msg, _ := s.ChannelMessageSendComplex(i.ChannelID, masterMessage())
	// log.Printf("%#v", msg)
	return msg
}

func masterMessage() *discordgo.MessageSend {
	return &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{infoEmbed()},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						CustomID:    "master",
						Placeholder: "Add or Remove Roles",
						MinValues:   0,
						MaxValues:   1,
						Options: []discordgo.SelectMenuOption{
							removeOption(),
						},
					},
				},
			},
		},
	}
}

func infoEmbed() *discordgo.MessageEmbed {
	return NewEmbed().
		SetTitle("Role Management").
		SetDescription("Select a category").
		SetFooter("Contact @qouesm#9558 with issues").
		MessageEmbed
}

func removeOption() discordgo.SelectMenuOption {
	return discordgo.SelectMenuOption{
		Label:       "Remove Roles",
		Value:       "remove",
		Description: "Remove any role added from any menu",
		Emoji: discordgo.ComponentEmoji{
			Name: "❌", // cross mark U+274C
		},
	}
}
