package discordcommands

import (
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"strings"
	"github.com/jmoiron/sqlx"
	"log"
)

const schema = `CREATE TABLE IF NOT EXISTS guild_%s (
	user_id VARCHAR(18) NOT NULL,
	mute_expires DATETIME NULL,
	experience BIGINT NOT NULL,
	PRIMARY KEY (user_id)
)`

const rolesSchema = `CREATE TABLE IF NOT EXISTS roles_%s (
	id VARCHAR(18) NOT NULL,
	experience BIGINT NOT NULL DEFAULT 0,
	PRIMARY KEY (id)
)`

func InitDB(user string, pass string, host string, database string) (*sqlx.DB, error) {
	return sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", user, pass, host, database))
}

func QueryAllData(db *sqlx.DB) (Guilds, error) {
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

func queryGuildData(tx *sqlx.Tx) (Guilds, error) {
	var length int
	err := tx.Get(&length, "SELECT COUNT(*) FROM guilds")
	if err != nil {
		return nil, err
	}
	xs := make([]*Guild, 0, length)
	err = tx.Select(&xs, "SELECT * FROM guilds")
	if err != nil {
		return nil, err
	}
	guilds := make(map[string]*Guild)
	for _, x := range xs {
		guilds[x.ID] = x
	}
	return guilds, nil
}

func queryMemberData(tx *sqlx.Tx, guildID string) (Members, error) {
	q := fmt.Sprintf("SELECT * FROM guild_%s", guildID)
	q2 := fmt.Sprintf("SELECT COUNT(*) FROM guild_%s", guildID)
	var length int
	err := tx.Get(&length, q2)
	if err != nil {
		return nil, err
	}
	xs := make([]*Member, 0, length)
	err = tx.Select(&xs, q)
	if err != nil {
		return nil, err
	}
	members := make(Members)
	for _, x := range xs {
		members[x.ID] = x
	}
	return members, nil
}

func queryRoleData(tx *sqlx.Tx, guildID string) (Roles, error) {
	q := fmt.Sprintf("SELECT * FROM roles_%s", guildID)
	q2 := fmt.Sprintf("SELECT COUNT(*) FROM roles_%s", guildID)
	var length int
	err := tx.Get(&length, q2)
	if err != nil {
		return nil, err
	}
	roles := make(Roles, 0, length)
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

func PostAllData(db *sqlx.DB, guilds Guilds) error {
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

func postGuildData(tx *sqlx.Tx, guilds Guilds) error {
	stmt, err := tx.PrepareNamed("INSERT INTO guilds VALUES (:id, :mute_role, :mod_role, :admin_role, :mod_log, :autorole, :exp_reload, :exp_gain_upper, :exp_gain_lower, :lottery_chance, :lottery_upper, :lottery_lower) " +
		"ON DUPLICATE KEY UPDATE mute_role=:mute_role, mod_role=:mod_role, admin_role=:admin_role, mod_log=:mod_log, autorole=:autorole, exp_reload=:exp_reload, " +
			"exp_gain_upper=:exp_gain_upper, exp_gain_lower=:exp_gain_lower, lottery_chance=:lottery_chance, lottery_upper=:lottery_upper, lottery_lower=:lottery_lower")
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

func postMemberData(tx *sqlx.Tx, guilds Guilds, guildID string) error {
	guild := guilds.Get(guildID)
	q := fmt.Sprintf("INSERT INTO guild_%s VALUES (:user_id, :mute_expires, :experience) "+
		"ON DUPLICATE KEY UPDATE user_id=:user_id, mute_expires=:mute_expires, experience=:experience", guildID)
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

func postRoleData(tx *sqlx.Tx, guilds Guilds, guildID string) error {
	guild := guilds.Get(guildID)
	q := fmt.Sprintf("INSERT INTO roles_%s VALUES (:id, :experience) "+
		"ON DUPLICATE KEY UPDATE id=:id, experience=:experience", guildID)
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
