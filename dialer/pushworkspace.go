package dialer

import (
	"fmt"
	"net/http"
)

// [ ] Do using Channels ->
//		[ ] Client Creation as Channel & go func ->
// 											[ ] Send Request as go func than use channel or global to count

//  [ ] Test Code
func PushToConnections(workspace_name string, connnection_ip []string) (int, error) {
	count := 0
	for _, ip := range connnection_ip {
		url := fmt.Sprintf("http://%s/new/workspace/%s", ip, workspace_name)
		res, err := http.Get(url)
		if err != nil {
			return -1, err
		}

		if res.StatusCode == 200 {
			count += 1
		}
	}
		
	return count, nil
}