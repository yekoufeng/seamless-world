package entity

import (
	"os"
	"path/filepath"

	log "github.com/cihub/seelog"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func init() {
	initDefs()

	initLogger()
}

var elogger log.LoggerInterface

func initLogger() {
	viper.SetConfigFile("../res/config/server.json")
	if err := viper.ReadInConfig(); err != nil {
		panic("加载配置文件失败")
	}

	logDir := viper.GetString("Config.LogDir")
	logLevel := viper.GetString("Config.LogLevel")
	load(logDir, logLevel)

	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Error(err)
			return
		}
		defer watcher.Close()

		if err = watcher.Add("../res/config/server.json"); err != nil {
			log.Error(err)
			return
		}

		for {
			select {
			case ev := <-watcher.Events:
				if ev.Op == fsnotify.Write {
					if err := viper.ReadInConfig(); err != nil {
						log.Error(err)
						continue
					}
					dir := viper.GetString("Config.LogDir")
					level := viper.GetString("Config.LogLevel")
					load(dir, level)
				}
			case err := <-watcher.Errors:
				log.Error(err)
			}
		}
	}()

}

func load(logdir string, loglevel string) {
	log.Info("Logdir:", logdir, " Loglevel:", loglevel)

	switch loglevel {
	case "prod":
		loglevel = "info"
	case "debug", "info", "warn", "error":
	default:
		log.Error("LogLevel unsupport! [debug/prod/info/warn/error]")
		return
	}

	// 	<filter levels="debug">
	// 	<console />
	// </filter>
	srvName := logdir + filepath.Base(os.Args[0])
	str1 := `<seelog minlevel="`
	str2 := `" maxlevel="error">
				<outputs formatid="main">

					<filter levels="debug,info,warn,error"> 
						<console />
						<buffered size="10000" flushperiod="1000">
							<rollingfile type="date" filename="`
	str3 := `.infolog" datepattern="2006.01.02.15" fullname="true" maxrolls="30"/>  
						</buffered>
						<filter levels="warn,error">
							<buffered size="10000" flushperiod="1000">
								<rollingfile type="date" filename="`
	str4 := `.warnlog" datepattern="2006.01.02.15" fullname="true" maxrolls="30"/>  
							</buffered>
							<filter levels="error">
								<buffered size="10000" flushperiod="1000">
									<rollingfile type="date" filename="`
	str5 := `.errorlog" datepattern="2006.01.02.15" fullname="true" maxrolls="30"/>  
								</buffered>
							</filter>
						</filter>
					</filter>
				</outputs>
				<formats>
					<format id="main" format="%Date(2006-01-02 15:04:05.999) [%LEV] [%File:%Line] %Msg%n"/>  
				</formats>
			</seelog>`
	path := str1 + loglevel + str2 + srvName + str3 + srvName + str4 + srvName + str5

	defer log.Flush()
	logger, err := log.LoggerFromConfigAsString(path)
	if err != nil {
		panic(err)
	}

	logger.SetAdditionalStackDepth(1)
	elogger = logger
}
