package utils

import (
	"context"
	"github.com/fzxiao233/Vtb_Record/config"
	"github.com/gogf/greuse"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"runtime"
	"runtime/debug"
	"time"
)

func PrintMemUsage() {
	bToMb := func(b uint64) uint64 {
		return b / 1024 / 1024
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	log.Warnf("Alloc = %v MiB\tTotalAlloc = %v MiB\tSys = %v MiB\tGoroutines = %v\tNumGC = %v",
		bToMb(m.Alloc),
		bToMb(m.TotalAlloc),
		bToMb(m.Sys),
		runtime.NumGoroutine(),
		m.NumGC)
}

var PprofServer *http.Server

func InitProfiling() {
	go func() {
		log.Warnf("Starting pprof server")
		ticker := time.NewTicker(time.Minute * 1)
		for {
			//go http.ListenAndServe("0.0.0.0:49314", nil)
			if PprofServer == nil || PprofServer.Addr != config.Config.PprofHost {
				if PprofServer != nil {
					go PprofServer.Shutdown(context.Background())
				}
				//PprofServer = &http.Server{Addr: config.Config.PprofHos5t, Handler: nil}
				listener, err := greuse.Listen("tcp", config.Config.PprofHost)
				if listener == nil {
					log.Warnf("Error creating reusable listener, creating a normal one instead!")
					listener, err = net.Listen("tcp", config.Config.PprofHost)
				}
				if err != nil {
					log.Warnf("Failed to reuse-listen addr: %s", err)
				}
				PprofServer = &http.Server{}
				//go PprofServer.ListenAndServe()
				go PprofServer.Serve(listener)
			}
			<-ticker.C
		}
	}()

	go func() {
		ticker := time.NewTicker(time.Minute * 1)
		for {
			PrintMemUsage()
			<-ticker.C
		}
	}()

	go func() {
		ticker := time.NewTicker(time.Millisecond * 3000)
		for {
			//start := time.Now()
			runtime.GC()
			//log.Debugf("G	C & scvg use %s", time.Now().Sub(start))
			<-ticker.C
		}
	}()

	ticker := time.NewTicker(time.Second * 5)
	for {

		start := time.Now()
		debug.FreeOSMemory()
		log.Debugf("scvg use %s", time.Now().Sub(start))
		<-ticker.C
	}
}
