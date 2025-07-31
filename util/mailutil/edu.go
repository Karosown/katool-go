package verify

import (
	"regexp"
	"strings"
)

var MicroStrategyEducationDomains = map[string][]string{
	// 欧洲微型国家
	"VA": {".va"},               // 梵蒂冈城国
	"MC": {".edu.mc", ".ac.mc"}, // 摩纳哥
	"SM": {".edu.sm", ".ac.sm"}, // 圣马力诺
	"AD": {".edu.ad", ".ac.ad"}, // 安道尔 [ref:7]()
	"LI": {".edu.li", ".ac.li"}, // 列支敦士登
	"MT": {".edu.mt", ".ac.mt"}, // 马耳他

	// 太平洋微型国家
	"NR": {".edu.nr", ".ac.nr"}, // 瑙鲁
	"TV": {".edu.tv", ".ac.tv"}, // 图瓦卢
	"PW": {".edu.pw", ".ac.pw"}, // 帕劳
	"MH": {".edu.mh", ".ac.mh"}, // 马绍尔群岛
	"FM": {".edu.fm", ".ac.fm"}, // 密克罗尼西亚联邦
	"KI": {".edu.ki", ".ac.ki"}, // 基里巴斯
	"NU": {".edu.nu", ".ac.nu"}, // 纽埃
	"CK": {".edu.ck", ".ac.ck"}, // 库克群岛
	"TK": {".edu.tk", ".ac.tk"}, // 托克劳

	// 加勒比海小国
	"KN": {".edu.kn", ".ac.kn"}, // 圣基茨和尼维斯
	"LC": {".edu.lc", ".ac.lc"}, // 圣卢西亚
	"VC": {".edu.vc", ".ac.vc"}, // 圣文森特和格林纳丁斯
	"GD": {".edu.gd", ".ac.gd"}, // 格林纳达
	"AG": {".edu.ag", ".ac.ag"}, // 安提瓜和巴布达
	"DM": {".edu.dm", ".ac.dm"}, // 多米尼克
	"BS": {".edu.bs", ".ac.bs"}, // 巴哈马
	"BB": {".edu.bb", ".ac.bb"}, // 巴巴多斯
	"JM": {".edu.jm", ".ac.jm"}, // 牙买加
	"TT": {".edu.tt", ".ac.tt"}, // 特立尼达和多巴哥
}
var AfricanSmallStatesEducationDomains = map[string][]string{
	// 西非小国
	"GM": {".edu.gm", ".ac.gm"}, // 冈比亚（您提到的）
	"GW": {".edu.gw", ".ac.gw"}, // 几内亚比绍
	"SL": {".edu.sl", ".ac.sl"}, // 塞拉利昂
	"LR": {".edu.lr", ".ac.lr"}, // 利比里亚
	"CV": {".edu.cv", ".ac.cv"}, // 佛得角
	"ST": {".edu.st", ".ac.st"}, // 圣多美和普林西比
	"GQ": {".edu.gq", ".ac.gq"}, // 赤道几内亚

	// 东非小国
	"DJ": {".edu.dj", ".ac.dj"}, // 吉布提
	"ER": {".edu.er", ".ac.er"}, // 厄立特里亚
	"BI": {".edu.bi", ".ac.bi"}, // 布隆迪
	"RW": {".edu.rw", ".ac.rw"}, // 卢旺达
	"SC": {".edu.sc", ".ac.sc"}, // 塞舌尔
	"KM": {".edu.km", ".ac.km"}, // 科摩罗
	"MU": {".edu.mu", ".ac.mu"}, // 毛里求斯
	"MV": {".edu.mv", ".ac.mv"}, // 马尔代夫

	// 南部非洲小国
	"SZ": {".edu.sz", ".ac.sz"}, // 斯威士兰
	"LS": {".edu.ls", ".ac.ls"}, // 莱索托
	"BW": {".edu.bw", ".ac.bw"}, // 博茨瓦纳
	"NA": {".edu.na", ".ac.na"}, // 纳米比亚
}
var AsianSmallStatesEducationDomains = map[string][]string{
	// 东南亚小国
	"BN": {".edu.bn", ".ac.bn"}, // 文莱
	"TL": {".edu.tl", ".ac.tl"}, // 东帝汶

	// 南亚小国
	"BT": {".edu.bt", ".ac.bt"}, // 不丹
	"MV": {".edu.mv", ".ac.mv"}, // 马尔代夫

	// 中亚小国
	"TJ": {".edu.tj", ".ac.tj"}, // 塔吉克斯坦
	"KG": {".edu.kg", ".ac.kg"}, // 吉尔吉斯斯坦
	"TM": {".edu.tm", ".ac.tm"}, // 土库曼斯坦

	// 特殊地区
	"HK": {".edu.hk", ".ac.hk"}, // 香港
	"MO": {".edu.mo", ".ac.mo"}, // 澳门
	"TW": {".edu.tw", ".ac.tw"}, // 台湾
}
var EuropeanSmallStatesEducationDomains = map[string][]string{
	// 波罗的海小国
	"EE": {".edu.ee", ".ac.ee"}, // 爱沙尼亚
	"LV": {".edu.lv", ".ac.lv"}, // 拉脱维亚
	"LT": {".edu.lt", ".ac.lt"}, // 立陶宛

	// 巴尔干小国
	"ME": {".edu.me", ".ac.me"}, // 黑山
	"MK": {".edu.mk", ".ac.mk"}, // 北马其顿
	"AL": {".edu.al", ".ac.al"}, // 阿尔巴尼亚
	"BA": {".edu.ba", ".ac.ba"}, // 波斯尼亚和黑塞哥维那
	"XK": {".edu.xk", ".ac.xk"}, // 科索沃
	"MD": {".edu.md", ".ac.md"}, // 摩尔多瓦

	// 其他欧洲小国
	"IS": {".edu.is", ".ac.is"}, // 冰岛
	"LU": {".edu.lu", ".ac.lu"}, // 卢森堡
	"CY": {".edu.cy", ".ac.cy"}, // 塞浦路斯
}
var SpecialTerritoriesEducationDomains = map[string][]string{
	// 英国海外领土
	"AI": {".edu.ai", ".ac.ai"}, // 安圭拉
	"BM": {".edu.bm", ".ac.bm"}, // 百慕大
	"VG": {".edu.vg", ".ac.vg"}, // 英属维尔京群岛
	"KY": {".edu.ky", ".ac.ky"}, // 开曼群岛
	"MS": {".edu.ms", ".ac.ms"}, // 蒙特塞拉特
	"TC": {".edu.tc", ".ac.tc"}, // 特克斯和凯科斯群岛
	"SH": {".edu.sh", ".ac.sh"}, // 圣赫勒拿
	"FK": {".edu.fk", ".ac.fk"}, // 福克兰群岛
	"GI": {".edu.gi", ".ac.gi"}, // 直布罗陀
	"IO": {".edu.io", ".ac.io"}, // 英属印度洋领土
	"PN": {".edu.pn", ".ac.pn"}, // 皮特凯恩群岛

	// 法国海外领土
	"NC": {".edu.nc", ".ac.nc"}, // 新喀里多尼亚
	"PF": {".edu.pf", ".ac.pf"}, // 法属波利尼西亚
	"WF": {".edu.wf", ".ac.wf"}, // 瓦利斯和富图纳
	"MQ": {".edu.mq", ".ac.mq"}, // 马提尼克
	"GP": {".edu.gp", ".ac.gp"}, // 瓜德罗普
	"GF": {".edu.gf", ".ac.gf"}, // 法属圭亚那
	"RE": {".edu.re", ".ac.re"}, // 留尼汪
	"YT": {".edu.yt", ".ac.yt"}, // 马约特
	"PM": {".edu.pm", ".ac.pm"}, // 圣皮埃尔和密克隆
	"TF": {".edu.tf", ".ac.tf"}, // 法属南部和南极洲领土

	// 美国领土
	"AS": {".edu.as", ".ac.as"}, // 美属萨摩亚
	"GU": {".edu.gu", ".ac.gu"}, // 关岛
	"MP": {".edu.mp", ".ac.mp"}, // 北马里亚纳群岛
	"PR": {".edu.pr", ".ac.pr"}, // 波多黎各
	"VI": {".edu.vi", ".ac.vi"}, // 美属维尔京群岛
	"UM": {".edu.um", ".ac.um"}, // 美国本土外小岛屿

	// 荷兰领土
	"AW": {".edu.aw", ".ac.aw"}, // 阿鲁巴
	"CW": {".edu.cw", ".ac.cw"}, // 库拉索
	"SX": {".edu.sx", ".ac.sx"}, // 荷属圣马丁
	"BQ": {".edu.bq", ".ac.bq"}, // 荷兰加勒比区

	// 丹麦领土
	"FO": {".edu.fo", ".ac.fo"}, // 法罗群岛
	"GL": {".edu.gl", ".ac.gl"}, // 格陵兰

	// 其他特殊地区
	"AC": {".edu.ac", ".ac.ac"}, // 阿森松岛
	"TA": {".edu.ta", ".ac.ta"}, // 特里斯坦-达库尼亚
	"CC": {".edu.cc", ".ac.cc"}, // 科科斯群岛
	"CX": {".edu.cx", ".ac.cx"}, // 圣诞岛
	"HM": {".edu.hm", ".ac.hm"}, // 赫德岛和麦克唐纳群岛
	"NF": {".edu.nf", ".ac.nf"}, // 诺福克岛
}

