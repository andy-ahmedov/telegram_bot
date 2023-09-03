package config

import (
	"github.com/andy-ahmedov/telegram_bot/pkg/logging"
	"github.com/spf13/viper"
)

type Config struct {
	TelegramToken string
	PasswordDB    string
	UserName      string
	DBname        string
	SMSapi        string
	Sslmode       string
	PathToXml     string `mapstructure:"path_to_xml"`
	CreateTable   string `mapstructure:"create_table"`
	ChangeTable   string `mapstructure:"change_request"`
	Port          string `mapstructure:"port"`
	Host          string `mapstructure:"host"`
	Delimiter     string `mapstructure:"delimiter"`
	Catalog       string `mapstructure:"catalog_file"`
	Balance       string `mapstructure:"balance"`
	Discounts_url string `mapstructure:"discounts_url"`
	ConnectDB     string `mapstructure:"connect_db"`
	SMSlink       string `mapstructure:"sms_link"`
	DBPath        string `mapstructure:"db_file"`

	Messages   Messages
	ReservWrds ReservedWords
}

type Messages struct {
	Errors
	Responses
	Addresses
}

type ReservedWords struct {
	StoreAddr
}

type StoreAddr struct {
	Mayakovskogo   string `mapstructure:"mayakovskogo"`
	Ryleyeva       string `mapstructure:"ryleyeva"`
	Pushkarovskoye string `mapstructure:"pushkarovskoye"`
	Optimus        string `mapstructure:"optimus"`
	Sovremennik    string `mapstructure:"sovremennik"`
	Promyshlennaya string `mapstructure:"promyshlennaya"`
	DiscountCenter string `mapstructure:"discount_center"`
}

type Addresses struct {
	Mayakovskogo   string `mapstructure:"mayakovskogo"`
	Ryleyeva       string `mapstructure:"ryleyeva"`
	Pushkarovskoye string `mapstructure:"pushkarovskoye"`
	Optimus        string `mapstructure:"optimus"`
	Sovremennik    string `mapstructure:"sovremennik"`
	Promyshlennaya string `mapstructure:"promyshlennaya"`
	DiscountCenter string `mapstructure:"discount_center"`
}

type Errors struct {
	NeedCode             string `mapstructure:"need_code"`
	AccessDenied         string `mapstructure:"access_denied"`
	Unauthorized         string `mapstructure:"unauthorized"`
	CantCreateToken      string `mapstructure:"cant_create_token"`
	NumberAlreadyAuth    string `mapstructure:"number_already_auth"`
	CantDeleteDataFromDB string `mapstructure:"cant_delete_data_from_DB"`
	CantSendSMSToPhone   string `mapstructure:"cant_send_sms_to_phone"`
	CantGetDataFromDB    string `mapstructure:"cant_get_data_from_DB"`
	CantSendMessage      string `mapstructure:"cant_send_message"`
	CheckUserStatus      string `mapstructure:"check_user_status"`
	CantSaveToDB         string `mapstructure:"cant_save_to_DB"`
	UnknownError         string `mapstructure:"unknown_error"`
}

type Responses struct {
	CodeSent      string `mapstructure:"code_sent"`
	NeedNumber    string `mapstructure:"need_number"`
	SuccessAuth   string `mapstructure:"success_auth"`
	TextStart     string `mapstructure:"text_start"`
	DeleteSuccess string `mapstructure:"delete_success"`
	AlreadyLogged string `mapstructure:"already_logged"`
	SelectStore   string `mapstructure:"select_store"`
	Novelties     string `mapstructure:"novelties"`
	Contacts      string `mapstructure:"contacts"`
	StartWhenAuth string `mapstructure:"start_already_logged"`
}

func Init(logger *logging.Logger) (*Config, error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		logger.Error("Ошибка в функции viper.ReadInConfig:", err)
		return nil, err
	}

	var cfg Config

	if err := viper.Unmarshal(&cfg); err != nil {
		logger.Error("Ошибка в функции viper.Unmarshal:", err)
		return nil, err
	}

	if err := viper.UnmarshalKey("messages.responses", &cfg.Messages.Responses); err != nil {
		logger.Error("Ошибка в функции viper.UnmarshalKey:", err)
		return nil, err
	}

	if err := viper.UnmarshalKey("messages.errors", &cfg.Messages.Errors); err != nil {
		logger.Error("Ошибка в функции viper.UnmarshalKey:", err)
		return nil, err
	}

	if err := viper.UnmarshalKey("messages.addresses", &cfg.Messages.Addresses); err != nil {
		logger.Error("Ошибка в функции viper.UnmarshalKey:", err)
		return nil, err
	}

	if err := viper.UnmarshalKey("reserved_words.store_addr", &cfg.ReservWrds.StoreAddr); err != nil {
		logger.Error("Ошибка в функции viper.UnmarshalKey:", err)
		return nil, err
	}

	if err := parseEnv(&cfg, logger); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func parseEnv(cfg *Config, logger *logging.Logger) error {
	if err := viper.BindEnv("token"); err != nil {
		logger.Error("Ошибка в функции viper.BindEnv:", err)
		return err
	}

	if err := viper.BindEnv("POSTGRES_PASSWORD"); err != nil {
		logger.Error("Ошибка в функции viper.BindEnv:", err)
		return err
	}

	if err := viper.BindEnv("POSTGRES_USER"); err != nil {
		logger.Error("Ошибка в функции viper.BindEnv:", err)
		return err
	}

	if err := viper.BindEnv("POSTGRES_DB"); err != nil {
		logger.Error("Ошибка в функции viper.BindEnv:", err)
		return err
	}

	if err := viper.BindEnv("sms_api"); err != nil {
		logger.Error("Ошибка в функции viper.BindEnv:", err)
		return err
	}

	if err := viper.BindEnv("path_to_xml"); err != nil {
		logger.Error("Ошибка в функции viper.BindEnv:", err)
		return err
	}

	if err := viper.BindEnv("ssl_mode"); err != nil {
		logger.Error("Ошибка в функции viper.BindEnv:", err)
		return err
	}

	cfg.TelegramToken = viper.GetString("TOKEN")
	cfg.PasswordDB = viper.GetString("POSTGRES_PASSWORD")
	cfg.UserName = viper.GetString("POSTGRES_USER")
	cfg.DBname = viper.GetString("POSTGRES_DB")
	cfg.SMSapi = viper.GetString("SMS_API")
	cfg.PathToXml = viper.GetString("PATH_TO_XML")
	cfg.Sslmode = viper.GetString("SSL_MODE")

	return nil
}
