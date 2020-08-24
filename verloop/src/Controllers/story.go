package Controllers

import (
	"Constants"
	"Models"
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"time"
	Sort "sort"
)

func FetchAllStories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var limit,offset int
	var sort,order string
	limit, _ = strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 50
	}
	offset, _ = strconv.Atoi(r.URL.Query().Get("offset"))
	sort = r.URL.Query().Get("sort")
	order = r.URL.Query().Get("order")
	if sort == ""{
		sort = Constants.SortCreatedAt
	}
	if order == ""{
		order = Constants.ASC
	}
	stories := fetchAllStories()
	if strings.EqualFold(sort,Constants.SortCreatedAt){
		if strings.EqualFold(order,Constants.ASC) {
			Sort.Slice(stories, func(i, j int) bool {
				return stories[i].CreatedAt.UnixNano() < stories[j].CreatedAt.UnixNano()
			})
		}else if strings.EqualFold(order,Constants.DESC){
			Sort.Slice(stories, func(i, j int) bool {
				return stories[i].CreatedAt.UnixNano() > stories[j].CreatedAt.UnixNano()
			})
		}
	}else if strings.EqualFold(sort,Constants.SortUpdatedAt){
		if strings.EqualFold(order,Constants.ASC) {
			Sort.Slice(stories, func(i, j int) bool {
				return stories[i].UpdatedAt.UnixNano() < stories[j].UpdatedAt.UnixNano()
			})
		}else if strings.EqualFold(order,Constants.DESC){
			Sort.Slice(stories, func(i, j int) bool {
				return stories[i].UpdatedAt.UnixNano() > stories[j].UpdatedAt.UnixNano()
			})
		}
	}else if strings.EqualFold(sort,Constants.SortTitle){
		if strings.EqualFold(order,Constants.ASC) {
			Sort.Slice(stories, func(i, j int) bool {
				return stories[i].Title < stories[j].Title
			})
		}else if strings.EqualFold(order,Constants.DESC){
			Sort.Slice(stories, func(i, j int) bool {
				return stories[i].Title > stories[j].Title
			})
		}
	}
	if len(stories)>offset {
		end := min(len(stories),offset+limit)
		stories = stories[offset:end]
	}else{
		stories = []Models.Stories{}
	}
	resp := Models.AllStoryAPIResponse{
		Limit:   limit,
		Offset:  offset,
		Count:   len(stories),
		Results: stories,
	}

	err := json.NewEncoder(w).Encode(resp)
	if err != nil{
		log.Error("error writing response into json")
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "error writing response into json"})
	}
}

func FetchStory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	storyId,err := strconv.ParseInt(vars["story_id"],10,64)
	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"story_id":storyId,
			"err":err,
		}).Error("story id is not present or incorrect value")
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "story id is not present or incorrect value"})
		return
	}
	story := fetchStoryById(storyId)
	err = json.NewEncoder(w).Encode(story)
	if err != nil{
		log.Error("error writing response into json")
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "error writing response into json"})
	}
}

func fetchAllStories() []Models.Stories {
	var stories []Models.Stories
	var id int64
	var title []string
	var createdAt,updatedAt time.Time
	query := "select story_id,created_at,updated_at,title from stories.story;"
	iter := Constants.CQL.Session.Query(query).Iter()
	for iter.Scan(&id,&createdAt,&updatedAt,&title){
		stories = append(stories, Models.Stories{
			StoryId:   id,
			Title:     strings.Join(title," "),
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
	}
	return stories
}

func fetchStoryById(storyId int64) Models.Story{
	var sentences []string
	var paraId int64
	var story Models.Story
	var title []string
	paraMap := make(map[int64][][]string)
	query := "select sentences,paragraph_id from stories.paragraph where story_id=?"
	iter := Constants.CQL.Session.Query(query,storyId).Iter()

	for iter.Scan(&sentences,&paraId){
		paraMap[paraId] = append(paraMap[paraId],sentences)
	}

	for _,values := range paraMap{
		//values are sentences in a paragraph
		//10 sentences in paragraph
		//15 words in a sentences
		story.Paragraphs = append(story.Paragraphs,Models.Paragraph{Sentences:values})
	}
	_ = iter.Close()

	query = "select story_id,created_at,updated_at,title from stories.story where story_id = ?"
	iter = Constants.CQL.Session.Query(query,storyId).Iter()
	iter.Scan(&story.StoryId,&story.CreatedAt,&story.UpdatedAt,&title)
	story.Title = strings.Join(title," ")
	return story
}

func min(a,b int) int{
	if a<b{
		return a
	}
	return b
}