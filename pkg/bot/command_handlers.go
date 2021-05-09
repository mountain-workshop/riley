package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"k8s.io/klog"
)

var (
	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"help": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			klog.V(4).Info("handling help interaction")
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionApplicationCommandResponseData{
					Content: `I'm here to help you track points. Here's the various commands I support:
					/help - This command`,
				},
			})
		},
		"create-team": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			klog.V(4).Info("handling create-team")
			margs := []interface{}{
				i.Data.Options[0].RoleValue(nil, "").ID,
			}

			msgformat := " We are going to register this team:\n"
			msgformat += "> role-option: <@&%s>\n"

			if len(i.Data.Options) >= 2 {
				margs = append(margs, i.Data.Options[1].ChannelValue(nil).ID)
				msgformat += "> team-name: %s\n"
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				// Ignore type for now, we'll discuss them in "responses" part
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionApplicationCommandResponseData{
					Content: fmt.Sprintf(
						msgformat,
						margs...,
					),
				},
			})
		},
		// "subcommands": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// 	content := ""

		// 	// As you can see, the name of subcommand (nested, top-level) or subcommand group
		// 	// is provided through arguments.
		// 	switch i.Data.Options[0].Name {
		// 	case "subcmd":
		// 		content =
		// 			"The top-level subcommand is executed. Now try to execute the nested one."
		// 	default:
		// 		if i.Data.Options[0].Name != "scmd-grp" {
		// 			return
		// 		}
		// 		switch i.Data.Options[0].Options[0].Name {
		// 		case "nst-subcmd":
		// 			content = "Nice, now you know how to execute nested commands too"
		// 		default:
		// 			// I added this in the case something might go wrong
		// 			content = "Oops, something gone wrong.\n" +
		// 				"Hol' up, you aren't supposed to see this message."
		// 		}
		// 	}
		// 	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		// 		Type: discordgo.InteractionResponseChannelMessageWithSource,
		// 		Data: &discordgo.InteractionApplicationCommandResponseData{
		// 			Content: content,
		// 		},
		// 	})
		// },
		// "responses": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// 	// Responses to a command are very important.
		// 	// First of all, because you need to react to the interaction
		// 	// by sending the response in 3 seconds after receiving, otherwise
		// 	// interaction will be considered invalid and you can no longer
		// 	// use the interaction token and ID for responding to the user's request

		// 	content := ""
		// 	// As you can see, the response type names used here are pretty self-explanatory,
		// 	// but for those who want more information see the official documentation
		// 	switch i.Data.Options[0].IntValue() {
		// 	case int64(discordgo.InteractionResponseChannelMessage):
		// 		content =
		// 			"Well, you just responded to an interaction, and sent a message.\n" +
		// 				"That's all what I wanted to say, yeah."
		// 		content +=
		// 			"\nAlso... you can edit your response, wait 5 seconds and this message will be changed"
		// 	case int64(discordgo.InteractionResponseChannelMessageWithSource):
		// 		content =
		// 			"You just responded to an interaction, sent a message and showed the original one. " +
		// 				"Congratulations!"
		// 		content +=
		// 			"\nAlso... you can edit your response, wait 5 seconds and this message will be changed"
		// 	default:
		// 		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		// 			Type: discordgo.InteractionResponseType(i.Data.Options[0].IntValue()),
		// 		})
		// 		if err != nil {
		// 			s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
		// 				Content: "Something went wrong",
		// 			})
		// 		}
		// 		return
		// 	}

		// 	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		// 		Type: discordgo.InteractionResponseType(i.Data.Options[0].IntValue()),
		// 		Data: &discordgo.InteractionApplicationCommandResponseData{
		// 			Content: content,
		// 		},
		// 	})
		// 	if err != nil {
		// 		s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
		// 			Content: "Something went wrong",
		// 		})
		// 		return
		// 	}
		// 	time.AfterFunc(time.Second*5, func() {
		// 		err = s.InteractionResponseEdit(s.State.User.ID, i.Interaction, &discordgo.WebhookEdit{
		// 			Content: content + "\n\nWell, now you know how to create and edit responses. " +
		// 				"But you still don't know how to delete them... so... wait 10 seconds and this " +
		// 				"message will be deleted.",
		// 		})
		// 		if err != nil {
		// 			s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
		// 				Content: "Something went wrong",
		// 			})
		// 			return
		// 		}
		// 		time.Sleep(time.Second * 10)
		// 		s.InteractionResponseDelete(s.State.User.ID, i.Interaction)
		// 	})
		// },
		// "followups": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// 	// Followup messages are basically regular messages (you can create as many of them as you wish)
		// 	// but work as they are created by webhooks and their functionality
		// 	// is for handling additional messages after sending a response.

		// 	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		// 		Type: discordgo.InteractionResponseChannelMessageWithSource,
		// 		Data: &discordgo.InteractionApplicationCommandResponseData{
		// 			// Note: this isn't documented, but you can use that if you want to.
		// 			// This flag just allows you to create messages visible only for the caller of the command
		// 			// (user who triggered the command)
		// 			Flags:   1 << 6,
		// 			Content: "Surprise!",
		// 		},
		// 	})
		// 	msg, err := s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
		// 		Content: "Followup message has been created, after 5 seconds it will be edited",
		// 	})
		// 	if err != nil {
		// 		s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
		// 			Content: "Something went wrong",
		// 		})
		// 		return
		// 	}
		// 	time.Sleep(time.Second * 5)

		// 	s.FollowupMessageEdit(s.State.User.ID, i.Interaction, msg.ID, &discordgo.WebhookEdit{
		// 		Content: "Now the original message is gone and after 10 seconds this message will ~~self-destruct~~ be deleted.",
		// 	})

		// 	time.Sleep(time.Second * 10)

		// 	s.FollowupMessageDelete(s.State.User.ID, i.Interaction, msg.ID)

		// 	s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
		// 		Content: "For those, who didn't skip anything and followed tutorial along fairly, " +
		// 			"take a unicorn :unicorn: as reward!\n" +
		// 			"Also, as bonus... look at the original interaction response :D",
		// 	})
		// },
	}
)
