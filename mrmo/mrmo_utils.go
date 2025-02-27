package mrmo

var regionToBasePathMap = map[string]string{
	"dca":            "inindca.com",
	"tca":            "inintca.com",
	"us-east-1":      "mypurecloud.com",
	"us-east-2":      "use2.us-gov-pure.cloud",
	"us-west-2":      "usw2.pure.cloud",
	"eu-west-1":      "mypurecloud.ie",
	"eu-west-2":      "euw2.pure.cloud",
	"ap-southeast-2": "mypurecloud.com.au",
	"ap-northeast-1": "mypurecloud.jp",
	"eu-central-1":   "mypurecloud.de",
	"ca-central-1":   "cac1.pure.cloud",
	"ap-northeast-2": "apne2.pure.cloud",
	"ap-south-1":     "aps1.pure.cloud",
	"sa-east-1":      "sae1.pure.cloud",
	"ap-northeast-3": "apne3.pure.cloud",
	"eu-central-2":   "euc2.pure.cloud",
	"me-central-1":   "mec1.pure.cloud",
}

func getRegionMap() map[string]string {
	return map[string]string{
		"dca":            "inindca.com",
		"tca":            "inintca.com",
		"us-east-1":      "mypurecloud.com",
		"us-east-2":      "use2.us-gov-pure.cloud",
		"us-west-2":      "usw2.pure.cloud",
		"eu-west-1":      "mypurecloud.ie",
		"eu-west-2":      "euw2.pure.cloud",
		"ap-southeast-2": "mypurecloud.com.au",
		"ap-northeast-1": "mypurecloud.jp",
		"eu-central-1":   "mypurecloud.de",
		"ca-central-1":   "cac1.pure.cloud",
		"ap-northeast-2": "apne2.pure.cloud",
		"ap-south-1":     "aps1.pure.cloud",
		"sa-east-1":      "sae1.pure.cloud",
		"ap-northeast-3": "apne3.pure.cloud",
		"eu-central-2":   "euc2.pure.cloud",
		"me-central-1":   "mec1.pure.cloud",
	}
}
