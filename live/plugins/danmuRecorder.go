package plugins

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/fzxiao233/Vtb_Record/config"
	"github.com/fzxiao233/Vtb_Record/live/videoworker"
	"github.com/fzxiao233/Vtb_Record/utils"
	log "github.com/sirupsen/logrus"
)

func callJsAPI(roomID string, status string, filename string) error {
	var err error
	host := "http://127.0.0.1"
	port := "3000"
	if config.Config.DanmuHost != "" {
		host = config.Config.DanmuHost
	}
	if config.Config.DanmuPort != "" {
		port = config.Config.DanmuPort
	}
	_, err = utils.HttpGet(nil, host+":"+port+"/api/live?roomID="+roomID+"&status="+status+"&filename="+url.QueryEscape(filename), map[string]string{})
	if err != nil {
		err = fmt.Errorf("call danmaku error %v", err)
		log.Warn(err)
		return err
	}
	log.Debugf("[Danmaku]%s: %s", roomID, status)
	return nil
}

func getRoomId(targetId string) string {
	var resp []byte
	var err error = nil
	for {
		resp, err = utils.HttpGet(nil, "https://api.vtbs.moe/v1/detail/"+targetId, map[string]string{})
		if err != nil {
			log.Errorf("cannot get roomid %v", err)
			continue
		}
		respJson, err := simplejson.NewJson(resp)
		if err != nil {
			log.Errorf("%s parse json error", targetId)
		}
		if respJson != nil {
			roomId := strconv.Itoa(respJson.Get("roomid").MustInt())
			return roomId
		}
	}
}

type PluginDanmuRecorder struct {
	path string
}

func (p *PluginDanmuRecorder) LiveStart(process *videoworker.ProcessVideo) error {
	return nil
}

func (p *PluginDanmuRecorder) DownloadStart(process *videoworker.ProcessVideo) error {
	video := process.LiveStatus.Video
	pathSlice := []string{config.Config.UploadDir, process.LiveStatus.Video.UsersConfig.Name, process.GetFullTitle()}
	p.path = strings.Join(pathSlice, "/")
	err := callJsAPI(getRoomId(video.UsersConfig.TargetId), "1", p.path)
	if err != nil {
		return err
	}
	return nil
}

func (p *PluginDanmuRecorder) LiveEnd(process *videoworker.ProcessVideo) error {
	video := process.LiveStatus.Video
	path := video.FilePath
	if config.Config.EnableTS2MP4 {
		path = strings.TrimSuffix(path, ".mp4")
	}
	err := callJsAPI(getRoomId(video.UsersConfig.TargetId), "0", path)
	if err != nil {
		return err
	}
	return nil
}
