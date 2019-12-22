package plugins

import (
	. "Vtb_Record/src/utils"
	"github.com/bitly/go-simplejson"
	"regexp"
)

type yfConfig struct {
	IsLive bool
	Title  string
	Target string
}
type Youtube struct {
	yfConfig
	Url  string
	Name string
}

func (y Youtube) getVideoInfo() yfConfig {
	htmlBody := HttpGet(y.Url)
	re, err := regexp.Compile(`ytplayer.config\s*=\s*([^\n]+?});`)
	CheckError(err, "cannot find yfconfig")
	jsonYtConfig := re.FindSubmatch(htmlBody)[1]
	ytConfigJson, _ := simplejson.NewJson(jsonYtConfig)
	playerResponse, _ := simplejson.NewJson([]byte(ytConfigJson.Get("args").Get("player_response").MustString()))
	videoDetails := playerResponse.Get("videoDetails")
	IsLive, err := videoDetails.Get("isLive").Bool()
	if err == nil {
		IsLive = true
	}
	y.Title = videoDetails.Get("title").MustString()
	y.Target = "https://www.youtube.com/watch?v=" + videoDetails.Get("videoId").MustString()
	y.IsLive = IsLive
	return y.yfConfig
}
func (y Youtube) createVideo() VideoInfo {
	return VideoInfo{
		Title:         y.Title,
		Date:          GetTimeNow(),
		Target:        y.Target,
		Provider:      "Youtube",
		Filename:      GenerateFilepath(y.Name, y.Title),
		StreamingLink: "",
	}
}
func YoutubeCheckLive(usersConfig UsersConfig) {
	y := new(Youtube)
	y.Url = "https://www.youtube.com/channel/" + usersConfig.TargetId + "/live"
	y.Name = usersConfig.Name
	yfConfig := y.getVideoInfo()
	if yfConfig.IsLive == true {
		ProcessVideo(y.createVideo())
	}
	NoLiving("Youtube", usersConfig.Name)
}
