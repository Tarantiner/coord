package loc

import (
	"fmt"
	"math"
)

// 提供不同坐标系转换方法
// 参考网站：https://artoria2e5.github.io/PRCoords/demo?lat=39&lon=114.429444

var x, y float64
const gcjA = 6378245
const gcjEE = 0.00669342162296594323 // f = 1/298.3; e^2 = 2*f - f**2
const BdDLAT = 0.0060
const BdDLON = 0.0065

type POS struct {
	LON, LAT float64
}

func Wgs2Gcj(pos POS) POS {
	x = pos.LON - 105
	y = pos.LAT - 35
	dlatM := -100 + 2*x + 3*y + 0.2*y*y + 0.1*x*y + 0.2*math.Sqrt(math.Abs(x)) + (
		2*math.Sin(x*6*math.Pi) + 2*math.Sin(x*2*math.Pi) +
			2*math.Sin(y*math.Pi) + 4*math.Sin(y/3*math.Pi) +
			16*math.Sin(y/12*math.Pi) + 32*math.Sin(y/30*math.Pi))*20/3

	dlonM := 300 + x + 2*y + 0.1*x*x + 0.1*x*y + 0.1*math.Sqrt(math.Abs(x)) + (
		2*math.Sin(x*6*math.Pi) + 2*math.Sin(x*2*math.Pi) +
			2*math.Sin(x*math.Pi) + 4*math.Sin(x/3*math.Pi) +
			15*math.Sin(x/12*math.Pi) + 30*math.Sin(x/30*math.Pi))*20/3

	radLAT := pos.LAT / 180 * math.Pi
	magic := 1 - gcjEE*math.Pow(math.Sin(radLAT), 2.0)
	latDegArclen := (math.Pi / 180) * (gcjA * (1 - gcjEE)) / math.Pow(magic, 1.5)
	lonDegArclen := (math.Pi / 180) * (gcjA * math.Cos(radLAT) / math.Sqrt(magic))
	return POS{
		LON: pos.LON + (dlonM / lonDegArclen),
		LAT: pos.LAT + (dlatM / latDegArclen),
	}
}

func Gcj2BD(pos POS) POS {
	x = pos.LON
	y = pos.LAT
	r := math.Sqrt(x*x+y*y) + 0.00002*math.Sin(y*math.Pi*3000/180)
	t := math.Atan2(y, x) + 0.000003*math.Cos(x*math.Pi*3000/180)
	return POS{
		LON: r*math.Cos(t) + BdDLON,
		LAT: r*math.Sin(t) + BdDLAT,
	}
}

func outOfChina(pos POS) bool {
	return !(pos.LON > 73.66 && pos.LON < 135.05 && pos.LAT > 3.86 && pos.LAT < 53.55)
}

func Wgs2BD(pos POS) POS {
	// wgs坐标系转百度坐标系
	if outOfChina(pos){
		fmt.Println("中国以外")
		return pos
	}
	pos = Wgs2Gcj(pos)
	return Gcj2BD(pos)
}

func BD2Gcj(pos POS) POS {
	x = pos.LON - BdDLON
	y = pos.LAT - BdDLAT
	r := math.Sqrt(x*x+y*y) - 0.00002*math.Sin(y*math.Pi*3000/180)
	t := math.Atan2(y, x) - 0.000003*math.Cos(x*math.Pi*3000/180)
	return POS{
		LON: r*math.Cos(t),
		LAT: r*math.Sin(t),
	}
}

func CoordDiff(pos1, pos2 POS) POS {
	return POS{
		LON: pos1.LON - pos2.LON,
		LAT: pos1.LAT - pos2.LAT,
	}
}


func Gcj2Wgs(pos POS) POS {
	return CoordDiff(pos, CoordDiff(Wgs2Gcj(pos), pos))
}

func BD2Wgs(pos POS) POS {
	if outOfChina(pos){
		fmt.Println("中国以外")
		return pos
	}
	pos = BD2Gcj(pos)
	return Gcj2Wgs(pos)
}

func TestPos() {
	lon := 114.429444
	lat := 39.0
	fmt.Println("正在测试坐标转换，114.429444|39.0")
	fmt.Println("/geo/api/&type=wgs2gcj", Wgs2Gcj(POS{LON: lon, LAT: lat}))
	fmt.Println("/geo/api/&type=wgs2bd", Wgs2BD(POS{LON: lon, LAT: lat}))
	fmt.Println("/geo/api/&type=gcj2wgs", Gcj2Wgs(POS{LON: lon, LAT: lat}))
	fmt.Println("/geo/api/&type=gcj2bd", Gcj2BD(POS{LON: lon, LAT: lat}))
	fmt.Println("/geo/api/&type=bd2wgs", BD2Wgs(POS{LON: lon, LAT: lat}))
	fmt.Println("/geo/api/&type=bd2gcj", BD2Gcj(POS{LON: lon, LAT: lat}))
}
