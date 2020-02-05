package quac

import (
	"fmt"
)

// for applications to receive content
func GetForApp(application string) string {
	tags := []string{"external-use", "app:" + application}
	content, found := ConcatAllContentFromTags(tags)
	if !found {
		fmt.Println("nothing found with those tags")
	}
	return string(content)
}

// get a group of images by tag
func GetImagesByTag(tags []string) (ideas Ideas) {
	return GetAllIdeas().WithImage().WithTags(tags)
}
