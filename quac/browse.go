package quac

import (
	"fmt"
	"log"
	"os"
	"sort"

	tui "github.com/marcusolsson/tui-go"
	"github.com/rigelrozanski/thranch/quac/idea"
)

type bList struct {
	*tui.List
	items      []string
	blankItems int
	maxWidth   int
}

func (b *bList) PrependBlanks(noBlanks int) {
	b.blankItems += noBlanks

	newList := tui.NewList()
	for i := 0; i < b.blankItems; i++ {
		newList.AddItems("")
	}
	newList.AddItems(b.items...)
	b.List = newList
}

func MaxWidth(strs []string) int {
	mw := 0
	for _, str := range strs {
		if len(str) > mw {
			mw = len(str)
		}
	}
	return mw
}

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (p PairList) Top(max int) []string {
	topItems := []string{}
	for i, item := range p {
		topItems = append(topItems, item.Key)
		if i > max {
			break
		}
	}
	return topItems
}

func rankByWordCount(wordFrequencies map[string]int) PairList {
	pl := make(PairList, len(wordFrequencies))
	i := 0
	for k, v := range wordFrequencies {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

func GetAssociations(idears idea.Ideas, tags []string) (associatedTags PairList) {
	subset := idears
	if len(tags) > 0 {
		subset = idears.WithTags(tags)
	}

	at := make(map[string]int)
	for _, idea := range subset {
		for _, tag := range idea.Tags {

			// don't count inputs or highlights
			inInputs := false
			for _, inputTags := range tags {
				if tag == inputTags {
					inInputs = true
					break
				}
			}
			for _, inputTags := range highlighted {
				if tag == inputTags || inInputs {
					inInputs = true
					break
				}
			}

			if !inInputs {
				at[tag] += 1
			}
		}
	}
	return rankByWordCount(at)
}

var highlighted []string

func AddHighlighted(tag string) {
	highlighted = append(highlighted, tag)
}

func RemoveHighlighted(tag string) {
	var newHighlighted []string
	for _, item := range highlighted {
		if item != tag {
			newHighlighted = append(newHighlighted, item)
		}
	}
	highlighted = newHighlighted
}

func Ls(clumpedTags string) {

	inputTags := idea.ParseClumpedTags(clumpedTags)
	idears := idea.GetAllIdeasNonConsuming()
	highlighted = inputTags

	listi := 0
	lists := []*bList{}
	l1Items := []string{}
	l1 := tui.NewList()

	initDrill := false
	if len(inputTags) == 0 {
		tagCounts := GetAssociations(idears, inputTags)
		if len(tagCounts) == 0 {
			fmt.Println("no associations found")
			os.Exit(1)
		}
		l1Items = tagCounts.Top(100)
	} else {
		l1.SetStyle("highlightedAllList")
		l1Items = inputTags
		initDrill = true
	}

	l1.SetFocused(true)
	l1.AddItems(l1Items...)
	l1.SetSelected(0)
	lists = append(lists, &bList{l1, l1Items, 0, MaxWidth(l1Items)})

	t := tui.NewTheme()
	t.SetStyle("list.item", tui.Style{Bg: tui.ColorDefault, Fg: tui.ColorWhite})
	t.SetStyle("list.item.selected", tui.Style{Bg: tui.ColorYellow, Fg: tui.ColorBlack})
	t.SetStyle("highlightedList", tui.Style{Bg: tui.ColorDefault, Fg: tui.ColorWhite})
	t.SetStyle("highlightedList.selected", tui.Style{Bg: tui.ColorRed, Fg: tui.ColorBlack})
	t.SetStyle("highlightedAllList", tui.Style{Bg: tui.ColorRed, Fg: tui.ColorBlack})
	t.SetStyle("highlightedAllList.selected", tui.Style{Bg: tui.ColorRed, Fg: tui.ColorBlack})

	hlists := tui.NewHBox(l1)
	s := tui.NewScrollArea(hlists)
	s.Scroll(-100, -10)
	ui, err := tui.New(s)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetTheme(t)

	setListFn := func() {
		l2 := tui.NewList()
		dudItems := lists[listi].Selected()
		for i := 0; i < dudItems; i++ {
			l2.AddItems("")
		}

		selectedTag := lists[listi].SelectedItem()
		newItems := GetAssociations(idears, append(highlighted, selectedTag)).Top(100)
		mw := MaxWidth(newItems)
		s.Scroll(mw, 0)
		l2.AddItems(newItems...)

		l2.SetSelected(dudItems)
		hlists.Append(l2)
		lists = append(lists, &bList{l2, newItems, dudItems, mw})
		listi++
	}

	if initDrill {
		setListFn()
	}

	ui.SetKeybinding("q", func() {
		ui.Quit()
	})

	ui.SetKeybinding("Enter", func() {
		if len(highlighted) > 0 {
			ui.Quit()
			fmt.Printf("debug highlighted: %v\n", highlighted)
			MultiOpenByTags(highlighted, false)
		}
	})

	ui.SetKeybinding("k", func() {
		if lists[listi].Selected()-1 >= lists[listi].blankItems {
			lists[listi].SetSelected(lists[listi].Selected() - 1)
			s.Scroll(0, -1)
		}
	})

	ui.SetKeybinding("j", func() {
		if lists[listi].Selected()+1 < lists[listi].Length() {
			lists[listi].SetSelected(lists[listi].Selected() + 1)
			s.Scroll(0, 1)
		}
	})

	ui.SetKeybinding("Ctrl+l", func() {
		lists[listi].SetStyle("highlightedList")
		AddHighlighted(lists[listi].SelectedItem())
		setListFn()
	})

	ui.SetKeybinding("l", func() {
		lists[listi].SetStyle("list.item")
		setListFn()
	})

	ui.SetKeybinding("h", func() {
		if len(lists) == 1 {
			return
		}
		heightScroll := lists[listi-1].Selected() - lists[listi].Selected()
		s.Scroll(-lists[listi].maxWidth, heightScroll)
		hlists.Remove(hlists.Length() - 1)
		lists = lists[:len(lists)-1]
		listi--
		lists[listi].SetStyle("list.item")
		if listi != 0 {
			RemoveHighlighted(lists[listi].SelectedItem())
		} else {
			for _, item := range lists[listi].items {
				RemoveHighlighted(item)
			}
		}
	})

	//ui.SetKeybinding("k", func() { s.Scroll(0, -1) })
	//ui.SetKeybinding("j", func() { s.Scroll(0, 1) })
	//ui.SetKeybinding("h", func() { s.Scroll(-1, 0) })
	//ui.SetKeybinding("l", func() { s.Scroll(1, 0) })

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
