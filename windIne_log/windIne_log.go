package windIne_log

/*
   Package windIne_log Log工具
*/

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/mxkcw/windIne/windIne_color"
	"github.com/mxkcw/windIne/windIne_time"
	"github.com/sirupsen/logrus"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	currentLogConfig *WindIneLogConf
	logConfigOnce    sync.Once
	setupComplete    bool
	mainLog          *WindIneLog
)

// WindIneLogStyle 日志样式
type WindIneLogStyle int

const (
	WindIneLogStyleDebug   WindIneLogStyle = iota // Debug
	WindIneLogStyleError                          // Error
	WindIneLogStyleWarning                        // Warning
	WindIneLogStyleInfo                           // Info
	WindIneLogStyleTrace                          // Trace
	WindIneLogStyleFatal                          // Fatal
	WindIneLogStylePanic                          // Panic
)

func (aStyle WindIneLogStyle) String() string {
	switch aStyle {
	case WindIneLogStyleFatal:
		return "fatal"
	case WindIneLogStyleTrace:
		return "trace"
	case WindIneLogStyleInfo:
		return "info"
	case WindIneLogStyleWarning:
		return "warning"
	case WindIneLogStyleError:
		return "error"
	case WindIneLogStyleDebug:
		return "debug"
	case WindIneLogStylePanic:
		return "panic"
	default:
		return "debug"
	}
}

// WindIneLogSaveType 日志分片类型
type WindIneLogSaveType int

const (
	WindIneLogSaveTypeDays WindIneLogSaveType = iota //按日分片
	WindIneLogSaveHours                              //按小时分片
)

func (aFlag WindIneLogSaveType) String() string {
	switch aFlag {
	case WindIneLogSaveTypeDays:
		return "Days"
	case WindIneLogSaveHours:
		return "Hours"
	default:
		return "Unknown"
	}
}

type WindIneLogConf struct {
	productName       string
	productLogDir     string
	enableSaveLogFile bool
	logLeve           WindIneLogStyle
	logMaxSaveDays    int64
	logSaveType       WindIneLogSaveType
}

func instanceConfig() *WindIneLogConf {
	logConfigOnce.Do(func() {
		currentLogConfig = &WindIneLogConf{}
	})
	return currentLogConfig
}

func setupDefaultLog() *WindIneLog {
	if setupComplete == false && mainLog == nil {
		mainLog = NewWindIneLog(strings.ToLower(instanceConfig().productName))
	}
	return mainLog
}

type WindIneLog struct {
	sync.RWMutex
	logger         *logrus.Logger // 添加这一行
	modelName      string
	logDir         string
	logDirWithDate string
	entryTime      time.Time // 日志初始化时间,留作后续比对使用
	lastCheckTime  time.Time // 记录最后一次检查时间,用作日志轮转

}

func GetProjectName() string {
	return instanceConfig().productName
}

func GetLogLevel() WindIneLogStyle {
	return instanceConfig().logLeve
}

func GetProductMainLogDir() string {
	return instanceConfig().productLogDir
}

