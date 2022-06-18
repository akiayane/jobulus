package config

var (
	Token     string //To store value of Token from config.json .
	BotPrefix string // To store value of BotPrefix from config.json.

	config *configStruct //To store value extracted from config.json.
)

type configStruct struct {
	Token     string `json : "Token"`
	BotPrefix string `json : "BotPrefix"`
}

func SetToken(token, botPrefix string) error {
	Token = token
	BotPrefix = botPrefix

	//If there isn't any error we will return nil.
	return nil
}
