package useragent

type sigFunc func(useragent, searchStr []byte, start, idx int, data *Data) bool

var (
	unknown = []byte("Unknown")

	mozillaPrefix = []byte("Mozilla/5.0 (")

	pltWindows  = []byte("Windows")
	pltLinux    = []byte("Linux")
	pltIOS      = []byte("iOS")
	pltIPadOS   = []byte("iPadOS")
	pltMacOS    = []byte("macOS")
	pltChromeOS = []byte("Chrome OS")

	strMsie            = []byte("MSIE ")
	strWindowsNT       = []byte("Windows NT ")
	strX11             = []byte("X11")
	strLinux           = []byte("Linux")
	strAndroid         = []byte("Android")
	strIPhone          = []byte("iPhone")
	strLikeMacOSX      = []byte(" like Mac OS X")
	strSpaceOSSpace    = []byte(" OS ")
	strIPad            = []byte("iPad")
	strIPod            = []byte("iPod")
	strIPodTouch       = []byte("iPod touch")
	strMacintosh       = []byte("Macintosh")
	strKHTMLWebLight   = []byte("(KHTML, like Gecko; googleweblight)")
	strCompatible      = []byte("compatible")
	strCompatibleComma = []byte("compatible,")
	strChromeSlash     = []byte("Chrome/")
	strSmartTV         = []byte("SmartTV")

	blockAfterPlatformU     = []byte("U;")
	blockAfterPlatformArm   = []byte("arm;")
	blockAfterPlatformArm64 = []byte("arm_64;")

	pltWin100   = []byte("10.0")
	pltWin81    = []byte("8.1")
	pltWin8     = []byte("8")
	pltWin7     = []byte("7")
	pltWinVista = []byte("Vista")
	pltWinXP    = []byte("XP")
	pltWin2000  = []byte("2000")
	pltAndroid  = []byte("Android")

	intelMacOSXLen = len("Intel Mac OS X ")

	sigsf = []struct {
		s []byte
		f sigFunc
	}{
		{[]byte("SmartTV"), fSmartTV},
		{[]byte("OPR/"), fOpera},
		{[]byte("OPT/"), fWrap(brwOperaTouch)},
		{[]byte("UCBrowser/"), fUCBrowser},
		//{[]byte("UCBrowser/"), fWrap(brwUC)},
		{[]byte("MiuiBrowser/"), fWrap(brwMIUI)},
		{[]byte("GSA/"), fWrap(brwGsa)},
		{[]byte("HeyTapBrowser/"), fWrap(brwHeyTap)},
		{[]byte("AlohaBrowser/"), fWrap(brwAloha)},
		//{[]byte("YaBrowser/"), fWrap(brwYaBrowser)},
		{[]byte("YaApp_Android/"), fWrap(brwYsa)},
		{[]byte("YaApp_iOS/"), fWrap(brwYsa)},
		{[]byte("YaBrowser/"), fYaBrowser},
		{[]byte("SamsungBrowser/"), fWrap(brwSamsungBrowser)},
		{[]byte("RakutenWebSearch/"), fWrapWebView(brwRakutenWebsearchApp)},
		{[]byte("CriOS/"), fWrap(brwChromeMobile)},
		{[]byte("TO-Browser/"), fTOBrowser},
		{[]byte("Silk/"), fWrap(brwSilk)},
		{[]byte("OppoBrowser/"), fWrap(brwOppo)},
		{[]byte("Edg/"), fWrap(brwEdge)},
		{[]byte("DuckDuckGo/"), fWrap(brwDuckDuckGo)},
		{[]byte("Mint Browser/"), fWrap(brwMint)},
		{[]byte("coc_coc_browser/"), fWrap(brwCocCoc)},
		{[]byte("QQBrowser/"), fWrap(brwQQ)},
		{[]byte("SogouMobileBrowser/"), fWrap(brwSogouMobile)},
		//{[]byte("FxiOS/"), fWrap(brwFirefoxMobile)},
		{[]byte("FxiOS/"), fFirefox},
		{[]byte("Amigo/"), fWrap(brwAmigo)},
		{[]byte("VivoBrowser/"), fWrap(brwVivo)},
		{[]byte("Edge/"), fWrap(brwEdge)},
		{[]byte("UBrowser/"), fWrap(brwU)},
		{[]byte("CoolNovo/"), fWrap(brwCoolNovo)},
		{[]byte("Iron"), fIron},
		{[]byte("Freeu/"), fWrap(brwFreeU)},
		{[]byte("Vivaldi/"), fWrap(brwVivaldi)},
		//{[]byte("Puffin/"), fWrap(brwPuffin)},
		{[]byte("Puffin/"), fPuffin},
		{[]byte("Elements Browser/"), fWrap(brwElements)},
		{[]byte("PaleMoon/"), fWrap(brwPaleMoon)},
		{[]byte("Atom/"), fWrap(brwAtom)},
		{[]byte("Photon/"), fWrap(brwPhoton)},
		{[]byte("GoBrowser/"), fWrap(brwGoBrowser)},
		{[]byte("Maxthon/"), fWrap(brwMaxthon)},
		{[]byte("LiteBrowser/"), fWrap(brwLiteBrowser)},
		{[]byte("Waterfox/"), fWrap(brwWaterfox)},
		{[]byte("SznProhlizec/"), fWrap(brwSznProhlizec)},
		{[]byte("Sleipnir/"), fWrap(brwSleipnir)},
		{[]byte("2345Explorer/"), fWrap(brw2345Explorer)},
		{[]byte("Avast/"), fWrap(brwAvast)},
		{[]byte("Iceweasel/"), fWrap(brwIceweasel)},
		{[]byte("Whale/"), fWrap(brwWhale)},
		{[]byte("115Browser/"), fWrap(brw115Browser)},
		{[]byte("baiduboxapp/"), fWrap(brwBaiduBoxApp)},
		{[]byte("MicroMessenger/"), fWrap(brwWeChatIOS)},
		{[]byte("Instagram "), fWrapWebView(brwInstagram)},
		{[]byte("FBSV/"), fWrapWebView(brwFacebookApp)},
		{[]byte("NAVER("), fNaver},
		{[]byte("[LinkedInApp]"), fLinkedIn},
		{[]byte("TopBuzz technology.mainspring.Babe/"), fWrapWebView(brwTopBuzzApp)},
		{[]byte("jp.co.yahoo.ipn.appli/"), fWrapWebView(brwYGApp)},
		{[]byte("KodeiOS/"), fWrapWebView(brwKodeApp)},
		{[]byte("Firefox/"), fFirefox},
		{[]byte("Chrome/"), fChrome},
		{[]byte("Version/"), fVersion},
	}

	brwFirefox       = []byte("Firefox")
	brwFirefoxIOS    = []byte("Firefox for iOS")
	brwFirefoxMobile = []byte("Firefox for Mobile")
	brwChrome        = []byte("Chrome")
	brwChromeIOS     = []byte("Chrome for iOS")
	brwChromeMobile  = []byte("Chrome Mobile")
	//brwOpera         = []byte("Opera")
	brwOperaTouch   = []byte("Opera Touch")
	brwOperaMobile  = []byte("Opera Mobile")
	brwOperaDesktop = []byte("Opera Desktop")
	brwSafari       = []byte("Safari")
	brwMobileSafari = []byte("Mobile Safari")

	brwPuffin        = []byte("Puffin")
	brwPuffinIOS     = []byte("Puffin for iOS")
	brwPuffinAndroid = []byte("Puffin for Android")

	brwEdge                = []byte("Edge")
	brwSilk                = []byte("Amazon Silk")
	brwTODEBrowser         = []byte("t-online.de Browser")
	brwYaBrowser           = []byte("Yandex.Browser")
	brwYaBrowserForIOS     = []byte("Yandex.Browser for iOS")
	brwYaBrowserForWindows = []byte("Yandex.Browser for Windows")
	brwYaBrowserForAndroid = []byte("Yandex.Browser for Android")
	brwYaBrowserForMacOS   = []byte("Yandex.Browser for MacOS")
	brwAloha               = []byte("Aloha Browser")
	brwSamsungBrowser      = []byte("Samsung Browser")
	brwUC                  = []byte("UCBrowser")
	brwUCAndroid           = []byte("UCBrowser for Android")
	brwUCMacOS             = []byte("UCBrowser for MacOS")
	brwUCIPhone            = []byte("UCBrowser for iPhone")
	brwMIUI                = []byte("MIUI Browser")
	brwGsa                 = []byte("Google Search App")
	brwHeyTap              = []byte("HeyTapBrowser")
	brwGoHTTPClient        = []byte("Go-http-client")
	brwAndroid             = []byte("Unknown Android App")
	brwWebLight            = []byte("Web Light")
	brwDuckDuckGo          = []byte("DuckDuckGo")
	brwMint                = []byte("Mint Browser")
	brwCocCoc              = []byte("Coc Coc")
	brwQQ                  = []byte("QQ Browser")
	brwInternetExplorer    = []byte("Internet Explorer")
	brwYsa                 = []byte("Yandex Search App")
	brwOppo                = []byte("Oppo Browser")
	brwSogouMobile         = []byte("Sogou Browser")
	brwAmigo               = []byte("Amigo Browser")
	brwVivo                = []byte("Vivo Browser")
	brwU                   = []byte("U Browser")
	brwCoolNovo            = []byte("Cool Novo Browser")
	brwIron                = []byte("SRWare Iron")
	brwFreeU               = []byte("FreeU")
	brwVivaldi             = []byte("Vivaldi")
	brwElements            = []byte("Elements Browser")
	brwSmartTV             = []byte("SmartTV")
	brwPaleMoon            = []byte("PaleMoon")
	brwAtom                = []byte("Atom")
	brwPhoton              = []byte("Photon")
	brwGoBrowser           = []byte("GoBrowser")
	brwMaxthon             = []byte("Maxthon")
	brwLiteBrowser         = []byte("LiteBrowser")
	brwWaterfox            = []byte("Waterfox")
	brwSznProhlizec        = []byte("Seznam Browser")
	brwSleipnir            = []byte("Sleipnir Browser")
	brw2345Explorer        = []byte("2345 Explorer")
	brwAvast               = []byte("Avast Browser")
	brwIceweasel           = []byte("Iceweasel Browser")
	brwWhale               = []byte("Whale Browser")
	brw115Browser          = []byte("115 Browser")
	brwBaiduBoxApp         = []byte("Baidu Box App")
	brwWeChatIOS           = []byte("WeChat for iOS")
	brwInstagram           = []byte("Instagram")
	brwFacebookApp         = []byte("Facebook App")
	brwNaverApp            = []byte("NAVER App")
	brwTopBuzzApp          = []byte("TopBuzz App")
	brwLinkedInApp         = []byte("LinkedIn App")
	brwYGApp               = []byte("YJ App")
	brwRakutenWebsearchApp = []byte("RakutenWebsearch App")
	brwIOSWebView          = []byte("iOS WebView")
	brwIPadOSWebView       = []byte("iPadOS WebView")
	brwKodeApp             = []byte("Kode App")
)
