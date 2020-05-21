package sync

// MergeErrorChannels takes n channels as parameters, attaches them to a new error channel and returns it
func MergeErrorChannels(channles ...<-chan error) <-chan error {
	out := make(chan error)

	for _, channel := range channles {
		go func(ch <-chan error) {
			for err := range ch {
				out <- err
			}
		}(channel)
	}
	return out
}
