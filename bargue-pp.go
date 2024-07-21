package main

/*
 * This is a "preprocessor", intended to be executed before
 * dtmpl(1), to help build the site more efficiently and decrease
 * data duplication.
 *
 * Read input/db/pages.json and input/db/groups.json; for each page:
 *	- if the page is a plate (name starts with "plate_"):
 *		- create & fill the corresponding input/db/imgs/$name.json
 *		- create & fill the corresponding input/$url.html.tmpl
 *	- If the page is a group of plates (is an index of groups.json):
 *		- create & fill the corresponding input/$url.html.tmpl
 */

import (
	"path"
	"log"
	"encoding/json"
	"path/filepath"
	"fmt"
	"strings"
	"os"
	"strconv"
	"text/template"
)

// input directory
var ind = "input/"

var dbImgsDir  = "db/imgs/"
var dbPagesFn  = "db/pages.json"
var dbGroupsFn = "db/groups.json"

type DB map[string]any

var dbPages  DB
var dbGroups DB

func isPlate(name string) bool {
	return strings.HasPrefix(name, "plate_")
}

func isGroup(name string) bool {
	_, ok := dbGroups[name]
	return ok
}

func createPath(fn string) (*os.File, error) {
	err := os.MkdirAll(filepath.Dir(fn), 0755)
	if err != nil {
		return nil, err
	}
	return os.Create(fn) //, os.O_RDWR|os.O_CREATE, 0644)
}

// Yo dawg.
var platePageTmpl = template.Must(template.New("").Delims("{{{", "}}}").Parse(""+
`{{< header "{{{ .lang }}}" >}}

{{< plate "{{{ .id }}}" "{{{ .lang }}}" >}}

{{< maybeparsefn "{{{ .url }}}/index.html.content" >}}

{{< footer "{{{ .lang }}}" >}}
`))

func mkPlatePage(id, lang, url string) error {
	fn := filepath.Join(ind, url, "index.html.tmpl")

//	var s strings.Builder
//	platePageTmpl.Execute(&s, map[string]any{

	f, err := createPath(fn)
	defer f.Close()
	if err != nil {
		return err
	}

	return platePageTmpl.Execute(f, map[string]any{
		"id"   : id,
		"url"  : url,
		"lang" : lang,
	})
//	println(s.String())
}

func mkPlatePages(id string, xpage any) error {
	// TODO/XXX why does .(DB) fails but this works?
	page, ok := xpage.(map[string]any)
	if !ok {
		panic("Invalid DB format for "+id)
	}
	for lang, xlpage := range page {
		lpage, ok := xlpage.(map[string]any)
		if !ok {
			panic("Invalid DB format (url) for "+id+", lang="+lang)
		}
		url, ok := lpage["url"].(string)
		if !ok {
			panic("Invalid DB format (url) for "+id+", lang="+lang)
		}
		if err := mkPlatePage(id, lang, url); err != nil {
			return err
		}
	}

	return nil
}

var plateImgTmpl = template.Must(template.New("").Parse(`{
	"file"   : "{{ .id }}.jpg",
	"author" : "charles-bargue",
	"descrs" : {
		"en" : "{{ .page.en.name }}",
		"fr" : "{{ .page.fr.name }}"
	},
	"urls"   : {
		"barguedrawing.wordpress.com" : "barguedrawing.wordpress.com-plate-I.{{ .n }}"
	}
}
`))

// id is a "plate_I-03_noses"; we want to fetch the "3"
func getNFromId(id string) int {
	n, err := strconv.Atoi(id[8:10])
	if err != nil {
		panic("Invalid plate ID: "+id)
	}
	return n
}

func mkPlateImg(id string) error {
	fn := filepath.Join(ind, dbImgsDir, id+".json")
	f, err := createPath(fn)
	defer f.Close()
	if err != nil {
		return err
	}
	return plateImgTmpl.Execute(f, map[string]any{
		"id"   : id,
		"page" : dbPages[id],
		"n"    : getNFromId(id),
	})
}

// Yo dawg.
var groupPageTmpl = template.Must(template.New("").Delims("{{{", "}}}").Parse(""+
`{{< header "{{{ .lang }}}" >}}

{{< group "{{{ .id }}}" "{{{ .lang }}}" >}}

{{< footer "{{{ .lang }}}" >}}
`))

func mkGroupPage(id, lang, url string) error {
	fn := filepath.Join(ind, url, "index.html.tmpl")
	f, err := createPath(fn)
	defer f.Close()
	if err != nil {
		return err
	}
	return groupPageTmpl.Execute(f, map[string]any{
		"id"   : id,
		"lang" : lang,
	})
}

func mkGroupPages(id string, xpage any) error {
	// TODO/XXX why does .(DB) fails but this works?
	page, ok := xpage.(map[string]any)
	if !ok {
		panic("Invalid DB format for "+id)
	}
	for lang, xlpage := range page {
		lpage, ok := xlpage.(map[string]any)
		if !ok {
			panic("Invalid DB format (url) for "+id+", lang="+lang)
		}
		url, ok := lpage["url"].(string)
		if !ok {
			panic("Invalid DB format (url) for "+id+", lang="+lang)
		}
		if err := mkGroupPage(id, lang, url); err != nil {
			return err
		}
	}

	return nil
}

func mkPages() error {
	for id := range dbPages {
		if isPlate(id) {
			if err := mkPlatePages(id, dbPages[id]); err != nil {
				return err
			}
			if err := mkPlateImg(id); err != nil {
				return err
			}
		} else if isGroup(id) {
			if err := mkGroupPages(id, dbPages[id]); err != nil {
				return err
			}
		}
	}
	return nil
}

func loadJSON(fn string) (DB, error) {
	var db DB
	raw, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(raw, &db)
	if err != nil {
		err = fmt.Errorf("%s: %s", fn, err)
	}
	return db, err
}

func loadDB(ind string) error {
	var err error
	dbPages, err = loadJSON(filepath.Join(ind, dbPagesFn))
	if err == nil {
		dbGroups, err = loadJSON(filepath.Join(ind, dbGroupsFn))
	}
	return err
}

func fails(err error) {
	argv0 := path.Base(os.Args[0])
	log.Fatal(argv0, ": ", err)
}

func init() {
	if err := loadDB(ind); err != nil {
		fails(err)
	}
}

func main() {
	if err := mkPages(); err != nil {
		fails(err)
	}
}
