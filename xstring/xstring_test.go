package xstring

import (
	"github.com/go-playground/assert/v2"
	"log"
	"testing"
	"time"
)

func TestCapitalize(t *testing.T) {
	assert.Equal(t, Capitalize("abc"), "Abc")
	assert.Equal(t, Capitalize("Abc"), "Abc")
	assert.Equal(t, Capitalize(""), "")
}

func TestUncapitalize(t *testing.T) {
	assert.Equal(t, Uncapitalize("Abc"), "abc")
	assert.Equal(t, Uncapitalize("abc"), "abc")
	assert.Equal(t, Uncapitalize(""), "")
}

func TestToRuneToByte(t *testing.T) {
	log.Printf("%T", 'a')         // int32
	log.Printf("%T", ToRune("a")) // int32
	log.Printf("%T", ToByte("a")) // uint8
	log.Printf("%T", "a"[0])      // uint8

	assert.Equal(t, ToRune("a"), 'a')
	assert.Equal(t, ToRune("bcd"), 'b')
	assert.Equal(t, ToRune(""), rune(0))

	assert.Equal(t, ToByte("a"), byte('a'))
	assert.Equal(t, ToByte("bcd"), byte('b'))
	assert.Equal(t, ToByte(""), byte(0))
}

func TestCurrentTimeUuid(t *testing.T) {
	log.Println(CurrentTimeUuid(5))
	log.Println(CurrentTimeUuid(24))
	log.Println(CurrentTimeUuid(30))
}

func TestRandLetterNumberString(t *testing.T) {
	log.Println(RandLetterString(20))
	log.Println(RandLetterString(20))
	log.Println(RandNumberString(20))
	log.Println(RandNumberString(20))
	log.Println(RandLetterNumberString(20))
	log.Println(RandLetterNumberString(20))

	log.Println(RandString(32, CapitalLetterRunes))
	log.Println(RandString(32, LowercaseLetterRunes))
	log.Println(RandString(32, NumberRunes))

	log.Println(RandString(32, LetterRunes))
	log.Println(RandString(32, LetterNumberRunes))
	log.Println(RandString(32, CapitalLetterNumberRunes))
	log.Println(RandString(32, LowercaseLetterNumberRunes))
}

func TestPrettifyJson(t *testing.T) {
	from := "{\"a\": \"b\", \"c\": {\"d\": \"e\", \"f\": 0}, \"g\": [{\"h\": 1}, {\"h\": 1}]}"
	to := "{\n" +
		"    \"a\": \"b\",\n" +
		"    \"c\": {\n" +
		"        \"d\": \"e\",\n" +
		"        \"f\": 0\n" +
		"    },\n" +
		"    \"g\": [\n" +
		"        {\n" +
		"            \"h\": 1\n" +
		"        },\n" +
		"        {\n" +
		"            \"h\": 1\n" +
		"        }\n" +
		"    ]\n" +
		"}"
	assert.Equal(t, PrettifyJson(from, 4, " "), to)
}

func TestToSnakeCase(t *testing.T) {
	assert.Equal(t, ToSnakeCase(""), "")
	assert.Equal(t, ToSnakeCase("AoiHosizora"), "aoi_hosizora")
	assert.Equal(t, ToSnakeCase("abc0d1EdF"), "abc0d1_ed_f")
	assert.Equal(t, ToSnakeCase("私達isわたしたち"), "私達isわたしたち")
	assert.Equal(t, ToSnakeCase("a bC"), "a_b_c")
}

func TestIsLowercase(t *testing.T) {
	assert.Equal(t, IsLowercase(ToRune("A")), false)
	assert.Equal(t, IsLowercase(ToRune("Z")), false)
	assert.Equal(t, IsLowercase(ToRune("a")), true)
	assert.Equal(t, IsLowercase(ToRune("z")), true)
	assert.Equal(t, IsLowercase(ToRune("0")), false)
	assert.Equal(t, IsLowercase(ToRune("")), false)
	assert.Equal(t, IsLowercase(ToRune("我")), false)
}

func TestIsUppercase(t *testing.T) {
	assert.Equal(t, IsUppercase(ToRune("A")), true)
	assert.Equal(t, IsUppercase(ToRune("Z")), true)
	assert.Equal(t, IsUppercase(ToRune("a")), false)
	assert.Equal(t, IsUppercase(ToRune("z")), false)
	assert.Equal(t, IsUppercase(ToRune("0")), false)
	assert.Equal(t, IsUppercase(ToRune("")), false)
	assert.Equal(t, IsUppercase(ToRune("我")), false)
}

func TestRemoveSpaces(t *testing.T) {
	assert.Equal(t, RemoveSpaces(""), "")
	assert.Equal(t, RemoveSpaces("a b  c 　d   e f"), "a b c d e f")
	assert.Equal(t, RemoveSpaces("a b 	 c d   e f"), "a b c d e f")
	assert.Equal(t, RemoveSpaces("a b \n	 c d   e f"), "a b c d e f")
	assert.Equal(t, RemoveSpaces("\n　"), "")
	assert.Equal(t, RemoveSpaces("　\n	"), "")
}

func TestMaskToken(t *testing.T) {
	assert.Equal(t, MaskToken(""), "")
	assert.Equal(t, MaskToken(" "), "*")
	assert.Equal(t, MaskToken("a"), "*")
	assert.Equal(t, MaskToken("aa"), "*a")
	assert.Equal(t, MaskToken("aaa"), "**a")
	assert.Equal(t, MaskToken("aaaa"), "a**a")
	assert.Equal(t, MaskToken("aaaaa"), "a***a")
	assert.Equal(t, MaskToken("aaaaaa"), "aa**aa")
}

func TestStringToBytes(t *testing.T) {
	assert.Equal(t, StringToBytes(""), []byte{})
	assert.Equal(t, StringToBytes("abcdefg"), []byte("abcdefg"))

	cnt := 200000000

	bs1 := make([]byte, cnt, cnt)
	bs2 := make([]byte, cnt, cnt)
	for i := 0; i < cnt; i++ {
		bs1[i] = 'A'
	}
	for i := 0; i < cnt; i++ {
		bs2[i] = 'B'
	}
	str1 := string(bs1)
	str2 := string(bs2)

	start := time.Now()
	bs01 := []byte(str1)
	log.Println(time.Now().Sub(start).String())

	start = time.Now()
	bs02 := StringToBytes(str2)
	log.Println(time.Now().Sub(start).String())

	assert.Equal(t, bs01, bs1)
	assert.Equal(t, bs02, bs2)
}

func TestBytesToString(t *testing.T) {
	assert.Equal(t, BytesToString(nil), "")
	assert.Equal(t, BytesToString([]byte{}), "")
	assert.Equal(t, BytesToString([]byte("abcdefg")), "abcdefg")

	cnt := 200000000

	bs1 := make([]byte, cnt, cnt)
	bs2 := make([]byte, cnt, cnt)
	for i := 0; i < cnt; i++ {
		bs1[i] = 'A'
	}
	for i := 0; i < cnt; i++ {
		bs2[i] = 'B'
	}
	str1 := string(bs1)
	str2 := string(bs2)

	start := time.Now()
	str01 := string(bs1)
	log.Println(time.Now().Sub(start).String())

	start = time.Now()
	str02 := BytesToString(bs2)
	log.Println(time.Now().Sub(start).String())

	assert.Equal(t, str01, str1)
	assert.Equal(t, str02, str2)
}