type EduEmailVerify struct {
	// 所有域名分类
	majorCountries      map[string][]string
	microstates         map[string][]string
	africanSmallStates  map[string][]string
	asianSmallStates    map[string][]string
	europeanSmallStates map[string][]string
	specialTerritories  map[string][]string

	// 编译后的正则表达式
	allPatterns     []*regexp.Regexp
	countryPatterns map[string][]*regexp.Regexp
}

// GlobalEducationDomains - 全球教育邮箱域名配置（完整版）
var GlobalEducationDomains = map[string][]string{
	// 北美洲
	"US": {".edu"},              // 美国
	"CA": {".ca", ".edu.ca"},    // 加拿大
	"MX": {".edu.mx"},           // 墨西哥
	"GT": {".edu.gt", ".ac.gt"}, // 危地马拉
	"BZ": {".edu.bz", ".ac.bz"}, // 伯利兹
	"SV": {".edu.sv", ".ac.sv"}, // 萨尔瓦多
	"HN": {".edu.hn", ".ac.hn"}, // 洪都拉斯
	"NI": {".edu.ni", ".ac.ni"}, // 尼加拉瓜
	"CR": {".edu.cr", ".ac.cr"}, // 哥斯达黎加
	"PA": {".edu.pa", ".ac.pa"}, // 巴拿马
	"CU": {".edu.cu", ".ac.cu"}, // 古巴
	"DO": {".edu.do", ".ac.do"}, // 多米尼加共和国
	"HT": {".edu.ht", ".ac.ht"}, // 海地

	// 欧洲主要国家
	"UK": {".ac.uk", ".edu.uk"},                      // 英国
	"DE": {".edu.de", ".uni.de", ".fh.de", ".tu.de"}, // 德国
	"FR": {".edu.fr", ".univ.fr", ".ens.fr"},         // 法国
	"IT": {".edu.it", ".univ.it"},                    // 意大利
	"ES": {".edu.es", ".univ.es"},                    // 西班牙
	"NL": {".edu.nl", ".uni.nl"},                     // 荷兰
	"SE": {".edu.se", ".uni.se"},                     // 瑞典
	"NO": {".edu.no", ".uni.no"},                     // 挪威
	"FI": {".edu.fi", ".uni.fi"},                     // 芬兰
	"DK": {".edu.dk", ".uni.dk"},                     // 丹麦
	"CH": {".edu.ch", ".ethz.ch", ".epfl.ch"},        // 瑞士
	"AT": {".edu.at", ".ac.at", ".uni.at"},           // 奥地利
	"BE": {".edu.be", ".ac.be", ".uni.be"},           // 比利时
	"PT": {".edu.pt", ".univ.pt"},                    // 葡萄牙
	"IE": {".edu.ie", ".ac.ie"},                      // 爱尔兰
	"GR": {".edu.gr", ".ac.gr"},                      // 希腊
	"PL": {".edu.pl", ".ac.pl"},                      // 波兰
	"CZ": {".edu.cz", ".ac.cz"},                      // 捷克 [ref:11]()
	"RU": {".edu.ru", ".ac.ru"},                      // 俄罗斯
	"UA": {".edu.ua", ".ac.ua"},                      // 乌克兰
	"BY": {".edu.by", ".ac.by"},                      // 白俄罗斯
	"SK": {".edu.sk", ".ac.sk"},                      // 斯洛伐克
	"HU": {".edu.hu", ".ac.hu"},                      // 匈牙利
	"RO": {".edu.ro", ".ac.ro"},                      // 罗马尼亚
	"BG": {".edu.bg", ".ac.bg"},                      // 保加利亚
	"HR": {".edu.hr", ".ac.hr"},                      // 克罗地亚
	"SI": {".edu.si", ".ac.si"},                      // 斯洛文尼亚
	"RS": {".edu.rs", ".ac.rs"},                      // 塞尔维亚
	"ME": {".edu.me", ".ac.me"},                      // 黑山
	"BA": {".edu.ba", ".ac.ba"},                      // 波斯尼亚和黑塞哥维那
	"MK": {".edu.mk", ".ac.mk"},                      // 北马其顿
	"AL": {".edu.al", ".ac.al"},                      // 阿尔巴尼亚
	"XK": {".edu.xk", ".ac.xk"},                      // 科索沃
	"MD": {".edu.md", ".ac.md"},                      // 摩尔多瓦
	"EE": {".edu.ee", ".ac.ee"},                      // 爱沙尼亚
	"LV": {".edu.lv", ".ac.lv"},                      // 拉脱维亚
	"LT": {".edu.lt", ".ac.lt"},                      // 立陶宛
	"IS": {".edu.is", ".ac.is"},                      // 冰岛
	"LU": {".edu.lu", ".ac.lu"},                      // 卢森堡
	"CY": {".edu.cy", ".ac.cy"},                      // 塞浦路斯
	"MT": {".edu.mt", ".ac.mt"},                      // 马耳他
	"MC": {".edu.mc", ".ac.mc"},                      // 摩纳哥
	"SM": {".edu.sm", ".ac.sm"},                      // 圣马力诺
	"AD": {".edu.ad", ".ac.ad"},                      // 安道尔
	"LI": {".edu.li", ".ac.li"},                      // 列支敦士登
	"VA": {".edu.va", ".ac.va"},                      // 梵蒂冈

	// 亚太地区
	"CN": {".edu.cn", ".ac.cn"},                     // 中国
	"JP": {".ac.jp", ".ed.jp"},                      // 日本
	"KR": {".ac.kr", ".edu.kr"},                     // 韩国
	"KP": {".edu.kp", ".ac.kp"},                     // 朝鲜
	"IN": {".ac.in", ".edu.in"},                     // 印度
	"PK": {".edu.pk", ".ac.pk"},                     // 巴基斯坦
	"BD": {".edu.bd", ".ac.bd"},                     // 孟加拉国
	"LK": {".edu.lk", ".ac.lk"},                     // 斯里兰卡
	"NP": {".edu.np", ".ac.np"},                     // 尼泊尔
	"BT": {".edu.bt", ".ac.bt"},                     // 不丹
	"MV": {".edu.mv", ".ac.mv"},                     // 马尔代夫
	"AF": {".edu.af", ".ac.af"},                     // 阿富汗
	"AU": {".edu.au", ".ac.au"},                     // 澳大利亚
	"NZ": {".ac.nz", ".edu.nz"},                     // 新西兰
	"SG": {".edu.sg", ".nus.edu.sg", ".ntu.edu.sg"}, // 新加坡
	"MY": {".edu.my", ".ac.my"},                     // 马来西亚
	"TH": {".ac.th", ".edu.th"},                     // 泰国
	"ID": {".ac.id", ".edu.id"},                     // 印度尼西亚
	"PH": {".edu.ph", ".ac.ph"},                     // 菲律宾
	"VN": {".edu.vn", ".ac.vn"},                     // 越南
	"KH": {".edu.kh", ".ac.kh"},                     // 柬埔寨
	"LA": {".edu.la", ".ac.la"},                     // 老挝
	"MM": {".edu.mm", ".ac.mm"},                     // 缅甸
	"BN": {".edu.bn", ".ac.bn"},                     // 文莱
	"TL": {".edu.tl", ".ac.tl"},                     // 东帝汶
	"HK": {".edu.hk", ".ac.hk"},                     // 香港
	"TW": {".edu.tw", ".ac.tw"},                     // 台湾
	"MO": {".edu.mo", ".ac.mo"},                     // 澳门
	"MN": {".edu.mn", ".ac.mn"},                     // 蒙古
	"KZ": {".edu.kz", ".ac.kz"},                     // 哈萨克斯坦
	"KG": {".edu.kg", ".ac.kg"},                     // 吉尔吉斯斯坦
	"TJ": {".edu.tj", ".ac.tj"},                     // 塔吉克斯坦
	"UZ": {".edu.uz", ".ac.uz"},                     // 乌兹别克斯坦
	"TM": {".edu.tm", ".ac.tm"},                     // 土库曼斯坦
	"GE": {".edu.ge", ".ac.ge"},                     // 格鲁吉亚
	"AM": {".edu.am", ".ac.am"},                     // 亚美尼亚
	"AZ": {".edu.az", ".ac.az"},                     // 阿塞拜疆

	// 太平洋岛国
	"FJ": {".edu.fj", ".ac.fj"}, // 斐济
	"PG": {".edu.pg", ".ac.pg"}, // 巴布亚新几内亚
	"SB": {".edu.sb", ".ac.sb"}, // 所罗门群岛
	"VU": {".edu.vu", ".ac.vu"}, // 瓦努阿图
	"NC": {".edu.nc", ".ac.nc"}, // 新喀里多尼亚
	"PF": {".edu.pf", ".ac.pf"}, // 法属波利尼西亚
	"TO": {".edu.to", ".ac.to"}, // 汤加
	"WS": {".edu.ws", ".ac.ws"}, // 萨摩亚
	"TV": {".edu.tv", ".ac.tv"}, // 图瓦卢
	"KI": {".edu.ki", ".ac.ki"}, // 基里巴斯
	"NR": {".edu.nr", ".ac.nr"}, // 瑙鲁
	"MH": {".edu.mh", ".ac.mh"}, // 马绍尔群岛
	"FM": {".edu.fm", ".ac.fm"}, // 密克罗尼西亚联邦
	"PW": {".edu.pw", ".ac.pw"}, // 帕劳
	"CK": {".edu.ck", ".ac.ck"}, // 库克群岛
	"NU": {".edu.nu", ".ac.nu"}, // 纽埃
	"TK": {".edu.tk", ".ac.tk"}, // 托克劳

	// 中东
	"IL": {".ac.il", ".edu.il"}, // 以色列
	"PS": {".edu.ps", ".ac.ps"}, // 巴勒斯坦
	"TR": {".edu.tr", ".ac.tr"}, // 土耳其
	"SA": {".edu.sa", ".ac.sa"}, // 沙特阿拉伯
	"AE": {".ac.ae", ".edu.ae"}, // 阿联酋
	"JO": {".edu.jo", ".ac.jo"}, // 约旦
	"LB": {".edu.lb", ".ac.lb"}, // 黎巴嫩
	"SY": {".edu.sy", ".ac.sy"}, // 叙利亚
	"IQ": {".edu.iq", ".ac.iq"}, // 伊拉克
	"IR": {".ac.ir", ".edu.ir"}, // 伊朗
	"KW": {".edu.kw", ".ac.kw"}, // 科威特
	"BH": {".edu.bh", ".ac.bh"}, // 巴林
	"QA": {".edu.qa", ".ac.qa"}, // 卡塔尔
	"OM": {".edu.om", ".ac.om"}, // 阿曼
	"YE": {".edu.ye", ".ac.ye"}, // 也门

	// 非洲
	"ZA": {".ac.za", ".edu.za"}, // 南非
	"EG": {".edu.eg", ".ac.eg"}, // 埃及
	"NG": {".edu.ng", ".ac.ng"}, // 尼日利亚
	"KE": {".ac.ke", ".edu.ke"}, // 肯尼亚
	"GH": {".edu.gh", ".ac.gh"}, // 加纳
	"MA": {".ac.ma", ".edu.ma"}, // 摩洛哥
	"TN": {".edu.tn", ".ac.tn"}, // 突尼斯
	"ET": {".edu.et", ".ac.et"}, // 埃塞俄比亚
	"DZ": {".edu.dz", ".ac.dz"}, // 阿尔及利亚
	"SD": {".edu.sd", ".ac.sd"}, // 苏丹
	"SS": {".edu.ss", ".ac.ss"}, // 南苏丹
	"UG": {".edu.ug", ".ac.ug"}, // 乌干达
	"TZ": {".edu.tz", ".ac.tz"}, // 坦桑尼亚
	"RW": {".edu.rw", ".ac.rw"}, // 卢旺达
	"BI": {".edu.bi", ".ac.bi"}, // 布隆迪
	"MW": {".edu.mw", ".ac.mw"}, // 马拉维
	"ZM": {".edu.zm", ".ac.zm"}, // 赞比亚
	"ZW": {".edu.zw", ".ac.zw"}, // 津巴布韦
	"BW": {".edu.bw", ".ac.bw"}, // 博茨瓦纳
	"NA": {".edu.na", ".ac.na"}, // 纳米比亚
	"SZ": {".edu.sz", ".ac.sz"}, // 斯威士兰
	"LS": {".edu.ls", ".ac.ls"}, // 莱索托
	"MZ": {".edu.mz", ".ac.mz"}, // 莫桑比克
	"AO": {".edu.ao", ".ac.ao"}, // 安哥拉
	"MG": {".edu.mg", ".ac.mg"}, // 马达加斯加
	"MU": {".edu.mu", ".ac.mu"}, // 毛里求斯
	"SC": {".edu.sc", ".ac.sc"}, // 塞舌尔
	"KM": {".edu.km", ".ac.km"}, // 科摩罗
	"DJ": {".edu.dj", ".ac.dj"}, // 吉布提
	"SO": {".edu.so", ".ac.so"}, // 索马里
	"ER": {".edu.er", ".ac.er"}, // 厄立特里亚
	"CI": {".edu.ci", ".ac.ci"}, // 科特迪瓦
	"BF": {".edu.bf", ".ac.bf"}, // 布基纳法索
	"ML": {".edu.ml", ".ac.ml"}, // 马里
	"NE": {".edu.ne", ".ac.ne"}, // 尼日尔
	"SN": {".edu.sn", ".ac.sn"}, // 塞内加尔
	"GM": {".edu.gm", ".ac.gm"}, // 冈比亚
	"GW": {".edu.gw", ".ac.gw"}, // 几内亚比绍
	"GN": {".edu.gn", ".ac.gn"}, // 几内亚
	"SL": {".edu.sl", ".ac.sl"}, // 塞拉利昂
	"LR": {".edu.lr", ".ac.lr"}, // 利比里亚
	"CV": {".edu.cv", ".ac.cv"}, // 佛得角
	"ST": {".edu.st", ".ac.st"}, // 圣多美和普林西比
	"GQ": {".edu.gq", ".ac.gq"}, // 赤道几内亚
	"GA": {".edu.ga", ".ac.ga"}, // 加蓬
	"CG": {".edu.cg", ".ac.cg"}, // 刚果共和国
	"CD": {".edu.cd", ".ac.cd"}, // 刚果民主共和国
	"CF": {".edu.cf", ".ac.cf"}, // 中非共和国
	"CM": {".edu.cm", ".ac.cm"}, // 喀麦隆
	"TD": {".edu.td", ".ac.td"}, // 乍得
	"LY": {".edu.ly", ".ac.ly"}, // 利比亚

	// 南美洲
	"BR": {".edu.br", ".ac.br"}, // 巴西
	"AR": {".edu.ar", ".ac.ar"}, // 阿根廷
	"CL": {".edu.cl", ".ac.cl"}, // 智利
	"CO": {".edu.co", ".ac.co"}, // 哥伦比亚
	"PE": {".edu.pe", ".ac.pe"}, // 秘鲁
	"VE": {".edu.ve", ".ac.ve"}, // 委内瑞拉
	"UY": {".edu.uy", ".ac.uy"}, // 乌拉圭
	"EC": {".edu.ec", ".ac.ec"}, // 厄瓜多尔
	"BO": {".edu.bo", ".ac.bo"}, // 玻利维亚
	"PY": {".edu.py", ".ac.py"}, // 巴拉圭
	"GY": {".edu.gy", ".ac.gy"}, // 圭亚那
	"SR": {".edu.sr", ".ac.sr"}, // 苏里南
	"GF": {".edu.gf", ".ac.gf"}, // 法属圭亚那
	"FK": {".edu.fk", ".ac.fk"}, // 福克兰群岛
	"GS": {".edu.gs", ".ac.gs"}, // 南乔治亚和南桑威奇群岛

	// 加勒比海地区
	"JM": {".edu.jm", ".ac.jm"}, // 牙买加
	"TT": {".edu.tt", ".ac.tt"}, // 特立尼达和多巴哥
	"BB": {".edu.bb", ".ac.bb"}, // 巴巴多斯
	"BS": {".edu.bs", ".ac.bs"}, // 巴哈马
	"AG": {".edu.ag", ".ac.ag"}, // 安提瓜和巴布达
	"DM": {".edu.dm", ".ac.dm"}, // 多米尼克
	"GD": {".edu.gd", ".ac.gd"}, // 格林纳达
	"KN": {".edu.kn", ".ac.kn"}, // 圣基茨和尼维斯
	"LC": {".edu.lc", ".ac.lc"}, // 圣卢西亚
	"VC": {".edu.vc", ".ac.vc"}, // 圣文森特和格林纳丁斯
	"AI": {".edu.ai", ".ac.ai"}, // 安圭拉
	"BM": {".edu.bm", ".ac.bm"}, // 百慕大
	"VG": {".edu.vg", ".ac.vg"}, // 英属维尔京群岛
	"VI": {".edu.vi", ".ac.vi"}, // 美属维尔京群岛
	"KY": {".edu.ky", ".ac.ky"}, // 开曼群岛
	"TC": {".edu.tc", ".ac.tc"}, // 特克斯和凯科斯群岛
	"MS": {".edu.ms", ".ac.ms"}, // 蒙特塞拉特
	"MQ": {".edu.mq", ".ac.mq"}, // 马提尼克
	"GP": {".edu.gp", ".ac.gp"}, // 瓜德罗普
	"AW": {".edu.aw", ".ac.aw"}, // 阿鲁巴
	"CW": {".edu.cw", ".ac.cw"}, // 库拉索
	"SX": {".edu.sx", ".ac.sx"}, // 荷属圣马丁
	"BQ": {".edu.bq", ".ac.bq"}, // 荷兰加勒比区
	"PR": {".edu.pr", ".ac.pr"}, // 波多黎各

	// 美国领土
	"AS": {".edu.as", ".ac.as"}, // 美属萨摩亚
	"GU": {".edu.gu", ".ac.gu"}, // 关岛
	"MP": {".edu.mp", ".ac.mp"}, // 北马里亚纳群岛
	"UM": {".edu.um", ".ac.um"}, // 美国本土外小岛屿

	// 其他地区和特殊领土
	"FO": {".edu.fo", ".ac.fo"}, // 法罗群岛
	"GL": {".edu.gl", ".ac.gl"}, // 格陵兰
	"SH": {".edu.sh", ".ac.sh"}, // 圣赫勒拿
	"AC": {".edu.ac", ".ac.ac"}, // 阿森松岛
	"TA": {".edu.ta", ".ac.ta"}, // 特里斯坦-达库尼亚
	"GI": {".edu.gi", ".ac.gi"}, // 直布罗陀
	"IO": {".edu.io", ".ac.io"}, // 英属印度洋领土
	"PN": {".edu.pn", ".ac.pn"}, // 皮特凯恩群岛
	"WF": {".edu.wf", ".ac.wf"}, // 瓦利斯和富图纳
	"YT": {".edu.yt", ".ac.yt"}, // 马约特
	"RE": {".edu.re", ".ac.re"}, // 留尼汪
	"PM": {".edu.pm", ".ac.pm"}, // 圣皮埃尔和密克隆
	"TF": {".edu.tf", ".ac.tf"}, // 法属南部和南极洲领土
	"BV": {".edu.bv", ".ac.bv"}, // 布韦岛
	"SJ": {".edu.sj", ".ac.sj"}, // 斯瓦尔巴和扬马延
	"HM": {".edu.hm", ".ac.hm"}, // 赫德岛和麦克唐纳群岛
	"CC": {".edu.cc", ".ac.cc"}, // 科科斯群岛
	"CX": {".edu.cx", ".ac.cx"}, // 圣诞岛
	"NF": {".edu.nf", ".ac.nf"}, // 诺福克岛
	"AQ": {".edu.aq", ".ac.aq"}, // 南极洲
}