// logF 快捷日志Function，含模块字段封装
// Params [style] log类型  fatal、trace、info、warning、error、debug
// Params [format] 模块名称：自定义字符串
// Params [args...] 模块名称：自定义字符串
func (aLog *WindIneLog) logF(style WindIneLogStyle, format string, args ...interface{}) {
	aLog.Lock()
	defer aLog.Unlock()

	colorFormat := format
	if instanceConfig().enableSaveLogFile != true {
		// TODO 每分钟检查一次是否需要更新日志文件路径
		now := time.Now().UTC()
		if now.Sub(aLog.lastCheckTime) > time.Minute {
			if windIne_time.WindIneDateEqualYearMoonDay(aLog.lastCheckTime, now) == false {
				aLog.logDirWithDate = fmt.Sprintf("%s/%s", aLog.logDir, now.Format("2006-01-02"))
				rLog := newLogSaveHandler(aLog)
				aLog.logger.SetOutput(rLog)
				aLog.lastCheckTime = now
			}
		}
		// 对每个占位符、非占位符片段和'['、']'进行迭代，为它们添加相应的颜色
		re := regexp.MustCompile(`(%[vTsdfqTbcdoxXUeEgGp]+)|(\[|\])|([^%\[\]]+)`)
		colorFormat = re.ReplaceAllStringFunc(format, func(s string) string {
			switch {
			case strings.HasPrefix(s, "%"):
				return fmt.Sprintf("%s%s%s", windIne_color.ANSIColorForegroundBrightYellow, s, windIne_color.ANSIColorReset)
			case s == "[" || s == "]":
				return s // 保持 `[` 和 `]` 的原始颜色
			default:
				if style == WindIneLogStyleError {
					return fmt.Sprintf("%s%s%s", windIne_color.ANSIColorForegroundBrightRed, s, windIne_color.ANSIColorReset)
				} else if style == WindIneLogStyleInfo {
					return fmt.Sprintf("%s%s%s", windIne_color.ANSIColorForegroundBrightGreen, s, windIne_color.ANSIColorReset)
				} else {
					return fmt.Sprintf("%s%s%s", windIne_color.ANSIColorForegroundBrightCyan, s, windIne_color.ANSIColorReset)
				}
			}
		})
	}

	if style != WindIneLogStyleInfo {
		pc, _, _, _ := runtime.Caller(2)
		fullName := runtime.FuncForPC(pc).Name()

		lastDot := strings.LastIndex(fullName, ".")
		if lastDot == -1 || lastDot == 0 || lastDot == len(fullName)-1 {
			return
		}
		callerClass := fullName[:lastDot]
		method := fullName[lastDot+1:]

		prefixFormat := fmt.Sprintf("[pkg--%s--][method--%s--] ", callerClass, method)
		colorFormat = prefixFormat + colorFormat
	}

	switch style {
	case WindIneLogStyleFatal:
		aLog.logger.Fatalf(colorFormat, args...)
	case WindIneLogStyleTrace:
		aLog.logger.Tracef(colorFormat, args...)
	case WindIneLogStyleInfo:
		aLog.logger.Infof(colorFormat, args...)
	case WindIneLogStyleWarning:
		aLog.logger.Warnf(colorFormat, args...)
	case WindIneLogStyleError:
		aLog.logger.Errorf(colorFormat, args...)
	case WindIneLogStyleDebug:
		aLog.logger.Debugf(colorFormat, args...)
	case WindIneLogStylePanic:
		aLog.logger.Panicf(colorFormat, args...)

	}
}

// LogInfof format格式化log--info信息
func (aLog *WindIneLog) LogInfof(format string, args ...interface{}) {
	aLog.logF(WindIneLogStyleInfo, format, args...)
}

// LogErrorf format格式化log--error信息
func (aLog *WindIneLog) LogErrorf(format string, args ...interface{}) {
	aLog.logF(WindIneLogStyleError, format, args...)
}

// LogDebugf format格式化log--debug信息
func (aLog *WindIneLog) LogDebugf(format string, args ...interface{}) {
	aLog.logF(WindIneLogStyleDebug, format, args...)
}

// LoWindIneracef format格式化log--Trace信息
func (aLog *WindIneLog) LoWindIneracef(format string, args ...interface{}) {
	aLog.logF(WindIneLogStyleTrace, format, args...)
}

// LogFatalf format格式化log--Fatal信息 !!!慎用，使用后程序会退出!!!
func (aLog *WindIneLog) LogFatalf(format string, args ...interface{}) {
	aLog.logF(WindIneLogStyleFatal, format, args...)
}

// LogWarnf format格式化log--Warning信息
func (aLog *WindIneLog) LogWarnf(format string, args ...interface{}) {
	aLog.logF(WindIneLogStyleWarning, format, args...)
}

// determineRotationTime 辅助函数：根据日志保存类型决定轮转时间
func determineRotationTime(logSaveType WindIneLogSaveType) time.Duration {
	switch logSaveType {
	case WindIneLogSaveHours:
		return time.Hour
	case WindIneLogSaveTypeDays:
		return time.Hour * 24
	default:
		return time.Hour * 24 // 默认按天轮转
	}
}

func newLogSaveHandler(WindIneLog *WindIneLog) (rotateLogger *rotatelogs.RotateLogs) {
	// 确保日志目录存在
	err := os.MkdirAll(WindIneLog.logDir, 0755)
	if err != nil {
		fmt.Printf("%s", err)
		return nil
	}
	logFilePath := fmt.Sprintf("%s/run", WindIneLog.logDirWithDate)
	linkLogFilePath := fmt.Sprintf("%s/run", WindIneLog.logDir)

	/* 日志轮转相关函数
	   `WithLinkName` 为最新的日志建立软连接
	   `WithRotationTime` 设置日志分割的时间，隔多久分割一次
	   WithMaxAge 和 WithRotationCount二者只能设置一个
	    `WithMaxAge` 设置文件清理前的最长保存时间
	    `WithRotationCount` 设置文件清理前最多保存的个数
	*/
	writer, err := rotatelogs.New(
		logFilePath+".%Y-%m-%d_%H",
		rotatelogs.WithLinkName(linkLogFilePath),
		rotatelogs.WithMaxAge(time.Duration(instanceConfig().logMaxSaveDays)*24*time.Hour),
		rotatelogs.WithRotationTime(determineRotationTime(instanceConfig().logSaveType)),
	)
	if err != nil {
		// 处理错误
		fmt.Println("Error setting up log writer:", err)
		return nil
	}
	return writer
}

