package bot

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"k8s.io/klog"
)

type Server struct {
	*discordgo.Session
	Guild         string
	CleanUpOnExit bool
}

func (s Server) Run() error {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.Data.Name]; ok {
			h(s, i)
		}
	})

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		klog.Infof("bot is up!")
	})
	err := s.Session.Open()
	if err != nil {
		return fmt.Errorf("cannot open the session: %v", err)
	}

	defer s.Close()

	if err := s.registerCommands(); err != nil {
		return err
	}
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	klog.Info("gracefully shutting down...")
	if s.CleanUpOnExit {
		if err := s.removeAllCommands(); err != nil {
			return err
		}
	}
	return nil
}

func (s Server) registerCommands() error {
	klog.Info("registering commands")
	for _, v := range commands {
		klog.Infof("registering command: %s", v.Name)
		_, err := s.ApplicationCommandCreate(s.State.User.ID, s.Guild, v)
		if err != nil {
			return fmt.Errorf("cannot create '%s' command: %v", v.Name, err)
		}
	}
	return nil
}

func (s Server) removeAllCommands() error {
	klog.Info("deleting all commands")
	commands, err := s.Session.ApplicationCommands(s.State.User.ID, s.Guild)
	if err != nil {
		return err
	}

	for _, c := range commands {
		klog.Infof("deleting command %s", c.Name)
		err := s.Session.ApplicationCommandDelete(s.State.User.ID, s.Guild, c.ID)
		if err != nil {
			klog.Errorf("Cannot delete '%v' command: %v", c.Name, err)
		}
	}

	return nil
}
