package main

func main() {

	alertRule := AlertRule{}
	notification := Notification{}
	tpsAlertHandler := TpsAlertHandler{alertHandler{
		rule:         alertRule,
		notification: notification,
	}}
	errAlertHandler := ErrAlertHandler{alertHandler{
		rule:         alertRule,
		notification: notification,
	}}
	timeoutAlertHandler := TimeoutAlertHandler{alertHandler{
		rule:         alertRule,
		notification: notification,
	}}
	alert := AlertChange{}
	alert.addHandler(tpsAlertHandler)
	alert.addHandler(errAlertHandler)
	alert.addHandler(timeoutAlertHandler)
	alert.check(ApiStatInfo{})
}

func Init() {

}

type AlertChange struct {
	handlers []AlertHandler
}

func (ar *AlertChange) addHandler(handler AlertHandler) {
	ar.handlers = append(ar.handlers, handler)
}

func (ar *AlertChange) check(apiStatInfo ApiStatInfo) {
	for _, handler := range ar.handlers {
		handler.check(apiStatInfo)
	}
}

type AlertHandler interface {
	check(apiStatInfo ApiStatInfo)
}

type alertHandler struct {
	rule         AlertRule
	notification Notification
}

type TpsAlertHandler struct {
	alertHandler
}

func (tah TpsAlertHandler) check(apiStatInfo ApiStatInfo) {
	tps := apiStatInfo.getRequestCount() / apiStatInfo.getDurationSeconds()
	if tps > tah.rule.getMaxTps() {
		tah.notification.notify("l1", "...")
	}
}

type ErrAlertHandler struct {
	alertHandler
}

func (tah ErrAlertHandler) check(apiStatInfo ApiStatInfo) {
	if tah.rule.getMaxErrorCount() < apiStatInfo.getErrCount() {
		tah.notification.notify("l2", "...")
	}
}

type TimeoutAlertHandler struct {
	alertHandler
}

func (tah TimeoutAlertHandler) check(apiStatInfo ApiStatInfo) {
	if tah.rule.getMaxTimeoutTps() < apiStatInfo.getTimeout() {
		tah.notification.notify("l3", "...")
	}
}

type ApiStatInfo struct {
	api             string
	requestCount    int
	errCount        int
	durationSeconds int
	timeouts        int
}

func (asi ApiStatInfo) getRequestCount() int {
	return asi.requestCount
}

func (asi ApiStatInfo) getErrCount() int {
	return asi.errCount
}

func (asi ApiStatInfo) getDurationSeconds() int {
	return asi.durationSeconds
}

func (asi ApiStatInfo) getTimeout() int {
	return asi.timeouts
}
