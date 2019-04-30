package common

import (
	"io/ioutil"
	"path/filepath"
)

func GetPcapAndSnapShot(path string) (string, string, []string) {
	dir := filepath.Dir(path)
	pcapPath := filepath.Join(dir, "dump.pcap")
	pcapSha1 := GetSha1ByPath(pcapPath)
	snapPath := filepath.Join(dir, "shots")
	reportJsonPath := filepath.Join(dir, "reports", "report.json")
	reportJsonSha1 := GetSha1ByPath(reportJsonPath)

	dir_list, _ := ioutil.ReadDir(snapPath)
	snapSha1s := make([]string, 0)
	for _, v := range dir_list {
		snapPath := filepath.Join(snapPath, v.Name())
		snapSha1 := GetSha1ByPath(snapPath)
		snapSha1s = append(snapSha1s, snapSha1)
	}
	return pcapSha1, reportJsonSha1, snapSha1s

}
