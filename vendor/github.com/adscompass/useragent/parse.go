package useragent

import (
	"bytes"
	"math"
)

// Parse the useragent and fill the Data struct. Returns false if parsing failed
func Parse(useragent []byte, data *Data) bool {
	var idx, start int

	if bytes.HasPrefix(useragent, mozillaPrefix) {
		idx = len(mozillaPrefix)
		start = idx
	} else if bytes.HasPrefix(useragent, []byte("Go-http-client/")) {
		data.PlatformName = append(data.PlatformName, unknown...)
		data.PlatformVersion = append(data.PlatformVersion, unknown...)
		data.BrowserName = append(data.BrowserName, brwGoHTTPClient...)
		data.BrowserVersion = append(data.BrowserVersion, unknown...)
		data.IsMobile = false
		data.IsCrawler = true
		return true
	} else if bytes.HasPrefix(useragent, []byte("Dalvik/2.1.0")) {
		data.PlatformName = append(data.PlatformName, pltAndroid...)
		idx := bytes.Index(useragent, strAndroid)
		if idx > 0 {
			idx++
			if len(useragent) < idx+len(strAndroid) {
				return false
			}
			i := bytes.IndexByte(useragent[idx+len(strAndroid):], ';')
			if i > 0 {
				data.PlatformVersion = append(data.PlatformVersion, useragent[idx+len(strAndroid):idx+len(strAndroid)+i]...)
				replaceUnderscoreToDot(data.PlatformVersion)
			}
		}
		if len(data.PlatformVersion) == 0 {
			data.PlatformVersion = append(data.PlatformVersion, unknown...)
		}
		data.BrowserName = append(data.BrowserName, brwAndroid...)
		if bytes.Contains(useragent, strSmartTV) {
			data.BrowserName = append(data.BrowserName[:0], brwSmartTV...)
		}
		data.BrowserVersion = append(data.BrowserVersion[:0], unknown...)
		data.IsMobile = true
		return true
	} else if bytes.HasPrefix(useragent, []byte("iOS/")) {
		// Календари iOS
		// iOS/11.1.2 (15B202) dataaccessd/1.0 для запроса данных
		// iOS/14.0 (18A373) accountsd/1.0 для подписки
		if !(bytes.HasSuffix(useragent, []byte("dataaccessd/1.0")) ||
			bytes.HasSuffix(useragent, []byte("accountsd/1.0"))) {
			return false
		}

		data.BrowserName = append(data.BrowserName, []byte("iOS Calendar")...)
		data.BrowserVersion = append(data.BrowserVersion, unknown...)
		data.PlatformName = append(data.PlatformName, pltIOS...)
		data.PlatformVersion = parseDottedVersion(useragent[len("iOS/"):], data.PlatformVersion)
		data.IsMobile = true

		return true
	} else {

		return false
	}

	endPlatformBlock := false

	for {
		idx++
		if len(useragent) <= idx {
			return false
		}
		if useragent[idx] == ';' || useragent[idx] == ',' {
			break
		}
		// если первый блок не содержит разделения групп (точка с запятой) - устанавливаем флаг endPlatformBlock
		if useragent[idx] == ')' {
			endPlatformBlock = true
			break
		}
	}

	if bytes.HasPrefix(useragent[start:idx], strWindowsNT) {
		data.PlatformName = append(data.PlatformName, pltWindows...)
		data.PlatformVersion = append(data.PlatformVersion, replaceWindowsVersion(useragent[start+len(strWindowsNT):idx])...)
	} else if bytes.Equal(useragent[start:idx], strX11) {
		data.PlatformName = append(data.PlatformName, pltLinux...)
		data.PlatformVersion = append(data.PlatformVersion, unknown...)
		if len(useragent) <= start+len(pltLinux) {
			return false
		}
		if bytes.HasPrefix(useragent[start+len(pltLinux):], []byte("CrOS ")) {
			data.PlatformName = append(data.PlatformName, pltChromeOS...)
		}
	} else if bytes.HasPrefix(useragent[start:idx], strAndroid) {
		// случай, когда Android идет сразу после открывающей скобки
		// Mozilla/5.0 (Android 10; Mobile; rv:86.0) Gecko/86.0 Firefox/86.0
		data.PlatformName = append(data.PlatformName, pltAndroid...)
		data.IsMobile = true

		i := start + len(pltAndroid) + 1

		if len(useragent) < i {
			data.PlatformVersion = append(data.PlatformVersion, unknown...)
		} else {
			data.PlatformVersion = append(data.PlatformVersion, parseDottedVersion(useragent[i:], data.PlatformVersion)...)
		}

	} else if bytes.Equal(useragent[start:idx], strLinux) {
		// Если у нас Linux, то не должно заканчиваться круглой скобкой, после точки с запятой должна быть версия платформы (например Android 8.1.5)
		if endPlatformBlock {
			return false
		}

		idx += 2 // '; '

		// Далее в этом месте может быть вставка U;, arm;, arm_64; Пропускаем ее

		if len(useragent) <= idx+idx+len(blockAfterPlatformArm64) {
			return false
		}
		if bytes.Equal(useragent[idx:idx+len(blockAfterPlatformU)], blockAfterPlatformU) {
			idx += len(blockAfterPlatformU) + 1
		}
		if bytes.Equal(useragent[idx:idx+len(blockAfterPlatformArm)], blockAfterPlatformArm) {
			idx += len(blockAfterPlatformArm) + 1
		}
		if bytes.Equal(useragent[idx:idx+len(blockAfterPlatformArm64)], blockAfterPlatformArm64) {
			idx += len(blockAfterPlatformArm64) + 1
		}

		// Парсим строку Android
		i := idx
		for {
			i++
			if len(useragent) <= i {
				return false
			}
			if useragent[i] == ';' {
				break
			}
			if useragent[i] == ')' {
				break
			}
		}

		if bytes.HasPrefix(useragent[idx:i], strAndroid) {
			if len(useragent[idx:i]) < 9 { // как минимум Android 1
				return false
			}
			data.PlatformName = append(data.PlatformName, pltAndroid...)
			data.PlatformVersion = append(data.PlatformVersion, useragent[idx+8:i]...)
			data.IsMobile = true
		} else {
			// Не смогли распарсить Mozilla/5.0 (Linux; Something...
			return false
		}
	} else if bytes.Equal(useragent[start:idx], strIPhone) { // iPhone
		//idx += 2 // '; '
		data.PlatformName = append(data.PlatformName, pltIOS...)
		data.IsMobile = true
		i := bytes.Index(useragent[start+len(strIPhone):], strSpaceOSSpace)
		if i == -1 {
			return false
		}
		idx += i + len(strSpaceOSSpace)
		//idx += cpuIPhoneOSLen
		if idx >= len(useragent) {
			return false
		}
		i = bytes.Index(useragent[idx:], strLikeMacOSX)
		if i == -1 {
			return false
		}
		data.PlatformVersion = append(data.PlatformVersion, useragent[idx:idx+i]...)
		replaceUnderscoreToDot(data.PlatformVersion)
	} else if bytes.Equal(useragent[start:idx], strIPad) { // iPad
		//idx += 2 // '; '
		data.PlatformName = append(data.PlatformName, pltIPadOS...)
		data.IsMobile = true
		//idx += cpuOSLen
		i := bytes.Index(useragent[start+len(strIPad):], strSpaceOSSpace)
		if i == -1 {
			return false
		}
		idx += i + len(strSpaceOSSpace)
		if idx >= len(useragent) {
			return false
		}
		i = bytes.Index(useragent[idx:], strLikeMacOSX)
		if i == -1 {
			return false
		}
		data.PlatformVersion = append(data.PlatformVersion, useragent[idx:idx+i]...)
		replaceUnderscoreToDot(data.PlatformVersion)
	} else if bytes.Equal(useragent[start:idx], strIPod) { // iPod
		//idx += 2 // '; '
		data.PlatformName = append(data.PlatformName, pltIOS...)
		data.IsMobile = true
		//idx += cpuOSLen
		i := bytes.Index(useragent[start+len(strIPod):], strSpaceOSSpace)
		if i == -1 {
			return false
		}
		idx += i + len(strSpaceOSSpace)
		if idx >= len(useragent) {
			return false
		}
		i = bytes.Index(useragent[idx:], strLikeMacOSX)
		if i == -1 {
			return false
		}
		data.PlatformVersion = append(data.PlatformVersion, useragent[idx:idx+i]...)
		replaceUnderscoreToDot(data.PlatformVersion)
	} else if bytes.Equal(useragent[start:idx], strIPodTouch) { // iPod touch
		//idx += 2 // '; '
		data.PlatformName = append(data.PlatformName, pltIOS...)
		data.IsMobile = true
		//idx += cpuOSLen
		i := bytes.Index(useragent[start+len(strIPodTouch):], strSpaceOSSpace)
		if i == -1 {
			return false
		}
		idx += i + len(strSpaceOSSpace)
		if idx >= len(useragent) {
			return false
		}
		i = bytes.Index(useragent[idx:], strLikeMacOSX)
		if i == -1 {
			return false
		}
		data.PlatformVersion = append(data.PlatformVersion, useragent[idx:idx+i]...)
		replaceUnderscoreToDot(data.PlatformVersion)
	} else if bytes.Equal(useragent[start:idx], strMacintosh) { // macintosh
		idx += 2 // '; '
		data.PlatformName = append(data.PlatformName, pltMacOS...)
		idx += intelMacOSXLen
		if idx >= len(useragent) {
			return false
		}
		i := bytes.IndexByte(useragent[idx:], ')')
		if i == -1 {
			return false
		}
		// Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:89.0) Gecko/20100101 Firefox/89.0
		//                                             ^ для случая, когда есть '; rv...'
		i2 := bytes.IndexByte(useragent[idx:], ';')
		if i2 > 0 && i2 < i {
			i = i2
		}

		data.PlatformVersion = append(data.PlatformVersion, useragent[idx:idx+i]...)
		replaceUnderscoreToDot(data.PlatformVersion)
		//data.PlatformVersion = dotZero(data.PlatformVersion)
	} else if bytes.Equal(useragent[start:idx], strCompatible) || bytes.Equal(useragent[start:idx], strCompatibleComma) {
		if len(useragent) < start+len(strCompatible)+2 {
			return false
		}
		return parseCompatible(useragent[start+len(strCompatible)+2:], data) // 2 = '; '
	} else {
		// в первом блоке круглых скобок не нашли ничего знакомого
		return false
	}

	if i := bytes.Index(useragent[idx:], []byte("Trident/")); i != -1 {
		data.BrowserName = append(data.BrowserName, brwInternetExplorer...)
		rvIdx := bytes.Index(useragent[idx+i:], []byte("; rv:"))
		if rvIdx == -1 {
			data.BrowserVersion = append(data.BrowserVersion, unknown...)
			return true
		}
		rvIdx += 5 // '; rv:'

		rvEndIdx := bytes.Index(useragent[idx+i+rvIdx:], []byte(")"))
		if rvEndIdx == -1 {
			data.BrowserVersion = append(data.BrowserVersion, unknown...)
			return true
		}

		data.BrowserVersion = append(data.BrowserVersion, useragent[idx+i+rvIdx:idx+i+rvIdx+rvEndIdx]...)

		return true
	}

	if i := bytes.Index(useragent[start:], strKHTMLWebLight); i != -1 {
		data.BrowserName = append(data.BrowserName, brwWebLight...)
		data.BrowserVersion = append(data.BrowserVersion, unknown...)
		return true
	}

	start = idx

	// Ищем вхождения строк
	for _, item := range sigsf {
		if i := bytes.Index(useragent[start:], item.s); i != -1 {
			return item.f(useragent, item.s, start, i, data)
		}
	}

	// не нашли ни одного вхождения. Это может быть, если это iOS/iPad webview,
	// либо что-то неизвестное

	data.BrowserVersion = append(data.BrowserVersion, unknown...)

	if bytes.Equal(data.PlatformName, pltIOS) {
		data.BrowserName = append(data.BrowserName, brwIOSWebView...)
		data.IsWebView = true
		return true
	}
	if bytes.Equal(data.PlatformName, pltIPadOS) {
		data.BrowserName = append(data.BrowserName, brwIPadOSWebView...)
		data.IsWebView = true
		return true
	}

	//if bytes.Equal(data.PlatformName, pltIOS) ||
	//	bytes.Equal(data.PlatformName, pltIPadOS) ||
	//	bytes.Equal(data.PlatformName, pltMacOS) {
	//
	//
	//
	//	if data.IsMobile {
	//		data.BrowserName = append(data.BrowserName, brwMobileSafari...)
	//	} else {
	//		data.BrowserName = append(data.BrowserName, brwSafari...)
	//	}
	//	data.BrowserVersion = append(data.BrowserVersion, unknown...)
	//	return true
	//}

	return false
}

