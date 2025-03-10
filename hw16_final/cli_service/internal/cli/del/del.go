package del

import (
	"errors"
	"fmt"
	"main/internal/cli/del/black"
	"main/internal/cli/del/white"
	"slices"
	"strings"

	cb "github.com/spf13/cobra"
)

var DelCmd = &cb.Command{
	Use:   "del",
	Short: "Удаление Сети/Подсети из белого или черного листа",
	Long:  "Команда удаления Сети/Подсети  из белого или черного списка",

	RunE: func(cmd *cb.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("отсутствует подкоманды для команды 'del'")
		}
		sub := []string{"white", "black"}
		if !slices.Contains(sub, args[0]) {
			return fmt.Errorf("подкомандами могут быть только: %s", strings.Join(sub, ", "))
		}
		return nil
	},
}

func init() {
	DelCmd.AddCommand(black.BlackCmd)
	DelCmd.AddCommand(white.WhiteCmd)
}