func NewEduEmailVerify() *EduEmailVerify {
	v := &EduEmailVerify{
		majorCountries:      GlobalEducationDomains, // 之前定义的主要国家
		microstates:         MicroStrategyEducationDomains,
		africanSmallStates:  AfricanSmallStatesEducationDomains,
		asianSmallStates:    AsianSmallStatesEducationDomains,
		europeanSmallStates: EuropeanSmallStatesEducationDomains,
		specialTerritories:  SpecialTerritoriesEducationDomains,
		countryPatterns:     make(map[string][]*regexp.Regexp),
	}

	v.compileAllPatterns()
	return v
}

func (v *EduEmailVerify) compileAllPatterns() {
	allDomainSets := []map[string][]string{
		v.majorCountries,
		v.microstates,
		v.africanSmallStates,
		v.asianSmallStates,
		v.europeanSmallStates,
		v.specialTerritories,
	}

	for _, domainSet := range allDomainSets {
		for country, suffixes := range domainSet {
			for _, suffix := range suffixes {
				pattern := v.buildPattern(suffix)
				compiled := regexp.MustCompile(pattern)
				v.allPatterns = append(v.allPatterns, compiled)
				v.countryPatterns[country] = append(v.countryPatterns[country], compiled)
			}
		}
	}
}

