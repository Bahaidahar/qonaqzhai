package ws

// HubPublisher adapts a Hub into the thread.Publisher port shape.
type HubPublisher struct{ Hub *Hub }

// Publish wraps payload into an envelope and pushes to each userID's connections.
func (p HubPublisher) Publish(event string, payload any, userIDs ...string) {
	if p.Hub == nil {
		return
	}
	p.Hub.SendToUsers(Envelope{Op: event, Data: payload}, userIDs...)
}
