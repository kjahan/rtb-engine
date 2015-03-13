package writer

import (
    "log"
    "fmt"
    _ "github.com/lib/pq"
    "database/sql"
)

func Write_pg(){
    db, err := sql.Open("postgres", "user=aioptify dbname=aioptify password=gocanadago sslmode=disable")
    if err != nil {
	log.Fatal("pq: ", err)
    }
    rows, err := db.Query(`SELECT COUNT(*) as cnt FROM ads`)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    for rows.Next() {
            var cnt int
            if err := rows.Scan(&cnt); err != nil {
                    log.Fatal(err)
            }
            fmt.Printf("cnt=%d\n", cnt)
    }
    if err := rows.Err(); err != nil {
            log.Fatal(err)
    }
}
