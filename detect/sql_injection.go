package detect

type SQLInjectionDetector struct{}

func (s *SQLInjectionDetector) Setup() error {
	return nil
}

func (s *SQLInjectionDetector) Detect() (int, error) {
	score := 0

	return score, nil
}