// берем версию до тех пор, пока у нас цифра или точка (ну или до конца строки)
// обрезаем версию, если встречается третья точка
// 10.20.30.40... -> 10.20.30
func parseDottedVersion(useragent []byte, dest []byte) []byte {
	var idx int
	var dots int
	for {
		if len(useragent) <= idx {
			break
		}
		if useragent[idx] != '1' &&
			useragent[idx] != '2' &&
			useragent[idx] != '3' &&
			useragent[idx] != '4' &&
			useragent[idx] != '5' &&
			useragent[idx] != '6' &&
			useragent[idx] != '7' &&
			useragent[idx] != '8' &&
			useragent[idx] != '9' &&
			useragent[idx] != '0' &&
			useragent[idx] != '.' {

			if len(dest) > 0 && dest[idx-1] == '.' {
				dest = append(dest, '0')
			}

			break
		}
		if useragent[idx] == '.' {
			dots++
		}
		if dots > 2 {
			break
		}
		dest = append(dest, useragent[idx])
		idx++
	}

	return dest
}

// старые версии Windows с Internet Explorer
func parseCompatible(useragent []byte, data *Data) bool {
	data.PlatformName = append(data.PlatformName, pltWindows...)
	data.BrowserName = append(data.BrowserName, brwInternetExplorer...)

	if !bytes.HasPrefix(useragent, strMsie) {
		return false
	}
	start := 5 // 'MSIE '

	// Ищем, где заканчивается версия MSIE. Это будет или ; или ,
	idxDC := bytes.IndexByte(useragent[start:], ';')
	idxC := bytes.IndexByte(useragent[start:], ',')
	if idxDC == -1 && idxC == -1 {
		return false
	}
	var idx int
	if idxDC == -1 {
		idx = idxC
	} else if idxC == -1 {
		idx = idxDC
	} else {
		idx = int(math.Min(float64(idxDC), float64(idxC)))
	}

	data.BrowserVersion = append(data.BrowserVersion, useragent[start:start+idx]...)
	start += idx + 2 // '; '
	if len(useragent) <= start {
		return false
	}
	//if !bytes.HasPrefix(useragent[start:], strWindowsNT) {
	//	return false
	//}
	start += len(strWindowsNT)
	idx = start
	for {
		if idx >= len(useragent) {
			return false
		}
		if useragent[idx] == ';' || useragent[idx] == ')' {
			break
		}
		idx++
	}
	data.PlatformVersion = append(data.PlatformVersion, replaceWindowsVersion(useragent[start:idx])...)
	return true

}

