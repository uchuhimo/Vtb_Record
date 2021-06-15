package provgo

import (
	"crypto/tls"
	"fmt"
	"github.com/fzxiao233/Vtb_Record/utils"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
)

func doDownloadHttp(entry *log.Entry, output string, httpURL string, headers map[string]string, needMove bool, proxy string) error {
	// Create the file
	/*out, err := os.Create(output)
	if err != nil {
		return err
	}
	if !needMove {
		defer func () {
			go out.Close()
		}()
	} else {
		defer out.Close()
	}*/
	out := utils.GetWriter(output)
	defer out.Close()

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	if proxy != "" {
		proxyURL, err := url.Parse("http://"+proxy)
		if err != nil {
			return err
		}
		entry.Debugf("use proxy %s for %s", proxy, httpURL)
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	client := &http.Client{
		Transport: transport,
	}
	// Get the data
	req, _ := http.NewRequest("GET", httpURL, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("downloader got bad status: %s", resp.Status)
	}

	buf := make([]byte, 1024*1024*3) // 1M buffer
	src := resp.Body
	dst := out
	for {
		// Writer the body to file
		written := int64(0)
		for {
			nr, er := src.Read(buf)
			if nr > 0 {
				nw, ew := dst.Write(buf[0:nr])
				if nw > 0 {
					written += int64(nw)
				}
				if ew != nil {
					err = ew
					break
				}
				if nr != nw {
					err = io.ErrShortWrite
					break
				}
			}
			if er != nil {
				err = er
				break
			}
		}

		//written, err := io.CopyBuffer(out, resp.Body, buf)
		entry.Infof("Wrote %d, err: %s", written, err)
		if err == nil {
			return nil
		} else if err == io.EOF {
			entry.Info("Stream ended")
			return nil
		} else {
			return err
		}
	}
}
