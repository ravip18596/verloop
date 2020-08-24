package Controllers

import (
	"Constants"
	"Models"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"sync"
	"time"
)

func AddWords(w http.ResponseWriter, r *http.Request) {
	var data Models.AddWordRequest
	decoder := json.NewDecoder(r.Body)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := decoder.Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error("error decoding post body. Err is ", err)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	//if there are more than one word in req then send error message
	if len(strings.Split(data.Word, " ")) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"request": data.Word,
		}).Error("multiple words sent")
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "multiple words sent"})
		return
	}
	//fetch current word count from global counter
	count, paraCnt, storyCnt, stmCnt := Constants.Add()
	var sentence string
	var title []string
	if count == 0 {
		//starting of a new story
		title = []string{data.Word}
		sentence = ""
		CreateStoriesTable(storyCnt, title, time.Now(), time.Now())
	} else {
		title = FetchCurrentStoryTitle(storyCnt)
		if count >= 2 {
			//now sentence writing starts
			sentence = FetchCurrentSentence(storyCnt, paraCnt, stmCnt)
			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				defer wg.Done()
				UpdateParagraph(data.Word, paraCnt, storyCnt, stmCnt)
			}()
			go func() {
				defer wg.Done()
				UpdateStoriesTable(storyCnt, title, time.Now())
			}()
			wg.Wait()
		} else {
			sentence = ""
			title = append(title,data.Word)
			UpdateStoriesTable(storyCnt, title, time.Now())
		}
	}
	log.Debug("current word count is ", count)
	response := Models.AddWordResponse{}
	response.Id = storyCnt
	response.Title = strings.Join(title," ")
	response.CurrentSentence = sentence
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error("error writing response into json")
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "error writing response into json"})
	}
}

func UpdateParagraph(word string, paraId, storyId int64, sentenceId int) {
	query := "update stories.paragraph set sentences = sentences + ? where story_id = ? and paragraph_id = ? and sentence_id = ?"
	err := Constants.CQL.Session.Query(query, []string{word}, storyId, paraId, sentenceId).Exec()
	if err != nil {
		log.WithFields(log.Fields{
			"err":          err,
			"query":        query,
			"word":         word,
			"story_id":     storyId,
			"paragraph_id": paraId,
			"sentence_id":  sentenceId,
		}).Error("error updating paragraph table")
	} else {
		log.WithFields(log.Fields{
			"word":         word,
			"story_id":     storyId,
			"paragraph_id": paraId,
			"sentence_id":  sentenceId,
		}).Debug("successfully updated paragraph table")
	}
}

func UpdateStoriesTable(storyId int64, title []string, updatedTime time.Time) {
	query := "insert into stories.story(story_id,updated_at,title) values (?,?,?)"
	err := Constants.CQL.Session.Query(query, storyId, updatedTime, title).Exec()
	if err != nil {
		log.WithFields(log.Fields{
			"err":        err,
			"query":      query,
			"updated_at": updatedTime,
			"story_id":   storyId,
			"title":      title,
		}).Error("error updating story table")
	} else {
		log.WithFields(log.Fields{
			"err":        err,
			"query":      query,
			"updated_at": updatedTime,
			"story_id":   storyId,
			"title":      title,
		}).Debug("successfully updated story table")
	}
}

func CreateStoriesTable(storyId int64, title []string, createdTime, updatedTime time.Time) {
	query := "insert into stories.story(story_id,created_at,updated_at,title) values (?,?,?,?)"
	err := Constants.CQL.Session.Query(query, storyId, createdTime, updatedTime, title).Exec()
	if err != nil {
		log.WithFields(log.Fields{
			"err":        err,
			"query":      query,
			"created_at": createdTime,
			"updated_at": updatedTime,
			"story_id":   storyId,
			"title":      title,
		}).Error("error creating story table")
	} else {
		log.WithFields(log.Fields{
			"err":        err,
			"query":      query,
			"created_at": createdTime,
			"updated_at": updatedTime,
			"story_id":   storyId,
			"title":      title,
		}).Debug("successfully created story table")
	}
}

func FetchCurrentSentence(storyId, paragraphId int64, sentenceId int) string {
	var sentences []string
	query := "select sentences from stories.paragraph where story_id=? and paragraph_id=? and sentence_id=?;"
	Constants.CQL.Session.Query(query, storyId, paragraphId, sentenceId).Iter().Scan(&sentences)
	return strings.Join(sentences, " ")
}

func FetchCurrentStoryTitle(storyId int64) []string {
	var title []string
	query := "select title from stories.story where story_id=? ;"
	Constants.CQL.Session.Query(query, storyId).Iter().Scan(&title)
	return title
}