func replaceWindowsVersion(v []byte) []byte {
	switch string(v) {
	case "10.0":
		return pltWin100
	case "6.3":
		return pltWin81
	case "6.2":
		return pltWin8
	case "6.1":
		return pltWin7
	case "6.0":
		return pltWinVista
	case "5.1":
		return pltWinXP
	case "5.0":
		return pltWin2000
	}
	return unknown
}

func replaceUnderscoreToDot(data []byte) {
	for i, b := range data {
		if b == '_' {
			data[i] = '.'
		}
	}
}

func fWrapWebView(browserName []byte) sigFunc {
	return func(useragent, searchStr []byte, start, idx int, data *Data) bool {
		data.BrowserName = append(data.BrowserName, browserName...)
		data.IsWebView = true

		start += idx + len(searchStr)

		if len(useragent) <= start {
			data.BrowserVersion = append(data.BrowserVersion, unknown...)
			return true
		}

		data.BrowserVersion = parseDottedVersion(useragent[start:], data.BrowserVersion)

		return true
	}
}

func fWrap(browserName []byte) sigFunc {
	return func(useragent, searchStr []byte, start, idx int, data *Data) bool {
		data.BrowserName = append(data.BrowserName, browserName...)

		start += idx + len(searchStr)

		if len(useragent) <= start {
			data.BrowserVersion = append(data.BrowserVersion, unknown...)
			return true
		}

		data.BrowserVersion = parseDottedVersion(useragent[start:], data.BrowserVersion)

		return true
	}
}

