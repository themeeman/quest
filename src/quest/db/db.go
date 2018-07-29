package db

import (
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"strings"
	"github.com/jmoiron/sqlx"
	"log"
	"../structures"
	"reflect"
	"bytes"
)

const schema = `CREATE TABLE guild_%s (
user_id VARCHAR(18) NOT NULL,
mute_expires DATETIME NULL DEFAULT NULL,
last_daily DATETIME NULL DEFAULT NULL,
experience BIGINT(20) NOT NULL,
chests JSON NOT NULL,
PRIMARY KEY (user_id)
);`

const rolesSchema = `CREATE TABLE roles_%s (
	id VARCHAR(18) NOT NULL,
	experience BIGINT(20) NOT NULL DEFAULT '0',
	PRIMARY KEY (id)
);`

func InitDB(user string, pass string, host string, database string) (*sqlx.DB, error) {
	return sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", user, pass, host, database))
}

func QueryAllData(db *sqlx.DB) (structures.Guilds, error) {
	fmt.Println(createOverwriteQuery("?", structures.Role{}))
	fmt.Println(createOverwriteQuery("?", structures.Member{}))
	fmt.Println(createOverwriteQuery("?", structures.Guild{}))
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			tx.Rollback()
		}
	}()
	defer tx.Commit()
	guilds, err := queryGuildData(tx)
	if err != nil {
		return nil, err
	}
	for id, g := range guilds {
	start:
		members, err := queryMemberData(tx, id)
		if err != nil {
			fmt.Println(err)
			if strings.HasPrefix(err.Error(), "Error 1146") {
				CreateAllData(tx, id)
				goto start
			}
		}
		g.Members = members
		roles, err := queryRoleData(tx, id)
		if err != nil {
			fmt.Println(err)
			if strings.HasPrefix(err.Error(), "Error 1146") {
				CreateAllData(tx, id)
				goto start
			}
		}
		g.Roles = roles
	}
	return guilds, nil
}

func queryGuildData(tx *sqlx.Tx) (structures.Guilds, error) {
	var length int
	err := tx.Get(&length, "SELECT COUNT(*) FROM guilds")
	if err != nil {
		return nil, err
	}
	xs := make([]*structures.Guild, 0, length)
	err = tx.Select(&xs, "SELECT * FROM guilds")
	if err != nil {
		return nil, err
	}
	guilds := make(structures.Guilds)
	for _, x := range xs {
		guilds[x.ID] = x
	}
	return guilds, nil
}

func queryMemberData(tx *sqlx.Tx, guildID string) (structures.Members, error) {
	q := fmt.Sprintf("SELECT * FROM guild_%s", guildID)
	q2 := fmt.Sprintf("SELECT COUNT(*) FROM guild_%s", guildID)
	var length int
	err := tx.Get(&length, q2)
	if err != nil {
		return nil, err
	}
	xs := make([]*structures.Member, 0, length)
	err = tx.Select(&xs, q)
	if err != nil {
		return nil, err
	}
	members := make(structures.Members)
	for _, x := range xs {
		members[x.ID] = x
	}
	return members, nil
}

func queryRoleData(tx *sqlx.Tx, guildID string) (structures.Roles, error) {
	q := fmt.Sprintf("SELECT * FROM roles_%s", guildID)
	q2 := fmt.Sprintf("SELECT COUNT(*) FROM roles_%s", guildID)
	var length int
	err := tx.Get(&length, q2)
	if err != nil {
		return nil, err
	}
	roles := make(structures.Roles, 0, length)
	err = tx.Select(&roles, q)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func CreateAllData(tx *sqlx.Tx, guildID string) {
	tx.MustExec(`INSERT IGNORE INTO guilds (id) VALUES (?);`, guildID)
	q := fmt.Sprintf(schema, guildID)
	tx.MustExec(q)
	q = fmt.Sprintf(rolesSchema, guildID)
	tx.MustExec(q)
}

func PostAllData(db *sqlx.DB, guilds structures.Guilds) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Commit()
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			tx.Rollback()
			PostAllData(db, guilds)
		}
	}()
	err = postGuildData(tx, guilds)
	if err != nil {
		return err
	}
	for id := range guilds {
	start:
		err = postMemberData(tx, guilds, id)
		if err != nil {
			fmt.Println(err)
			if strings.HasPrefix(err.Error(), "Error 1146") {
				CreateAllData(tx, id)
				goto start
			}
		}
		err = postRoleData(tx, guilds, id)
		if err != nil {
			fmt.Println(err)
			if strings.HasPrefix(err.Error(), "Error 1146") {
				CreateAllData(tx, id)
				goto start
			}
		}
	}
	return nil
}

func postGuildData(tx *sqlx.Tx, guilds structures.Guilds) error {
	stmt, err := tx.PrepareNamed(createOverwriteQuery("guilds", structures.Guild{}))
	if err != nil {
		return err
	}
	for _, guild := range guilds {
		_, err = stmt.Exec(guild)
		if err != nil {
			return err
		}
	}
	return nil
}

func postMemberData(tx *sqlx.Tx, guilds structures.Guilds, guildID string) error {
	guild := guilds.Get(guildID)
	q := fmt.Sprintf(createOverwriteQuery("guild_%s", structures.Member{}), guildID)
	stmt, err := tx.PrepareNamed(q)
	if err != nil {
		return err
	}
	for _, member := range guild.Members {
		_, err = stmt.Exec(member)
		if err != nil {
			return err
		}
	}
	return nil
}

func postRoleData(tx *sqlx.Tx, guilds structures.Guilds, guildID string) error {
	guild := guilds.Get(guildID)
	q := fmt.Sprintf(createOverwriteQuery("roles_%s", structures.Role{}), guildID)
	stmt, err := tx.PrepareNamed(q)
	if err != nil {
		return err
	}
	for _, role := range guild.Roles {
		_, err = stmt.Exec(role)
		if err != nil {
			return err
		}
	}
	return nil
}

func createOverwriteQuery(table string, v interface{}) string {
	t := reflect.TypeOf(v)
	start := fmt.Sprintf("INSERT INTO %s VALUES ", table)
	cols := make([]string, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		if s, ok := t.Field(i).Tag.Lookup("db"); ok {
			cols = append(cols, s)
		}
	}
	ccols := make([]string, len(cols))
	for i, s := range cols {
		ccols[i] = ":" + s
	}
	values := fmt.Sprintf("(%s)", strings.Join(ccols, ", "))
	var end bytes.Buffer
	end.WriteString(" ON DUPLICATE KEY UPDATE ")
	sss := make([]string, len(cols))
	for i, s := range cols {
		sss[i] = fmt.Sprintf("%s=:%s", s, s)
	}
	end.WriteString(strings.Join(sss, ", "))
	return start + values + end.String() + ";"
}
