// https://habr.com/ru/post/499164/
// https://yandex.ru/support/adfox-sites/banners/specs/video.html

package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"time"

	"github.com/adscompass/vast"
)

func main() {
	//fmt.Println("vast1:")
	//fmt.Println(vast1())
	//fmt.Println()
	//fmt.Println("vast2:")
	//fmt.Println(vast2())

	fmt.Println("buildAdm:")
	fmt.Println(buildAdm())

	fmt.Println()
	fmt.Println("Done...")
}

func buildAdm() string {
	x := vast.VAST{
		Version: "3.0",
		Ads: []vast.Ad{
			{
				ID: "01",
				Wrapper: &vast.Wrapper{
					AdSystem: &vast.AdSystem{Version: "1.0", Name: "Ivan"},
					VASTAdTagURI: vast.CDATAString{
						CDATA: "vast.xsd",
					},
					Impressions: []vast.Impression{
						{
							ID:  "2999",
							URI: "127.0.0.1:2999",
						},
					},
				},
			},
		},
	}

	buf := bytes.NewBuffer(nil)

	errEncode := xml.NewEncoder(buf).Encode(x)
	if errEncode != nil {
		log.Printf("error encode, %v", errEncode)
	}

	res := xml.Header + string(buf.Bytes())

	return res
}

func vast2() string {
	x := vast.VAST{
		Version: "3.0",
		Ads: []vast.Ad{
			{
				ID: "1",
				Wrapper: &vast.Wrapper{
					AdSystem: &vast.AdSystem{Name: "Adscompass"},
					VASTAdTagURI: vast.CDATAString{
						CDATA: "https://2000.200.adscompass.ru/1/vast",
					},
					Impressions: []vast.Impression{
						{
							ID:  "1",
							URI: "https://2000.200.adscompass.ru/2/imp",
						},
					},
					Creatives: []vast.CreativeWrapper{
						{
							ID:       "1",
							Sequence: 0,
							AdID:     "1",
							Linear: &vast.LinearWrapper{
								Icons: nil,
								TrackingEvents: []vast.Tracking{
									{
										Event: "start",
										URI:   "https://2000.200.adscompass.ru/2/event/start",
									},
									{
										Event: "skip",
										URI:   "https://2000.200.adscompass.ru/2/event/skip",
									},
									{
										Event: "complete",
										URI:   "https://2000.200.adscompass.ru/2/event/complete",
									},
								},
								VideoClicks: nil,
							},
							CompanionAds: nil,
							NonLinearAds: nil,
						},
					},
				},
			},
		},
	}

	buf := bytes.NewBuffer(nil)

	errEncode := xml.NewEncoder(buf).Encode(x)
	if errEncode != nil {
		log.Printf("error encode, %v", errEncode)
	}

	return string(buf.Bytes())
}

func vast1() string {
	d := vast.Duration(time.Second * 3)

	x := vast.VAST{
		Version: "3.0",
		Ads: []vast.Ad{
			{
				ID: "1",
				InLine: &vast.InLine{
					AdSystem: &vast.AdSystem{Name: "Adscompass"},
					AdTitle:  vast.CDATAString{CDATA: "Title1"},
					Impressions: []vast.Impression{
						{
							URI: "https://2000.200.adscompass.ru/1/imp",
						},
					},
					Creatives: []vast.Creative{
						{
							Linear: &vast.Linear{
								SkipOffset: &vast.Offset{
									Duration: &d,
								},
								Duration: vast.Duration(time.Second * 10),
								TrackingEvents: []vast.Tracking{
									{
										Event: "start",
										URI:   "https://2000.200.adscompass.ru/1/event/start",
									},
									{
										Event: "skip",
										URI:   "https://2000.200.adscompass.ru/1/event/skip",
									},
									{
										Event: "progress",
										Offset: &vast.Offset{
											Duration: &d,
										},
										URI: "https://2000.200.adscompass.ru/1/event/progress3sec",
									},
								},
								VideoClicks: &vast.VideoClicks{
									ClickThroughs: []vast.VideoClick{
										{
											ID:  "1",
											URI: "https://2000.200.adscompass.ru/1/clickthrough?n=1",
										},
										//{
										//	ID:  "2",
										//	URI: "https://2000.200.adscompass.ru/1/clickthrough?n=2",
										//},
									},
									ClickTrackings: []vast.VideoClick{
										{
											ID:  "1",
											URI: "https://2000.200.adscompass.ru/1/clicktrack?n=1",
										},
										//{
										//	ID:  "2",
										//	URI: "https://2000.200.adscompass.ru/1/clicktrack?n=2",
										//},
									},
									CustomClicks: nil,
								},
								MediaFiles: []vast.MediaFile{
									{
										Delivery: "progressive",
										Type:     "video/mp4",
										Width:    1280,
										Height:   720,
										URI:      "https://arts.fra1.digitaloceanspaces.com/file.mp4",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	buf := bytes.NewBuffer(nil)

	errEncode := xml.NewEncoder(buf).Encode(x)
	if errEncode != nil {
		log.Printf("error encode, %v", errEncode)
	}

	return string(buf.Bytes())
}
