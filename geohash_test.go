package geohash

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func TestLatErr(t *testing.T) {
	tr := []struct {
		Input  int
		Out    float64
		Errstr string
	}{
		{1, 22.5, "error at bit 1"},
		{2, 2.81, "error at bit 2"},
		{3, 0.703, "error at bit 3"},
		{4, 0.0879, "error at bit 4"},
		{5, 0.02197, "error at bit 5"},
		{6, 0.002747, "error at bit 6"},
		{7, 0.0006866, "error at bit 7"},
		{8, 8.583e-05, "error at bit 8"},
		{9, 2.1458e-05, "error at bit 9"},
		{10, 2.6822e-06, "error at bit 10"},
		{11, 6.7055e-07, "error at bit 11"},
		{12, 8.3819e-08, "error at bit 12"},
	}

	for _, v := range tr {
		if r := latErr(v.Input); v.Out != r {
			fmt.Println(v.Errstr, v.Out, " != ", r)
			t.FailNow()
		}
	}
}

func TestLngErr(t *testing.T) {
	tr := []struct {
		Input  int
		Out    float64
		Errstr string
	}{
		{1, 22.5, "error at len 1"},
		{2, 5.63, "error at len 2"},
		{3, 0.703, "error at len 3"},
		{4, 0.1758, "error at len 4"},
		{5, 0.02197, "error at len 5"},
		{6, 0.005493, "error at len 6"},
		{7, 0.0006866, "error at len 7"},
		{8, 0.00017166, "error at len 8"},
		{9, 2.1458e-05, "error at len 9"},
		{10, 5.3644e-06, "error at len 10"},
		{11, 6.7055e-07, "error at len 11"},
		{12, 1.67638e-07, "error at len 12"},
	}

	for _, v := range tr {
		if r := lngErr(v.Input); v.Out != r {
			fmt.Println(v.Errstr, v.Out, " != ", r)
			t.FailNow()
		}
	}
}

func TestCryptorInf(t *testing.T) {
	c1 := NewDefaultGeoHash()
	c2 := NewGeoHash("abcdefghijklmnopqrstuvwxyz123456")
	switch c1.(type) {
	case GeoCryptor:
	default:
		fmt.Printf("c1 is not type of GeoCryptor\n")
		t.FailNow()
	}
	switch c2.(type) {
	case GeoCryptor:
	default:
		fmt.Printf("c2 is not type of GeoCryptor\n")
		t.FailNow()
	}
}

func TestEncodeGeoHash(t *testing.T) {
	cryptor := NewDefaultGeoHash()
	actual, expected := "", ""

	actual, expected = cryptor.Encode(12.04512315, 118.20385763, 9), "wdhh9b9rv"
	if actual != expected {
		_, file, line, _ := runtime.Caller(0)
		fmt.Printf("%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\n\n", filepath.Base(file), line, expected, actual)
		t.FailNow()
	}

	actual, expected = cryptor.Encode(-2, -3, 1), "7"
	if actual != expected {
		_, file, line, _ := runtime.Caller(0)
		fmt.Printf("%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\n\n", filepath.Base(file), line, expected, actual)
		t.FailNow()
	}

	actual, expected = cryptor.Encode(-2, -3, 6), "7ztuee"
	if actual != expected {
		_, file, line, _ := runtime.Caller(0)
		fmt.Printf("%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\n\n", filepath.Base(file), line, expected, actual)
		t.FailNow()
	}

	latErr, lngErr, expLatErr, expLngErr := 0.0, 0.0, 0.0, 0.0

	actual, latErr, lngErr = cryptor.EncodeWithErr(12.04512315, 118.20385763, 9)
	expected, expLatErr, expLngErr = "wdhh9b9rv", 2.1458e-05, 2.1458e-05
	if actual != expected || latErr != expLatErr || lngErr != expLngErr {
		_, file, line, _ := runtime.Caller(0)
		fmt.Printf("%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\n\n", filepath.Base(file), line, expected, actual)
		t.FailNow()
	}

	actual, latErr, lngErr = cryptor.EncodeWithErr(-2, -3, 1)
	expected, expLatErr, expLngErr = "7", 22.5, 22.5
	if actual != expected || latErr != expLatErr || lngErr != expLngErr {
		_, file, line, _ := runtime.Caller(0)
		fmt.Printf("%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\n\n", filepath.Base(file), line, expected, actual)
		t.FailNow()
	}

	actual, latErr, lngErr = cryptor.EncodeWithErr(-2, -3, 6)
	expected, expLatErr, expLngErr = "7ztuee", 0.002747, 0.005493
	if actual != expected || latErr != expLatErr || lngErr != expLngErr {
		_, file, line, _ := runtime.Caller(0)
		fmt.Printf("%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\n\n", filepath.Base(file), line, expected, actual)
		t.FailNow()
	}

	r := cryptor.EncodeAsBox(-2, -3, 6)
	latErr, lngErr = r.ErrPair()
	if h, err := r.Geohash(); err != nil || h != expected || latErr != expLatErr || lngErr != expLngErr {
		_, file, line, _ := runtime.Caller(0)
		fmt.Printf("%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\n\n", filepath.Base(file), line, expected, actual)
		t.FailNow()
	}

}

