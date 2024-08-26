package main

type Alert struct {
	rule         AlertRule
	notification Notification
}

type AlertRule struct {
}

func (ar *AlertRule) getMatchedRule(api string) *AlertRule {
	return nil
}

func (ar *AlertRule) getMaxTps() int {
	return 0
}

func (ar *AlertRule) getMaxErrorCount() int {
	return 0
}

func (ar *AlertRule) getMaxTimeoutTps() int {
	return 0
}

type Notification struct {
}

func (n *Notification) notify(level, msg string) {

}

func NewAlert(rule AlertRule, notification Notification) *Alert {
	return &Alert{
		rule:         rule,
		notification: notification,
	}
}

func (a *Alert) Check(api string, requestCount, errCount, durationSeconds, timeouts int) {
	tps := requestCount / durationSeconds
	if tps > a.rule.getMatchedRule(api).getMaxTps() {
		a.notification.notify("l1", "...")
	}
	if errCount > a.rule.getMatchedRule(api).getMaxErrorCount() {
		a.notification.notify("l2", "...")
	}
	timeoutTps := timeouts / durationSeconds
	if timeoutTps > a.rule.getMatchedRule(api).getMaxTimeoutTps() {
		a.notification.notify("l3", "...")
	}
}
