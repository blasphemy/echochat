package main

import (
	"time"
)

var (
	bancounter = 1
)

type Ban struct {
	mask   string
	whoset string
	ts     time.Time
	id     int
}

func NewBan(bmask string, bwhoset string) *Ban {
	bid := bancounter
	bancounter++
	return &Ban{
		mask:   bmask,
		whoset: bwhoset,
		ts:     time.Now(),
		id:     bid,
	}
}

func CheckIfBanExists(channel *Channel, b string) bool {
	for _, k := range channel.banlist {
		if k.mask == b {
			return true
		}
	}
	return false
}
