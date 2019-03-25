package build

// 3rd party mocks
//go:generate mockgen -destination=../mocks/mock_k8s.io_apimachinery_pkg_watch.go -package=mocks k8s.io/apimachinery/pkg/watch Interface
//go:generate mockgen -destination=../mocks/mock_k8s.io_client-go_kubernetes_typed_core_v1.go -package=mocks k8s.io/client-go/kubernetes/typed/core/v1 CoreV1Interface,NamespaceInterface
//go:generate mockgen -destination=../mocks/mock_k8s.io_client-go_kubernetes_typed_extensions_v1beta1.go -package=mocks k8s.io/client-go/kubernetes/typed/extensions/v1beta1 ExtensionsV1beta1Interface
//go:generate mockgen -destination=../mocks/mock_k8s.io_kubernetes_pkg_kubectl_cmd_util.go -package=mocks k8s.io/kubernetes/pkg/kubectl/cmd/util Factory

// local mocks
//go:generate mockgen -destination=../mocks/mocks_kubernetes_client_clientFactory.go -package=mocks -source=../kubernetes/client/clientFactory.go
//go:generate mockgen -destination=../mocks/mocks_kubernetes_core_watcherFetcher.go -package=mocks -source=../kubernetes/core/watcherFetcher.go
//go:generate mockgen -destination=../mocks/mocks_kubernetes_core_watcher.go -package=mocks -source=../kubernetes/core/watcher.go
