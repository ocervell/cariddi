/*
==========
Cariddi
==========
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.
You should have received a copy of the GNU General Public License
along with this program.  If not, see http://www.gnu.org/licenses/.

	@Repository:  https://github.com/ocervell/cariddi

	@Author:      edoardottt, https://www.edoardoottavianelli.it

	@License: https://github.com/ocervell/cariddi/blob/main/LICENSE

*/

package scanner

// Error struct.
// ErrorName = the name that identifies the error.
// Regex = The regular expression to be matched.
type Error struct {
	ErrorName string
	Regex     []string
}

// ErrorMatched struct.
// Error = Error struct.
// Url = url in which the error is found.
// Match = the string matching the regex.
type ErrorMatched struct {
	Error Error
	URL   string
	Match string
}

// GetErrorRegexes returns all the error structs.
func GetErrorRegexes() []Error {
	var regexes = []Error{
		{
			"PHP error",
			[]string{
				`(?i)php (warning|error)`,
				`(?i)include_path`,
				`(?i)undefined (variable|index)`,
				`(?i)expect(s*) parameter [A-Za-z0-9-_]{1,30}`,
				`(?i)call to undefined method`,
				`(?i)failed to open stream`,
				`(?i)cannot modify header information`,
				`(?i)(syntax|parse|fatal) error`,
				`(?i)safe mode restriction in effect`,
				`(?i)uncaught exception`},
		},
		{
			"General error",
			[]string{
				`(?i)(fatal|critical|severe|high|medium|low) error`},
		},
		{
			"Debug information",
			[]string{
				`(?i)Debug trace`,
				`(?i)stack trace\:`},
		},
		{
			"MySQL error",
			[]string{
				`(?i)valid MySQL result`,
				`(?i)check the manual that (fits|corresponds to) your MySQL server version`,
				`(?i)MySQLSyntaxErrorException`,
				`(?i)MySqlException`,
				`(?i)MySql error`,
				`(?i)Unknown column `,
				`(?i)SQL syntax.*?MySQL`,
				`(?i)Warning.*?mysqli?`,
				`(?i)com\.mysql\.jdbc`,
				`(?i)Zend_Db_(Adapter|Statement)_Mysqli_Exception`,
				`(?i)Syntax error or access violation`},
		},
		{
			"MariaDB error",
			[]string{
				`(?i)check the manual that (fits|corresponds to) your MariaDB server version`,
				`(?i)MariaDB error`},
		},
		{
			"PostgreSQL error",
			[]string{
				`(?i)valid PostgreSQL result`,
				`(?i)PG::SyntaxError:`,
				`(?i)PSQLException`,
				`(?i)PostgreSQL query failed`,
				`(?i)ERROR: parser: parse error at or near`,
				`(?i)PostgreSQL error`,
				`(?i)PostgreSQL.*?ERROR`,
				`(?i)Warning.*?\\Wpg_`,
				`(?i)Npgsql\.`,
				`(?i)org\.postgresql\.util\.PSQLException`,
				`(?i)ERROR:\s\ssyntax error at or near`,
				`(?i)org\.postgresql\.jdbc`},
		},
		{
			"MSSQL error",
			[]string{
				`(?i)Microsoft SQL error`,
				`(?i)Microsoft SQL Native Client error`,
				`(?i)ODBC SQL Server Driver`,
				`(?i)Unclosed quotation mark after the character string`,
				`(?i)SQLServer JDBC Driver`,
				`(?i)Driver.*? SQL[\-\_\ ]*Server`,
				`(?i)OLE DB.*? SQL Server`,
				`(?i)\bSQL Server[^&lt;&quot;]+Driver`,
				`(?i)Warning.*?\W(mssql|sqlsrv)_`,
				`(?i)\bSQL Server[^&lt;&quot;]+[0-9a-fA-F]{8}`,
				`(?i)System\.Data\.SqlClient\.SqlException`,
				`(?is)Exception.*?\bRoadhouse\.Cms\.`,
				`(?i)\[SQL Server\]`,
				`(?i)ODBC Driver \d+ for SQL Server`,
				`(?i)com\.jnetdirect\.jsql`,
				`(?i)macromedia\.jdbc\.sqlserver`,
				`(?i)Zend_Db_(Adapter|Statement)_Sqlsrv_Exception`,
				`(?i)com\.microsoft\.sqlserver\.jdbc`,
				`(?i)SQL(Srv|Server)Exception`},
		},
		{
			"OracleDB error",
			[]string{
				`(?i)(\bORA-\d{5}|Oracle error|Oracle.*Driver|Warning.*\Woci_.*|Warning.*\Wora_.*)`},
		},
		{
			"IBMDB2 error",
			[]string{
				`(?i)(CLI Driver.*DB2|DB2 SQL error|\bdb2_\w+\(|SQLSTATE.+SQLCODE)`},
		},
		{
			"SQLite error",
			[]string{
				`(?i)(SQLite\/JDBCDriver|SQLite.Exception|System.Data.SQLite.SQLiteException` +
					`|Warning.*sqlite_.*|Warning.*SQLite3::|\[SQLITE_ERROR\])`},
		},
	}

	return regexes
}

// RemoveDuplicateErrors removes duplicates from Errors found.
func RemoveDuplicateErrors(input []ErrorMatched) []ErrorMatched {
	keys := make(map[string]bool)
	list := []ErrorMatched{}

	for _, entry := range input {
		if _, value := keys[entry.Match+entry.URL]; !value {
			keys[entry.Match+entry.URL] = true
			list = append(list, entry)
		}
	}

	return list
}
