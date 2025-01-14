package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var algorithm string

func main() {
	var rootCmd = &cobra.Command{
		Use:   "rate_limiter",
		Short: "Rate Limiter Server",
		Run: func(cmd *cobra.Command, args []string) {
			startServer()
		},
	}

	rootCmd.Flags().StringVarP(&algorithm, "algorithm", "a", "tokenbucket", "Rate limiting algorithm (tokenbucket or leakybucket)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func startServer() {
	if algorithm != "" {
		fmt.Printf("Server started with %s algorithm", algorithm)
		err := http.ListenAndServe(":8080", Route(algorithm))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("Invalid algorithm specified")
	}
}

func Unlimited(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Unlimited! Let's Go!")
}

func Limited(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Cosmos Allows you ! \n")
}
