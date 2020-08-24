package Constants

import (
	"Redis"
	log "github.com/sirupsen/logrus"
)

// 2 + (10*15)*7 = 1052
func Add() (int, int64, int64, int) {
	Redis.ExecuteQueryLock(Redis.LockKey)
	defer func() {
		success := Redis.ReleaseLock(Redis.LockKey)
		if success {
			log.Debug("successfully release redis lock")
		} else {
			log.Error("error release redis lock.")
		}
	}()
	count := Redis.IncrementWordCounter()
	var paragraphCount, storyCount, stmCnt int64
	paragraphCount = Redis.LoadCount(Redis.ParagraphCounterKey)
	storyCount = Redis.LoadCount(Redis.StoryCounterKey)
	stmCnt = Redis.LoadCount(Redis.SentenceCounterKey)
	if count >= 1052 {
		//new story
		Redis.SetCounterCount(Redis.WordCounterKey, "0")
		storyCount = Redis.IncrementStoryCounter()
		Redis.SetCounterCount(Redis.ParagraphCounterKey, "0")
		Redis.SetCounterCount(Redis.SentenceCounterKey, "0")
		paragraphCount = 0
		stmCnt = 0
		count = 0
	}
	if count > 2 && (count-2)%150 == 0 {
		//new paragraph
		paragraphCount = Redis.IncrementParagraphCounter()
	}
	if count > 2 && (count-2)%15 == 0 {
		//new sentence
		stmCnt = Redis.IncrementSentenceCounter()
	}
	if count == 2 {
		//title finished.
		paragraphCount = Redis.IncrementParagraphCounter()
		stmCnt = Redis.IncrementSentenceCounter()
	}
	return int(count), paragraphCount, storyCount, int(stmCnt)
}

func SetInitialCount() {
	if !Redis.CheckKeyExistance(Redis.WordCounterKey) {
		Redis.SetCounterCount(Redis.WordCounterKey, "-1")
	}
	if !Redis.CheckKeyExistance(Redis.ParagraphCounterKey) {
		Redis.SetCounterCount(Redis.ParagraphCounterKey, "0")
	}
	if !Redis.CheckKeyExistance(Redis.SentenceCounterKey) {
		Redis.SetCounterCount(Redis.SentenceCounterKey, "0")
	}
	if !Redis.CheckKeyExistance(Redis.StoryCounterKey) {
		Redis.SetCounterCount(Redis.StoryCounterKey, "1")
	}

}
