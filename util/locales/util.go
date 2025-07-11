package locales

func List() []string {
	return NativeLocalesMap.ToStream().KeySetStream().ToList()
}
