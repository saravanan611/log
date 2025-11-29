package log

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/google/uuid"
)

type ownStr string
/* -----------------------Sid Key Name------------------- */
const GateKey ownStr = "G-key"

/* ----------------------- log session------------------- */
type LogStruct struct {
	Uid       string
	lInfoFlag bool
}

/* -----------------------custom error ------------------- */
type ownErr struct {
	lFileInfo string
	lErr      string
}

/* -----------------------err to string-------------------- */
func (pErr *ownErr) Error() string {
	return pErr.lErr
}

/* --------------------Log Ination ---------------- */
func Init() *LogStruct {
	return &LogStruct{
		Uid:       strings.ReplaceAll(uuid.New().String(), "-", ""),
		lInfoFlag: strings.EqualFold(os.Getenv("InfoFlog"), "Y"),
	}
}

/*-----------------------------read req id -----------------------  */

func ReqInit(pReq *http.Request) *LogStruct {
	lUid, lOk := pReq.Context().Value(GateKey).(string)
	if !lOk || lUid == "" {
		lUid = strings.ReplaceAll(uuid.New().String(), "-", "")
	}
	return &LogStruct{
		Uid:       lUid,
		lInfoFlag: strings.EqualFold(os.Getenv("InfoFlog"), "Y"),
	}
}

/* --------------------Basic log variable ---------------- */
var (
	err    = log.New(os.Stderr, "ERROR: ", log.LstdFlags|log.Lshortfile)
	info   = log.New(os.Stdout, "INFO: ", log.LstdFlags|log.Lshortfile)
	lOnErr = log.New(os.Stdout, "ERROR: ", log.LstdFlags)
)

// ---------- INFO LOGGER ----------
func (lId *LogStruct) Info(pMsg ...any) {
	if !lId.lInfoFlag {
		return
	}
	msg := fmt.Sprint(pMsg...) // cleaner than fmt.Sprintf("%v")
	info.Output(2, fmt.Sprintf("[%s] %s", lId.Uid, msg))
}

// ---------- ERROR LOGGER ----------
func (lId *LogStruct) Err(pErr any) {
	if lErr, lok := pErr.(*ownErr); lok {
		lOnErr.Printf("%s [%s] %s", lErr.lFileInfo, lId.Uid, lErr.Error())
		return
	}
	// use Output(counter+2, ...) so call depth aligns properly
	err.Output(2, fmt.Sprintf("[%s] %v", lId.Uid, pErr))
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