func TestDecode(t *testing.T) {
	cryptor := NewDefaultGeoHash()

	lat, lng := cryptor.Decode("wdhh9b9rv", 8)
	explat, explng := 12.04511404, 118.20385695
	if lat != explat && lng != explng {
		_, file, line, _ := runtime.Caller(0)
		fmt.Printf("%s:%d:\n\n\texp: (%v, %v)\n\n\tgot: (%v, %v)\n\n", filepath.Base(file), line, explat, explng, lat, lng)
		t.FailNow()
	}

	lat, lng = cryptor.Decode("7", 1)
	explat, explng = -22.5, -22.5
	if lat != explat && lng != explng {
		_, file, line, _ := runtime.Caller(0)
		fmt.Printf("%s:%d:\n\n\texp: (%v, %v)\n\n\tgot: (%v, %v)\n\n", filepath.Base(file), line, explat, explng, lat, lng)
		t.FailNow()
	}

	lat, lng = cryptor.Decode("7ztuee", 6)
	explat, explng = -2.002258, -3.004761
	if lat != explat && lng != explng {
		_, file, line, _ := runtime.Caller(0)
		fmt.Printf("%s:%d:\n\n\texp: (%v, %v)\n\n\tgot: (%v, %v)\n\n", filepath.Base(file), line, explat, explng, lat, lng)
		t.FailNow()
	}

	latErr, lngErr, expLatErr, expLngErr := 0.0, 0.0, 0.0, 0.0
	var err error

	lat, lng, latErr, lngErr = cryptor.DecodeWithErr("wdhh9b9rv", 8)
	explat, explng, expLatErr, expLngErr = 12.04511404, 118.20385695, 8.583e-05, 0.00017166
	if lat != explat || lng != explng || latErr != expLatErr || lngErr != expLngErr {
		_, file, line, _ := runtime.Caller(0)
		fmt.Printf(
			"%s:%d:\n\n\texp: (%v, %v) and err pair (%v, %v)\n\n\tgot: (%v, %v) and err pair (%v, %v)\n\n",
			filepath.Base(file), line, explat, explng, expLatErr, expLngErr, lat, lng, latErr, lngErr)
		t.FailNow()
	}

	lat, lng, latErr, lngErr = cryptor.DecodeWithErr("7", 1)
	explat, explng, expLatErr, expLngErr = -22.5, -22.5, 22.5, 22.5
	if lat != explat || lng != explng || latErr != expLatErr || lngErr != expLngErr {
		_, file, line, _ := runtime.Caller(0)
		fmt.Printf("%s:%d:\n\n\texp: (%v, %v)\n\n\tgot: (%v, %v)\n\n", filepath.Base(file), line, explat, explng, lat, lng)
		t.FailNow()
	}

	lat, lng, latErr, lngErr = cryptor.DecodeWithErr("7ztuee", 6)
	explat, explng, expLatErr, expLngErr = -2.002258, -3.004761, 0.002747, 0.005493
	if lat != explat || lng != explng || latErr != expLatErr || lngErr != expLngErr {
		_, file, line, _ := runtime.Caller(0)
		fmt.Printf("%s:%d:\n\n\texp: (%v, %v)\n\n\tgot: (%v, %v)\n\n", filepath.Base(file), line, explat, explng, lat, lng)
		t.FailNow()
	}

	r := cryptor.DecodeAsBox("7ztuee", 6)
	latErr, lngErr = r.ErrPair()
	lat, lng, err = r.GetCenter()
	if err != nil || lat != explat || lng != explng || latErr != expLatErr || lngErr != expLngErr {
		_, file, line, _ := runtime.Caller(0)
		fmt.Printf("%s:%d:\n\n\texp: (%v, %v)\n\n\tgot: (%v, %v)\n\n", filepath.Base(file), line, explat, explng, lat, lng)
		t.FailNow()
	}
}

