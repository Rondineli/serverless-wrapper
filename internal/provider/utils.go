package provider

import (
	"os"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"strings"
    "archive/zip"
    "path/filepath"
)

// TODO: Managing s3 uploads zips and returning its location
//func s3Uploader(filePath string) (string, err) {
//	var output string
//	sess := session.Must(session.NewSession())
//
//	uploader := s3manager.NewUploader(sess)
//
//	f, err  := os.Open(filePath)
//	if err != nil {
//	    return output, fmt.Errorf("failed to open file %q, %v", filePath, err)
//	}
//
//	// Upload the file to S3.
//	result, err := uploader.Upload(&s3manager.UploadInput{
//	    Bucket: aws.String(myBucket),
//	    Key:    aws.String(myString),
//	    Body:   f,
//	})
//	if err != nil {
//	    return output, fmt.Errorf("failed to upload file, %v", err)
//	}
//	output = aws.StringValue(result.Location)
//	return output, nil
//}

func hash_file_sha1(filePath string) (string, error) {
	//Initialize variable returnMD5String now in case an error has to be returned
	var returnSHA1String string

	file, err := os.Open(filePath)
	if err != nil {
		return returnSHA1String, err
	}

	defer file.Close()

	hash := sha1.New()

	if _, err := io.Copy(hash, file); err != nil {
		return returnSHA1String, err
	}
	hashInBytes := hash.Sum(nil)[:20]
	returnSHA1String = hex.EncodeToString(hashInBytes)
	return returnSHA1String, nil
}

func zip_wrapper(inputDir string, outputFile string) (string, error) {
	// initiating empty string to return
	var output string

    destinationFile, err := os.Create(outputFile)
    if err != nil {
        return output, err
    }
    myZip := zip.NewWriter(destinationFile)
    err = filepath.Walk(inputDir, func(filePath string, info os.FileInfo, err error) error {
        if info.IsDir() {
            return nil
        }
        if err != nil {
            return err
        }
        relPath := strings.TrimPrefix(filePath, filepath.Dir(inputDir))
        zipFile, err := myZip.Create(relPath)
        if err != nil {
            return err
        }
        fsFile, err := os.Open(filePath)
        if err != nil {
            return err
        }
        _, err = io.Copy(zipFile, fsFile)
        if err != nil {
            return err
        }
        return nil
    })
    if err != nil {
        return output, err
    }
    err = myZip.Close()
    if err != nil {
        return output, err
    }
    return output, nil
}