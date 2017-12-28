package caches

import "github.com/AzureRelease/boiler-server/models"

type BoilerRuntimeCacheHistory struct {
	models.MyIdObject

	Boiler			*models.Boiler			`orm:"rel(fk);index"`

	P1001			float32
	P1002			float32
	P1003			float32
	P1004			float32
	P1005			float32
	P1006			float32
	P1007			float32
	P1008			float32
	P1009			float32
	P1010			float32
	P1011			float32
	P1012			float32
	P1013			float32
	P1014			float32
	P1015			float32
	P1016			float32
	P1017			float32
	P1018			float32
	P1019			float32
	P1020			float32
	P1021			float32
	P1022			float32
	P1023			float32
	P1024			float32
	P1025			float32
	P1026			float32
	P1027			float32
	P1028			float32
	P1029			float32
	P1030			float32
	P1031			float32
	P1032			float32
	P1033			float32
	P1034			float32
	P1035			float32
	P1036			float32
	P1037			float32
	P1038			float32
	P1039			float32
	P1040			float32
	P1041			float32
	P1042			float32
	P1043			float32
	P1044			float32
	P1045			float32
	P1046			float32
	P1047			float32
	P1048			float32
	P1049			float32
	P1050			float32
	P1051			float32
	P1052			float32
	P1053			float32
	P1054			float32
	P1055			float32
	P1056			float32
	P1057			float32
	P1058			float32
	P1059			float32
	P1060			float32
	P1061			float32
	P1062			float32
	P1063			float32
	P1064			float32
	P1065			float32
	P1066			float32
	P1067			float32
	P1068			float32
	P1069			float32
	P1070			float32
	P1071			float32
	P1072			float32
	P1073			float32
	P1080			float32
	P1090			float32
	P1091			float32
	P1092			float32
	P1093			float32
	P1094			float32
	P1095			float32
	P1096			float32
	P1097			float32
	P1098			float32
	P1099			float32

	P1101			float32
	P1102			float32
	P1103			float32
	P1104			float32
	P1105			float32
	P1106			float32
	P1107			float32
	P1108			float32
	P1109			float32
	P1110			float32

	P1201			float32
	P1202			float32
	P1203			float32
	P1204			float32
	P1205			float32
	P1206			float32
	P1207			float32

	A1001			int
	A1002			int
	A1003			int
	A1004			int
	A1005			int
	A1006			int
	A1007			int
	A1008			int
	A1009			int
	A1010			int
	A1011			int
	A1012			int
	A1013			int
	A1014			int
	A1015			int
	A1016			int
	A1017			int
	A1018			int
	A1019			int
	A1020			int
	A1021			int
	A1022			int
	A1023			int
	A1024			int
	A1025			int
	A1026			int
	A1027			int
	A1028			int
	A1029			int
	A1030			int
	A1031			int
	A1032			int
	A1033			int
	A1034			int
	A1035			int
	A1036			int
	A1037			int
	A1038			int
	A1039			int
	A1040			int
	A1041			int
	A1042			int
	A1043			int
	A1044			int
	A1045			int
	A1046			int
	A1047			int
	A1048			int
	A1049			int
	A1050			int
	A1051			int
	A1052			int
	A1053			int
	A1054			int
	A1055			int
	A1056			int
	A1057			int
	A1058			int
	A1059			int
	A1060			int
	A1061			int
	A1062			int
	A1063			int
	A1064			int
	A1065			int
	A1066			int
	A1067			int
	A1068			int
	A1069			int
	A1070			int
	A1071			int
	A1072			int
	A1073			int
	A1080			int
	A1090			int
	A1091			int
	A1092			int
	A1093			int
	A1094			int
	A1095			int
	A1096			int
	A1097			int
	A1098			int
	A1099			int

	A1101			int
	A1102			int
	A1103			int
	A1104			int
	A1105			int
	A1106			int
	A1107			int
	A1108			int
	A1109			int
	A1110			int

	A1201			int
	A1202			int
	A1203			int
	A1204			int
	A1205			int
	A1206			int
	A1207			int
}

func (history *BoilerRuntimeCacheHistory) TableUnique() [][]string {
	return [][]string {
		[]string{"Boiler", "UpdatedDate"},
	}
}