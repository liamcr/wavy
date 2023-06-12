package wav

// SpeedUp speeds up the wav file by a specified factor.
// This will not only increase the speed of the audio, but
// will increase the pitch as well.
func (w *Wav) SpeedUp(factor float32) {
	w.SampleRate = uint32(factor * float32(w.SampleRate))
}

// SlowDown slows up the wav file by a specified factor.
// This will not only decrease the speed of the audio, but
// will lower the pitch as well.
func (w *Wav) SlowDown(factor float32) {
	w.SampleRate = uint32(float32(w.SampleRate) / factor)
}