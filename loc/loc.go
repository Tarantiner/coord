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
	LNG, LAT float64
}

func Wgs2Gcj(pos POS) POS {
	x = pos.LNG - 105
	y = pos.LAT - 35
	dlatM := -100 + 2*x + 3*y + 0.2*y*y + 0.1*x*y + 0.2*math.Sqrt(math.Abs(x)) + (
		2*math.Sin(x*6*math.Pi) + 2*math.Sin(x*2*math.Pi) +
			2*math.Sin(y*math.Pi) + 4*math.Sin(y/3*math.Pi) +
			16*math.Sin(y/12*math.Pi) + 32*math.Sin(y/30*math.Pi))*20/3

	dlngM := 300 + x + 2*y + 0.1*x*x + 0.1*x*y + 0.1*math.Sqrt(math.Abs(x)) + (
		2*math.Sin(x*6*math.Pi) + 2*math.Sin(x*2*math.Pi) +
			2*math.Sin(x*math.Pi) + 4*math.Sin(x/3*math.Pi) +
			15*math.Sin(x/12*math.Pi) + 30*math.Sin(x/30*math.Pi))*20/3

	radLAT := pos.LAT / 180 * math.Pi
	magic := 1 - gcjEE*math.Pow(math.Sin(radLAT), 2.0)
	latDegArclen := (math.Pi / 180) * (gcjA * (1 - gcjEE)) / math.Pow(magic, 1.5)
	lngDegArclen := (math.Pi / 180) * (gcjA * math.Cos(radLAT) / math.Sqrt(magic))
	return POS{
		LNG: pos.LNG + (dlngM / lngDegArclen),
		LAT: pos.LAT + (dlatM / latDegArclen),
	}
}

func Gcj2BD(pos POS) POS {
	x = pos.LNG
	y = pos.LAT
	r := math.Sqrt(x*x+y*y) + 0.00002*math.Sin(y*math.Pi*3000/180)
	t := math.Atan2(y, x) + 0.000003*math.Cos(x*math.Pi*3000/180)
	return POS{
		LNG: r*math.Cos(t) + BdDLON,
		LAT: r*math.Sin(t) + BdDLAT,
	}
}

func Wgs2BD(pos POS) POS {
	// wgs坐标系转百度坐标系
	pos = Wgs2Gcj(pos)
	return Gcj2BD(pos)
}

func BD2Gcj(pos POS) POS {
	x = pos.LNG - BdDLON
	y = pos.LAT - BdDLAT
	r := math.Sqrt(x*x+y*y) + 0.00002*math.Sin(y*math.Pi*3000/180)
	t := math.Atan2(y, x) + 0.000003*math.Cos(x*math.Pi*3000/180)
	return POS{
		LNG: r*math.Cos(t),
		LAT: r*math.Sin(t),
	}
}

func CoordDiff(pos1, pos2 POS) POS {
	return POS{
		LNG: pos1.LNG - pos2.LNG,
		LAT: pos1.LAT - pos2.LAT,
	}
}


func Gcj2Wgs(pos POS) POS {
	return CoordDiff(pos, CoordDiff(Wgs2Gcj(pos), pos))
}

func BD2Wgs(pos POS) POS {
	pos = BD2Gcj(pos)
	return Gcj2Wgs(pos)
}

func TestPos() {
	lng := 114.429444
	lat := 39.0
	fmt.Println(Wgs2Gcj(POS{LNG: lng, LAT: lat}))
	fmt.Println(Wgs2BD(POS{LNG: lng, LAT: lat}))
	fmt.Println(Gcj2Wgs(POS{LNG: lng, LAT: lat}))
	fmt.Println(BD2Wgs(POS{LNG: lng, LAT: lat}))
	fmt.Println(Gcj2BD(POS{LNG: lng, LAT: lat}))
	fmt.Println(BD2Gcj(POS{LNG: lng, LAT: lat}))
}