func (v *EduEmailVerify) buildPattern(suffix string) string {
	escaped := strings.ReplaceAll(suffix, ".", `\.`)
	return `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+` + escaped + `$`
}

// 判断是否为教育邮箱（超级全面版）
func (v *EduEmailVerify) IsEducationEmail(email string) bool {
	email = strings.ToLower(strings.TrimSpace(email))

	if !v.isValidEmailFormat(email) {
		return false
	}

	// 检查所有模式
	for _, pattern := range v.allPatterns {
		if pattern.MatchString(email) {
			return true
		}
	}

	return false
}

// 获取完整信息
func (v *EduEmailVerify) GetCompleteInfo(email string) CompleteEmailInfo {
	email = strings.ToLower(strings.TrimSpace(email))

	info := CompleteEmailInfo{
		Email:       email,
		IsEducation: false,
		Country:     "",
		CountryType: "",
		Domain:      v.extractDomain(email),
		Institution: v.extractInstitution(email),
		Suffix:      v.extractSuffix(email),
	}

	if !v.isValidEmailFormat(email) {
		return info
	}

	// 检查各类别
	categories := map[string]map[string][]string{
		"Major Country":        v.majorCountries,
		"Microstate":           v.microstates,
		"African Small State":  v.africanSmallStates,
		"Asian Small State":    v.asianSmallStates,
		"European Small State": v.europeanSmallStates,
		"Special Territory":    v.specialTerritories,
	}

	for category, domains := range categories {
		for country, suffixes := range domains {
			for _, suffix := range suffixes {
				if strings.HasSuffix(email, suffix) {
					info.IsEducation = true
					info.Country = country
					info.CountryType = category
					return info
				}
			}
		}
	}

	return info
}

