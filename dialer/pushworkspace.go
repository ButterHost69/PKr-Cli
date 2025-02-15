package dialer

import (
	"context"
	"time"

	"github.com/ButterHost69/PKr-cli/logger"
	"github.com/ButterHost69/PKr-cli/pb"
	"google.golang.org/grpc"
)

// [ ] Do using Channels ->
//		[ ] Client Creation as Channel & go func ->
// 											[ ] Send Request as go func than use channel or global to count

// [ ] Test Code
// HTTP Handler
// func PushToConnections(workspace_name string, connnection_ip []string) (int, error) {
// 	count := 0
// 	for _, ip := range connnection_ip {
// 		url := fmt.Sprintf("http://%s/new/workspace/%s", ip, workspace_name)
// 		res, err := http.Get(url)
// 		if err != nil {
// 			return -1, err
// 		}

// 		if res.StatusCode == 200 {
// 			count += 1
// 		}
// 	}

// 	return count, nil
// }

// [ ] Test Function
func PushToConnections(workspace_name string, connnection_ip []string, workspace_logger *logger.WorkspaceLogger) int {
	successful_count := 0

	for _, workspace_ip := range connnection_ip {
		conn, err := grpc.NewClient(workspace_ip, grpc.WithInsecure())
		if err != nil {
			err_log := "Failed to Establish Connection with " + workspace_ip + " while sending Push Notification"
			// config.AddLogEntry(workspace_name, err_log)
			workspace_logger.Info(workspace_name, err_log)
			continue
		}

		client := pb.NewBackgroundServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
		defer cancel()

		response, err := client.NotifyPush(ctx, &pb.NotifyPushRequest{
			WorkspaceName: workspace_name,
		})

		if err != nil {
			err_log := "Error in Response from " + workspace_ip + " while sending Push Notification"
			// config.AddLogEntry(workspace_name, err_log)
			workspace_logger.Info(workspace_name, err_log)
			continue
		}

		if response.Response == 200 {
			successful_count += 1
		}
		conn.Close()
	}
	return successful_count
}
