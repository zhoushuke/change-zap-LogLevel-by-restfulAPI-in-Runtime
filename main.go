package main

import (
    "time"
    "./src/logger"
)


func main() {

//    go another.run()
	a := "you are me"
    for i := 0; i < 10000; i++ {
        time.Sleep(5 * time.Second)
    	logger.Debug("log Debub:")
	logger.Infof("%s", a)
    	logger.Info("log Info")
    	logger.Warn("log Warn")
    	logger.Error("log Error")
    	logger.Infow("无法获取网址",
    		"url", "http://www.baidu.com",
    		"attempt", 3,
    		"backoff", time.Second,
		)
	}
}
