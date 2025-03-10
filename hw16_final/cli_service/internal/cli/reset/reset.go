package reset

import (
	"fmt"
	"net/http"
	"net/url"

	cm "main/internal/common"
	pr "main/internal/params"

	cb "github.com/spf13/cobra"
)

var group, backet string

var ResetCmd = &cb.Command{
	Use:   "reset",
	Short: "reset the backet",
	Long:  "Команда сброса данных корзины (удаления корзины)",

	RunE: func(cmd *cb.Command, args []string) error {
		if group == "" {
			return fmt.Errorf("необходим параметр 'group'")
		}
		if backet == "" {
			return fmt.Errorf("необходим параметр 'backet'")
		}

		url, err := url.Parse(fmt.Sprintf("%s://%s:%s/reset", pr.ConfigHttp.Scheme, pr.ConfigHttp.Host, pr.ConfigHttp.Port))
		if err != nil {
			return fmt.Errorf("ошибка при парсинге URL: %v", err)
		}
		query := url.Query()
		query.Set("params", fmt.Sprintf("%s^%s", group, backet))
		url.RawQuery = query.Encode()
		if err := cm.NewRequest(http.MethodPatch, url.String(), nil); err != nil {
			return fmt.Errorf("ошибка при отправке запроса %s: %v", url.String(), err)
		}
		return nil
	},
}

func init() {
	ResetCmd.Flags().StringVarP(&group, "group", "g", "", "группы корзин, например: ip/login/pass и пр.")
	if err := ResetCmd.MarkFlagRequired("group"); err != nil {
		panic(fmt.Errorf("необходим параметр 'group': %v", err))
	}
	ResetCmd.Flags().StringVarP(&backet, "backet", "b", "", "наименование корзины например: 137.23.45.167/Volga/foo и пр.")
	if err := ResetCmd.MarkFlagRequired("backet"); err != nil {
		panic(fmt.Errorf("необходим параметр 'backet': %v", err))
	}
}