func fOpera(useragent, searchStr []byte, start, idx int, data *Data) bool {
	if data.IsMobile {
		data.BrowserName = append(data.BrowserName, brwOperaMobile...)
	} else {
		data.BrowserName = append(data.BrowserName, brwOperaDesktop...)
	}

	if len(useragent) <= start+idx+len(searchStr) {
		data.BrowserVersion = append(data.BrowserVersion, unknown...)
		return true
	}

	data.BrowserVersion = parseDottedVersion(useragent[start+idx+len(searchStr):], data.BrowserVersion)

	return true
}

func fFirefox(useragent, searchStr []byte, start, idx int, data *Data) bool {
	if data.IsMobile {
		data.BrowserName = append(data.BrowserName, brwFirefoxMobile...)
		if bytes.Equal(data.PlatformName, pltIOS) {
			data.BrowserName = append(data.BrowserName[:0], brwFirefoxIOS...)
		}
	} else {
		data.BrowserName = append(data.BrowserName, brwFirefox...)
	}

	i := start + idx + len(searchStr)

	if len(useragent) <= i {
		data.BrowserVersion = append(data.BrowserVersion, unknown...)
		return true
	}

	data.BrowserVersion = parseDottedVersion(useragent[i:], data.BrowserVersion)

	return true
}

