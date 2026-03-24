# 🪵 `log` — Reference Documentation

A minimal, goroutine-aware logging and error wrapping package for Go.

---

## 📦 Import

```go
import "your-module/log"
```

---

## ⚙️ Configuration

### Enable Info & Debug Logs

```go
log.EnableInfo()
```

* By default, **info/debug logs are disabled**
* Calling `EnableInfo()` enables both

---

## 🔑 Request ID Management

### Set Request ID (per goroutine)

```go
log.SetRequestID("req-123")
```

### Get Current Request ID

```go
id := log.GetRequestID()
```

Returns:

* Assigned ID → `"req-123"`
* If not set → `"UNKNOWN"`

### Clear Request ID

```go
log.ClearRequestID()
```

---

## 🧩 Logging Functions

### Info

```go
log.Info(format string, args ...any)
```

* Printed only when enabled
* Output: `stdout`

---

### Debug

```go
log.Debug(format string, args ...any)
```

* Printed only when enabled
* Output: `stdout`

---

### Error Log

```go
log.Err(err any)
```

* Always printed
* Output: `stderr`
* Handles wrapped errors automatically

---

## ⚡ Error Handling

### Wrap Error

```go
err := log.Error(err)
```

### Create Error

```go
err := log.Error("message")
```

### Format Error

```go
err := log.Error(fmt.Sprintf("failed: %v", err))
```

---

## 🧠 Behavior

### Error Wrapping

* Adds:

  * File name
  * Line number
* Prevents double wrapping

---

### Request ID Binding

* Stored per goroutine using internal map
* Automatically included in logs

---

## 🧾 Log Format

### Info / Debug

```
INFO: 2026/03/24 file.go:10: [ReqID: req-123] message
DEBUG: 2026/03/24 file.go:10: [ReqID: req-123] message
```

### Error (wrapped)

```
ERROR: file.go:25 [ReqID: req-123] message
```

### Error (normal)

```
ERROR: 2026/03/24 file.go:25: [ReqID: req-123] message
```

---

## 📊 Summary

| Function           | Description                    |
| ------------------ | ------------------------------ |
| `EnableInfo()`     | Enable info & debug logs       |
| `Info()`           | Print info log                 |
| `Debug()`          | Print debug log                |
| `Err()`            | Print error log                |
| `Error()`          | Wrap or create error           |
| `SetRequestID()`   | Assign request ID to goroutine |
| `GetRequestID()`   | Retrieve current request ID    |
| `ClearRequestID()` | Remove request ID              |