func TestDecodeError(t *testing.T) {
	cryptor := NewDefaultGeoHash()

	lat, lng := cryptor.Decode("aaaaaaa", 7)
	explat, explng := 89.9993134, 179.9993134
	if lat != explat && lng != explng {
		_, file, line, _ := runtime.Caller(0)
		fmt.Printf("%s:%d:\n\n\texp: (%v, %v)\n\n\tgot: (%v, %v)\n\n", filepath.Base(file), line, explat, explng, lat, lng)
		t.FailNow()
	}
}

func TestNeighbors(t *testing.T) {
	cryptor := NewDefaultGeoHash()
	got, exp := []string{}, []string{}
	neighbors := cryptor.Neighbors("7ztuee", 6)

	exp = []string{"7ztue6", "7ztued", "7ztuef", "7ztue7", "7ztueg", "7ztuek", "7ztues", "7ztueu"}
	for _, v := range neighbors {
		if v.(*LocationBox) == nil {
			continue
		}
		h, err := v.Geohash()
		if err == nil {
			got = append(got, h)
		}
	}
	if !reflect.DeepEqual(exp, got) {
		_, file, line, _ := runtime.Caller(0)
		fmt.Printf("%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\n\n", filepath.Base(file), line, exp, got)
		t.FailNow()
	}

	neighbors = cryptor.Neighbors("eb", 2)
	got, exp = []string{}, []string{"7x", "7z", "kp", "e8", "s0", "e9", "ec", "s1"}
	for _, v := range neighbors {
		if v.(*LocationBox) == nil {
			continue
		}
		h, err := v.Geohash()
		if err == nil {
			got = append(got, h)
		}
	}
	if !reflect.DeepEqual(exp, got) {
		_, file, line, _ := runtime.Caller(0)
		fmt.Printf("%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\n\n", filepath.Base(file), line, exp, got)
		t.FailNow()
	}

	neighbors = cryptor.Neighbors("gz", 2)
	got, exp = []string{}, []string{"gw", "gy", "un", "gx", "up"}
	for _, v := range neighbors {
		if v.(*LocationBox) == nil {
			continue
		}
		h, err := v.Geohash()
		if err == nil {
			got = append(got, h)
		}
	}
	if !reflect.DeepEqual(exp, got) {
		_, file, line, _ := runtime.Caller(0)
		fmt.Printf("%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\n\n", filepath.Base(file), line, exp, got)
		t.FailNow()
	}
}

func BenchmarkEncode(b *testing.B) {
	b.ReportAllocs()
	cryptor := NewDefaultGeoHash()
	for i := 0; i < b.N; i++ {
		cryptor.Encode(12.04512315, 118.20385763, 12)
	}
}

func BenchmarkDecode(b *testing.B) {
	b.ReportAllocs()
	cryptor := NewDefaultGeoHash()
	for i := 0; i < b.N; i++ {
		cryptor.Decode("wdhh9b9rv", 12)
	}
}

func BenchmarkNeighbors(b *testing.B) {
	b.ReportAllocs()
	cryptor := NewDefaultGeoHash()
	for i := 0; i < b.N; i++ {
		cryptor.Neighbors("7ztuee", 6)
	}
}
