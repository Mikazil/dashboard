package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Weather  WeatherConfig  `yaml:"weather"`
	Currency CurrencyConfig `yaml:"currency"`
	RSS     RSSConfig      `yaml:"rss"`
	Theme   ThemeConfig    `yaml:"theme"`
}

type WeatherConfig struct {
	City     string `yaml:"city"`
	Interval string `yaml:"interval"`
	APIKey   string `yaml:"api_key"`
}

type CurrencyConfig struct {
	Codes    []string `yaml:"codes"`
	Crypto   []string `yaml:"crypto"`
	Interval string   `yaml:"interval"`
}

type RSSConfig struct {
	Feeds         []RSSFeed `yaml:"feeds"`
	ScrollSpeed   string    `yaml:"scroll_speed"`
	UpdateInterval string   `yaml:"update_interval"`
}

type RSSFeed struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type ThemeConfig struct {
	Primary   string `yaml:"primary"`
	Secondary string `yaml:"secondary"`
	Dim       string `yaml:"dim"`
	Bg        string `yaml:"bg"`
}

func Default() Config {
	return Config{
		Weather: WeatherConfig{
			City:     "Moscow",
			Interval: "5m",
			APIKey:   "",
		},
		Currency: CurrencyConfig{
			Codes:    []string{"USD", "EUR", "CNY"},
			Crypto:   []string{"bitcoin"},
			Interval: "5m",
		},
		RSS: RSSConfig{
			Feeds: []RSSFeed{
				{Name: "Habr", URL: "https://habr.com/ru/rss/hubs/all/"},
				{Name: "Lenta", URL: "https://lenta.ru/rss/news"},
			},
			ScrollSpeed:    "50ms",
			UpdateInterval: "1h",
		},
		Theme: ThemeConfig{
			Primary:   "#FFB000",
			Secondary: "#E89900",
			Dim:       "#886600",
			Bg:        "#000000",
		},
	}
}

func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Default(), nil
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Default(), err
	}

	return cfg, nil
}
