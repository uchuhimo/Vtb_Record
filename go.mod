module github.com/fzxiao233/Vtb_Record

go 1.14

require (
	github.com/bitly/go-simplejson v0.5.0
	github.com/etherlabsio/go-m3u8 v0.1.2
	github.com/fsnotify/fsnotify v1.4.9
	github.com/fzxiao233/Go-Emoji-Utils v0.0.0-20200305114615-005e99b02c2f
	github.com/go-redis/redis v6.15.8+incompatible
	github.com/go-redis/redis/v8 v8.7.1
	github.com/gogf/greuse v1.1.0
	github.com/gopherjs/gopherjs v0.0.0-20200217142428-fce0ec30dd00 // indirect
	github.com/hashicorp/golang-lru v0.5.4
	github.com/knq/sdhook v0.0.0-20190801142816-0b7fa827d09a
	github.com/mitchellh/mapstructure v1.1.2
	github.com/orandin/lumberjackrus v1.0.1
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/rclone/rclone v1.52.2
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/viper v1.6.2
	github.com/tidwall/gjson v1.6.0
	github.com/tidwall/pretty v1.0.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0
	github.com/xanzy/ssh-agent v0.3.0 // indirect
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/ratelimit v0.1.0
	golang.org/x/sync v0.0.0-20201020160332-67f06af15bc9
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	google.golang.org/protobuf v1.23.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	honnef.co/go/tools v0.0.1-2020.1.6 // indirect
	storj.io/common v0.0.0-20201030140758-31112c1cc750 // indirect
	storj.io/uplink v1.3.1 // indirect
)

replace github.com/rclone/rclone v1.52.2 => github.com/NyaMisty/rclone v1.52.2-mod

//replace github.com/smallnest/ringbuffer => ../../smallnest/ringbuffer
