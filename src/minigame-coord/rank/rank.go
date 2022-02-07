package rank

//定时清理排行榜

func Register() {
	go schedule_clear()
}
