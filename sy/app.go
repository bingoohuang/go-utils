package sy

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/bingoohuang/gou/file"
	"github.com/spf13/pflag"

	"github.com/spf13/viper"

	"github.com/bingoohuang/gou/cnf"
	"github.com/bingoohuang/gou/htt"
	"github.com/bingoohuang/gou/lo"
)

// AppOption defines the application options.
type AppOption struct {
	EnvPrefix   string
	LogLevel    string
	ConfigBeans []interface{}

	LogWriter io.Writer
	CnfTpl    string
	CtlTpl    string
}

// SetupApp setup the application.
func SetupApp(appOption *AppOption) {
	var ipo *bool

	if appOption.CnfTpl != "" || appOption.CtlTpl != "" {
		ipo = pflag.BoolP("init", "", false, "init to create demo config file and ctl.sh")
	}

	lo.DeclareLogPFlags()
	cnf.DeclarePflags()
	cnf.DeclarePflagsByStruct(appOption.ConfigBeans...)

	pprofAddr := htt.PprofAddrPflag()

	if err := cnf.ParsePflags(appOption.EnvPrefix); err != nil {
		panic(err)
	}

	if ipo != nil && *ipo {
		initIPO(appOption)
	}

	htt.StartPprof(*pprofAddr)

	if appOption.LogLevel != "" {
		viper.Set(lo.LoglevelKey, appOption.LogLevel)
	}

	appOption.LogWriter = lo.SetupLog()

	_ = UpdatePidFile()

	cnf.LoadByPflag(appOption.ConfigBeans...)
}

func initIPO(appOption *AppOption) {
	if appOption.CnfTpl != "" {
		if err := initCfgFile(appOption.CnfTpl, "cnf.toml"); err != nil {
			fmt.Println(err.Error())
		}
	}

	if appOption.CtlTpl != "" {
		if err := initCtl(appOption.CtlTpl, "ctl"); err != nil {
			fmt.Println(err.Error())
		}
	}

	os.Exit(0)
}

// initCfgFile initializes the cfg file.
func initCfgFile(conf, configFileName string) error {
	if file.Stat(configFileName) == file.Exists {
		fmt.Printf("%s already exists, ignored!\n", configFileName)
		return nil
	}

	if err := file.WriteValue(configFileName, conf); err != nil {
		return err
	}

	fmt.Println(configFileName, "created!")

	return nil
}

// initCtl initializes the ctl file.
func initCtl(ctl, ctlFilename string) error {
	if file.Stat(ctlFilename) == file.Exists {
		fmt.Println(ctlFilename + " already exists, ignored!")
		return nil
	}

	tpl, err := template.New(ctlFilename).Parse(ctl)
	if err != nil {
		return err
	}

	binArgs := argsExcludeInit()
	m := map[string]string{
		"BinName": os.Args[0],
		"BinArgs": strings.Join(binArgs, " ")}

	var content bytes.Buffer
	if err := tpl.Execute(&content, m); err != nil {
		return err
	}

	if err := file.WriteValue(ctlFilename, content.String()); err != nil {
		return err
	}

	// 0755->即用户具有读/写/执行权限，组用户和其它用户具有读写权限；
	if err := os.Chmod(ctlFilename, 0755); err != nil {
		return err
	}

	fmt.Println(ctlFilename, "created!")

	return nil
}

func argsExcludeInit() []string {
	binArgs := make([]string, 0, len(os.Args)-2) // nolint gomnd

	for _, arg := range os.Args[1:] {
		if strings.Index(arg, "--init") == 0 {
			continue
		}

		if strings.Index(arg, "-") != 0 {
			arg = strconv.Quote(arg)
		}

		binArgs = append(binArgs, arg)
	}

	return binArgs
}
