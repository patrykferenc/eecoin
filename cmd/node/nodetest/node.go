package nodetest

import "os"

// TODO: implement test harness

func NewPeersFile(peers []string) string {
	testDir, err := os.MkdirTemp(os.TempDir(), "eecoin_node_test_peers")
	if err != nil {
		panic(err)
	}

	file, err := os.Create(testDir + "/peers")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for _, peer := range peers {
		_, err := file.WriteString(peer + "\n")
		if err != nil {
			panic(err)
		}
	}

	return file.Name()
}
