package quac

import (
	"fmt"
	"path"

	"github.com/rigelrozanski/common"
)

// for applications to receive content
func GetForApp(application string) string {
	tags := []string{"external-use", "app=" + application}
	content, found := ConcatAllContentFromTags(tags)
	if !found {
		fmt.Println("nothing found with those tags")
	}
	return string(content)
}

func AppendLineForApp(application, appendLine string) {
	tags := []string{"external-use", "app=" + application}
	ideas := GetAllIdeasNonConsuming()
	subset := ideas.WithTags(tags).WithText()
	if len(subset) != 1 {
		return
	}
	idea := subset[0]
	path := path.Join(IdeasDir, idea.Filename)
	content, err := common.ReadLines(path)
	if err != nil {
		return
	}
	content = append(content, appendLine)
	_ = common.WriteLines(content, path)
}

// get a group of images by tag
func GetImagesByTag(tags []string) (ideas Ideas) {
	return GetAllIdeas().WithImage().WithTags(tags)
}
