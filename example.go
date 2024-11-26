package main

// import (
// 	"context"
// 	"fmt"
// 	"io"
// 	"log"
// 	"os"

// 	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
// )

// func Example_blob_Client_Download() {
// 	// From the Azure portal, get your Storage account blob service URL endpoint.
// 	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

// 	// Create a blobClient object to a blob in the container (we assume the container & blob already exist).
// 	blobURL := fmt.Sprintf("https://%s.blob.core.windows.net/mycontainer/BigBlob.bin", accountName)
// 	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
// 	handleError(err)
// 	blobClient, err := blob.NewClientWithSharedKeyCredential(blobURL, credential, nil)
// 	handleError(err)

// 	contentLength := int64(0) // Used for progress reporting to report the total number of bytes being downloaded.

// 	// Download returns an intelligent retryable stream around a blob; it returns an io.ReadCloser.
// 	dr, err := blobClient.DownloadStream(context.TODO(), nil)
// 	handleError(err)
// 	rs := dr.Body

// 	// NewResponseBodyProgress wraps the GetRetryStream with progress reporting; it returns an io.ReadCloser.
// 	stream := streaming.NewResponseProgress(
// 		rs,
// 		func(bytesTransferred int64) {
// 			fmt.Printf("Downloaded %d of %d bytes.\n", bytesTransferred, contentLength)
// 		},
// 	)
// 	defer func(stream io.ReadCloser) {
// 		err := stream.Close()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}(stream) // The client must close the response body when finished with it

// 	file, err := os.Create("BigFile.bin") // Create the file to hold the downloaded blob contents.
// 	handleError(err)
// 	defer func(file *os.File) {
// 		err := file.Close()
// 		if err != nil {

// 		}
// 	}(file)

// 	written, err := io.Copy(file, stream) // Write to the file by reading from the blob (with intelligent retries).
// 	handleError(err)
// 	fmt.Printf("Wrote %d bytes.\n", written)
// }
