package covidtracker

import "fmt"

var popByDep = map[string]int{
	"01":  655171,
	"02":  549587,
	"03":  349336,
	"04":  167331,
	"05":  146148,
	"06":  1098539,
	"07":  334591,
	"08":  283004,
	"09":  158205,
	"10":  316639,
	"11":  377580,
	"12":  289481,
	"13":  2047433,
	"14":  709715,
	"15":  151615,
	"16":  365697,
	"17":  660458,
	"18":  315100,
	"19":  249707,
	"2A":  156958,
	"2B":  179037,
	"21":  546466,
	"22":  618478,
	"23":  123500,
	"24":  426667,
	"25":  552619,
	"26":  522276,
	"27":  620046,
	"28":  445281,
	"29":  936432,
	"30":  757564,
	"31":  1373626,
	"32":  197851,
	"33":  1595903,
	"34":  1152125,
	"35":  1079333,
	"36":  229772,
	"37":  620671,
	"38":  1279514,
	"39":  270142,
	"40":  418200,
	"41":  343026,
	"42":  778211,
	"43":  234613,
	"44":  1415805,
	"45":  691942,
	"46":  179390,
	"47":  342358,
	"48":  80141,
	"49":  833602,
	"50":  516010,
	"51":  584108,
	"52":  183720,
	"53":  317742,
	"54":  747614,
	"55":  195047,
	"56":  769772,
	"57":  1064905,
	"58":  216182,
	"59":  2639070,
	"60":  842804,
	"61":  294421,
	"62":  1494330,
	"63":  667365,
	"64":  694279,
	"65":  235131,
	"66":  482567,
	"67":  1139258,
	"68":  777734,
	"69":  1864962,
	"70":  244305,
	"71":  572527,
	"72":  582211,
	"73":  442775,
	"74":  823928,
	"75":  2210875,
	"76":  1280803,
	"77":  1419206,
	"78":  1458275,
	"79":  385495,
	"80":  584797,
	"81":  397929,
	"82":  263125,
	"83":  1073201,
	"84":  570921,
	"85":  689496,
	"86":  447026,
	"87":  383215,
	"88":  382328,
	"89":  350970,
	"90":  147347,
	"91":  1305061,
	"92":  1622143,
	"93":  1616311,
	"94":  1389336,
	"95":  1237218,
	"971": 400170,
	"972": 382294,
	"973": 271829,
	"974": 862814,
}

func PopulationOfDepartment(dep string) (int, error) {
	if pop, ok := popByDep[dep]; ok {
		return pop, nil
	}
	return 0, fmt.Errorf("invalid department %q", dep)
}
