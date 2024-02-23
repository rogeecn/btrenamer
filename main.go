package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "btrenamer",
	Short: "A bt rule base rename tool",
	RunE:  run,
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.btrenamer.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func run(cmd *cobra.Command, args []string) error {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	if len(cfgFile) > 0 {
		viper.SetConfigFile(cfgFile)
	} else {
		pwd, _ := os.Getwd()
		viper.AddConfigPath(pwd)
		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err != nil {
		return errors.Wrap(err, "read config failed")
	}

	if err := viper.Unmarshal(&rule); err != nil {
		return errors.Wrap(err, "unmarshal failed")
	}

	if len(args) < 1 {
		return errors.New("no path provided")
	}

	log.Println("rules:", rule.Rules)
	log.Println("match: ", args)
	args = []string{
		"./test/【高清影视之家发布 www.HDBTHD.com】飞鸭向前冲[高码版][国英多音轨+中文字幕].Migration.2023.2160p.HQ.WEB-DL.H265.DDP5.1.2Audio-DreamHD",
		"./test/【高清剧集网发布 www.DDHDTV.com】猎冰[第04-05集][国语音轨+简繁英字幕].The.Hunter.S01.2024.1080p.WeTV.WEB-DL.H264.AAC-BlackTV",
		"./test/【高清剧集网发布 www.DDHDTV.com】猎冰[第01-03集][国语音轨+简繁英字幕].The.Hunter.S01.2024.1080p.WeTV.WEB-DL.H264.AAC-BlackTV",
		"./test/【高清剧集网发布 www.DDHDTV.com】猎冰[第04-05集][国语音轨+简繁英字幕].The.Hunter.S02.2024.1080p.WeTV.WEB-DL.H264.AAC-BlackTV",
		"./test/【高清剧集网发布 www.DDHDTV.com】猎冰[第01-03集][国语音轨+简繁英字幕].The.Hunter.S02.2024.1080p.WeTV.WEB-DL.H264.AAC-BlackTV",
	}

	for _, rawPath := range args {
		rawPath = strings.TrimSuffix(rawPath, "/")
		for _, r := range rule.Rules {
			baseName := filepath.Base(rawPath)
			result, match, err := matchAndReplace(baseName, r)
			if !match {
				continue
			}
			if err != nil {
				log.Println("[ERR] ", err, " match: ", r.Match)
				continue
			}
			newPath := filepath.Join(filepath.Dir(rawPath), result)

			if !dirExists(newPath) {
				if err := os.MkdirAll(newPath, os.ModePerm); err != nil {
					return err
				}
			}

			if err := moveFiles(rawPath, newPath, r.Junk); err != nil {
				return err
			}

			if err := os.RemoveAll(rawPath); err != nil {
				return err
			}
			break
		}
	}

	return nil
}
