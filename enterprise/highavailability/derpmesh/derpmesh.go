package derpmesh

import (
	"context"
	"sync"

	"golang.org/x/xerrors"
	"tailscale.com/derp"
	"tailscale.com/derp/derphttp"
	"tailscale.com/types/key"

	"github.com/coder/coder/tailnet"

	"cdr.dev/slog"
)

// New constructs a new mesh for DERP servers.
func New(logger slog.Logger, server *derp.Server) *Mesh {
	return &Mesh{
		logger: logger,
		server: server,
		ctx:    context.Background(),
		closed: make(chan struct{}),
		active: make(map[string]context.CancelFunc),
	}
}

type Mesh struct {
	logger slog.Logger
	server *derp.Server
	ctx    context.Context

	mutex  sync.Mutex
	closed chan struct{}
	active map[string]context.CancelFunc
}

// SetAddresses performs a diff of the incoming addresses and adds
// or removes DERP clients from the mesh.
func (m *Mesh) SetAddresses(addresses []string) {
	total := make(map[string]struct{}, 0)
	for _, address := range addresses {
		total[address] = struct{}{}
		added, err := m.addAddress(address)
		if err != nil {
			m.logger.Error(m.ctx, "failed to add address", slog.F("address", address), slog.Error(err))
			continue
		}
		if added {
			m.logger.Debug(m.ctx, "added mesh address", slog.F("address", address))
		}
	}

	m.mutex.Lock()
	for address := range m.active {
		_, found := total[address]
		if found {
			continue
		}
		removed := m.removeAddress(address)
		if removed {
			m.logger.Debug(m.ctx, "removed mesh address", slog.F("address", address))
		}
	}
	m.mutex.Unlock()
}

// addAddress begins meshing with a new address.
// It's expected that this is a full HTTP address with a path.
// e.g. http://127.0.0.1:8080/derp
func (m *Mesh) addAddress(address string) (bool, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	_, isActive := m.active[address]
	if isActive {
		return false, nil
	}
	client, err := derphttp.NewClient(m.server.PrivateKey(), address, tailnet.Logger(m.logger))
	if err != nil {
		return false, xerrors.Errorf("create derp client: %w", err)
	}
	client.MeshKey = m.server.MeshKey()
	ctx, cancelFunc := context.WithCancel(m.ctx)
	closed := make(chan struct{})
	closeFunc := func() {
		cancelFunc()
		_ = client.Close()
		<-closed
	}
	m.active[address] = closeFunc
	go func() {
		defer close(closed)
		client.RunWatchConnectionLoop(ctx, m.server.PublicKey(), tailnet.Logger(m.logger), func(np key.NodePublic) {
			m.server.AddPacketForwarder(np, client)
		}, func(np key.NodePublic) {
			m.server.RemovePacketForwarder(np, client)
		})
	}()
	return true, nil
}

// removeAddress stops meshing with a given address.
func (m *Mesh) removeAddress(address string) bool {
	cancelFunc, isActive := m.active[address]
	if isActive {
		cancelFunc()
	}
	return isActive
}

// Close ends all active meshes with the DERP server.
func (m *Mesh) Close() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	select {
	case <-m.closed:
		return nil
	default:
	}
	close(m.closed)
	for _, cancelFunc := range m.active {
		cancelFunc()
	}
	return nil
}
