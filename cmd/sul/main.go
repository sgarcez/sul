package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync"

	"github.com/sgarcez/sul/pkg/auth"
	"github.com/sgarcez/sul/pkg/uploader"
	"github.com/spf13/cobra"
	strava "github.com/strava/go.strava"
)

// new creates the application command tree
func new() *cobra.Command {

	rootCmd := &cobra.Command{
		Use:   "sul",
		Short: "A Strava activity file uploader",
	}

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Sul",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Sul - Strava Uploader v0.0.1")
		},
	}
	rootCmd.AddCommand(versionCmd)

	uploadCmd := &cobra.Command{
		Use:   "upload",
		Short: "Uploads activity files from directory",
		Run: func(cmd *cobra.Command, args []string) {
			accessToken := cmd.Flag("token").Value.String()
			u := uploader.NewUploader(accessToken)

			inputDir := cmd.Flag("dir").Value.String()
			files, err := ioutil.ReadDir(inputDir)
			if err != nil {
				log.Fatal(err)
			}

			var wg sync.WaitGroup
			wg.Add(len(files))
			log.Printf("Processing %d files\n", len(files))
			for _, f := range files {
				go func(fname string) {
					defer wg.Done()
					f, err := os.Open(path.Join(inputDir, fname))
					if err != nil {
						log.Print(err)
						return
					}
					u.Upload(fname, f)
				}(f.Name())
			}
			wg.Wait()
			log.Print("Done")
		},
	}
	uploadCmd.Flags().StringP("token", "t", "", "Access token")
	uploadCmd.MarkFlagRequired("token")
	uploadCmd.Flags().StringP("dir", "d", "", "Input directory")
	uploadCmd.MarkFlagRequired("dir")
	rootCmd.AddCommand(uploadCmd)

	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Retrieves access token with write permissions via OAuth flow",
		Run: func(cmd *cobra.Command, args []string) {
			port := cmd.Flag("port").Value.String()
			auth.Start(port)
		},
	}
	authCmd.Flags().IntVarP(&strava.ClientId, "id", "i", 0, "Strava Client ID")
	authCmd.MarkFlagRequired("id")
	authCmd.Flags().StringVarP(&strava.ClientSecret, "secret", "s", "", "Strava Client Secret")
	authCmd.MarkFlagRequired("secret")
	authCmd.Flags().IntP("port", "p", 8080, "Port for temp server")
	rootCmd.AddCommand(authCmd)

	return rootCmd
}

func main() {
	if err := new().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
