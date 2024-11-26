package main

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func main() {
	var concurrency int

	var downloadCmd = &cobra.Command{
		Use:   "download URI [target]",
		Short: "Download from blob store",
		Long:  `download a blob form the blog store with a concurrency`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			target := filepath.Base(args[0])
			if len(args) == 2 {
				target = args[1]
			}
			startDownload(args[0], target, concurrency)
		},
	}
	downloadCmd.PersistentFlags().IntVarP(&concurrency, "concurrency", "c", runtime.NumCPU(), "the number of concurrent downloads")
	var app = &cobra.Command{Use: "app"}
	app.AddCommand(downloadCmd)
	app.Execute()
}

func startDownload(url string, target string, concurrency int) {
	if concurrency < 1 {
		log.Fatal().Msg("Concurrency must be greater than 0")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatal().Err(err).Msg("could not auth")
	}

	token, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{
		Scopes: []string{"https://storage.azure.com/.default"},
	})
	if err != nil {
		log.Fatal().Err(err).Msg("could not auth")
	}

	resp, err := request(token.Token, "HEAD", url)
	if err != nil {
		log.Fatal().Err(err).Msg("could not make request")
	}

	if err != nil {
		log.Fatal().Err(err).Msg("could not make request")
	}
	if resp.StatusCode > 299 {
		log.Fatal().
			Any("statusCode", resp.StatusCode).
			Any("status", resp.Status).
			Any("responseHeader", resp.Header).
			// Any("requestHeader", request.Header).
			Msg("bad status code")
	}
	contentLength := resp.Header.Get("Content-Length")

	fileSize, err := strconv.ParseInt(contentLength, 10, 64)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("invalid content length")
	}

	println("file size", fileSize, "thread count", concurrency)

	file, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("target", target).
			Msg("Cannot open file")
	}
	defer file.Close()

	var wg sync.WaitGroup
	pieceSize := fileSize / int64(concurrency)
	remainder := fileSize - pieceSize*int64(concurrency)
	for i := 0; i < concurrency; i++ {
		offset := pieceSize * int64(i)
		length := pieceSize
		if i == concurrency-1 {
			length = pieceSize + remainder
		}
		go download(&wg, file, token.Token, url, offset, length)
	}
	wg.Wait()
}

func download(wg *sync.WaitGroup, file *os.File, token string, url string, offset int64, length int64) {
	wg.Add(1)
	defer wg.Done()
	println("offset", offset, "length", length)
}

func seekWriter(file *os.File, input chan(string), wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	select input {

	}
}
