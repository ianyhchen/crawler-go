package ptt

const (
	url_gossip          = "https://www.ptt.cc/bbs/Gossiping/index.html"
	url_lifeismoney     = "https://www.ptt.cc/bbs/Lifeismoney/index.html"
	DefaultGetDataLimit = 20
)

type Board struct {
	Name string
	URL  string
}

var BoardInfo = []Board{
	{Name: "gossip", URL: url_gossip},
	{Name: "lifeismoney", URL: url_lifeismoney},
}

func GetBoardURL(name string) (returnURL string) {
	for _, info := range BoardInfo {
		if name == info.Name {
			returnURL = info.URL
			break
		}
	}
	return returnURL
}

func ValidateBoardName(input string) bool {
	for _, info := range BoardInfo {
		if input == info.Name {
			return true
		}
	}
	return false
}
