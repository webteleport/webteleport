package edge

// Subscribe to incoming requests
type Subscriber interface {
	Subscribe(u Upgrader)
}