func fChrome(useragent, searchStr []byte, start, idx int, data *Data) bool {
	if data.IsMobile {
		data.BrowserName = append(data.BrowserName, brwChromeMobile...)
		if bytes.Equal(data.PlatformName, pltIOS) {
			data.BrowserName = append(data.BrowserName[:0], brwChromeIOS...)
		}
	} else {
		data.BrowserName = append(data.BrowserName, brwChrome...)
	}

	if len(useragent) <= start+idx+len(searchStr) {
		data.BrowserVersion = append(data.BrowserVersion, unknown...)
		return true
	}

	data.BrowserVersion = parseDottedVersion(useragent[start+idx+len(searchStr):], data.BrowserVersion)

	if bytes.Equal(data.PlatformName, pltChromeOS) {
		data.PlatformVersion = append(data.PlatformVersion[:0], data.BrowserVersion...)
	}

	return true
}

func fSmartTV(_, _ []byte, _, _ int, data *Data) bool {
	data.BrowserName = append(data.BrowserName, brwSmartTV...)
	data.BrowserVersion = append(data.BrowserVersion, unknown...)
	return true
}

func fUCBrowser(useragent, searchStr []byte, start, idx int, data *Data) bool {
	if bytes.Equal(data.PlatformName, pltMacOS) {
		data.BrowserName = append(data.BrowserName, brwUCMacOS...)
	} else if bytes.Equal(data.PlatformName, pltIOS) {
		data.BrowserName = append(data.BrowserName, brwUCIPhone...)
	} else if bytes.Equal(data.PlatformName, pltIPadOS) {
		data.BrowserName = append(data.BrowserName, brwUCIPhone...)
	} else if bytes.Equal(data.PlatformName, pltAndroid) {
		data.BrowserName = append(data.BrowserName, brwUCAndroid...)
	} else {
		data.BrowserName = append(data.BrowserName, brwUC...)
	}

	i := start + idx + len(searchStr)

	if len(useragent) <= i {
		data.BrowserVersion = append(data.BrowserVersion, unknown...)
		return true
	}

	data.BrowserVersion = parseDottedVersion(useragent[start+idx+len(searchStr):], data.BrowserVersion)

	return true
}

func fLinkedIn(_, _ []byte, _, _ int, data *Data) bool {
	data.BrowserName = append(data.BrowserName, brwLinkedInApp...)
	data.BrowserVersion = append(data.BrowserVersion, unknown...)
	data.IsWebView = true

	return true
}

func fVersion(useragent, searchStr []byte, start, idx int, data *Data) bool {
	if bytes.Equal(data.PlatformName, pltMacOS) {
		data.BrowserName = append(data.BrowserName, brwSafari...)
	} else if bytes.Equal(data.PlatformName, pltIOS) {
		data.BrowserName = append(data.BrowserName, brwMobileSafari...)
	} else if bytes.Equal(data.PlatformName, pltIPadOS) {
		data.BrowserName = append(data.BrowserName, brwMobileSafari...)
	} else {
		data.BrowserName = append(data.BrowserName, unknown...)
	}

	i := start + idx + len(searchStr)

	if len(useragent) <= i {
		data.BrowserVersion = append(data.BrowserVersion, unknown...)
		return true
	}

	data.BrowserVersion = parseDottedVersion(useragent[start+idx+len(searchStr):], data.BrowserVersion)

	return true
}

