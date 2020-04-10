/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package twqrp

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	re3Digit  *regexp.Regexp
	reAccount *regexp.Regexp
)

func init() {
	re3Digit = regexp.MustCompile("^[0-9]{3}$")
	reAccount = regexp.MustCompile("^[0-9]{1,16}$")
}

// G 代表 Generator，用來產生 TWQRP 內容
//
// 預設國別是台灣 (158)
type G struct {
	Name    string // 廠商自訂的服務名稱
	country int
	typ     int
	params  map[string]string
	Mutable bool // 是否可以讓付款人更改內容
}

// NewEmpty 回傳一個空白的 TWQRP 產生器
//
// 所有參數都必須自行透過 SetParam() 手動設定，也不會檢查
func NewEmpty(serviceType int) (ret *G) {
	return &G{
		country: 158,
		typ:     serviceType,
		params:  map[string]string{},
	}
}

// NewTransfer 回傳一個用來產生轉帳 TWQRP 的產生器
//
// 這個函式會檢查以下情況
//
//   - bankCode 必須是三碼數字
//   - account 必須是 1 ~ 16 碼數字
func NewTransfer(bankCode string, account string) (ret *G, err error) {
	if !re3Digit.MatchString(bankCode) {
		err = errors.New("twqrp: invalid bank code")
		return
	}

	if !reAccount.MatchString(account) {
		err = errors.New("twqrp: invalid account")
	}

	ret = &G{
		country: 158,
		typ:     2,
		params:  map[string]string{},
	}
	ret.SetParam(5, bankCode)
	ret.add(6, "%016s", account)

	return
}

// SetParam 手動設定相關參數
func (g *G) SetParam(key int, val string) {
	g.params[strconv.Itoa(key)] = val
}

func (g *G) add(key int, tmpl string, args ...interface{}) {
	g.SetParam(key, fmt.Sprintf(tmpl, args...))
}

// SortedString 實際產生 TWQRP 內容，順序固定
func (g *G) SortedString() (ret string) {
	prefix := "D"
	if g.Mutable {
		prefix = "M"
	}

	// 固定 order
	keys := make([]string, 0, len(g.params))
	for k := range g.params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	arr := make([]string, 0, len(g.params))
	for _, k := range keys {
		arr = append(arr, prefix+k+"="+g.params[k])
	}

	return fmt.Sprintf(
		"TWQRP://%s/%03d/%02d/V1?%s",
		g.Name,
		g.country,
		g.typ,
		strings.Join(arr, "&"),
	)
}

// String 實際產生 TWQRP 內容，順序不固定
func (g *G) String() (ret string) {
	prefix := "D"
	if g.Mutable {
		prefix = "M"
	}

	buf := &bytes.Buffer{}
	buf.WriteString(fmt.Sprintf(
		"TWQRP://%s/%03d/%02d/V1?",
		g.Name,
		g.country,
		g.typ,
	))

	for k, v := range g.params {
		buf.WriteString(prefix + k + "=" + v + "&")
	}

	return string(buf.Bytes()[:buf.Len()-1])
}

func (g *G) must(err error) (ret *G) {
	if err != nil {
		panic(err)
	}

	return g
}

// TryCountry 用來設定國別代碼，會檢查是否為三位數
func (g *G) TryCountry(code int) (err error) {
	if code < 1 || code > 999 {
		return errors.New("twqrp: invalid country code")
	}
	g.country = code
	return
}

// Country 用來設定國別代碼，有錯會直接 panic
func (g *G) Country(code int) (ret *G) {
	return g.must(g.TryCountry(code))
}

// TryAmount 用來設定交易金額，會檢查是否為五位正整數
func (g *G) TryAmount(a int) (err error) {
	if a < 1 || a > 99999 {
		return errors.New("twqrp: invalid amount")
	}
	g.add(1, "%d", a*100)
	return
}

// Amount 用來設定交易金額，有錯就直接 panic
func (g *G) Amount(a int) (ret *G) {
	return g.must(g.TryAmount(a))
}

// TryNote 用來設定備註，最多 19 字
func (g *G) TryNote(n string) (err error) {
	if len(n) > 19 {
		return errors.New("twqrp: invalid note")
	}
	g.SetParam(9, n)
	return
}

// Note 用來設定備註，有錯直接 panic
func (g *G) Note(n string) (ret *G) {
	return g.must(g.TryNote(n))
}

// TryCurrency 用來設定幣別，會檢查是否為三位數，預設為台幣 901
func (g *G) TryCurrency(n string) (err error) {
	if !re3Digit.MatchString(n) {
		return errors.New("twqrp: invalid currency code")
	}

	g.SetParam(10, n)
	return
}

// Currency 用來設定幣別，有錯直接 panic
func (g *G) Currency(n string) (ret *G) {
	return g.must(g.TryCurrency(n))
}

// TryQRDue 用來設定 QRCode 有效期限，請自行注意時區問題
func (g *G) TryQRDue(t time.Time) (err error) {
	if t.IsZero() {
		return errors.New("twqrp: invalid qrcode due time")
	}

	g.SetParam(12, t.Format("20060102150405"))
	return
}

// QRDue 用來設定 QRCode 有效期限，有錯直接 panic
func (g *G) QRDue(t time.Time) (ret *G) {
	return g.must(g.TryQRDue(t))
}
