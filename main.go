package main

import (
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"

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
	log.Println("match: ", args[0])
	rawPath := args[0]
	// rawPath = "./test/【高清影视之家发布 www.HDBTHD.com】飞鸭向前冲[高码版][国英多音轨+中文字幕].Migration.2023.2160p.HQ.WEB-DL.H265.DDP5.1.2Audio-DreamHD"

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
		log.Println("[BINGO] ", result)
		// rename old to new path
		newFullPath := filepath.Join(filepath.Dir(rawPath), result)
		os.Rename(rawPath, newFullPath)

		filepath.Walk(newFullPath, func(filePath string, info fs.FileInfo, err error) error {
			if filePath == newFullPath {
				return nil
			}
			baseName := info.Name()[:len(info.Name())-len(path.Ext(info.Name()))]
			log.Println(">>>> ", baseName)
			for _, junk := range r.Junk {
				if baseName == junk {
					os.Remove(filePath)
				}
			}
			return nil
		})
		break
	}

	return nil
}
