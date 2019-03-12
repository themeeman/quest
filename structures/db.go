package structures

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/tomvanwoow/quest/modlog"
	"reflect"
	"strings"
)

const schema = `CREATE TABLE guilds (
	id VARCHAR(18) NOT NULL,
	mute_role VARCHAR(18) NULL DEFAULT NULL,
	mod_role VARCHAR(18) NULL DEFAULT NULL,
	admin_role VARCHAR(18) NULL DEFAULT NULL,
	mod_log VARCHAR(18) NULL DEFAULT NULL,
	autorole VARCHAR(18) NULL DEFAULT NULL,
	exp_reload SMALLINT(5) UNSIGNED NOT NULL DEFAULT '60',
	exp_gain_upper SMALLINT(6) UNSIGNED NOT NULL DEFAULT '0',
	exp_gain_lower SMALLINT(6) UNSIGNED NOT NULL DEFAULT '0',
	lottery_chance TINYINT(3) UNSIGNED NOT NULL DEFAULT '0',
	lottery_upper INT(10) UNSIGNED NOT NULL DEFAULT '0',
	lottery_lower INT(10) UNSIGNED NOT NULL DEFAULT '0',
	cases JSON NOT NULL,
	PRIMARY KEY (id),
	INDEX id (id)
);`

const membersSchema = `CREATE TABLE members (
	guild_id VARCHAR(18) NOT NULL,
	id VARCHAR(18) NOT NULL,
	mute_expires DATETIME NULL DEFAULT NULL,
	last_daily DATETIME NULL DEFAULT NULL,
	experience BIGINT(20) NOT NULL DEFAULT '0',
	chests JSON NOT NULL,
	PRIMARY KEY (guild_id, id),
	INDEX guild_id (guild_id),
	INDEX user_id (id),
	CONSTRAINT FK_members_guilds FOREIGN KEY (guild_id) REFERENCES guilds (id)
);`

const rolesSchema = `CREATE TABLE roles (
	guild_id VARCHAR(18) NOT NULL,
	id VARCHAR(18) NOT NULL,
	experience BIGINT(20) NOT NULL DEFAULT '0',
	PRIMARY KEY (guild_id, id),
	INDEX guild_id (guild_id),
	INDEX id (id),
	CONSTRAINT FK_roles_guilds FOREIGN KEY (guild_id) REFERENCES guilds (id)
);`

const guildsInsert = "INSERT INTO guilds VALUES (:id, :mute_role, :mod_role, :admin_role, :mod_log, :autorole, :exp_reload, :exp_gain_upper, :exp_gain_lower, :lottery_chance, :lottery_upper, :lottery_lower, :cases);"

type memberWrapper struct {
	GuildID string `db:"guild_id"`
	*Member
}

type roleWrapper struct {
	GuildID string `db:"guild_id"`
	*Role
}

func caseType(s string) modlog.Case {
	var a modlog.Case
	switch s {
	case "ban":
		a = &modlog.CaseBan{}
	case "kick":
		a = &modlog.CaseKick{}
	case "mute":
		a = &modlog.CaseMute{}
	case "purge":
		a = &modlog.CasePurge{}
	case "unban":
		a = &modlog.CaseUnban{}
	case "unmute":
		a = &modlog.CaseUnmute{}
	case "warn":
		a = &modlog.CaseWarn{}
	case "set":
		a = &modlog.CaseSet{}
	case "addexp":
		a = &modlog.CaseAddExp{}
	}
	return a
}

func CaseQuery(T string) string {
	structType := reflect.TypeOf(caseType(T)).Elem()
	if structType == reflect.TypeOf(nil) {
		return ""
	}
	xs := make([]string, 0, structType.NumField())
	for i := 0; i < structType.NumField(); i++ {
		if key, ok := structType.Field(i).Tag.Lookup("db"); ok {
			xs = append(xs, key)
		}
	}
	return fmt.Sprintf("SELECT (%s) FROM cases WHERE guild_id=? AND id=?", strings.Join(xs, ", "))
}

func InitDB(user string, pass string, host string, database string) (*sqlx.DB, error) {
	return sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", user, pass, host, database))
}

func FetchGuild(db *sqlx.DB, id string) (*Guild, error) {
	guild := NewGuild(id)
	err := db.Get(guild, "SELECT * FROM guilds WHERE id=?;", id)
	if err != nil {
		return nil, err
	}
	return guild, nil
}

func FetchMember(db *sqlx.DB, guildID string, id string) (*Member, error) {
	var member Member
	err := db.Get(&member, "SELECT * FROM members WHERE guild_id=? AND id=?;", guildID, id)
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func FetchRole(db *sqlx.DB, guildID string, id string) (*Role, error) {
	var role Role
	err := db.Get(&role, "SELECT * FROM roles WHERE guild_id=? AND id=?;", guildID, id)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func FetchCase(db *sqlx.DB, guildID string, id uint) (modlog.Case, error) {
	var T string
	err := db.Get(&T, "SELECT `type` FROM cases WHERE `guild_id`=? AND `id`=?;", guildID, id)
	if err != nil {
		return nil, err
	}
	c := caseType(T)
	err = db.Get(&c, CaseQuery(T), guildID, id)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func SaveGuild(db *sqlx.DB, guild *Guild) error {
	if guild == nil {
		return errors.New("Can't save nil guild")
	}
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Commit()
	_, _ = tx.Exec("DELETE FROM guilds WHERE id=?", guild.ID)
	stmt, err := tx.PrepareNamed(guildsInsert)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(guild)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func SaveMember(db *sqlx.DB, guildID string, member *Member) error {
	if member == nil {
		return errors.New("Can't save nil member")
	}
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Commit()
	_, _ = tx.Exec("DELETE FROM guilds WHERE guild_id=? AND id=?", guildID, member.ID)
	stmt, err := tx.PrepareNamed("INSERT INTO members VALUES (:guild_id, :id, :mute_expires, :last_daily, :experience, :chests)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(memberWrapper{GuildID: guildID, Member: member})
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func SaveRole(db *sqlx.DB, guildID string, role *Role) error {
	if role == nil {
		return errors.New("Can't save nil role")
	}
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Commit()
	_, _ = tx.Exec("DELETE FROM guilds WHERE guild_id=? AND id=?", guildID, role.ID)
	stmt, err := tx.PrepareNamed("INSERT INTO roles VALUES (:guild_id, :id, :mute_expires, :last_daily, :experience, :chests)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(roleWrapper{GuildID: guildID, Role: role})
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil

}
