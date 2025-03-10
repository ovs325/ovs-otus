package white

import (
	"fmt"
	"net"
	"net/http"
	"net/url"

	cm "main/internal/common"
	pr "main/internal/params"

	cb "github.com/spf13/cobra"
)

var netw string

var WhiteCmd = &cb.Command{
	Use:   "white",
	Short: "Удаление Сети/Подсети из белого листа",
	Long:  "Команда удаления Сети/Подсети из белого списка",

	RunE: func(cmd *cb.Command, args []string) error {
		if netw == "" {
			return fmt.Errorf("необходим параметр 'network'")
		}
		if _, _, err := net.ParseCIDR(netw); err != nil {
			return fmt.Errorf("параметр 'network' %s не может быть распознан (пример: 137.23.45.167/25)", netw)
		}

		url, err := url.Parse(fmt.Sprintf("%s://%s:%s/del/white", pr.ConfigHttp.Scheme, pr.ConfigHttp.Host, pr.ConfigHttp.Port))
		if err != nil {
			return fmt.Errorf("ошибка при парсинге URL: %v", err)
		}
		query := url.Query()
		query.Set("network", netw)
		url.RawQuery = query.Encode()
		if err := cm.NewRequest(http.MethodDelete, url.String(), nil); err != nil {
			return fmt.Errorf("ошибка при отправке запроса %s: %v", url.String(), err)
		}
		return nil
	},
}

func init() {
	WhiteCmd.Flags().StringVarP(&netw, "network", "n", "", "адрес с подсетью, например: 137.23.45.167/25")
	if err := WhiteCmd.MarkFlagRequired("network"); err != nil {
		panic(fmt.Errorf("необходим параметр 'network': %v", err))
	}
}
