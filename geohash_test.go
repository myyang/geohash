package geohash

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func TestEncode(t *testing.T) {
	actual, expected := "", ""
	actual, expected = Encode(12.04512315, 118.20385763, 9), "wdhh9b9rv"
	if !reflect.DeepEqual(expected, actual) {
		_, file, line, _ := runtime.Caller(0)
		fmt.Printf("%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\n\n", filepath.Base(file), line, expected, actual)
		t.FailNow()
	}

	actual, expected = Encode(-2, -3, 1), "7"
	if !reflect.DeepEqual(expected, actual) {
		_, file, line, _ := runtime.Caller(0)
		fmt.Printf("%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\n\n", filepath.Base(file), line, expected, actual)
		t.FailNow()
	}

	actual, expected = Encode(-2, -3, 6), "7ztuee"
	if !reflect.DeepEqual(expected, actual) {
		_, file, line, _ := runtime.Caller(0)
		fmt.Printf("%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\n\n", filepath.Base(file), line, expected, actual)
		t.FailNow()
	}
}

func TestDecode(t *testing.T) {
	lat, lng := Decode("wdhh9b9rv", 8)
	explat, explng := 12.04511404, 118.20385695
	if lat != explat && lng != explng {
		_, file, line, _ := runtime.Caller(0)
		fmt.Printf("%s:%d:\n\n\texp: (%v, %v)\n\n\tgot: (%v, %v)\n\n", filepath.Base(file), line, explat, explng, lat, lng)
		t.FailNow()
	}

	lat, lng = Decode("7", 1)
	explat, explng = -22.5, -22.5
	if lat != explat && lng != explng {
		_, file, line, _ := runtime.Caller(0)
		fmt.Printf("%s:%d:\n\n\texp: (%v, %v)\n\n\tgot: (%v, %v)\n\n", filepath.Base(file), line, explat, explng, lat, lng)
		t.FailNow()
	}

	lat, lng = Decode("7ztuee", 6)
	explat, explng = -2.002258, -3.004761
	if lat != explat && lng != explng {
		_, file, line, _ := runtime.Caller(0)
		fmt.Printf("%s:%d:\n\n\texp: (%v, %v)\n\n\tgot: (%v, %v)\n\n", filepath.Base(file), line, explat, explng, lat, lng)
		t.FailNow()
	}

}

func TestDecodeError(t *testing.T) {
	lat, lng := Decode("aaaaaaa", 7)
	explat, explng := 89.9993134, 179.9993134
	if lat != explat && lng != explng {
		_, file, line, _ := runtime.Caller(0)
		fmt.Printf("%s:%d:\n\n\texp: (%v, %v)\n\n\tgot: (%v, %v)\n\n", filepath.Base(file), line, explat, explng, lat, lng)
		t.FailNow()
	}
}

func BenchmarkEncode(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Encode(12.04512315, 118.20385763, 9)
	}
}

func BenchmarkDecode(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Decode("wdhh9b9rv", 8)
	}
}
