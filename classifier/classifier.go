package classifier

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/peterbourgon/diskv"
	"github.com/sourcegraph/httpcache/diskcache"
	"github.com/sourcegraph/thesrc"
)

func Classify(post *thesrc.Post) (string, error) {
	if post.LinkURL == "" {
		return "", nil
	}

	resp, err := httpClient.Get(post.LinkURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non-200 HTTP response status: %d", resp.StatusCode)
	}

	page, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	allLen := len(page.Find("body").Text())
	codeLen := len(page.Find("code").Text() + page.Find("pre").Text()) // might double-count

	prop := float64(codeLen) / float64(allLen)
	summary := fmt.Sprintf("%.1f%% code (%d/%d)", prop, codeLen, allLen)

	var classif string
	if prop > 0.07 || codeLen > 300 {
		classif = "CODE"
	} else {
		classif = "NOTCODE"
	}

	return classif + " " + summary, nil
}

var (
	// httpCacheDir is the directory used for caching HTTP responses. It can be reused
	// executions (it is not necessary to create a new random temp dir upon
	// startup).
	httpCacheDir = "/tmp/thesrc-http-cache"

	localCache = diskcache.NewWithDiskv(diskv.New(diskv.Options{
		BasePath:     httpCacheDir,
		CacheSizeMax: 50 * 1024 * 1024 * 1024, // 50 GB
	}))

	httpClient = &http.Client{
		//Transport: &httpcache.Transport{Cache: localCache},
		Timeout: time.Second * 3,
		// TODO(sqs): add timeout
	}
)

func init() {
	if err := os.Mkdir(httpCacheDir, 0700); err != nil && !os.IsExist(err) {
		log.Fatalf("Mkdir(%s) failed: %s", httpCacheDir, err)
	}
}
