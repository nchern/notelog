package searcher

import (
	"bufio"
	"os"
)

// GetLastNthResult returns nth result from last saved search results
func GetLastNthResult(notes Notes, n int) (string, error) {
	f, err := os.Open(notes.MetadataFilename(lastResultsFile))
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil // just an empty result, same if we ask for non-existing item
		}
		return "", err
	}
	defer f.Close()

	i := 1
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if i == n {
			return scanner.Text(), nil
		}
		i++
	}
	return "", scanner.Err()
}
