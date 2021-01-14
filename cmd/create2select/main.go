package main

import (
	//"fmt"

	"log"
	"pgutil"
	// "github.com/koplec/pgutil"
)

func main() {
	log.Printf("BEGIN create2select")
	// sql := `
	// CREATE TABLE public.orders (
	// 	id integer NOT NULL,
	// 	donor_id integer,
	// 	ordered_at timestamp with time zone DEFAULT now(),
	// );
	// `

	sql := `CREATE TABLE public.orders (
		id integer NOT NULL,
		donor_id integer,
		ordered_at timestamp with time zone DEFAULT now(),
		order_no character varying(255),
		canceled_at timestamp with time zone,
		canceled_by_id integer,
		schedule_date date,
		ordered_by_code character varying(255),
		ordered_by_name character varying(255),
		order_section_name character varying(255),
		infection text,
		created_at timestamp with time zone DEFAULT now(),
		updated_at timestamp with time zone DEFAULT now(),
		version_no integer DEFAULT 0,
		nyugai_kbn integer,
		byoto_code character varying(16),
		byoto_name character varying(32),
		note text,
		order_type integer
	);`
	log.Printf("sql:%s", sql)
	tokens, err := pgutil.ParseToken(sql)
	if err != nil {
		log.Fatalf("ptoken failed:%v\n", err)
	}
	log.Printf("debug: tokens:%v len(tokens):%v\n", tokens, len(tokens))
	for i, t := range tokens {
		log.Printf("%d:%s\n", i, t)
	}

	table, err := pgutil.ParseCreateTokens(tokens)
	if err != nil {
		log.Fatalf("parseCreateStmt failed : %v\n", err)
	}
	log.Printf("table:%v\n", table)

	log.Printf("select stmt")
	log.Printf("%s", table.CreateSelectStmt())
}
