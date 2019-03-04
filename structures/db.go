package structures

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
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

func InitDB(user string, pass string, host string, database string) (*sqlx.DB, error) {
	return sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", user, pass, host, database))
}

func FetchGuild(db *sqlx.DB, id string) (*Guild, error) {
	var guild Guild
	err := db.Get(&guild, "SELECT * FROM guilds WHERE id=?;", id)
	if err != nil {
		return nil, err
	}
	return &guild, nil
}
