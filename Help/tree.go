package Help

import (
	"blog/models"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

type CategoryTree struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	Title     string             `json:"title"`
	Parent_Id string             `json:"parent_id"`
	Checked   bool               `json:"checked"`
	Children  []*CategoryTree    `orm:"-" json:"children"`
}

func GetTree(categoryIds map[string]string) []*CategoryTree {
	db := models.ConnectDb()
	var people []*CategoryTree
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cursor, err := db.Collection("category").Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person CategoryTree
		cursor.Decode(&person)

		people = append(people, &person)
	}
	tree := tree(people, categoryIds)
	return tree
}
func tree(list []*CategoryTree, categoryIds map[string]string) []*CategoryTree {
	data := buildData(list)
	result := makeTreeCore("0", data, categoryIds)
	body := result
	return body

}
func buildData(list []*CategoryTree) map[string]map[string]*CategoryTree {
	var data map[string]map[string]*CategoryTree = make(map[string]map[string]*CategoryTree)
	for _, v := range list {
		id := v.Id.Hex()
		parent_id := v.Parent_Id
		if _, ok := data[parent_id]; !ok {
			data[parent_id] = make(map[string]*CategoryTree)
		}
		data[parent_id][id] = v

	}
	return data
}

func makeTreeCore(index string, data map[string]map[string]*CategoryTree, categoryIds map[string]string) []*CategoryTree {
	tmp := make([]*CategoryTree, 0)
	for id, item := range data[index] {
		fmt.Println(data[id])
		if data[id] != nil {
			item.Children = makeTreeCore(id, data, categoryIds)
		}
		if categoryIds != nil {
			if categoryIds[id] != "" {
				item.Checked = true
			}
		}
		tmp = append(tmp, item)
	}
	return tmp
}
