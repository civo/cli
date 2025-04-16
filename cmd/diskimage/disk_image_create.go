package diskimage

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var (
	createDiskImageName         string
	createDiskImageDistribution string
	createDiskImageVersion      string
	createDiskImagePath         string
	createOS                    string
	createLogoPath              string
)

var diskImageCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"upload", "new"},
	Short:   "Create a new disk image",
	Example: "civo diskimage create --name ubuntu-22.04 --distribution ubuntu --version 22.04 --path ./image.qcow2",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate required flags
		if createDiskImageName == "" || createDiskImageDistribution == "" || createDiskImageVersion == "" || createDiskImagePath == "" {
			return errors.New("missing required flags: --name, --distribution, --version, --path")
		}

		// Validate disk image file
		ext := strings.ToLower(filepath.Ext(createDiskImagePath))
		if ext != ".raw" && ext != ".qcow2" {
			return errors.New("invalid disk image format, must be .raw or .qcow2")
		}

		// Calculate checksums
		file, err := os.Open(createDiskImagePath)
		if err != nil {
			return fmt.Errorf("failed to open disk image: %s", err)
		}
		defer file.Close()

		hashSHA256 := sha256.New()
		hashMD5 := md5.New()

		if _, err := io.Copy(io.MultiWriter(hashSHA256, hashMD5), file); err != nil {
			return fmt.Errorf("failed to calculate checksums: %s", err)
		}

		sha256Sum := hex.EncodeToString(hashSHA256.Sum(nil))
		md5Sum := hex.EncodeToString(hashMD5.Sum(nil))

		// Get file size
		fileInfo, err := file.Stat()
		if err != nil {
			return fmt.Errorf("failed to get file info: %s", err)
		}
		size := fileInfo.Size()

		// Handle logo
		var logoBase64 string
		if createLogoPath != "" {
			if strings.ToLower(filepath.Ext(createLogoPath)) != ".svg" {
				return errors.New("logo must be an SVG file")
			}
			logoBytes, err := os.ReadFile(createLogoPath)
			if err != nil {
				return fmt.Errorf("failed to read logo file: %s", err)
			}
			logoBase64 = base64.StdEncoding.EncodeToString(logoBytes)
		}

		// Create API request
		client, err := config.CivoAPIClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %s", err)
		}

		params := &civogo.CreateDiskImageParams{
			Name:         createDiskImageName,
			Distribution: createDiskImageDistribution,
			Version:      createDiskImageVersion,
			OS:           createOS,
			ImageSHA256:  sha256Sum,
			ImageMD5:     md5Sum,
			ImageSize:    size,
			LogoBase64:   logoBase64,
		}

		resp, err := client.CreateDiskImage(params)
		if err != nil {
			return fmt.Errorf("failed to create disk image: %s", err)
		}

		fmt.Println("DISK URL: ", resp.DiskImageURL)

		// Upload disk image
		err = uploadUsingPresigned(resp.DiskImageURL, createDiskImagePath)
		if err != nil {
			return fmt.Errorf("failed to upload disk image: %s", err)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{
			"id":     resp.ID,
			"name":   resp.Name,
			"status": resp.Status,
		})

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			ow.WriteTable()
		}

		return nil
	},
}

func uploadUsingPresigned(url string, filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get file size for progress bar and content-length header
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	// Create progress bar
	bar := progressbar.NewOptions64(
		fileInfo.Size(),
		progressbar.OptionSetDescription("Uploading disk image..."),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(40),
		progressbar.OptionThrottle(100*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
	)

	// Create a reader that updates the progress bar
	reader := io.TeeReader(file, bar)

	req, err := http.NewRequest("PUT", url, reader)
	if err != nil {
		return err
	}

	// Set required headers
	req.ContentLength = fileInfo.Size()
	req.Header.Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upload failed with status: %s", resp.Status)
	}

	return nil
}
