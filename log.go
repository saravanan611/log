package log

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

// ==================================================
// Global map to store requestID for each goroutine
// ==================================================
var (
	GEnvLogVal    bool
	gRequestIDMap = sync.Map{}
	gInfo         = log.New(os.Stdout, "INFO: ", log.LstdFlags|log.Lshortfile)
	gDebug        = log.New(os.Stdout, "DEBUG: ", log.LstdFlags|log.Lshortfile)
	gErr          = log.New(os.Stderr, "ERROR: ", log.LstdFlags|log.Lshortfile)
	gOwnErr       = log.New(os.Stderr, "ERROR: ", log.LstdFlags)
)

/* -----------------------custom error ------------------- */
type ownErr struct {
	lFileInfo string
	lErr      string
}

/* -----------------------err to string-------------------- */
func (pErr *ownErr) Error() string {
	return pErr.lErr
}

/* ---------------- enable info and debug log -------------- */
func EnableInfo() {
	GEnvLogVal = true
}

// ---------- ERROR WRAPPER ----------
func Error(pErr any) error {

	if lErr, lOk := pErr.(*ownErr); lOk {
		return lErr
	}

	_, lFile, lLine, _ := runtime.Caller(1)
	lStrArray := strings.Split(lFile, "/")
	lFilename := lStrArray[len(lStrArray)-2] + "/" + lStrArray[len(lStrArray)-1]
	return &ownErr{lFileInfo: fmt.Sprintf("%s:%d", lFilename, lLine), lErr: fmt.Sprintf("%v", pErr)}

}

// ============================================
// Logging functions - NO PARAMETERS NEEDED!
// ============================================

// ---------- INFO LOGGER ----------
func Info(format string, args ...any) {
	if !GEnvLogVal {
		return
	}
	requestID := GetRequestID()
	msg := fmt.Sprintf(format, args...)
	gInfo.Output(2, fmt.Sprintf("[ReqID: %s] %s", requestID, msg))
}

// ---------- Debug LOGGER ----------
func Debug(format string, args ...any) {
	if !GEnvLogVal {
		return
	}
	requestID := GetRequestID()
	msg := fmt.Sprintf(format, args...)
	gDebug.Output(2, fmt.Sprintf("[ReqID: %s] %s", requestID, msg))
}

// ---------- Err LOGGER ----------
func Err(pErr any) {
	requestID := GetRequestID()
	if lErr, lok := pErr.(*ownErr); lok {
		gOwnErr.Printf("%s [ReqID: %s] %s", lErr.lFileInfo, requestID, lErr.Error())
		return
	}
	// use Output(counter+2, ...) so call depth aligns properly
	gErr.Output(2, fmt.Sprintf("[ReqID: %s] %v", requestID, pErr))
}

// Get current goroutine ID
func getGoroutineID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = b[len("goroutine "):]
	b = b[:strings.IndexByte(string(b), ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

// Set request ID for current goroutine
func SetRequestID(requestID string) {
	gid := getGoroutineID()
	gRequestIDMap.Store(gid, requestID)
}

// Get request ID for current goroutine
func GetRequestID() string {
	gid := getGoroutineID()
	// log.Println("gid :", gid)
	if id, ok := gRequestIDMap.Load(gid); ok {
		return id.(string)
	}
	return "UNKNOWN"
}

// Clear request ID when done
func ClearRequestID() {
	gid := getGoroutineID()
	gRequestIDMap.Delete(gid)
}
