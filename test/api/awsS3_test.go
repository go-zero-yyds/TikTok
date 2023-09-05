package test

import (
	"TikTok/pkg/filesystem"
	"bytes"
	"io"
	"net/http"
	"testing"
)

func TestS3Utility(t *testing.T) {
	// Replace the following values with your S3 server configuration for testing
	URL := "http://127.0.0.1:9000"
	Bucket := "test"
	AwsAccessKeyId := "7KDhgxEm3G5sQLj4ePHv"
	AwsSecretAccessKey := "zaB97osHgENMqdc970FfOL9PGiA1UTiosa7CBLQ8"

	// Create a new S3 client for testing
	var s3Client filesystem.FileSystem
	s3Client = filesystem.NewS3(URL, Bucket, "", AwsAccessKeyId, AwsSecretAccessKey)

	// Test Upload and GetDownloadLink
	data := []byte("This is a test data.")
	key := "testfile.txt"

	// Upload the data to S3
	err := s3Client.Upload(bytes.NewReader(data), key)
	if err != nil {
		t.Fatalf("Failed to upload file to S3: %v", err)
	}

	// Get the download link for the uploaded file
	downloadLink, err := s3Client.GetDownloadLink(key)
	if err != nil {
		t.Fatalf("Failed to get download link: %v", err)
	}
	t.Logf("Download Link: %s", downloadLink)

	// Test FileExists
	exists, err := s3Client.FileExists(key)
	if err != nil {
		t.Fatalf("Failed to check file existence: %v", err)
	}
	if !exists {
		t.Errorf("Expected file %s to exist, but it doesn't.", key)
	}

	// Make an HTTP request to the download link
	resp, err := http.Get(downloadLink)
	if err != nil {
		t.Fatalf("Failed to make HTTP request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}
	// Create a buffer to store the downloaded data
	var downloadedData bytes.Buffer

	// Copy the response body to the buffer
	_, err = io.Copy(&downloadedData, resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Check the downloaded data with the original data
	if !bytes.Equal(data, downloadedData.Bytes()) {
		t.Error("Response data does not match expected data")
	}

	// Test Delete
	err = s3Client.Delete(key)
	if err != nil {
		t.Fatalf("Failed to delete file: %v", err)
	}

	// Check if the file still exists after deletion
	exists, err = s3Client.FileExists(key)
	if err != nil {
		t.Fatalf("Failed to check file existence after deletion: %v", err)
	}
	if exists {
		t.Errorf("Expected file %s to be deleted, but it still exists.", key)
	}
}
