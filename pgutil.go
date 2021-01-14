package pgutil

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

type Table struct {
	TableName string
	Columns   []Column
}

func (t *Table) NormalizeName() (str string) {
	name := t.TableName
	ary := strings.Split(name, ".")
	return ary[len(ary)-1]
}

func (t *Table) CreateSelectStmt() (sql string) {
	tblName := t.NormalizeName()
	columns := t.Columns
	columnStrs := make([]string, 0)
	for _, c := range columns {
		columnStrs = append(columnStrs, fmt.Sprintf("%s.%s as %s_%s", tblName, c.Name, tblName, c.Name))
	}

	return fmt.Sprintf(`SELECT
%s
FROM %s`, strings.Join(columnStrs, ",\n"), tblName)
}

type Column struct {
	Name string
	Type string
}

const (
	TOKEN_CREATE      = "create"
	TOKEN_TABLE       = "table"
	TOKEN_INT         = "integer"
	TOKEN_LEFT_PAREN  = "("
	TOKEN_RIGHT_PAREN = ")"
	TOKEN_COMMA       = ","
)

/**
 *  create文をtokenに分ける
 */
func ParseToken(sql string) (tokens []string, err error) {
	sql = strings.ToLower(sql)
	sql = strings.Replace(sql, "\n", "", -1)
	sql = strings.Replace(sql, "\t", "", -1)

	_tokens := strings.Split(sql, " ")
	for _, t := range _tokens {
		log.Printf("debug t:%s\n", t)
		//tに"(", ","が含まれていたら分ける
		if strings.ContainsAny(t, "(,)") {
			//1文字ずつ見る
			start := 0
			for pos, c := range t {
				if c == '(' || c == ',' || c == ')' || c == ';' {
					if start < pos {
						log.Printf("t[start:pos]:%s\n", t[start:pos])
						tokens = append(tokens, t[start:pos])
					}
					start = pos + 1
					log.Printf("string([]rune{c}):%s\n", string([]rune{c}))
					tokens = append(tokens, string([]rune{c}))
				}
			}
			//終わりの処理
			if start < len(t) {
				end := t[start:]
				log.Printf("end:%s\n", end)
				tokens = append(tokens, end)
			}
		} else {
			tokens = append(tokens, t)
		}

	}
	return tokens, nil
}

const (
	CreateStmtParserState_START                    = iota
	CreateStmtParserState_FIND_TABLE_TOKEN         //"table"トークンを探している状態
	CreateStmtParserState_PARSING_TABLE_NAME_TOKEN //テーブルの名称を探している状態
	CreateStmtParserState_FIND_COLUMN_START        //カラムの開始を探している状態
	CreateStmtParserState_PARSING_COLUMNS_TOKEN
	CreateStmtParserState_END
)

func ParseCreateTokens(tokens []string) (table *Table, err error) {
	table = &Table{}
	columns := make([]Column, 0)
	state := CreateStmtParserState_START

	idx := -1 //0startだとidx++で最初のtokenを無視してしまう
	// PARSE_LOOP:
	for idx < len(tokens) {
		idx++
		t := tokens[idx]

		if state == CreateStmtParserState_START {
			if t != TOKEN_CREATE {
				err = errors.New(fmt.Sprintf("not create tokens:%s\n", t))
				return nil, err
			}
			state = CreateStmtParserState_FIND_TABLE_TOKEN
			continue
		}
		if state == CreateStmtParserState_FIND_TABLE_TOKEN {
			//"table"トークンか確認
			//次はテーブル名
			state = CreateStmtParserState_PARSING_TABLE_NAME_TOKEN
			continue
		}
		//テーブル名
		if state == CreateStmtParserState_PARSING_TABLE_NAME_TOKEN {
			table.TableName = t
			state = CreateStmtParserState_FIND_COLUMN_START
			continue
		}
		//カラムを探す状態
		if state == CreateStmtParserState_FIND_COLUMN_START {
			if t == TOKEN_LEFT_PAREN {
				state = CreateStmtParserState_PARSING_COLUMNS_TOKEN
				continue
			}
		}
		//カラム解析中
		if state == CreateStmtParserState_PARSING_COLUMNS_TOKEN {
			parenCount := 0
		COLUMN_LOOP:
			for idx < len(tokens) {
				columnName := tokens[idx]
				idx++
				columnType := tokens[idx]
				column := Column{
					Name: columnName, Type: columnType,
				}
				columns = append(columns, column)

				for idx < len(tokens) {
					idx++
					t := tokens[idx]
					if t == TOKEN_COMMA {
						idx++
						continue COLUMN_LOOP
					}
					if t == TOKEN_LEFT_PAREN {
						parenCount++
						continue
					}
					if t == TOKEN_RIGHT_PAREN && parenCount > 0 {
						parenCount--
						continue
					}
					if t == TOKEN_RIGHT_PAREN && parenCount == 0 {
						state = CreateStmtParserState_END
						idx++
						break COLUMN_LOOP
					}
				}
			}
		}
		if state == CreateStmtParserState_END {
			table.Columns = columns
			return table, nil
		}
		err = errors.New(fmt.Sprintf("invalid parse current token:%s\n", t))
		return
	}
	err = errors.New("never reached")
	return
}
