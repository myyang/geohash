package geohash

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func TestEncode36(t *testing.T) {
	exp, got := "", ""

	exp = "bdrdC26BqH"
	got, _, _ = Encode36(51.504444, -0.086666, 10)

	if !reflect.DeepEqual(exp, got) {
		_, file, line, _ := runtime.Caller(0)
		fmt.Printf(
			"%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\n\n",
			filepath.Base(file), line, exp, got)
		t.FailNow()
	}

	exp = "bdrdC26"
	got, _, _ = Encode36(51.504444, -0.086666, 7)

	if !reflect.DeepEqual(exp, got) {
		_, file, line, _ := runtime.Caller(0)
		fmt.Printf(
			"%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\n\n",
			filepath.Base(file), line, exp, got)
		t.FailNow()
	}

	exp = "H2RXqLHNG6"
	got, _, _ = Encode36(25.03297033, 121.56542031, 10)

	if !reflect.DeepEqual(exp, got) {
		_, file, line, _ := runtime.Caller(0)
		fmt.Printf(
			"%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\n\n",
			filepath.Base(file), line, exp, got)
		t.FailNow()
	}
}

func TestDecode36(t *testing.T) {
	lat, lng := 0.0, 0.0
	explat, explng := 0.0, 0.0

	lat, lng = Decode36("bdrdC26BqH", 6)
	explat, explng = 51.504444, -0.086666

	if lat != explat || lng != explng {
		_, file, line, _ := runtime.Caller(0)
		fmt.Printf(
			"%s:%d:\n\n\texp: (%#v, %#v)\n\n\tgot: (%#v, %#v)\n\n",
			filepath.Base(file), line, explat, explng, lat, lng)
		t.FailNow()
	}

	lat, lng = Decode36("H2RXqLHNG6", 8)
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
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Encode36(25.03297033, 121.56542031, 10)
	}
}

func BenchmarkDecode36(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Decode36("H2RXqLHNG6", 8)
	}
}