type CompleteEmailInfo struct {
	Email       string
	IsEducation bool
	Country     string
	CountryType string // Major Country, Microstate, etc.
	Domain      string
	Institution string
	Suffix      string
}

func (v *EduEmailVerify) extractDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

func (v *EduEmailVerify) extractSuffix(email string) string {
	domain := v.extractDomain(email)
	if domain == "" {
		return ""
	}

	// 找到最长匹配的后缀
	allSuffixes := v.getAllSuffixes()

	// 按长度排序，先匹配长的后缀
	for i := 0; i < len(allSuffixes); i++ {
		for j := i + 1; j < len(allSuffixes); j++ {
			if len(allSuffixes[i]) < len(allSuffixes[j]) {
				allSuffixes[i], allSuffixes[j] = allSuffixes[j], allSuffixes[i]
			}
		}
	}

	for _, suffix := range allSuffixes {
		if strings.HasSuffix(domain, suffix) {
			return suffix
		}
	}

	return ""
}

func (v *EduEmailVerify) getAllSuffixes() []string {
	var allSuffixes []string

	allDomainSets := []map[string][]string{
		v.majorCountries,
		v.microstates,
		v.africanSmallStates,
		v.asianSmallStates,
		v.europeanSmallStates,
		v.specialTerritories,
	}

	for _, domainSet := range allDomainSets {
		for _, suffixes := range domainSet {
			allSuffixes = append(allSuffixes, suffixes...)
		}
	}

	return allSuffixes
}

func (v *EduEmailVerify) extractInstitution(email string) string {
	domain := v.extractDomain(email)
	suffix := v.extractSuffix(email)

	if domain == "" || suffix == "" {
		return ""
	}

	return strings.TrimSuffix(domain, suffix)
}

func (v *EduEmailVerify) isValidEmailFormat(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// 统计功能
func (v *EduEmailVerify) GetStatistics() map[string]int {
	stats := make(map[string]int)

	categories := map[string]map[string][]string{
		"Major Countries":       v.majorCountries,
		"Microstates":           v.microstates,
		"African Small States":  v.africanSmallStates,
		"Asian Small States":    v.asianSmallStates,
		"European Small States": v.europeanSmallStates,
		"Special Territories":   v.specialTerritories,
	}

	totalCountries := 0
	totalDomains := 0

	for category, domains := range categories {
		countryCount := len(domains)
		domainCount := 0
		for _, suffixes := range domains {
			domainCount += len(suffixes)
		}

		stats[category+" Countries"] = countryCount
		stats[category+" Domains"] = domainCount

		totalCountries += countryCount
		totalDomains += domainCount
	}

	stats["Total Countries"] = totalCountries
	stats["Total Domain Patterns"] = totalDomains

	return stats
}