func fIron(useragent, _ []byte, start, _ int, data *Data) bool {
	// После строки Iron может быть версия (тогда берем ее), либо не быть -тогда ищем Chrome/ и берем ее версию
	// Mozilla/5.0 (Windows; U; Windows NT 5.2; en-US) AppleWebKit/532.0 (KHTML, like Gecko) Chrome/0.0.0 Safari/532.0 Iron/3.0.197.0
	// Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3750.0 Iron Safari/537.36
	data.BrowserName = append(data.BrowserName, brwIron...)
	// Если после Iron есть слеш, то идем далее, стандартный путь по определению версии
	//if len(useragent) <= start+idx+4 && useragent[start+idx+5] == '/' {
	//	break
	//}

	if i := bytes.Index(useragent[start:], strChromeSlash); i > -1 {
		if len(useragent) > start+i+7 { // 7  = Chrome/
			data.BrowserVersion = parseDottedVersion(useragent[start+i+7:], data.BrowserVersion)
			return true
		}
	}

	data.BrowserVersion = append(data.BrowserVersion, unknown...)
	return true
}

func fYaBrowser(useragent, searchStr []byte, start, idx int, data *Data) bool {
	if bytes.Equal(data.PlatformName, pltWindows) {
		data.BrowserName = append(data.BrowserName, brwYaBrowserForWindows...)
	} else if bytes.Equal(data.PlatformName, pltAndroid) {
		data.BrowserName = append(data.BrowserName, brwYaBrowserForAndroid...)
	} else if bytes.Equal(data.PlatformName, pltIOS) {
		data.BrowserName = append(data.BrowserName, brwYaBrowserForIOS...)
	} else if bytes.Equal(data.PlatformName, pltIPadOS) {
		data.BrowserName = append(data.BrowserName, brwYaBrowserForIOS...)
	} else if bytes.Equal(data.PlatformName, pltMacOS) {
		data.BrowserName = append(data.BrowserName, brwYaBrowserForMacOS...)
	} else {
		data.BrowserName = append(data.BrowserName, brwYaBrowser...)
	}

	if len(useragent) <= start+idx+len(searchStr) {
		data.BrowserVersion = append(data.BrowserVersion, unknown...)
		return true
	}

	data.BrowserVersion = parseDottedVersion(useragent[start+idx+len(searchStr):], data.BrowserVersion)

	return true
}

func fPuffin(useragent, searchStr []byte, start, idx int, data *Data) bool {
	if bytes.Equal(data.PlatformName, pltAndroid) {
		data.BrowserName = append(data.BrowserName, brwPuffinAndroid...)
	} else if bytes.Equal(data.PlatformName, pltIOS) {
		data.BrowserName = append(data.BrowserName, brwPuffinIOS...)
	} else if bytes.Equal(data.PlatformName, pltIPadOS) {
		data.BrowserName = append(data.BrowserName, brwPuffinIOS...)
	} else {
		data.BrowserName = append(data.BrowserName, brwPuffin...)
	}

	if len(useragent) <= start+idx+len(searchStr) {
		data.BrowserVersion = append(data.BrowserVersion, unknown...)
		return true
	}

	data.BrowserVersion = parseDottedVersion(useragent[start+idx+len(searchStr):], data.BrowserVersion)

	return true
}

func fNaver(useragent, searchStr []byte, start, idx int, data *Data) bool {
	// Mozilla/5.0 (iPhone; CPU iPhone OS 14_4_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/605.1 NAVER(inapp; search; 900; 11.2.1; XR)

	data.BrowserName = append(data.BrowserName, brwNaverApp...)
	data.IsWebView = true

	i := start + idx + len(searchStr) + 15 // 14 = len(inapp; search;), +1 для установки позиции ЗА ;

	j := bytes.IndexByte(useragent[i:], ';')
	if j == -1 {
		data.BrowserVersion = append(data.BrowserVersion, unknown...)
		return true
	}

	data.BrowserVersion = parseDottedVersion(useragent[i+j+2:], data.BrowserVersion) // 2 = len(; )

	return true
}

func fTOBrowser(useragent, searchStr []byte, start, idx int, data *Data) bool {
	data.BrowserName = append(data.BrowserName, brwTODEBrowser...)
	// В случае tod браузера, версия выглядит так: TO-Browser/TOB7.85.0.301_01
	// то есть, для парсинга версии помимо слеша надо пропустить еще 3 символа: /TOB
	i := start + idx + len(searchStr) + 3
	if len(useragent) <= i {
		data.BrowserVersion = append(data.BrowserVersion, unknown...)
		return true
	}

	data.BrowserVersion = parseDottedVersion(useragent[i:], data.BrowserVersion)
	return true

}
