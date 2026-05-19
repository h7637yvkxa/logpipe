package filter

// Pipeline reads normalized log entries from in, applies the filter,
// and forwards matching entries to the returned channel.
func Pipeline(in <-chan map[string]interface{}, f *Filter) <-chan map[string]interface{} {
	out := make(chan map[string]interface{}, 64)
	go func() {
		defer close(out)
		for entry := range in {
			if f.Match(entry) {
				out <- entry
			}
		}
	}()
	return out
}
