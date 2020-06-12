package neoproxy

func NotifyDosage(flow *Flow) {
	flow.UpdateDosage()
	flow.SendDosage()
}

func NotifyNews(flow *Flow) {
	flow.CrawlLastNews()
	flow.SendLastNews()
}
