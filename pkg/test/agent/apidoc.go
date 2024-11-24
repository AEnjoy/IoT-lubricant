package agent

import (
	"context"
	"fmt"

	testMeta "github.com/AEnjoy/IoT-lubricant/pkg/test"
	"github.com/AEnjoy/IoT-lubricant/protobuf/agent"
)

func TestGetOpenApiDoc(cli agent.EdgeServiceClient, docType *agent.OpenapiDocType, content *[]byte) *testMeta.Result {
	fmt.Println("Test_GetOpenApiDoc:")
	if docType == nil {
		// test all
		fmt.Print("--Test get all doc(should be null because of no setting):")
		doc, err := cli.GetOpenapiDoc(context.Background(),
			&agent.GetOpenapiDocRequest{
				AgentID: testMeta.AgentID,
				DocType: agent.OpenapiDocType_All,
			})
		if err != nil {
			return &testMeta.Result{Success: false, Message: err.Error()}
		}
		if len(doc.GetEnableFile()) == 0 && len(doc.GetOriginalFile()) == 0 {
			fmt.Println("Success")
			return &testMeta.Result{Success: true}
		} else {
			return &testMeta.Result{Success: false, Message: "get the content when no doc is set"}
		}
	}
	fmt.Print("--Test get doc: ")
	switch *docType {
	case agent.OpenapiDocType_enableFile:
		doc, err := cli.GetOpenapiDoc(context.Background(), &agent.GetOpenapiDocRequest{
			AgentID: testMeta.AgentID,
			DocType: agent.OpenapiDocType_enableFile,
		})
		if err != nil {
			return &testMeta.Result{Success: false, Message: err.Error()}
		}
		if len(doc.GetEnableFile()) == 0 {
			return &testMeta.Result{Success: false, Message: "can't get the content when doc has been set"}
		} else if content != nil && len(doc.GetEnableFile()) == len(*content) {
			fmt.Println("Success")
			return &testMeta.Result{Success: true}
		} else {
			return &testMeta.Result{Success: false, Message: "Inconsistency between actual and expected"}
		}
	case agent.OpenapiDocType_originalFile:
		doc, err := cli.GetOpenapiDoc(context.Background(), &agent.GetOpenapiDocRequest{
			AgentID: testMeta.AgentID,
			DocType: agent.OpenapiDocType_originalFile,
		})
		if err != nil {
			return &testMeta.Result{Success: false, Message: err.Error()}
		}
		if len(doc.GetOriginalFile()) == 0 {
			return &testMeta.Result{Success: false, Message: "can't get the content when doc has been set"}
		} else if content != nil && len(doc.GetOriginalFile()) == len(*content) {
			fmt.Println("Success")
			return &testMeta.Result{Success: true}
		} else {
			return &testMeta.Result{Success: false, Message: "Inconsistency between actual and expected"}
		}
	default:
		return &testMeta.Result{Success: false, Message: "unknown doc type"}
	}
}
