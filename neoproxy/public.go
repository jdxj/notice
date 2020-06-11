package neoproxy

var (
	flow = NewFlow()
)

func NotifyDosage() {
	flow.UpdateDosage()
	flow.SendDosage()
}

func NotifyNews() {
	flow.CrawlLastNews()
	flow.SendLastNews()
}
