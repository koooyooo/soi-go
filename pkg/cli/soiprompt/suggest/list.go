package suggest

import (
	"log"
	"strings"

	"github.com/koooyooo/soi-go/pkg/cli/soiprompt/utils"

	"github.com/koooyooo/soi-go/pkg/cli/soiprompt/suggest/common"

	"github.com/c-bata/go-prompt"
	"github.com/koooyooo/soi-go/pkg/cli/constant"
)

// ListCmd はlistコマンド系のSuggestを提示します
func ListCmd(d prompt.Document) []prompt.Suggest {
	// option探索
	if strings.HasPrefix(d.GetWordBeforeCursor(), "-") {
		return browserOptSuggests
	}
	soisDir, err := constant.LocalBucket.Path()
	if err != nil {
		log.Fatal(err)
	}
	files, err := utils.ListFilesRecursively(soisDir)
	if err != nil {
		log.Fatal(err)
	}
	swp, err := common.LoadSoiDataArray(files)
	if err != nil {
		panic(err)
	}
	var sgs []prompt.Suggest
	for _, s := range swp {
		sgs = append(sgs, prompt.Suggest{
			Text:        common.CreateHeader(s.SoiData) + " " + strings.TrimPrefix(s.Path, soisDir+"/"),
			Description: "",
		})
	}
	return prompt.FilterContains(sgs, d.GetWordBeforeCursor(), true)
}

var browserOptSuggests = []prompt.Suggest{
	{Text: "-f", Description: "open w/ firefox"},
	{Text: "-s", Description: "open w/ safari"},
}
