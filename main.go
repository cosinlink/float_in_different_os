package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const (
	GenesisTimestamp = uint32(1561910400) //2019-07-01 00:00:00
)

func calVoteWeight(now uint32, stakeWeight uint64) uint64 {
	precision := math.Pow10(4)
	weight := math.Floor(float64(now-GenesisTimestamp)/(7*24*3600*52)*precision) / precision
	return stakeWeight * uint64(math.Pow(2, float64(weight)))
}

func calVoteWeightFloat(now uint32) float64 {
	precision := math.Pow10(4)
	weight := math.Floor(float64(now-GenesisTimestamp)/(7*24*3600*52)*precision) / precision
	return math.Pow(2, float64(weight))
}

func main() {
	fmt.Println("---------------")
	log.Println("------test calVoteWeight in different OS ----")

	okUnix, _ := PathExists("print_float_unix.log")
	okWindows, _ := PathExists("print_float_windows.log")

	if okUnix && okWindows {
		log.Println("------print_float_unix.log and print_float_windows.log already exist, just diff them ----")
		c := "diff print_float_unix.log print_float_windows.log"
		cmd := exec.Command("sh", "-c", c)
		out, err := cmd.Output()
		if err != nil {
			log.Printf("%s\n", out)
		} else {
			log.Println(err)
		}
		return
	}

	calFloatToLogFile()
}

func func_log2file() {
	//创建日志文件
	f, err := os.OpenFile("test.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	//完成后，延迟关闭
	defer f.Close()
	// 设置日志输出到文件
	log.SetOutput(f)
	// 写入日志内容
	log.Println("check to make sure it works")
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func calFloatToLogFile() {
	fileName := "print_float_"
	osStr := runtime.GOOS
	if strings.Contains(osStr, "windows") {
		fileName += "windows.log"
	} else if strings.Contains(osStr, "darwin") {
		fileName += "unix.log"
	}

	ok, _ := PathExists(fileName)
	if ok {
		log.Println(fileName + " already exits")
		return
	}

	//创建日志文件
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	//完成后，延迟关闭
	defer f.Close()
	// 设置日志输出到文件
	// 定义多个写入器
	writers := []io.Writer{
		f,
		os.Stdout}
	fileAndStdoutWriter := io.MultiWriter(writers...)
	// 创建新的log对象
	//logger := log.New(fileAndStdoutWriter, "", log.Ldate|log.Ltime|log.Lshortfile)
	logger := log.New(fileAndStdoutWriter, "", 0)
	// 使用新的log对象，写入日志内容
	//logger.Println("--> logger :  check to make sure it works")
	testCalToFile(logger)
}

func testCalToFile(logger *log.Logger) {
	timeMaxLength := uint32(3600 * 24 * 6000)
	var now, last float64
	for i := uint32(1); i <= timeMaxLength; i += 1 {

		now = calVoteWeightFloat(i + GenesisTimestamp)
		if last != now {
			logger.Println(i, now)
		}
		last = now
	}

}
