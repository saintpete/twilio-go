package twilio

import (
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type MediaService struct {
	client *Client
}

func mediaPathPart(messageSid string) string {
	return "Messages/" + messageSid + "/Media"
}

type MediaPage struct {
	Page
	MediaList []*Media `json:"media_list"`
}

type Media struct {
	Sid         string     `json:"sid"`
	ContentType string     `json:"content_type"`
	AccountSid  string     `json:"account_sid"`
	DateCreated TwilioTime `json:"date_created"`
	DateUpdated TwilioTime `json:"date_updated"`
	ParentSid   string     `json:"parent_sid"`
	URI         string     `json:"uri"`
}

// MediaClient is used for fetching images and does not follow redirects.
var MediaClient = http.Client{
	Timeout: defaultTimeout,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

func (m *MediaService) GetPage(messageSid string, data url.Values) (*MediaPage, error) {
	mp := new(MediaPage)
	err := m.client.ListResource(mediaPathPart(messageSid), data, mp)
	return mp, err
}

// Get returns a Media struct representing a Media instance, or an error.
func (m *MediaService) Get(messageSid string, sid string) (*Media, error) {
	me := new(Media)
	err := m.client.GetResource(mediaPathPart(messageSid), sid, me)
	return me, err
}

// GetURL returns a URL that can be retrieved to download the given image.
func (m *MediaService) GetURL(messageSid string, sid string) (*url.URL, error) {
	uriEnd := strings.Join([]string{mediaPathPart(messageSid), sid}, "/")
	path := m.client.FullPath(uriEnd)
	// We want the media, not the .json representation
	if strings.HasSuffix(path, ".json") {
		path = path[:len(path)-len(".json")]
	}
	urlStr := m.client.Client.Base + path
	count := 0
	for {
		req, err := http.NewRequest("GET", urlStr, nil)
		if err != nil {
			return nil, err
		}
		req.SetBasicAuth(m.client.AccountSid, m.client.AuthToken)
		req.Header.Set("User-Agent", userAgent)
		resp, err := MediaClient.Do(req)
		if err != nil {
			return nil, err
		}
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
		// This is brittle because we need to detect/rewrite the S3 URL.
		// I don't want to hard code a S3 URL but we have to do some
		// substitution.
		location := resp.Header.Get("Location")
		if location == "" {
			return nil, errors.New("twilio: Couldn't follow redirect")
		}
		u, err := url.Parse(location)
		if err != nil {
			return nil, err
		}
		if strings.Contains(u.Host, "media.twiliocdn.com.") && strings.Contains(u.Host, "amazonaws") {
			// This is the URL we can use to download the content. The URL that
			// Twilio gives us back is insecure and uses HTTP. Rewrite it to
			// use the HTTPS path-based URL scheme.
			//
			// https://docs.aws.amazon.com/AmazonS3/latest/dev/UsingBucket.html
			if u.Scheme == "http" {
				u.Host = strings.Replace(u.Host, "media.twiliocdn.com.", "", 1)
				u.Path = "/media.twiliocdn.com" + u.Path
				u.Scheme = "https"
			}
			return u, nil
		}
		count++
		if count > 5 {
			return nil, errors.New("twilio: too many redirects")
		}
		urlStr = location
	}
}

// GetImage downloads a Media object and returns an image.Image. The
// documentation isn't great on what happens - as of October 2016, we make a
// request to the Twilio API, then to media.twiliocdn.com, then to a S3 URL. We
// then download that image and decode it based on the provided content-type.
func (m *MediaService) GetImage(messageSid string, sid string) (image.Image, error) {
	u, err := m.GetURL(messageSid, sid)
	if err != nil {
		return nil, err
	}
	if u.Scheme == "http" {
		return nil, fmt.Errorf("Attempted to download image over insecure URL: %s", u.String())
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := MediaClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// https://www.twilio.com/docs/api/rest/accepted-mime-types#supported
	ctype := resp.Header.Get("Content-Type")
	switch ctype {
	case "image/jpeg":
		return jpeg.Decode(resp.Body)
	case "image/gif":
		return gif.Decode(resp.Body)
	case "image/png":
		return png.Decode(resp.Body)
	default:
		io.Copy(ioutil.Discard, resp.Body)
		return nil, fmt.Errorf("twilio: Unknown content-type %s", ctype)
	}
}
