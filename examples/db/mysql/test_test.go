package mysql

import (
	"github.com/v587-zyf/gc/db/mysql"
	"testing"
	"time"
)

func TestInsertOne(t *testing.T) {
	d := &TestModel{
		Time: time.Now(),
	}
	if err := mysql.CreateModel[TestModel](d, nil); err != nil {
		t.Error(err)
		return
	}

	data := &TestModel{
		ModelBase: ModelBase{
			ID: d.ID,
		},
	}
	if err := mysql.LoadModel[TestModel](data); err != nil {
		t.Error(err)
		return
	}
	t.Log(data)
}
