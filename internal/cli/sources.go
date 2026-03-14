package cli

import (
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

var sourcesCmd = &cobra.Command{
	Use:   "sources",
	Short: "Check accessibility of all configured sources",
	RunE: func(cmd *cobra.Command, args []string) error {
		httpClient := &http.Client{Timeout: 10 * time.Second}

		fmt.Printf("%-30s %-10s %-50s %s\n", "NAME", "TYPE", "URL", "STATUS")
		fmt.Println("----------------------------  ---------- -------------------------------------------------- --------")

		for _, s := range cfg.Sources.RSS {
			status := checkURL(httpClient, s.URL)
			fmt.Printf("%-30s %-10s %-50s %s\n", s.Name, "rss", s.URL, status)
		}

		for _, s := range cfg.Sources.YouTube {
			url := fmt.Sprintf("https://www.youtube.com/channel/%s", s.ChannelID)
			status := checkURL(httpClient, url)
			fmt.Printf("%-30s %-10s %-50s %s\n", s.Name, "youtube", url, status)
		}

		for _, s := range cfg.Sources.Podcasts {
			status := checkURL(httpClient, s.URL)
			fmt.Printf("%-30s %-10s %-50s %s\n", s.Name, "podcast", s.URL, status)
		}

		return nil
	},
}

func checkURL(client *http.Client, url string) string {
	resp, err := client.Head(url)
	if err != nil {
		return "UNREACHABLE"
	}
	_ = resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return "OK"
	}
	return fmt.Sprintf("HTTP %d", resp.StatusCode)
}

func init() {
	rootCmd.AddCommand(sourcesCmd)
}
