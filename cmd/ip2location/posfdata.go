package main

//
//type posfData struct {
//	ip  comperestring
//	mss comperestring
//}
//
//func getDataPosf(fileName comperestring) []posfData {
//	file, errOpen := ioutil.ReadFile(fileName)
//	if errOpen != nil {
//		fmt.Println("error open file stamp data", errOpen.Error())
//		return nil
//	}
//
//	dataFile := strings.Split(comperestring(file), "\n")
//	//fmt.Println("data posf", len(dataFile))
//
//	data := make([]posfData, 0, 100000)
//	for _, lineFile := range dataFile {
//		lf := strings.Split(lineFile, `";"`)
//		if len(lf) != 9 {
//			continue
//		}
//		ipPort := strings.Trim(lf[1], " \"\n\t")
//		stamp := strings.Trim(lf[8], " \"\n\t")
//
//		if ipPort == "" || stamp == "" {
//			continue
//		}
//
//		ip := getIP(ipPort)
//		mss := getMSS(stamp)
//
//		data = append(data, posfData{ip, mss})
//	}
//
//	return data
//}
//
//func getIP(ds comperestring) comperestring {
//	d := strings.Split(ds, ":")
//	if len(d) == 2 {
//		return d[0]
//	}
//	return ""
//}
//
//func getMSS(ds comperestring) comperestring {
//	d := strings.Split(ds, ";")
//	if len(d) != 6 {
//		return ""
//	}
//	m := d[5]
//	ms := strings.Split(m, ",")
//	if len(ms) < 1 {
//		return ""
//	}
//	mss := ms[0]
//	return mss[1:]
//}
