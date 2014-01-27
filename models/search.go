package models

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego/orm"

	"github.com/beego/wetalk/utils"
)

var searchEscapePattern = []string{
	`\\`, `(`, `)`, `|`, `-`, `!`, `@`, `~`, `'`, `&`, `/`, `^`, `$`, `=`,
	`\\\\`, `\(`, `\)`, `\|`, `\-`, `\!`, `\@`, `\~`, `\'`, `\&`, `\/`, `\^`, `\$`, `\=`,
}

func filterSearchQ(q string) string {
	q = strings.TrimSpace(q)
	replacer := strings.NewReplacer(searchEscapePattern...)
	return replacer.Replace(q)
}

func SearchPost(q string, page int) ([]*Post, *utils.SphinxMeta, error) {
	q = filterSearchQ(q)
	if len(q) == 0 {
		return nil, nil, fmt.Errorf("empty query")
	}

	sdb, err := utils.SphinxPools.GetConn()
	if err != nil {
		return nil, nil, err
	}
	defer sdb.Close()

	pers := 20
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * pers

	var pids orm.ParamsList
	query := fmt.Sprintf(`SELECT @id AS pid, updated
		FROM `+utils.SphinxIndex+`
		WHERE MATCH('@title `+q+`')
		ORDER BY @weight DESC, updated DESC
		LIMIT %d, %d OPTION ranker=bm25`, offset, pers)

	if _, err = sdb.RawValuesFlat(query, &pids, "pid"); err != nil {
		return nil, nil, err
	}

	var meta *utils.SphinxMeta
	if meta, err = sdb.ShowMeta(); err != nil {
		return nil, nil, err
	}
	sdb.Close()

	if len(pids) == 0 {
		return nil, meta, nil
	}

	var posts []*Post
	_, err = Posts().Filter("Id__in", pids).RelatedSel().All(&posts)
	if err != nil {
		return nil, nil, err
	}

	return posts, meta, nil
}