// NewWindIneLog 添加WindIneLog模块
func NewWindIneLog(modelName string) *WindIneLog {
	currentTime := time.Now().UTC()

	WindIneLog := &WindIneLog{
		modelName:      modelName,
		logDir:         fmt.Sprintf("%s/%s", instanceConfig().productLogDir, modelName),
		logDirWithDate: fmt.Sprintf("%s/%s/%s", instanceConfig().productLogDir, modelName, currentTime.Format("2006-01-02")),
		logger:         logrus.New(),
		entryTime:      currentTime,
		lastCheckTime:  currentTime,
	}

	// 初始化日志设置（代码简化，具体初始化逻辑可以根据需要调整）
	WindIneLog.logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	WindIneLog.logger.SetLevel(logrus.TraceLevel)

	// 设置默认日志输出为控制台
	WindIneLog.logger.SetOutput(os.Stdout)

	// 根据LogLevel设置logrus的日志级别
	switch currentLogConfig.logLeve {
	case WindIneLogStyleFatal:
		WindIneLog.logger.SetLevel(logrus.FatalLevel)
	case WindIneLogStyleTrace:
		WindIneLog.logger.SetLevel(logrus.TraceLevel)
	case WindIneLogStyleInfo:
		WindIneLog.logger.SetLevel(logrus.InfoLevel)
	case WindIneLogStyleWarning:
		WindIneLog.logger.SetLevel(logrus.WarnLevel)
	case WindIneLogStyleError:
		WindIneLog.logger.SetLevel(logrus.ErrorLevel)
	case WindIneLogStyleDebug:
		WindIneLog.logger.SetLevel(logrus.DebugLevel)
	default:
		WindIneLog.logger.SetLevel(logrus.InfoLevel)
	}

	// 设置日志输出，可以根据EnableSaveLogFile和其他参数来配置
	// （省略了日志轮转和文件输出的设置，可以直接使用SetupLoWindIneools中相关的代码）
	//	设置Log
	if instanceConfig().enableSaveLogFile == true {
		rLog := newLogSaveHandler(WindIneLog)
		WindIneLog.logger.SetOutput(rLog)
	}
	return WindIneLog
}

// LogInfof format格式化log--info信息
func LogInfof(format string, args ...interface{}) {
	setupDefaultLog().logF(WindIneLogStyleInfo, format, args...)
}

// LogErrorf format格式化log--error信息
func LogErrorf(format string, args ...interface{}) {
	setupDefaultLog().logF(WindIneLogStyleError, format, args...)
}

// LogDebugf format格式化log--debug信息
func LogDebugf(format string, args ...interface{}) {
	setupDefaultLog().logF(WindIneLogStyleDebug, format, args...)
}

// LoWindIneracef format格式化log--Trace信息
func LoWindIneracef(format string, args ...interface{}) {
	setupDefaultLog().logF(WindIneLogStyleTrace, format, args...)
}

// LogFatalf format格式化log--Fatal信息 !!!慎用，使用后程序会退出!!!
func LogFatalf(format string, args ...interface{}) {
	setupDefaultLog().logF(WindIneLogStyleFatal, format, args...)
}

// LogWarnf format格式化log--Warning信息
func LogWarnf(format string, args ...interface{}) {

	setupDefaultLog().logF(WindIneLogStyleWarning, format, args...)
}

// SetupLoWindIneools 初始化日志
func SetupLoWindIneools(productName string, enableSaveLogFile bool, logLeve WindIneLogStyle, logMaxSaveDays int64, logSaveType WindIneLogSaveType, productLogDir string) {
	setupComplete = false

	instanceConfig().productName = productName
	instanceConfig().enableSaveLogFile = enableSaveLogFile
	instanceConfig().logLeve = logLeve
	instanceConfig().logMaxSaveDays = logMaxSaveDays
	instanceConfig().logSaveType = logSaveType
	instanceConfig().productLogDir = productLogDir

	if productLogDir == "" {
		if runtime.GOOS == "linux" {
			instanceConfig().productLogDir = fmt.Sprintf("%s/%s", "/var/log", strings.ToLower(instanceConfig().productName))
		} else {
			instanceConfig().productLogDir = "./logs"
		}
	}

	if mainLog == nil {
		setupDefaultLog()
	}
}
