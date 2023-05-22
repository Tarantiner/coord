package main

import (
	"coord/loc"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strconv"
)

var port string

func init() {
	flag.StringVar(&port, "port", "8000", "web监听端口")
	flag.Parse()
}

type MsgBox struct {
	Type    string  `json:"type"`
	Success int     `json:"success"`
	ERR     string  `json:"error"`
	LON     float64 `json:"longitude"`
	LAT     float64 `json:"latitude"`
}

func HandleResult(mg *MsgBox, pos loc.POS, w http.ResponseWriter) {
	mg.Success = 1
	mg.LON, _ = strconv.ParseFloat(fmt.Sprintf("%.6f", pos.LON), 64)
	mg.LAT, _ = strconv.ParseFloat(fmt.Sprintf("%.6f", pos.LAT), 64)
	b, _ := json.Marshal(mg)
	fmt.Println(string(b))
	w.Write(b)
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
	mg.Type = tp
	lonStr := qMap.Get("lon")
	latStr := qMap.Get("lat")
	lon, err := strconv.ParseFloat(lonStr, 64)
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

	if -180.0 <= lon && lon <= 180.0 && -90.0 <= lat && lat <= 90.0 {
	} else {
		mg.ERR = "请提供有效经纬度"
		b, _ := json.Marshal(mg)
		fmt.Println(string(b), mg)
		w.Write(b)
		return
	}

	if tp == "wgs2gcj" {
		pos := loc.Wgs2Gcj(loc.POS{LON: lon, LAT: lat})
		HandleResult(&mg, pos, w)
		return
	}

	if tp == "wgs2bd" {
		pos := loc.Wgs2BD(loc.POS{LON: lon, LAT: lat})
		HandleResult(&mg, pos, w)
		return
	}

	if tp == "gcj2bd" {
		pos := loc.Gcj2BD(loc.POS{LON: lon, LAT: lat})
		HandleResult(&mg, pos, w)
		return
	}

	if tp == "bd2gcj" {
		pos := loc.BD2Gcj(loc.POS{LON: lon, LAT: lat})
		HandleResult(&mg, pos, w)
		return
	}

	if tp == "gcj2wgs" {
		pos := loc.Gcj2Wgs(loc.POS{LON: lon, LAT: lat})
		HandleResult(&mg, pos, w)
		return
	}

	if tp == "bd2wgs" {
		pos := loc.BD2Wgs(loc.POS{LON: lon, LAT: lat})
		HandleResult(&mg, pos, w)
		return
	}

	mg.ERR = "请有效转换类型"
	b, _ := json.Marshal(mg)
	fmt.Println(string(b), mg)
	w.Write(b)
	return
}

func main() {
	loc.TestPos()
	http.HandleFunc("/geo/api", GeoAPI)
	http.ListenAndServe(":"+port, nil)
}
