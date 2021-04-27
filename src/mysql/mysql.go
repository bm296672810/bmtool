package mysql

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var regexMysqlType string = "(bigint|binary|bit|blob|char|date|datetime|decimal|double|enum|float|tinyint|int|integer|linestring|longblob|longtext|text|tinytext|varchar|time|timestamp|year)"
var regexKeyName string = "(`|')?\\w+(`|')?"
var regexComment string = "COMMENT\\s+'[ \\S]*'"
var regexCreateEnd string = "\\)[ \\S]*;"

var suffix string = "Data"

func ProcessSql(sqlfile string, packageName string) error {
	data, err := os.ReadFile(sqlfile)
	if err != nil {
		log.Printf("ReadConfigure error:%v\n", err)
		log.Panicln("Read file error:", err)
	}

	sd := string(data)
	// remove the comments
	reg0 := regexp.MustCompile(regexComment)
	strs := reg0.FindAllString(sd, -1)
	for _, s := range strs {
		sd = strings.ReplaceAll(sd, s, "")
	}
	// ioutil.WriteFile("noComment.log", []byte(sd), 0755)
	reg1 := regexp.MustCompile(regexCreateEnd)
	ends := reg1.FindAllString(sd, -1)
	for _, s := range ends {
		sd = strings.ReplaceAll(sd, s, ");")
	}

	// ioutil.WriteFile("noComment.log", []byte(sd), 0755)

	reg2 := regexp.MustCompile("(create|CREATE)\\s+(table|TABLE)\\s+(`|')?\\w+(`|')?")

	is := reg2.FindAllStringIndex(sd, -1)
	log.Println(is)
	for i := 0; i < len(is); i++ {
		in := is[i]
		log.Println(sd[in[0]:in[1]])
	}

	parseSql(sd, is, packageName)
	return nil
}

func sqlTypeToGo(sqlType string) string {
	var r string
	switch sqlType {
	case "bigint":
		r = "int64"
	case "binary", "bit", "blob", "char", "linestring", "longblob", "longtext", "text", "tinytext", "varchar":
		r = "string"
	case "date", "datetime", "time", "timestamp", "year", "decimal":
		r = "string"
	case "double":
		r = "double"
	case "float":
		r = "float"
	case "enum", "tinyint", "int", "integer":
		r = "int"
	}

	return r
}

func firstToUpper(src string) string {
	if len(src) > 0 {
		result := strings.ToUpper(src[0:1]) + src[1:]
		return result
	}

	return ""
}

func underlineToUpper(src string) string {
	con := firstToUpper(src)
	index := strings.Index(con, "_")
	for index != -1 {
		if index == len(con)-1 {
			con = con[:index]
		} else {
			con = con[:index] + firstToUpper(con[index+1:])
		}
		index = strings.Index(con, "_")
	}

	return con
}

