package main

import (
	"context"
	"strings"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"sync"

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
			fmt.Println("Sul - Strava Uploader v0.0.2")
		},
	}
	rootCmd.AddCommand(versionCmd)

	uploadCmd := &cobra.Command{
		Use:   "upload",
		Short: "Uploads activity files from directory",
		Run: func(cmd *cobra.Command, args []string) {
			accessToken := cmd.Flag("token").Value.String()
			u := newUploader(accessToken)

			inputDir := cmd.Flag("dir").Value.String()
			files, err := ioutil.ReadDir(inputDir)
			if err != nil {
				log.Fatal(err)
			}

			var wg sync.WaitGroup
			log.Printf("Processing %d files\n", len(files))

			for _, f := range files {
				if f.IsDir() || strings.HasPrefix(f.Name(), ".") || !strings.HasSuffix(strings.ToLower(f.Name()), ".fit") {
					log.Printf("Ignoring %s\n", f.Name())
					continue
				}

				wg.Add(1)
				go func(fname string) {
					defer wg.Done()
					f, err := os.Open(path.Join(inputDir, fname))
					if err != nil {
						log.Printf("%s - open: %s", fname, err)
						return
					}
					aid, err := u.Upload(fname, f)
					if err != nil {
						log.Printf("%s - upload: %s", fname, err)
						return
					}
					log.Printf(
						"%s - Activity created, you can view it at http://www.strava.com/activities/%d",
						fname, aid)
				}(f.Name())
			}
			wg.Wait()
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

			m := http.NewServeMux()
			s := &http.Server{Addr: fmt.Sprintf(":%s", port), Handler: m}

			authURL, path, handler := AuthHandler(port)

			m.HandleFunc(path, handleAndKill(s, handler))

			fmt.Printf("-------------------------------\n")
			fmt.Printf("Use this URL to authorise your application:\n\n%s\n", authURL)
			fmt.Printf("-------------------------------\n")

			if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatal(err)
			}
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

func handleAndKill(s *http.Server, in http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in(w, r)
		fmt.Fprintf(w, "\nYou can close this window")
		go func() {
			if err := s.Shutdown(context.Background()); err != nil {
				log.Fatal(err)
			}
		}()
	}
}