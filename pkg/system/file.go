package system

func (ssi *SystemServiceImpl) OverwriteFile(path string, data []byte) error {
	return ssi.writeFile(path, data, 0644)
}

func (ssi *SystemServiceImpl) ReadFile(path string) ([]byte, error) {
	return ssi.readFile(path)
}
