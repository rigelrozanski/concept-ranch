package quac

import (
	"errors"
	"fmt"
	"path"

	"github.com/rigelrozanski/common"
	"github.com/rigelrozanski/thranch/quac/idea"
)

func appTags(application string) []idea.Tag {
	return []idea.Tag{
		idea.MustNewTagReg("external-use"),
		idea.MustNewTagRegWithValue("app", application),
	}
}

// for applications to receive content
func GetForApp(application string) string {
	content, found := ConcatAllContentFromTags(appTags(application))
	if !found {
		fmt.Println("nothing found with those tags")
	}
	return string(content)
}

func AppendLineForApp(application, appendLine string) error {
	ideas := GetAllIdeasNonConsuming()
	subset := ideas.WithTags(appTags(application)).WithText()
	if len(subset) != 1 {
		return errors.New("nothing found with those tags")
	}
	idea := subset[0]
	path := path.Join(IdeasDir, idea.Filename)
	content, err := common.ReadLines(path)
	if err != nil {
		return err
	}
	content = append(content, appendLine)
	_ = common.WriteLines(content, path)
	return nil
}
