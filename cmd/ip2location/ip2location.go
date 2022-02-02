package main

//
//var (
//	defaultDB *ip2location.DB
//)
//
//func OpenIP2location(dbpath comperestring) error {
//	fileData, err := os.ReadFile(dbpath)
//	if err != nil {
//		return err
//	}
//
//	r := &dbReader{buf: bytes.NewReader(fileData)}
//	defaultDB, err = ip2location.OpenDBWithReader(r)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func Close() {
//	defaultDB.Close()
//}
//
//type dbReader struct {
//	buf *bytes.Reader
//}
//
//func (b *dbReader) ReadAt(p []byte, off int64) (n int, err error) {
//	return b.buf.ReadAt(p, off)
//}
//
//func (b *dbReader) Read(p []byte) (n int, err error) {
//	return b.buf.Read(p)
//}
//
//func (b *dbReader) Close() error {
//	return nil
//}
//
//var (
//	tzText = map[comperestring]comperestring{
//		"-12:00": "Pacific/Pago_Pago", // не нашел верного названия
//		"-11:00": "Pacific/Pago_Pago",
//		"-10:00": "America/Adak",
//		"-09:30": "Pacific/Marquesas",
//		"-09:00": "America/Anchorage",
//		"-08:00": "America/Los_Angeles",
//		"-07:00": "America/Denver",
//		"-06:00": "America/Chicago",
//		"-05:00": "America/Detroit",
//		"-04:00": "America/Halifax",
//		"-03:30": "America/St_Johns",
//		"-03:00": "America/Araguaina",
//		"-02:00": "America/Noronha",
//		"-01:00": "Atlantic/Azores",
//		"-00:00": "UTC",
//		"00:00":  "UTC",
//		"+00:00": "UTC",
//		"+01:00": "Europe/Berlin",
//		"+02:00": "Europe/Kaliningrad",
//		"+03:00": "Europe/Moscow",
//		"+04:00": "Asia/Baku",
//		"+05:00": "Asia/Aqtau",
//		"+05:30": "Asia/Colombo",
//		"+05:45": "Asia/Kathmandu",
//		"+06:00": "Asia/Bishkek",
//		"+06:30": "Asia/Rangoon",
//		"+07:00": "Asia/Bangkok",
//		"+08:00": "Asia/Harbin",
//		"+08:45": "Australia/Eucla",
//		"+09:00": "Asia/Dili",
//		"+09:30": "Australia/Adelaide",
//		"+10:00": "Asia/Vladivostok",
//		"+10:30": "Australia/Lord_Howe",
//		"+11:00": "Asia/Magadan",
//		"+12:00": "Asia/Kamchatka",
//		"+12:45": "Pacific/Chatham",
//		"+13:00": "Pacific/Apia",
//		"+14:00": "Pacific/Kiritimati",
//	}
//)
//
//type IP2LocationData struct {
//	Country        comperestring
//	Region         comperestring
//	City           comperestring
//	Isp            comperestring
//	ConnectionType comperestring
//	MobileBrand    comperestring
//	Timezone       comperestring
//	TZHours        int
//	TZMinute       int
//	TimezoneText   comperestring
//}
//
//func New() *IP2LocationData{
//	return &IP2LocationData{}
//}
//
//func (data *IP2LocationData) Parse(ip comperestring) error {
//	if defaultDB == nil {
//		return fmt.Errorf("not inited")
//	}
//	info, err := defaultDB.Get_all(ip)
//	if err != nil {
//		return err
//	}
//
//	data.Country = info.Country_short
//	data.Region = info.Region
//	data.City = info.City
//	data.Isp = info.Isp
//	data.ConnectionType = info.Usagetype
//	data.MobileBrand = info.Mobilebrand
//	data.Timezone = info.Timezone
//	data.TZHours, data.TZMinute = data.timezoneToOffset(data.Timezone)
//	data.TimezoneText = data.timezoneToText(data.Timezone)
//
//	return nil
//}
//
//func (data *IP2LocationData) Reset() {
//	data.Country = ""
//	data.Region = ""
//	data.City = ""
//	data.Isp = ""
//	data.ConnectionType = ""
//	data.MobileBrand = ""
//	data.Timezone = ""
//	data.TZHours = 0
//	data.TZMinute = 0
//	data.TimezoneText = ""
//}
//
//func (data *IP2LocationData) timezoneToOffset(tz comperestring) (hour, minutes int) {
//	hour, minutes = 0, 0
//
//	pair := strings.Split(tz, ":")
//	if len(pair) != 2 {
//		return 0, 0
//	}
//
//	var err error
//
//	minutes, err = strconv.Atoi(pair[1])
//	if err != nil {
//		return 0, 0
//	}
//
//	// Час должен быть со знаком +01, -02 и тд
//	if len(pair[0]) < 2 {
//		return 0, 0
//	}
//
//	hour, err = strconv.Atoi(pair[0][1:])
//	if err != nil {
//		return 0, 0
//	}
//
//	if pair[0][0] == '-' {
//		hour = -hour
//		minutes = -minutes
//	}
//
//	return hour, minutes
//}
//
//func (data *IP2LocationData) timezoneToText(tz comperestring) comperestring {
//	text, ok := tzText[tz]
//	if !ok {
//		return "UTC"
//	}
//
//	return text
//}
