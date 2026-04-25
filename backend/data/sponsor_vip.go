package data

import (
	"encoding/hex"
	"encoding/json"
	"strings"
	"time"

	"github.com/duke-git/lancet/v2/convertor"
	"github.com/duke-git/lancet/v2/cryptor"
)

const DefaultSponsorAESKeyHex = ""

var SponsorDecryptKeyHex string

func EffectiveSponsorVipLevel() (level int, active bool) {
	return 99, true
}

func EffectiveSponsorVipLevelOriginal() (level int, active bool) {
	keyHex := strings.TrimSpace(SponsorDecryptKeyHex)
	if keyHex == "" {
		keyHex = DefaultSponsorAESKeyHex
	}
	sponsorCode := strings.TrimSpace(GetSettingConfig().SponsorCode)
	if sponsorCode == "" {
		return 0, false
	}
	encrypted, err := hex.DecodeString(sponsorCode)
	if err != nil {
		return 0, false
	}
	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return 0, false
	}
	raw := cryptor.AesEcbDecrypt(encrypted, key)
	if len(raw) == 0 {
		return 0, false
	}
	var info map[string]any
	if err := json.Unmarshal(raw, &info); err != nil {
		return 0, false
	}
	lvl64, _ := convertor.ToInt(info["vipLevel"])
	lvl := int(lvl64)
	vipStartTime, err1 := time.ParseInLocation("2006-01-02 15:04:05", convertor.ToString(info["vipStartTime"]), time.Local)
	vipEndTime, err2 := time.ParseInLocation("2006-01-02 15:04:05", convertor.ToString(info["vipEndTime"]), time.Local)
	vipAuthTime, err3 := time.ParseInLocation("2006-01-02 15:04:05", convertor.ToString(info["vipAuthTime"]), time.Local)
	if err1 != nil || err2 != nil || err3 != nil {
		return lvl, false
	}
	now := time.Now()
	if now.After(vipAuthTime) && now.After(vipStartTime) && now.Before(vipEndTime) {
		return lvl, true
	}
	return lvl, false
}
