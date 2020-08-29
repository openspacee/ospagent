package container

import (
	"github.com/openspacee/ospagent/pkg/container/resource"
	"github.com/openspacee/ospagent/pkg/kubernetes"
	"github.com/openspacee/ospagent/pkg/utils"
	"github.com/openspacee/ospagent/pkg/websocket"
)

const (
	LIST     = "list"
	GET      = "get"
	DELETE   = "delete"
	UPDATE   = "update"
	EXEC     = "exec"
	STDIN    = "stdin"
	OPENLOG  = "openLog"
	CLOSELOG = "closeLog"
)

type Handler func(interface{}) *utils.Response

type ActionHandler map[string]Handler

type ResourceActions struct {
	KubeClient            *kubernetes.KubeClient
	ResourceActionHandler map[string]ActionHandler
}

func NewResourceActions(kubeClient *kubernetes.KubeClient, sendResponse websocket.SendResponse) *ResourceActions {
	actionHandlers := make(map[string]ActionHandler)

	watch := resource.NewWatchResource(sendResponse)
	watchActions := ActionHandler{
		GET: watch.WatchAction,
	}
	actionHandlers["watch"] = watchActions

	pod := resource.NewPod(kubeClient, sendResponse, watch)
	podActions := ActionHandler{
		LIST:     pod.List,
		GET:      pod.Get,
		EXEC:     pod.Exec,
		STDIN:    pod.ExecStdIn,
		OPENLOG:  pod.OpenLog,
		CLOSELOG: pod.CloseLog,
		DELETE:   pod.Delete,
		UPDATE:   pod.Update,
	}
	actionHandlers["pod"] = podActions

	ns := resource.NewNamespace(kubeClient, sendResponse, watch)
	nsActions := ActionHandler{
		LIST: ns.List,
	}
	actionHandlers["namespace"] = nsActions

	node := resource.NewNode(kubeClient, sendResponse)
	nodeActions := ActionHandler{
		LIST: node.List,
	}
	actionHandlers["node"] = nodeActions

	configMap := resource.NewConfigMap(kubeClient, sendResponse)
	configMapActions := ActionHandler{
		LIST: configMap.List,
	}
	actionHandlers["configMap"] = configMapActions

	return &ResourceActions{
		KubeClient:            kubeClient,
		ResourceActionHandler: actionHandlers,
	}
}

func (r *ResourceActions) GetRequestHandler(resource string, action string) Handler {
	return r.ResourceActionHandler[resource][action]
}
