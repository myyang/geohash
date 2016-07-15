package geohash

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
)

var degreeErr = []struct {
	LatErr float64
	LngErr float64
}{
	{0, 0},
	{30.0, 60.0},
	{5.0, 10.0},
	{0.833, 1.667},
	{0.1389, 0.2778},
	{0.02315, 0.0463},
	{0.003858, 0.007716},
	{0.000643, 0.001286},
	{0.00010717, 0.00021433},
	{1.7861e-05, 3.5722e-05},
	{2.9769e-06, 5.9537e-06},
	{4.9615e-07, 9.9229e-07},
	{8.2691e-08, 1.65382e-07},
}

func TestEncode36(t *testing.T) {
	cryptor := NewGeoHash36WithDefaultKey()
	exp, got := "", ""

	exp = "bdrdC26BqH"
	got = cryptor.Encode(51.504444, -0.086666, 10)

	if exp != got {
		_, file, line, _ := runtime.Caller(0)
		fmt.Printf(
			"%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\n\n",
			filepath.Base(file), line, exp, got)
		t.FailNow()
	}

	exp = "bdrdC26"
	got = cryptor.Encode(51.504444, -0.086666, 7)

	if exp != got {
		_, file, line, _ := runtime.Caller(0)
		fmt.Printf(
			"%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\n\n",
			filepath.Base(file), line, exp, got)
		t.FailNow()
	}

	exp = "H2RXqLHNG6"
	got = cryptor.Encode(25.03297033, 121.56542031, 10)

	if exp != got {
		_, file, line, _ := runtime.Caller(0)
		fmt.Printf(
			"%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\n\n",
			filepath.Base(file), line, exp, got)
		t.FailNow()
	}

}

func TestDecode36(t *testing.T) {
	cryptor := NewGeoHash36WithDefaultKey()
	lat, lng := 0.0, 0.0
	explat, explng := 0.0, 0.0

	lat, lng = cryptor.Decode("bdrdC26BqH", 6)
	explat, explng = 51.504444, -0.086666

	if lat != explat || lng != explng {
		_, file, line, _ := runtime.Caller(0)
		fmt.Printf(
			"%s:%d:\n\n\texp: (%#v, %#v)\n\n\tgot: (%#v, %#v)\n\n",
			filepath.Base(file), line, explat, explng, lat, lng)
		t.FailNow()
	}

	lat, lng = cryptor.Decode("H2RXqLHNG6", 8)
	explat, explng = 25.03297033, 121.56542031

	if lat != explat || lng != explng {
		_, file, line, _ := runtime.Caller(0)
		fmt.Printf(
			"%s:%d:\n\n\texp: (%#v, %#v)\n\n\tgot: (%#v, %#v)\n\n",
			filepath.Base(file), line, explat, explng, lat, lng)
		t.FailNow()
	}
}

func BenchmarkEncode36(b *testing.B) {
	cryptor := NewGeoHash36WithDefaultKey()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cryptor.Encode(25.03297033, 121.56542031, 10)
	}
}

func BenchmarkDecode36(b *testing.B) {
	cryptor := NewGeoHash36WithDefaultKey()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cryptor.Decode("H2RXqLHNG6", 8)
	}
}
