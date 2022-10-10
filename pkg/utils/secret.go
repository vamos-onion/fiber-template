package utils

import (
	"crypto/sha256"
	"fmt"
	mrand "math/rand"
	"strconv"
	"time"

	"github.com/google/uuid"
)

const digits = uint64(1_000_000)

// Same logic as MakeUintUniqueID().
// But is not perfectly same.
// Its returning value is definitely unique.
// But no one can expect the value.
func MakeUintUniqueUserId() (uint64, error) {
	now := time.Now()
	ts := fmt.Sprintf("%d", now.UnixMilli())
	mrand.Seed(time.Now().UnixNano())
	var randNum string = fmt.Sprintf("%d", mrand.Intn(1_000))
	for len(randNum) < 3 {
		randNum += fmt.Sprintf("%d", mrand.Intn(1_000))
	}
	randNum = randNum[:3]
	result, err := strconv.ParseUint(ts+randNum, 10, 64)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func HashPassword(pwd *string) {
	hash := sha256.New()
	_, err := hash.Write([]byte(*pwd))
	if err != nil {
		pwd = nil
	}
	*pwd = fmt.Sprintf("%x", hash.Sum(nil))
}

func MakeUID() uuid.UUID {
	return uuid.New()
}

// string(Unix timestamp) + randomNumericID[:9].
//
// The biggest value of type uint64 is 18446744073709551615.
// 	- It has 20 digits.
//
// For example, current timestamp in 2022 is 1647235843. (10 digits)
// 	So after 2028, the range of timestamp will be over 1844674407.
// 	It's the first 10 digits of the biggest value of uint64's maximum value.
//
// If this UUID doesn't contain "0" at the first place of 20 digits, overflow will be occured in 2028. (1844******)
// 	So we should contain "0" at the first place of 20 digits.
//
func MakeUintUniqueID() (uint64, error) {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	random_id := mrand.New(mrand.NewSource(time.Now().UnixNano()))
	var randIntToStr string = fmt.Sprintf("%d", random_id.Uint64()%digits)
	for len(randIntToStr) < 6 {
		randIntToStr += fmt.Sprintf("%d", random_id.Uint64()%digits)
	}
	randIntToStr = randIntToStr[:6]
	result, err := strconv.ParseUint(timestamp+randIntToStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return result, nil
}
