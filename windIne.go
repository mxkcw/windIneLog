package windIne

import (
	"fmt"
	"github.com/mxkcw/windIneLog/config"
	"github.com/mxkcw/windIneLog/windIne_http"
	"github.com/mxkcw/windIneLog/windIne_log"
)

type RunMode int

const (
	RunModeUnknown RunMode = iota
	RunModeDebug
	RunModeRelease
	RunModeTest
)

type WindIneAppSignalInfo struct {
	SigCode string
	Msg     string
}

var (
	currentRunMode RunMode
)

func (rm RunMode) String() string {
	switch rm {
	case RunModeDebug:
		return "Debug"
	case RunModeRelease:
		return "Release"
	case RunModeTest:
		return "Test"
	case RunModeUnknown:
		return "Unknown"
	default:
		return "Unknown"
	}
}

func GetCurrentRunMode() RunMode {
	return currentRunMode
}

// SetupWindIneBox
// 必须--YES
// 必须使用此方法初始化工具库,未使用此方法初始化，无法使用完整功能，亦存在兼容性问题
// logMaxSaveDays Log是否开启文件存储模式
// log_dir 自定义日志目录,默认为:/usr/logs/${projectName},如果传"" 即使用默认值
// httpRequestTimeOut 网络请求超时时间
// projectName--项目名称，
// run_mode 运行模式 debug
// logLevel--日志等级，
// logMaxSaveTime--默认365天,
// logSaveType--日志分片格式，默认按天分片，可选按小时分片
func SetupWindIneBox(projectName string, runMode RunMode, productLogDir string, logMaxSaveDays int64, logSaveType windIne_log.WindIneLogSaveType, httpRequestTimeOut int) {
	enableSaveLogFile := false
	logLevel := windIne_log.WindIneLogStyleDebug
	currentRunMode = runMode
	switch runMode {
	case RunModeDebug:
		enableSaveLogFile = false
		logLevel = windIne_log.WindIneLogStyleDebug
	case RunModeTest:
		enableSaveLogFile = true
		logLevel = windIne_log.WindIneLogStyleDebug
	case RunModeRelease:
		enableSaveLogFile = true
		logLevel = windIne_log.WindIneLogStyleInfo
	}

	windIne_log.SetupLoWindIneools(projectName, enableSaveLogFile, logLevel, logMaxSaveDays, logSaveType, productLogDir)
	windIne_http.DefaultTimeout = httpRequestTimeOut
	config.IsSetup = true
	fmt.Printf("Tools Setup End\nProjcetName=[%s]\nrunMode=[%s]\nlogLeve=[%s]\nproduct main logdir=[%s]\nlogCutType=[%s]\nlogSaveDays=[%d]\nhttpRequestTimeout=[%d Second]\n",
		windIne_log.GetProjectName(),
		runMode.String(),
		windIne_log.GetLogLevel().String(),
		windIne_log.GetProductMainLogDir(),
		logSaveType.String(),
		logMaxSaveDays,
		windIne_http.DefaultTimeout,
	)
}
