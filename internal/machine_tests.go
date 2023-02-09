package internal

import (
	"context"

	v1 "github.com/tinkerbell/pbnj/api/v1"
	v1Client "github.com/tinkerbell/pbnj/client"
	"google.golang.org/grpc"
)

func (t *Tester) powerStatus(ctx context.Context, conn *grpc.ClientConn) ([]byte, error) {
	mc := v1.NewMachineClient(conn)
	//	bc := v1.NewBMCClient(conn)
	tc := v1.NewTaskClient(conn)

	resp, err := v1Client.MachinePower(ctx, mc, tc, &v1.PowerRequest{
		Authn: &v1.Authn{
			Authn: &v1.Authn_DirectAuthn{
				DirectAuthn: &v1.DirectAuthn{
					Host: &v1.Host{
						Host: t.bmcHost,
					},
					Username: t.bmcUser,
					Password: t.bmcPass,
				},
			},
		},
		Vendor: &v1.Vendor{
			Name: "",
		},
		PowerAction: v1.PowerAction_POWER_ACTION_STATUS,
	})
	if err != nil {
		t.logger.Error(err, "error calling")
		return nil, err
	}

	t.logger.V(2).Info("resp", "resp", []interface{}{resp})

	return nil, nil
}

func (t *Tester) powerOn(ctx context.Context, conn *grpc.ClientConn) ([]byte, error) {
	mc := v1.NewMachineClient(conn)
	tc := v1.NewTaskClient(conn)

	resp, err := v1Client.MachinePower(ctx, mc, tc, &v1.PowerRequest{
		Authn: &v1.Authn{
			Authn: &v1.Authn_DirectAuthn{
				DirectAuthn: &v1.DirectAuthn{
					Host: &v1.Host{
						Host: t.bmcHost,
					},
					Username: t.bmcUser,
					Password: t.bmcPass,
				},
			},
		},
		Vendor: &v1.Vendor{
			Name: "",
		},
		PowerAction: v1.PowerAction_POWER_ACTION_ON,
	})
	if err != nil {
		t.logger.Error(err, "error calling")
		return nil, err
	}

	t.logger.V(2).Info("resp", "resp", []interface{}{resp})

	return nil, nil
}

func (t *Tester) powerOff(ctx context.Context, conn *grpc.ClientConn) ([]byte, error) {
	mc := v1.NewMachineClient(conn)
	tc := v1.NewTaskClient(conn)

	resp, err := v1Client.MachinePower(ctx, mc, tc, &v1.PowerRequest{
		Authn: &v1.Authn{
			Authn: &v1.Authn_DirectAuthn{
				DirectAuthn: &v1.DirectAuthn{
					Host: &v1.Host{
						Host: t.bmcHost,
					},
					Username: t.bmcUser,
					Password: t.bmcPass,
				},
			},
		},
		Vendor: &v1.Vendor{
			Name: "",
		},
		PowerAction: v1.PowerAction_POWER_ACTION_OFF,
	})
	if err != nil {
		t.logger.Error(err, "error calling")
		return nil, err
	}

	t.logger.V(2).Info("resp", "resp", []interface{}{resp})

	return nil, nil
}

func (t *Tester) pxeBoot(ctx context.Context, conn *grpc.ClientConn) ([]byte, error) {
	return nil, nil
}

func (t *Tester) powerCycle(ctx context.Context, conn *grpc.ClientConn) ([]byte, error) {
	mc := v1.NewMachineClient(conn)
	tc := v1.NewTaskClient(conn)

	resp, err := v1Client.MachinePower(ctx, mc, tc, &v1.PowerRequest{
		Authn: &v1.Authn{
			Authn: &v1.Authn_DirectAuthn{
				DirectAuthn: &v1.DirectAuthn{
					Host: &v1.Host{
						Host: t.bmcHost,
					},
					Username: t.bmcUser,
					Password: t.bmcPass,
				},
			},
		},
		Vendor: &v1.Vendor{
			Name: "",
		},
		PowerAction: v1.PowerAction_POWER_ACTION_CYCLE,
	})
	if err != nil {
		t.logger.Error(err, "error calling")
		return nil, err
	}

	t.logger.V(2).Info("resp", "resp", []interface{}{resp})
	return nil, nil
}
