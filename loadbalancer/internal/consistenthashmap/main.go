package consistenthashmap

const (
	Slots = 512
	K     = 9
)

type ConsistentHashMap struct {
	virtualServers [Slots]int
}

func hashRequest(i int) int {
	return i*i + 2*i + 17
}

func hashVirtualServer(i, j int) int {
	return i*i + j*j + 2*j + 25
}

func (hm *ConsistentHashMap) Init() {
	for i := 0; i < Slots; i++ {
		hm.virtualServers[i] = -1
	}
}

func (hm *ConsistentHashMap) findEmptyServerSlot(hashValue int) int {
	slot := hashValue % Slots
	for hm.virtualServers[slot] != -1 {
		slot = (slot + 1) % Slots
	}
	return slot
}

func (hm *ConsistentHashMap) AddServer(serverID int) {
	for j := 0; j < K; j++ {
		slot := hm.findEmptyServerSlot(hashVirtualServer(serverID, j))
		hm.virtualServers[slot] = slot
	}
}

func (hm *ConsistentHashMap) GetServerForRequest(requestID int) int {
	slot := hashRequest(requestID) % Slots
	for hm.virtualServers[slot] == -1 {
		slot = (slot + 1) % Slots
	}
	return hm.virtualServers[slot]
}
