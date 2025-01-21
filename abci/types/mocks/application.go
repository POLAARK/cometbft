// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	v1 "github.com/cometbft/cometbft/api/cometbft/abci/v1"
)

// Application is an autogenerated mock type for the Application type
type Application struct {
	mock.Mock
}

// ApplySnapshotChunk provides a mock function with given fields: ctx, req
func (_m *Application) ApplySnapshotChunk(ctx context.Context, req *v1.ApplySnapshotChunkRequest) (*v1.ApplySnapshotChunkResponse, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for ApplySnapshotChunk")
	}

	var r0 *v1.ApplySnapshotChunkResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ApplySnapshotChunkRequest) (*v1.ApplySnapshotChunkResponse, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ApplySnapshotChunkRequest) *v1.ApplySnapshotChunkResponse); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.ApplySnapshotChunkResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.ApplySnapshotChunkRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CheckTx provides a mock function with given fields: ctx, req
func (_m *Application) CheckTx(ctx context.Context, req *v1.CheckTxRequest) (*v1.CheckTxResponse, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for CheckTx")
	}

	var r0 *v1.CheckTxResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.CheckTxRequest) (*v1.CheckTxResponse, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.CheckTxRequest) *v1.CheckTxResponse); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.CheckTxResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.CheckTxRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Commit provides a mock function with given fields: ctx, req
func (_m *Application) Commit(ctx context.Context, req *v1.CommitRequest) (*v1.CommitResponse, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for Commit")
	}

	var r0 *v1.CommitResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.CommitRequest) (*v1.CommitResponse, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.CommitRequest) *v1.CommitResponse); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.CommitResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.CommitRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ExtendVote provides a mock function with given fields: ctx, req
func (_m *Application) ExtendVote(ctx context.Context, req *v1.ExtendVoteRequest) (*v1.ExtendVoteResponse, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for ExtendVote")
	}

	var r0 *v1.ExtendVoteResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ExtendVoteRequest) (*v1.ExtendVoteResponse, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ExtendVoteRequest) *v1.ExtendVoteResponse); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.ExtendVoteResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.ExtendVoteRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FinalizeBlock provides a mock function with given fields: ctx, req
func (_m *Application) FinalizeBlock(ctx context.Context, req *v1.FinalizeBlockRequest) (*v1.FinalizeBlockResponse, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for FinalizeBlock")
	}

	var r0 *v1.FinalizeBlockResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.FinalizeBlockRequest) (*v1.FinalizeBlockResponse, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.FinalizeBlockRequest) *v1.FinalizeBlockResponse); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.FinalizeBlockResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.FinalizeBlockRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Info provides a mock function with given fields: ctx, req
func (_m *Application) Info(ctx context.Context, req *v1.InfoRequest) (*v1.InfoResponse, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for Info")
	}

	var r0 *v1.InfoResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.InfoRequest) (*v1.InfoResponse, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.InfoRequest) *v1.InfoResponse); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.InfoResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.InfoRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InitChain provides a mock function with given fields: ctx, req
func (_m *Application) InitChain(ctx context.Context, req *v1.InitChainRequest) (*v1.InitChainResponse, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for InitChain")
	}

	var r0 *v1.InitChainResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.InitChainRequest) (*v1.InitChainResponse, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.InitChainRequest) *v1.InitChainResponse); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.InitChainResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.InitChainRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListSnapshots provides a mock function with given fields: ctx, req
func (_m *Application) ListSnapshots(ctx context.Context, req *v1.ListSnapshotsRequest) (*v1.ListSnapshotsResponse, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for ListSnapshots")
	}

	var r0 *v1.ListSnapshotsResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ListSnapshotsRequest) (*v1.ListSnapshotsResponse, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ListSnapshotsRequest) *v1.ListSnapshotsResponse); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.ListSnapshotsResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.ListSnapshotsRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LoadSnapshotChunk provides a mock function with given fields: ctx, req
func (_m *Application) LoadSnapshotChunk(ctx context.Context, req *v1.LoadSnapshotChunkRequest) (*v1.LoadSnapshotChunkResponse, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for LoadSnapshotChunk")
	}

	var r0 *v1.LoadSnapshotChunkResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.LoadSnapshotChunkRequest) (*v1.LoadSnapshotChunkResponse, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.LoadSnapshotChunkRequest) *v1.LoadSnapshotChunkResponse); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.LoadSnapshotChunkResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.LoadSnapshotChunkRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OfferSnapshot provides a mock function with given fields: ctx, req
func (_m *Application) OfferSnapshot(ctx context.Context, req *v1.OfferSnapshotRequest) (*v1.OfferSnapshotResponse, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for OfferSnapshot")
	}

	var r0 *v1.OfferSnapshotResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.OfferSnapshotRequest) (*v1.OfferSnapshotResponse, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.OfferSnapshotRequest) *v1.OfferSnapshotResponse); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.OfferSnapshotResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.OfferSnapshotRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PrepareProposal provides a mock function with given fields: ctx, req
func (_m *Application) PrepareProposal(ctx context.Context, req *v1.PrepareProposalRequest) (*v1.PrepareProposalResponse, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for PrepareProposal")
	}

	var r0 *v1.PrepareProposalResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.PrepareProposalRequest) (*v1.PrepareProposalResponse, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.PrepareProposalRequest) *v1.PrepareProposalResponse); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.PrepareProposalResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.PrepareProposalRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProcessProposal provides a mock function with given fields: ctx, req
func (_m *Application) ProcessProposal(ctx context.Context, req *v1.ProcessProposalRequest) (*v1.ProcessProposalResponse, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for ProcessProposal")
	}

	var r0 *v1.ProcessProposalResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProcessProposalRequest) (*v1.ProcessProposalResponse, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.ProcessProposalRequest) *v1.ProcessProposalResponse); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.ProcessProposalResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.ProcessProposalRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Query provides a mock function with given fields: ctx, req
func (_m *Application) Query(ctx context.Context, req *v1.QueryRequest) (*v1.QueryResponse, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for Query")
	}

	var r0 *v1.QueryResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.QueryRequest) (*v1.QueryResponse, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.QueryRequest) *v1.QueryResponse); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.QueryResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.QueryRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// VerifyVoteExtension provides a mock function with given fields: ctx, req
func (_m *Application) VerifyVoteExtension(ctx context.Context, req *v1.VerifyVoteExtensionRequest) (*v1.VerifyVoteExtensionResponse, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for VerifyVoteExtension")
	}

	var r0 *v1.VerifyVoteExtensionResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.VerifyVoteExtensionRequest) (*v1.VerifyVoteExtensionResponse, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.VerifyVoteExtensionRequest) *v1.VerifyVoteExtensionResponse); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.VerifyVoteExtensionResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.VerifyVoteExtensionRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewApplication creates a new instance of Application. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewApplication(t interface {
	mock.TestingT
	Cleanup(func())
}) *Application {
	mock := &Application{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
