package lang

var cs = map[string]string{
	"year":       "rok|:count roky|:count let",
	"month":      "měsíc|:count měsíce|:count měsíců",
	"week":       "týden|:count týdny|:count týdnů",
	"day":        "den|:count dny|:count dní",
	"hour":       "hodinu|:count hodiny|:count hodin",
	"minute":     "minutu|:count minuty|:count minut",
	"second":     "sekundu|:count sekundy|:count sekund",
	"ago":        "před :time",
	"from_now":   "za :time",
	"after":      ":time později",
	"before":     ":time předtím",
	"year_ago":   "rokem|[2,Inf]:count lety",
	"month_ago":  "měsícem|[2,Inf]:count měsíci",
	"week_ago":   "týdnem|[2,Inf]:count týdny",
	"day_ago":    "dnem|[2,Inf]:count dny",
	"hour_ago":   "hodinou|[2,Inf]:count hodinami",
	"minute_ago": "minutou|[2,Inf]:count minutami",
	"second_ago": "sekundou|[2,Inf]:count sekundami",
}
