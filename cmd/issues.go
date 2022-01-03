package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/titolins/jira_scrapper/config"
	"github.com/titolins/jira_scrapper/internal/cache"
	"github.com/titolins/jira_scrapper/internal/scrapper"
	"github.com/titolins/jira_scrapper/internal/service"
)

var issuesCmd = &cobra.Command{
	Use:   "issues",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		s := buildScrapper()
		for _, p := range getProjects() {
			for _, sp := range loadCachedSprints(p).Sprints {
				if err := service.NewIssues(s, sp, p.Key).Call(); err != nil {
					log.Printf("failed to fetch issues for sprint %d: %v", sp.ID, err)
				}
			}
		}
	},
}

func loadCachedSprints(p config.Project) (d scrapper.GetSprintsResult) {
	f := fmt.Sprintf("%s_%d", p.Key, p.BoardID)
	c := cache.New("sprints", f)
	if _, exists := c.Exists(); !exists {
		log.Printf("cache for project %q board %d doesn't exist\n", p.Key, p.BoardID)
	}
	if err := c.Load(&d); err != nil {
		log.Printf("failed to load cache for project %q board %d: %v", p.Key, p.BoardID, err)
	}

	return
}