package main

import (
	"coord/loc"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type MsgBox struct {
	Success int     `json:"success"`
	ERR     string  `json:"error"`
	LNG     float64 `json:"longitude"`
	LAT     float64 `json:"latitude"`
}

func GeoAPI(w http.ResponseWriter, r *http.Request) {
	qMap := r.URL.Query()
	tp := qMap.Get("type")
	var mg MsgBox
	if tp == "" {
		mg.ERR = "请提供转换类型"
		b, _ := json.Marshal(mg)
		fmt.Println(string(b))
		w.Write(b)
		return
	}
	lngStr := qMap.Get("lng")
	latStr := qMap.Get("lat")
	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		mg.ERR = "请提供有效经纬度"
		b, _ := json.Marshal(mg)
		fmt.Println(string(b), mg)
		w.Write(b)
		return
	}
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		mg.ERR = "请提供有效经纬度"
		b, _ := json.Marshal(mg)
		fmt.Println(string(b), mg)
		w.Write(b)
		return
	}

	if -180.0 <= lng && lng <= 180.0 && -90.0 <= lat && lat <= 90.0 {
	}else{
		mg.ERR = "请提供有效经纬度"
		b, _ := json.Marshal(mg)
		fmt.Println(string(b), mg)
		w.Write(b)
		return
	}

	if tp == "wgs2gcj" {
		pos := loc.Wgs2Gcj(loc.POS{LNG: lng, LAT: lat})
		mg.Success = 1
		mg.LNG = pos.LNG
		mg.LAT = pos.LAT
		b, _ := json.Marshal(mg)
		fmt.Println(string(b))
		w.Write(b)
		return
	}

	if tp == "gcj2bd" {
		pos := loc.Gcj2BD(loc.POS{LNG: lng, LAT: lat})
		mg.Success = 1
		mg.LNG = pos.LNG
		mg.LAT = pos.LAT
		b, _ := json.Marshal(mg)
		fmt.Println(string(b))
		w.Write(b)
		return
	}

	if tp == "bd2gcj" {
		pos := loc.BD2Gcj(loc.POS{LNG: lng, LAT: lat})
		mg.Success = 1
		mg.LNG = pos.LNG
		mg.LAT = pos.LAT
		b, _ := json.Marshal(mg)
		fmt.Println(string(b))
		w.Write(b)
		return
	}

	if tp == "gcj2wgs" {
		pos := loc.Gcj2Wgs(loc.POS{LNG: lng, LAT: lat})
		mg.Success = 1
		mg.LNG = pos.LNG
		mg.LAT = pos.LAT
		b, _ := json.Marshal(mg)
		fmt.Println(string(b))
		w.Write(b)
		return
	}

	if tp == "bd2wgs" {
		pos := loc.BD2Wgs(loc.POS{LNG: lng, LAT: lat})
		mg.Success = 1
		mg.LNG = pos.LNG
		mg.LAT = pos.LAT
		b, _ := json.Marshal(mg)
		fmt.Println(string(b))
		w.Write(b)
		return
	}

	mg.ERR = "请有效转换类型"
	b, _ := json.Marshal(mg)
	fmt.Println(string(b), mg)
	w.Write(b)
	return
}

func main() {
	http.HandleFunc("/geo/api", GeoAPI)
	http.ListenAndServe(":8000", nil)
}
