package models

type BoilerM163 struct {
	TerminalCode			string		`orm:"column(Boiler_term_id);size(6)"`
	TerminalSetId			string		`orm:"column(Boiler_boiler_id);size(2)"`
	TerminalSystemTime		string		`orm:"column(Term_sys_time);size(19)"`
	DataFormatVersion		string		`orm:"column(Boiler_data_fmt_ver);size(2)"`
	SerialNumber			int32		`orm:"column(Boiler_sn)"`
}
