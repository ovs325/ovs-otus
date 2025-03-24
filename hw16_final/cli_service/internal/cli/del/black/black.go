package black

import (
	"fmt"
	"net"
	"net/http"
	"net/url"

	cb "github.com/spf13/cobra"
	cm "main/internal/common"
	pr "main/internal/params"
)

var netw string

var BlackCmd = &cb.Command{
	Use:   "black",
	Short: "Удаление Сети/Подсети из черного листа",
	Long:  "Команда удаления Сети/Подсети из черного списка",

	RunE: func(_ *cb.Command, _ []string) error {
		if netw == "" {
			return fmt.Errorf("необходим параметр 'network'")
		}
		if _, _, err := net.ParseCIDR(netw); err != nil {
			return fmt.Errorf("параметр 'network' %s не может быть распознан (пример: 137.23.45.167/25)", netw)
		}

		url, err := url.Parse(
			fmt.Sprintf("%s://%s:%s/del/black", pr.ConfigHTTP.Scheme, pr.ConfigHTTP.Host, pr.ConfigHTTP.Port),
		)
		if err != nil {
			return fmt.Errorf("ошибка при парсинге URL: %w", err)
		}
		query := url.Query()
		query.Set("network", netw)
		url.RawQuery = query.Encode()
		if err := cm.NewRequest(http.MethodDelete, url.String(), nil); err != nil {
			return fmt.Errorf("ошибка при отправке запроса %s: %w", url.String(), err)
		}
		return nil
	},
}

func init() {
	BlackCmd.Flags().StringVarP(&netw, "network", "n", "", "адрес с подсетью, например: 137.23.45.167/25")
	if err := BlackCmd.MarkFlagRequired("network"); err != nil {
		panic(fmt.Errorf("необходим параметр 'network': %w", err))
	}
}