func parseSql(sql string, index [][]int, packageName string) error {
	dataCon := "package " + packageName + "\n" +
		`import (
		"database/sql"
		"errors"
		)` + "\n"
	funcCon := "\n"

	source := sql

	var head string
	var before []int
	for i := 0; i < len(index); i++ {
		in := index[i]
		thead := source[in[0]:in[1]]
		fun := func(head string, body string) (string, string) {
			var con, fcon string

			tabname := parseTableName(head)
			keys := parseKey(body)
			log.Println("tabname:", tabname)
			log.Println("keys:", keys)
			if len(tabname) <= 0 || keys == nil {
				return "", ""
			}
			/*
			   func AddArea(d *ud.AreaData) (int64, error) {
			   	if d == nil {
			   		ud.ELogger.Println("param nil")
			   		return 0, errors.New("param nil")
			   	}

			   	r, err := db.Exec("INSERT INTO parking_area (name,parking_id,type) VALUES(?,?,?)", d.Name, d.ParkingId, d.Type)
			   	if err != nil {
			   		ud.ELogger.Println("AddArea error:", err)
			   		return 0, err
			   	}

			   	id, err := r.LastInsertId()
			   	if err != nil {
			   		ud.ELogger.Println("LastInsertId error:", err)
			   		return 0, err
			   	}
			   	ud.ILogger.Println("AddArea success")

			   	return id, nil
			   }
			*/
			typeName := underlineToUpper(tabname) + suffix

			fcon = fcon + "func Insert" + underlineToUpper(tabname) + "(d " + typeName + ",db *sql.DB) (int64, error) {\n" +
				"if db == nil { return 0, errors.New(\"db unused!\")}\n" +
				"r, err := db.Exec(\"INSERT INTO " + tabname + " ("
			fconValues := " VALUES("
			fconParams := ","

			con = con + "type " + typeName + " struct {\n"
			for _, s := range keys {
				tp := sqlTypeToGo(s[1])
				con = con + underlineToUpper(s[0]) + " " + tp + " `json:\"" + s[0] + "," + tp + "\"`\n"

				if s[0] != "id" && s[0] != "create_time" && s[0] != "time_version" {
					fcon = fcon + s[0] + ","
					fconValues = fconValues + "?,"
					fconParams = fconParams + "d." + underlineToUpper(s[0]) + ","
				}
			}

			fcon = fcon[:len(fcon)-1] + ")"
			fconValues = fconValues[:len(fconValues)-1] + ");\""
			fconParams = fconParams[:len(fconParams)-1] + ")\n"
			fcon = fcon + fconValues + fconParams
			fcon = fcon + "if err != nil { return 0, err }\n"
			fcon = fcon + "id, err := r.LastInsertId()\n" +
				"if err != nil { return 0,err }\n return id, nil"
			fcon = fcon + "}\n"

			con = con + "}\n"
			// }

			return con, fcon
		}

		if i != len(index)-1 && len(before) > 0 {
			dc, fc := fun(head, source[before[1]:in[0]])
			dataCon = dataCon + dc
			funcCon = funcCon + fc
		} else if i == len(index)-1 {
			dc, fc := fun(head, source[before[1]:in[0]])
			dataCon = dataCon + dc
			funcCon = funcCon + fc

			subStr := source[in[1]:]
			ind := strings.Index(subStr, ";")
			c := subStr[:ind]
			log.Println("c:", c)

			dc, fc = fun(thead, c)
			dataCon = dataCon + dc
			funcCon = funcCon + fc
		}

		head = thead
		before = in
	}

	log.Println("content:", dataCon)
	filename := packageName + ".go"

	err := ioutil.WriteFile(filename, []byte(dataCon+funcCon), 0755)
	if err != nil {
		log.Println("write file error:", err)
		return err
	}

	func(filename string) {
		cmdStr := "gofmt -l -w " + filename
		cmd := exec.Command("/bin/sh", "-c", cmdStr)
		if _, err := cmd.Output(); err != nil {
			log.Panicln(err)
		}
	}(filename)

	return nil
}

func parseTableName(head string) string {
	reg := regexp.MustCompile("(create|CREATE)\\s+(table|TABLE)\\s+")
	st := reg.FindAllString(head, -1)
	if len(st) == 1 {
		s := strings.Replace(head, st[0], "", -1)
		s = strings.ReplaceAll(s, "`", "")
		s = strings.ReplaceAll(s, "'", "")

		return s
	}

	return ""
}

func parseKey(body string) [][]string {
	var result [][]string

	values := strings.Split(body, ",")

	for i := 0; i < len(values); i++ {
		sourceStr := values[i]
		lowerStr := strings.ToLower(sourceStr)

		regStr := regexKeyName + "\\s+" + regexMysqlType

		reg := regexp.MustCompile(regStr)
		v := reg.FindAllString(sourceStr, -1)
		var regResult string
		if v == nil {
			ins := reg.FindAllStringIndex(lowerStr, -1)
			if ins != nil {
				regResult = sourceStr[ins[0][0]:ins[0][1]]
			} else {
				return result
			}
		} else {
			regResult = v[0]
		}

		reg1 := regexp.MustCompile(regexKeyName)
		keyNames := reg1.FindAllString(regResult, -1)
		if keyNames == nil {
			log.Panicln("key find error!")
		}

		reg2 := regexp.MustCompile(regexMysqlType)
		regResult = strings.ToLower(regResult)
		types := reg2.FindAllString(regResult, -1)
		if types == nil {
			log.Panicln("type find error!")
		}

		var vs []string
		key := keyNames[0]
		tp := types[0]
		key = strings.ReplaceAll(key, "`", "")
		key = strings.ReplaceAll(key, "'", "")

		vs = append(vs, key, tp)
		result = append(result, vs)

		log.Println("regResult:", regResult)
	}

	return result
}
