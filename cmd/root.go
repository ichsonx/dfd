package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var keywords string    // 要删除文件所包含的字符串
var rootpath = "./"    // 当前目录的相对路径
var wd, _ = os.Getwd() // 当前目录的绝对路径

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dfd",
	Short: "一个简单，递归删除所有空文件夹、指定包含字符串文件的工具",
	Long: `一个简单，递归删除所有空文件夹、指定包含字符串文件的工具. For example:
1. 删除所有空文件夹（包括子目录里的）：     dfd
2. 删除包含关键字的文件（包括子目录里的）：  dfd -f 关键字1,关键字2,关键字3,....
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// 如果使用了 -f，并且带有有效参数，删除文件
		// 否则，删除空文件夹
		if len(keywords) > 0 {
			k := strings.Split(strings.Trim(keywords, " "), ",")
			if len(k) == 0 {
				fmt.Println("关键字不能为空！")
				os.Exit(0)
			} else {
				loopDelFile(k)
			}
		} else {
			loopDelDir()
		}
		time.Sleep(30 * time.Second)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().StringVarP(&keywords, "file", "f", "", "删除文件所包含的关键字，模糊搜索，多个关键词可用逗号隔开")
}

// 递归删除所有指定文件
func loopDelFile(keywords []string) {
	fmt.Println("开始删除对应文件")
	total := 0
	for {
		c := 0
		err := filepath.Walk(rootpath, func(path string, info fs.FileInfo, err error) error {
			if !info.IsDir() {
				for _, k := range keywords {
					if strings.Contains(info.Name(), k) {
						fp := filepath.Join(wd, path)
						err := os.Remove(fp)
						if err == nil {
							c += 1
						}
					}
				}
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
		total += c
		if c == 0 {
			break
		}
	}
	fmt.Printf("本次总共删除文件：%d \n", total)
}

// 递归删除所有空目录
func loopDelDir() {
	fmt.Println("开始删除所有空目录...")
	total := 0
	for {
		c := 0
		err := filepath.Walk(rootpath, func(path string, info fs.FileInfo, err error) error {
			if info.IsDir() {
				fp := filepath.Join(wd, path)
				err := os.Remove(fp)
				if err == nil {
					c += 1
				}
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
		total += c
		if c == 0 {
			break
		}
	}
	fmt.Printf("本次总共删除空文件夹：%d \n", total)
}

// initConfig reads in config file and ENV variables if set.
//func initConfig() {
//	if cfgFile != "" {
//		// Use config file from the flag.
//		viper.SetConfigFile(cfgFile)
//	} else {
//		// Find home directory.
//		home, err := os.UserHomeDir()
//		cobra.CheckErr(err)
//
//		// Search config in home directory with name ".dfd" (without extension).
//		viper.AddConfigPath(home)
//		viper.SetConfigType("yaml")
//		viper.SetConfigName(".dfd")
//	}
//
//	viper.AutomaticEnv() // read in environment variables that match
//
//	// If a config file is found, read it in.
//	if err := viper.ReadInConfig(); err == nil {
//		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
//	}
//}
