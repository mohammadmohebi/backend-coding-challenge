package geonames

var CanAdminDivision = map[string]string{
	"01": "AB",
	"02": "BC",
	"03": "MB",
	"04": "NB",
	"05": "NL",
	"07": "NS",
	"13": "NS",
	"14": "NU",
	"08": "ON",
	"09": "PE",
	"10": "QC",
	"11": "SK",
	"12": "YT",
}

var CountryCode = map[string]string{
	"CA": "Canada",
	"US": "USA",
}

func GetCountryName(countryCode string) string {
	return CountryCode[countryCode]
}

func GetAdmin1Code(code string, countryCode string) string {
	if countryCode == "CA" {
		return CanAdminDivision[code]
	} else {
		return code
	}
}
