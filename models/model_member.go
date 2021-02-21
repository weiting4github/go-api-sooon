package models

import "database/sql"

// NewMemberPrepare member.NewMember()使用
// "INSERT IGNORE INTO `sooon_db`.`member`(`email`, `pwd`, `salt`, `ip_field`, `ipv4v6`, `create_ts`) VALUES(?, ?, ?, ?, INET6_ATON(?), ?)"
func (m *DBManager) NewMemberPrepare() (*sql.Stmt, error) {
	stmt, err := m.DB.Prepare("INSERT IGNORE INTO `sooon_db`.`member`(`email`, `pwd`, `salt`, `ip_field`, `ipv4v6`, `create_ts`) VALUES(?, ?, ?, ?, INET6_ATON(?), ?)")
	return stmt, err
}

// SelMemberEmail member.Login()使用
// "SELECT `member_id`, `email`, `pwd`, `salt` FROM `sooon_db`.`member` WHERE `email` = ?"
func (m *DBManager) SelMemberEmail() (*sql.Stmt, error) {
	stmt, err := m.DB.Prepare("SELECT `member_id`, `email`, `pwd`, `salt` FROM `sooon_db`.`member` WHERE `email` = ?")
	return stmt, err
}

// NewMemberloginLog member.Login()使用
// "INSERT INTO `sooon_db`.`member_login_log`(`member_id`, `client_device`, `login_ts`) VALUES (?, ?, ?)"
func (m *DBManager) NewMemberloginLog() (*sql.Stmt, error) {
	stmt, err := m.DB.Prepare("INSERT INTO `sooon_db`.`member_login_log`(`member_id`, `client_device`, `login_ts`) VALUES (?, ?, ?)")
	return stmt, err
}

// SelMemberLoginLog member.Do()使用
//"SELECT * FROM `sooon_db`.`member_login_log` WHERE `member_id` = ?"
func (m *DBManager) SelMemberLoginLog() (*sql.Stmt, error) {
	stmt, err := m.DB.Prepare("SELECT * FROM `sooon_db`.`member_login_log` WHERE `member_id` = ?")
	return stmt, err
}
