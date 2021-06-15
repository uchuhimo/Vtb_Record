package bilibili

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/fzxiao233/Vtb_Record/config"
	"github.com/fzxiao233/Vtb_Record/live/interfaces"
	"github.com/fzxiao233/Vtb_Record/live/monitor/base"
	"github.com/fzxiao233/Vtb_Record/utils"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Bilibili struct {
	base.BaseMonitor
	TargetId      string
	Title         string
	isLive        bool
	streamingLink string
}

type BilibiliPoller struct {
	LivingUids map[int]base.LiveInfo
	lock       sync.Mutex
}

var Poller BilibiliPoller

func (b *BilibiliPoller) getStatusUseBatch() error {
	ctx := base.GetCtx("Bilibili")
	url := ctx.ApiHostUrl
	if url == "" {
		url = "https://api.live.bilibili.com"
	}

	biliMod := base.GetMod("Bilibili")
	allUids := make([]string, 0)
	for _, u := range biliMod.Users {
		allUids = append(allUids, u.TargetId)
	}
	rand.Shuffle(len(allUids), func(i, j int) { allUids[i], allUids[j] = allUids[j], allUids[i] })
	livingUids := make(map[int]base.LiveInfo)
	for i := 0; i < len(allUids); i += 200 {
		payload := fmt.Sprintf("%s/room/v1/Room/get_status_info_by_uids?uids[]=%s",
			url,
			strings.Join(allUids[i:utils.Min(i+200, len(allUids))], "&uids[]="))
		rawInfoJSON, err := ctx.HttpGet(payload, map[string]string{})
		if err != nil {
			return err
		}
		infoJson, _ := simplejson.NewJson(rawInfoJSON)
		users := infoJson.Get("data")
		userMap := users.MustMap()
		for uid := range userMap {
			user := users.Get(uid)
			if user.Get("live_status").MustInt() == 1 {
				liveUrl := fmt.Sprintf("https://live.bilibili.com/%d", user.Get("room_id").MustInt())
				livingUids[user.Get("uid").MustInt()] = base.LiveInfo{
					Title:         user.Get("title").MustString(),
					StreamingLink: liveUrl,
				}
			}
		}
	}
	b.LivingUids = livingUids
	log.Tracef("Parsed uids: %v", b.LivingUids)
	return nil
}

func (b *BilibiliPoller) GetStatus() error {
	return b.getStatusUseBatch()
}

func (b *BilibiliPoller) StartPoll() error {
	err := b.GetStatus()
	if err != nil {
		return err
	}
	go func() {
		for {
			biliMod := base.GetMod("Bilibili")
			_interval := biliMod.PollInterval
			interval := time.Duration(config.Config.CriticalCheckSec) * time.Second
			if _interval != 0 {
				interval = time.Duration(_interval) * time.Second
			}
			time.Sleep(interval)
			err := b.GetStatus()
			if err != nil {
				log.WithError(err).Warnf("Error during polling GetStatus")
			}
		}
	}()
	return nil
}

func (b *BilibiliPoller) IsLiving(uid int) *base.LiveInfo {
	b.lock.Lock()
	if b.LivingUids == nil {
		err := b.StartPoll()
		if err != nil {
			log.WithError(err).Warnf("Failed to poll from bilibili")
		}
	}
	b.lock.Unlock()
	info, ok := b.LivingUids[uid]
	if ok {
		return &info
	} else {
		return nil
	}
}

func (b *Bilibili) getVideoInfoByPoll() error {
	uid, err := strconv.Atoi(b.TargetId)
	if err != nil {
		return err
	}
	ret := Poller.IsLiving(uid)
	b.isLive = ret != nil
	if !b.isLive {
		return nil
	}

	b.streamingLink = ret.StreamingLink
	b.Title = ret.Title
	return nil
}

func (b *Bilibili) getVideoInfoByRoom() error {
	url := b.Ctx.ApiHostUrl
	if url == "" {
		url = "https://api.live.bilibili.com"
	}
	rawInfoJSON, err := b.Ctx.HttpGet(url+"/room/v1/Room/getRoomInfoOld?mid="+b.TargetId, map[string]string{})
	if err != nil {
		return err
	}
	infoJson, _ := simplejson.NewJson(rawInfoJSON)
	livestatus := infoJson.Get("data").Get("liveStatus").MustInt()
	b.isLive = livestatus == 1
	b.streamingLink = infoJson.Get("data").Get("url").MustString("")
	b.Title = infoJson.Get("data").Get("title").MustString("")
	return nil
}

func (b *Bilibili) CreateVideo(usersConfig config.UsersConfig) *interfaces.VideoInfo {
	v := &interfaces.VideoInfo{
		Title:       b.Title,
		Date:        utils.GetTimeNow(),
		Target:      b.streamingLink,
		Provider:    "Bilibili",
		UsersConfig: usersConfig,
	}
	return v
}

func (b *Bilibili) CheckLive(usersConfig config.UsersConfig) bool {
	b.TargetId = usersConfig.TargetId
	var err error
	if b.Ctx.UseFollowPolling {
		err = b.getVideoInfoByPoll()
	} else {
		err = b.getVideoInfoByRoom()
	}

	if err != nil {
		b.isLive = false
		log.WithField("user", fmt.Sprintf("%s|%s", "Bilibili", usersConfig.Name)).WithError(err).Errorf("GetVideoInfo error")
	}
	if !b.isLive {
		base.NoLiving("Bilibili", usersConfig.Name)
	}
	return b.isLive
}
