package gomemcache

// 根据key选择响应的服务节点
type PeerPicker interface {
	PickPeer(key string) (PeerGetter, bool)
}

// 每个节点都必须要实现 PeerGetter 接口
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}