package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	debug   bool
)
var tvFileRegExp = `(.*?\.S\d{2}E\d{2}\.\d{4}\..*?)\..*\.(.*?)$`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "btrenamer",
	Short: "A bt rule base rename tool",
	RunE:  run,
}

var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "rename dir season files",
	RunE:  renameSeasonFiles,
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.btrenamer.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().BoolVarP(&debug, "debug", "d", false, "is debug mode")

	rootCmd.AddCommand(renameCmd)
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

	if debug {
		args = []string{
			"./test/【高清影视之家发布 www.HDBTHD.com】飞鸭向前冲[高码版][国英多音轨+中文字幕].Migration.2023.2160p.HQ.WEB-DL.H265.DDP5.1.2Audio-DreamHD",
			"./test/【高清剧集网发布 www.DDHDTV.com】猎冰[第04-05集][国语音轨+简繁英字幕].The.Hunter.S01.2024.1080p.WeTV.WEB-DL.H264.AAC-BlackTV",
			"./test/【高清剧集网发布 www.DDHDTV.com】猎冰[第01-03集][国语音轨+简繁英字幕].The.Hunter.S01.2024.1080p.WeTV.WEB-DL.H264.AAC-BlackTV",
			"./test/【高清剧集网发布 www.DDHDTV.com】猎冰[第04-05集][国语音轨+简繁英字幕].The.Hunter.S02.2024.1080p.WeTV.WEB-DL.H264.AAC-BlackTV",
			"./test/【高清剧集网发布 www.DDHDTV.com】猎冰[第01-03集][国语音轨+简繁英字幕].The.Hunter.S02.2024.1080p.WeTV.WEB-DL.H264.AAC-BlackTV",
		}
	}

	for _, rawPath := range args {
		rawPath = strings.TrimSuffix(rawPath, "/")
		for i, r := range rule.Rules {
			baseName := filepath.Base(rawPath)
			path, match, err := matchAndReplace(baseName, r)
			if !match {
				continue
			}
			if err != nil {
				log.Println("[ERR] ", err, " match: ", r.Match)
				continue
			}

			if err := moveFiles(rawPath, path, rule, i); err != nil {
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

func renameSeasonFiles(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("need dir params")
	}

	path := args[0]
	if !dirExists(path) {
		return errors.New("dir not exists: " + path)
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return errors.Wrap(err, "read dir failed")
	}

	r, err := regexp.Compile(tvFileRegExp)
	if err != nil {
		return err
	}

	for _, file := range files {
		if r.MatchString(file.Name()) {
			matches := r.FindStringSubmatch(file.Name())
			if len(matches) != 3 {
				continue
			}
			newFilename := fmt.Sprintf("%s.%s", matches[1], matches[2])

			log.Println("rename: from ", file.Name(), " to: ", newFilename)
			if err := os.Rename(filepath.Join(path, file.Name()), filepath.Join(path, newFilename)); err != nil {
				return err
			}
		}
	}

	return nil
}
