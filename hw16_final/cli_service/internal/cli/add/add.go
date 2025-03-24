package add

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	cb "github.com/spf13/cobra"
	"main/internal/cli/add/black"
	"main/internal/cli/add/white"
)

var AddCmd = &cb.Command{
	Use:   "add",
	Short: "Добавление Сети/Подсети в белый или черный лист",
	Long:  "Команда добавление Сети/Подсети в белый или черный список",

	RunE: func(_ *cb.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("отсутствует подкоманды для команды 'add'")
		}
		sub := []string{"white", "black"}
		if !slices.Contains(sub, args[0]) {
			return fmt.Errorf("подкомандами могут быть только: %s", strings.Join(sub, ", "))
		}
		return nil
	},
}

func init() {
	AddCmd.AddCommand(black.BlackCmd)
	AddCmd.AddCommand(white.WhiteCmd)
}
