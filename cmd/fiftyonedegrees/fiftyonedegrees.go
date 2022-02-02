package main

import (
	"fmt"
	"github.com/adscompass/useragent"
	"io/ioutil"
	"strconv"
	"strings"

	FiftyOneDegreesTrieV3 "gitlab.com/adscompass/f51degrees"
	"go.uber.org/zap"
)

const (
	versionUnknown       = "Unknown"
	versionPlatformXP    = "XP"
	versionPlatformVista = "Vista"
)

type FiftyOneDegreesData struct {
	BrowserName,
	BrowserVersion,
	PlatformName,
	PlatformVersion,
	HardwareModel,
	HardwareFamily,
	HardwareVendor,
	DeviceType string

	NumericBrowserVersion,
	NumericPlatformVersion []int32

	IsMobile,
	IsCrawler bool
}

func main() {
	var err error
	provider := FiftyOneDegreesTrieV3.NewProvider("/home/ivan/data/FIFTYONEDEGREES.BIN_V3")
	d := &FiftyOneDegreesData{}

	openFileName := "cmd/fiftyonedegrees/files/uaAndroidSort.txt"
	f, err := ioutil.ReadFile(openFileName)
	if err != nil {
		fmt.Println(err)
	}
	if len(f) == 0 {
		fmt.Println("data = 0")
		return
	}

	dataFile := strings.Split(string(f), "\n")

	uas := make([]string, 0, len(dataFile))
	for _, data := range dataFile {
		if strings.Contains(data, "\"") {
			data = data[1 : len(data)-1]
		}
		fmt.Println(data)
		uas = append(uas, data)
	}

	if len(uas) == 0 {
		fmt.Println("len(uas) == 0")
		return
	}

	//uas := []comperestring{
	//	"Mozilla/5.0 (Android 8.0.0; Mobile; rv:91.0) Gecko/91.0 Firefox/91.0",
	//}
	//

	fmt.Printf("%24s %14s %24s %14s \n", "DeviceModel", "Vendor", "Name", "DeviceType")
	fmt.Println("-----------------------------------------------------------------------------------------")

	var fs string
	divisor := ","
	for _, ua := range uas {
		d.Parse(ua, provider)
		fs += d.HardwareModel + divisor + d.HardwareVendor + divisor + d.HardwareFamily + divisor + d.DeviceType + "\n"
		fmt.Printf("%24s, %14s, %24s, %14s \n", d.HardwareModel, d.HardwareVendor, d.HardwareFamily, d.DeviceType)
	}

	saveDataFile := []byte(fs)
	saveFileName := "cmd/fiftyonedegrees/files/devicemodels.txt"

	err = ioutil.WriteFile(saveFileName, saveDataFile, 0644)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("------------------------------------------------------------------------------------------------")
	fmt.Println("saved file: ", saveFileName)
	fmt.Println("Done")

}

func (data *FiftyOneDegreesData) Parse(ua string, provider FiftyOneDegreesTrieV3.Provider) {
	match := provider.GetMatch(ua)

	data.BrowserName = match.GetValue("BrowserName")
	data.BrowserVersion = match.GetValue("BrowserVersion")
	data.PlatformName = match.GetValue("PlatformName")
	data.PlatformVersion = match.GetValue("PlatformVersion")
	data.HardwareFamily = match.GetValue("HardwareFamily")
	data.HardwareModel = match.GetValue("HardwareModel")
	data.HardwareVendor = match.GetValue("HardwareVendor")
	data.DeviceType = match.GetValue("DeviceType")
	data.IsMobile = match.GetValueAsBool("IsMobile")
	data.IsCrawler = match.GetValueAsBool("IsCrawler")

	if data.BrowserVersion != versionUnknown {
		versionSegments := strings.Split(data.BrowserVersion, ".")

		for _, v := range versionSegments {
			n, err := strconv.Atoi(v)

			if err != nil {
				fmt.Println("error parse browser version to numeric, browaser version", data.BrowserVersion, zap.Error(err))
				data.NumericBrowserVersion = data.NumericBrowserVersion[:0]
				continue
			}

			data.NumericBrowserVersion = append(data.NumericBrowserVersion, int32(n))
		}
	}

	if data.PlatformVersion != versionUnknown {
		if data.PlatformVersion == versionPlatformXP {
			data.NumericPlatformVersion = append(data.NumericPlatformVersion, 5)
		} else if data.PlatformVersion == versionPlatformVista {
			data.NumericPlatformVersion = append(data.NumericPlatformVersion, 6)
		} else {
			versionSegments := strings.Split(data.PlatformVersion, ".")

			for _, v := range versionSegments {
				n, err := strconv.Atoi(v)

				if err != nil {
					fmt.Println("error parse platform version to numeric, platform version", data.PlatformVersion, zap.Error(err))
					data.NumericPlatformVersion = data.NumericPlatformVersion[:0]
					continue
				}

				data.NumericPlatformVersion = append(data.NumericPlatformVersion, int32(n))
			}
		}
	}

	FiftyOneDegreesTrieV3.DeleteMatch(match)

	//if data.BrowserName == versionUnknown {
	//	data.parseUserAgent(ua)
	//}
}

func (data *FiftyOneDegreesData) Reset() {
	data.BrowserName = ""
	data.BrowserVersion = ""
	data.NumericBrowserVersion = data.NumericBrowserVersion[:0]
	data.PlatformName = ""
	data.PlatformVersion = ""
	data.NumericPlatformVersion = data.NumericPlatformVersion[:0]
	data.HardwareModel = ""
	data.HardwareFamily = ""
	data.HardwareVendor = ""
	data.DeviceType = ""

	data.IsMobile = false
	data.IsCrawler = false
}

func (data *FiftyOneDegreesData) parseUserAgent(ua string) {
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		fmt.Printf("USERAGENT PARSE PANIC: [%s] %v", ua, e)
	}()

	dataUserAgent := useragent.AcquireData()
	defer useragent.ReleaseData(dataUserAgent)

	if !useragent.Parse([]byte(ua), dataUserAgent) {
		return
	}

	data.Reset()

	data.BrowserName = string(dataUserAgent.BrowserName)
	data.BrowserVersion = string(dataUserAgent.BrowserVersion)
	data.PlatformName = string(dataUserAgent.PlatformName)
	data.PlatformVersion = string(dataUserAgent.PlatformVersion)
	data.IsMobile = dataUserAgent.IsMobile
	data.IsCrawler = dataUserAgent.IsCrawler
}
