package controllers

type ArticleClass struct {
	Id       int             `json:"id"`
	ParentID int             `json:"parent_id"`
	Name     string          `json:"name"`
	List     []*ArticleClass `json:"list,omitempty"`
}
