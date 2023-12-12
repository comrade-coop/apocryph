// SPDX-License-Identifier: GPL-3.0

// Package podmanagement provides functionality for managing pods (create pods, upload package to ipfs).
package podmanagement

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/spf13/viper"
)

func extractImageNames(manifest string) []string {
	vp := viper.New()
	vp.SetConfigFile(manifest)
	err := vp.ReadInConfig()
	if err != nil {
		fmt.Printf("Could not read manifest file %v\n", err)
		return nil
	}
	containers := vp.Get("containers")
	yaml := fmt.Sprintf("%v", containers)

	pattern := `image:([^ \] ]+)`

	regex := regexp.MustCompile(pattern)

	matches := regex.FindAllStringSubmatch(yaml, -1)
	var images []string
	if len(matches) > 0 {
		for _, match := range matches {
			imageName := match[1]
			images = append(images, imageName)
		}
	} else {
		log.Println("No image names found in the manifest file.")
	}
	return images
}

func saveImagesLocally(manifest string, savePath string) error {
	images := extractImageNames(manifest)
	os.MkdirAll(savePath, os.ModePerm)
	for _, v := range images {
		name := fmt.Sprintf("%v.tar", v)
		fileName := strings.Replace(name, "/", "_", -1)
		tarPath := fmt.Sprintf("%v/%v", savePath, fileName)
		cmd := exec.Command("docker", "save", "-o", tarPath, v)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("could not save image:", v, err)
			fmt.Println("Command output:", string(output))
			return err
		}
	}
	return nil
}
