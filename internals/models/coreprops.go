package models

type CoreProps struct {
	Address            string `json:"address"`
	EncryptedAddress   string `json:"encryptedAddress"`
	GgEncryptedAddress string `json:"ggEncryptedAddress"`
}

type SubApps struct {
	SubApps struct {
		Sonar struct {
			Metadata struct {
				WebServerAddress string `json:"webServerAddress"`
			} `json:"metadata"`
		} `json:"sonar"`
	} `json:"subApps"`
}
